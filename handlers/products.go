package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/alisilver78/goWebAPI/dbiface"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/bson"
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

// findProducts finds a list of product
func findProducts(ctx context.Context, collection dbiface.CollectionAPI, q url.Values) ([]Product, error) {
	var products []Product
	filter := make(map[string]interface{})
	for k, v := range q {
		filter[k] = v[0]
	}
	if filter["_id"] != nil {
		docID, err := primitive.ObjectIDFromHex(filter["_id"].(string))
		if err != nil {
			return products, err
		}
		filter["_id"] = docID
	}

	cursor, err := collection.Find(ctx, bson.M(filter))
	if err != nil {
		log.Errorf("Unable to find products: %v", err)
	}
	if err := cursor.All(ctx, &products); err != nil {
		log.Errorf("Unable to read cursor: %v", err)
	}
	return products, nil
}

// findproduct finds a single product
func findProduct(c context.Context, id string, collection dbiface.CollectionAPI) (Product, error) {
	var product Product
	docId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Errorf("Unable to convert id param to objectID: %v", err)
		return product, err
	}
	res := collection.FindOne(c, bson.M{"_id": docId})
	if err := res.Decode(&product); err != nil {
		log.Errorf("Unable to decode result or product not found: %v", err)
		return product, err
	}
	return product, nil
}

// GetProducts gets a list of products
func (h ProductHandler) GetProducts(c echo.Context) error {
	products, err := findProducts(context.Background(), h.Col, c.QueryParams())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, products)
}

func insertProducts(ctx context.Context, products []Product, collection dbiface.CollectionAPI) ([]interface{}, error) {
	var insertedIds []interface{}
	for _, product := range products {
		product.ID = primitive.NewObjectID()
		InsertID, err := collection.InsertOne(context.Background(), product)
		if err != nil {
			log.Errorf("Unable to insert: %v", err)

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
		log.Errorf("Unable to bind: %v", err)
	}
	for _, product := range products {
		if err := c.Validate(product); err != nil {
			log.Errorf("Unable to validate %+v: %v", product, err)
			return err
		}
	}
	IDs, err := insertProducts(context.Background(), products, h.Col)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, IDs)
}

// UpdateProduct updates a product
func (h *ProductHandler) UpdateProduct(c echo.Context) error {
	var product Product
	c.Echo().Validator = &ProductValidator{validator: v}
	product, err := findProduct(context.Background(), c.Param("id"), h.Col)
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}
	//decode request payload
	if err := json.NewDecoder(c.Request().Body).Decode(&product); err != nil {
		log.Errorf("Unable to decode request payload: %v", err)
		return err
	}
	//validating product
	if err := v.Struct(product); err != nil {
		log.Errorf("Unable to validate request payload: %v", err)
		return err
	}
	//updating database
	_, err = h.Col.UpdateOne(context.Background(), bson.M{"_id": c.Param("id")}, bson.M{"$set": product})
	if err != nil {
		log.Errorf("Unable to update database: %v", err)
		return err
	}

	return c.JSON(http.StatusOK, product)
}
