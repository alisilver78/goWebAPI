package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/alisilver78/goWebAPI/config"
	"github.com/alisilver78/goWebAPI/dbiface"
	"github.com/golang-jwt/jwt/v4"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var (
	props config.Properties
)

// errorMassage is an error struct
type errorMassage struct {
	Message string `json:"message"`
}

// User represents a user
type User struct {
	Email    string `json:"username" bson:"username" validate:"required,email"`
	Password string `json:"password,omitempty" bson:"password" validate:"required,min=8,max=300"`
	IsAdmin  bool   `json:"isadmin,omitempty" bson:"isadmin"`
}

// UsersHndler user handler struct
type UsersHandler struct {
	Col dbiface.CollectionAPI
}

// InsertUser inserts a user
func insertUser(ctx context.Context, user User, collection dbiface.CollectionAPI) (interface{}, *echo.HTTPError) {
	var newuser User
	res := collection.FindOne(ctx, bson.M{"username": user.Email})
	if err := res.Decode(&newuser); err != nil && err != mongo.ErrNoDocuments {
		log.Errorf("Unable to decode retrived user: %v", err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Unable to decode retrived user")
	}
	if newuser.Email != "" {
		log.Errorf("User by %v already exists", newuser.Email)
		return nil, echo.NewHTTPError(http.StatusBadRequest, "User already exists")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	if err != nil {
		log.Errorf("Unable to process the error: %v", err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Unable to process the error")
	}
	user.Password = string(hashedPassword)
	insertRes, err := collection.InsertOne(ctx, user)
	if err != nil {
		log.Errorf("Unable to create user: %v", err)
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Unable to create user")
	}
	return insertRes.InsertedID, nil
}

// CreateUser is handler POST method for user endpoint
func (h *UsersHandler) CreateUser(c echo.Context) error {
	var user User
	c.Echo().Validator = &userValidator{validator: v}
	if err := c.Bind(&user); err != nil {
		log.Errorf("Unable to bind to user struct: %v", err)
		return c.JSON(http.StatusBadRequest, errorMassage{Message: "Unable to parse the request payload"})
	}
	if err := c.Validate(&user); err != nil {
		log.Errorf("Unable to validate user struct: %v", err)
		return c.JSON(http.StatusBadRequest, errorMassage{Message: "Unable to validate request payload"})
	}
	_, err := insertUser(context.Background(), user, h.Col)
	if err != nil {
		log.Errorf("Unable to insert user: %v", err)
		return c.JSON(err.Code, errorMassage{Message: "Unable to insert user"})
	}
	token, err := user.createToken()
	if err != nil {
		log.Errorf("Unable to create token: %v", err)
		return c.JSON(err.Code, errorMassage{Message: "Unable to create a token"})
	}
	c.Response().Header().Set("x-auth-token", token)

	return c.JSON(http.StatusCreated, user.Email)
}

func isCredValid(su, ru string) error {
	err := bcrypt.CompareHashAndPassword([]byte(su), []byte(ru))
	//if did not match return error
	if err != nil {
		log.Errorf("Password did not match: %v", err)
		return err
	}
	return nil
}

// authenticateUser athenticates a user
func authenticateUser(ctx context.Context, user User, collection dbiface.CollectionAPI) (User, *echo.HTTPError) {
	var storedUser User
	res := collection.FindOne(ctx, bson.M{"username": user.Email})
	err := res.Decode(&storedUser)
	if err != nil && err != mongo.ErrNoDocuments {
		log.Errorf("Unable to decode retrived user: %v", err)
		return storedUser, echo.NewHTTPError(http.StatusUnprocessableEntity, "Unable to decode retrived user. ")
	}
	if err == mongo.ErrNoDocuments {
		log.Errorf("User does not exist: %v", user.Email)
		return storedUser, echo.NewHTTPError(http.StatusNotFound, "User does not exist. ")
	}
	if err := isCredValid(storedUser.Password, user.Password); err != nil {
		log.Errorf("Credendtials not valid: %v", err)
		return storedUser, echo.NewHTTPError(http.StatusUnauthorized, "Credendtials not valid")
	}
	return storedUser, nil
}

// createToken creates a jwt token
func (u User) createToken() (string, *echo.HTTPError) {
	if err := cleanenv.ReadEnv(&props); err != nil {
		log.Fatalf("Unable to read configuration: %v", err)
	}
	claims := jwt.MapClaims{}
	claims["authorized"] = u.IsAdmin
	claims["user_id"] = u.Email
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := at.SignedString([]byte(props.JwtTokenSecret))
	if err != nil {
		log.Errorf("Unable to generate token: %v", err)
		return "", echo.NewHTTPError(http.StatusInternalServerError, "Unable to generate token. ")
	}
	return token, nil
}

// AuthnUser is handler of /authn endpoint
func (h *UsersHandler) AthnUser(c echo.Context) error {
	var user User
	c.Echo().Validator = &userValidator{validator: v}
	if err := c.Bind(&user); err != nil {
		log.Errorf("Unable to bind request payload. ")
		return c.JSON(http.StatusUnprocessableEntity, errorMassage{Message: "Unable to bind request payload"})
	}
	if err := c.Validate(user); err != nil {
		log.Errorf("Unable to validate request payload. ")
		return c.JSON(http.StatusBadRequest, errorMassage{Message: "Unable to validate request payload"})
	}
	user, err := authenticateUser(context.Background(), user, h.Col)
	if err != nil {
		log.Errorf("Unable to athenticate user: %v", err)
		return c.JSON(err.Code, errorMassage{Message: "Unable to athenticate user"})
	}
	token, err := user.createToken()
	if err != nil {
		log.Errorf("Unable to create token: %v", err)
		return c.JSON(http.StatusInternalServerError, errorMassage{Message: "Unable to create a token"})
	}
	c.Response().Header().Set("x-auth-token", token)

	return c.JSON(http.StatusOK, User{Email: user.Email})
}
