package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/alisilver78/goWebAPI/config"
	"github.com/alisilver78/goWebAPI/handlers"
	"github.com/golang-jwt/jwt/v4"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/labstack/gommon/random"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	CorrelationID = "X-Correlation-ID"
)

var (
	c        *mongo.Client
	db       *mongo.Database
	prodcol  *mongo.Collection
	userscol *mongo.Collection
	cfg      config.Properties
)

func init() {
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("Configuration cannot be read: %v", err)
	}
	connectURI := fmt.Sprintf("mongodb://%s:%s", cfg.DBHost, cfg.DPPort)
	var err error
	c, err = mongo.Connect(context.Background(), options.Client().ApplyURI(connectURI))
	if err != nil {
		log.Fatalf("Unable connect to database: %v", err)
	}

	db = c.Database(cfg.DBName)
	prodcol = db.Collection(cfg.ProductsCollection)
	userscol = db.Collection(cfg.UsersCollection)

	IsUserIndexUnique := true
	indexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "username", Value: 1}},
		Options: &options.IndexOptions{
			Unique: &IsUserIndexUnique,
		},
	}

	_, err = userscol.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		log.Fatalf("Unable to create an index: %v", err)
	}
}

func addCorrelationID(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		//generate correlation id
		var newID string
		id := c.Request().Header.Get(CorrelationID)
		if id == "" {
			//generate id
			newID = random.String(12)
		} else {
			//assigen id
			newID = id
		}
		c.Request().Header.Set(CorrelationID, newID)
		c.Response().Header().Set(CorrelationID, newID)

		return next(c)
	}
}

// adminMiddleware checks is user admin or not
func adminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("x-auth-token")
		claims := jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(token, claims, func(*jwt.Token) (interface{}, error) {
			return []byte(cfg.JwtTokenSecret), nil
		})
		if err != nil {
			log.Errorf("Unable to parse jwt token with claims: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Unable to parse jwt token with claims")
		}
		boolClaims := claims["authorized"].(bool)
		if !boolClaims {
			return echo.NewHTTPError(http.StatusForbidden, "This user do not have required permissions")
		}
		return next(c)
	}
}

func main() {
	e := echo.New()
	h := &handlers.ProductHandler{Col: prodcol}
	uh := &handlers.UsersHandler{Col: userscol}
	jwtMiddleware := middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:  []byte(cfg.JwtTokenSecret),
		TokenLookup: "header:x-auth-token",
	})

	e.Logger.SetLevel(log.ERROR)
	e.Pre(middleware.RemoveTrailingSlash())
	e.Pre(addCorrelationID)
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{Format: config.LoggerConfigFormat}))
	e.GET("/products", h.GetProducts)
	e.GET("/products/:id", h.GetProduct)

	e.DELETE("/products/:id", h.DeleteProduct, jwtMiddleware, adminMiddleware)
	e.POST("/products", h.CreateProducts, middleware.BodyLimit("1M"), jwtMiddleware)
	e.PUT("/products/:id", h.UpdateProduct, middleware.BodyLimit("1M"), jwtMiddleware)

	e.POST("/users", uh.CreateUser, middleware.BodyLimit("1M"))
	e.POST("/auth", uh.AthnUser)
	e.Logger.Infof("listening on: %s:%s", cfg.Host, cfg.Port)
	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)))
}
