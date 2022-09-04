package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/alisilver78/goWebAPI/dbiface"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/go-playground/validator.v9"
)

var (
	v = validator.New()
)

// product define an electronic product
type Product struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"product_name" bson:"product_name"`
	Price       int                `json:"price" bson:"price"`
	Currency    string             `json:"currency" bson:"currency"`
	Quantity    int                `json:"quantity" bson:"quantity"`
	Discount    int                `json:"discount,omitempty" bson:"discount,omitempty"`
	Vendor      string             `json:"vendor" bson:"vendor"`
	Accessories []string           `json:"accessories,omitempty" bson:"accessories,omitempty"`
}

// ProductHandler is a product handler
type ProductHandler struct {
	Col dbiface.CollectionAPI
}

type ProductValidator struct {
	validator *validator.Validate
}

func (p *ProductValidator) Validate(i interface{}) error {
	return p.validator.Struct(i)
}

func insertProducts(ctx context.Context, products []Product, collection dbiface.CollectionAPI) ([]interface{}, error) {
	var insertedIds []interface{}
	for _, product := range products {
		product.ID = primitive.NewObjectID()
		InsertID, err := collection.InsertOne(context.Background(), product)
		if err != nil {
			log.Printf("Unable to insert: %v", err)
			return nil, err
		}
		insertedIds = append(insertedIds, InsertID.InsertedID)
	}
	return insertedIds, nil
}

// CreateProducts creates products
func (h *ProductHandler) CreateProducts(c echo.Context) error {
	var products []Product
	c.Echo().Validator = &ProductValidator{validator: v}
	if err := c.Bind(&products); err != nil {
		log.Printf("Unable to bind: %v", err)
	}
	for _, product := range products {
		if err := c.Validate(product); err != nil {
			log.Printf("Unable to validate %+v: %v", product, err)
			return err
		}
	}
	IDs, err := insertProducts(context.Background(), products, h.Col)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, IDs)
}
