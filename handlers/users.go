package handlers

import (
	"context"
	"net/http"

	"github.com/alisilver78/goWebAPI/dbiface"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"
)

// User represents a user
type User struct {
	Email    string `json:"username" bson:"username" validate:"required,email"`
	Password string `json:"password" bson:"password" validate:"required,min=8,max=300"`
}

// UsersHndler user handler struct
type UsersHandler struct {
	Col dbiface.CollectionAPI
}

// userValidator validate struct
type userValidator struct {
	validator *validator.Validate
}

func (u *userValidator) Validate(i interface{}) error {
	return u.validator.Struct(i)
}

// InsetUser inserts a user
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

// CreateUser is POST method handler
func (h *UsersHandler) CreateUser(c echo.Context) error {
	var user User
	c.Echo().Validator = &userValidator{validator: v}
	if err := c.Bind(&user); err != nil {
		log.Errorf("Unable to bind to user struct: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Unable to parse the request payload")
	}
	if err := c.Validate(&user); err != nil {
		log.Errorf("Unable to bind to user struct: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Unable to validate request payload")
	}
	insertedID, err := insertUser(context.Background(), user, h.Col)
	if err != nil {
		log.Errorf("Unable to insert user: %v", err)
		return err
	}
	return c.JSON(http.StatusCreated, insertedID)
}
