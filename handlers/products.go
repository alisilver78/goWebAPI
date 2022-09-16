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

// findProducts finds a list of product
func findProducts(ctx context.Context, collection dbiface.CollectionAPI, q url.Values) ([]Product, *echo.HTTPError) {
	var products []Product
	filter := make(map[string]interface{})
	for k, v := range q {
		filter[k] = v[0]
	}
	if filter["_id"] != nil {
		docID, err := primitive.ObjectIDFromHex(filter["_id"].(string))
		if err != nil {
			log.Errorf("Unable to convert string id to objectID: %v", err)
			return products, echo.NewHTTPError(http.StatusBadRequest, "Unable to convert string id to objectID.")
		}
		filter["_id"] = docID
	}

	cursor, err := collection.Find(ctx, bson.M(filter))
	if err != nil {
		log.Errorf("Unable to find products: %v", err)
		return products, echo.NewHTTPError(http.StatusNotFound, "Unable to find products.")
	}
	if err := cursor.All(ctx, &products); err != nil {
		log.Errorf("Unable to read cursor: %v", err)
		return products, echo.NewHTTPError(http.StatusInternalServerError, "Unable to read cursor.")
	}
	return products, nil
}

// findproduct finds a single product
func findProduct(c context.Context, id string, collection dbiface.CollectionAPI) (Product, *echo.HTTPError) {
	var product Product
	docId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Errorf("Unable to convert id param to objectID: %v", err)
		return product, echo.NewHTTPError(http.StatusBadRequest, "Unable to convert id param to objectID.")
	}
	res := collection.FindOne(c, bson.M{"_id": docId})
	if err := res.Decode(&product); err != nil {
		log.Errorf("Unable to decode result or product not found: %v", err)
		return product, echo.NewHTTPError(http.StatusBadRequest, "Unable to decode result or product not found.")
	}
	return product, nil
}

// GetProducts gets list of products
func (h ProductHandler) GetProducts(c echo.Context) error {
	products, err := findProducts(context.Background(), h.Col, c.QueryParams())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, products)
}

// GetProduct gets a single product
func (h *ProductHandler) GetProduct(c echo.Context) error {
	product, err := findProduct(context.Background(), c.Param("id"), h.Col)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, product)
}

// deleteProduct deletes a single product for DeleteProduct
func deleteProduct(c context.Context, id string, collection dbiface.CollectionAPI) (int64, *echo.HTTPError) {
	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Errorf("Unable to convert string id to objectID: %v", err)
		return 0, echo.NewHTTPError(http.StatusBadRequest, "Unable to convert string id to objectID.")
	}
	delCount, err := collection.DeleteOne(context.Background(), bson.M{"_id": docID})
	if err != nil {
		log.Errorf("Unable to delete product: %v", err)
		return 0, echo.NewHTTPError(http.StatusNotFound, "Unable to delete product.")

	}
	return delCount.DeletedCount, nil
}

// DeleteProduct deletes a single product
func (h *ProductHandler) DeleteProduct(c echo.Context) error {
	delCount, err := deleteProduct(context.Background(), c.Param("id"), h.Col)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, delCount)
}

func insertProducts(ctx context.Context, products []Product, collection dbiface.CollectionAPI) ([]interface{}, *echo.HTTPError) {
	var insertedIds []interface{}
	for _, product := range products {
		product.ID = primitive.NewObjectID()
		InsertID, err := collection.InsertOne(context.Background(), product)
		if err != nil {
			log.Errorf("Unable to insert: %v", err)
			return nil, echo.NewHTTPError(http.StatusInternalServerError, "Unable to insert.")
		}
		insertedIds = append(insertedIds, InsertID.InsertedID)
	}
	return insertedIds, nil
}

// CreateProducts creates products
func (h *ProductHandler) CreateProducts(c echo.Context) error {
	var products []Product
	c.Echo().Validator = &productValidator{validator: v}
	if err := c.Bind(&products); err != nil {
		log.Errorf("Unable to bind: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Unable to bind.")

	}
	for _, product := range products {
		if err := c.Validate(product); err != nil {
			log.Errorf("Unable to validate %+v: %v", product, err)
			return echo.NewHTTPError(http.StatusBadRequest, "Unable to validate.")
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
	c.Echo().Validator = &productValidator{validator: v}
	//finding product
	product, err := findProduct(context.Background(), c.Param("id"), h.Col)
	if err != nil {
		log.Errorf("Unable to find the product: %v", err)
		return err
	}
	//decode request payload
	if err := json.NewDecoder(c.Request().Body).Decode(&product); err != nil {
		log.Errorf("Unable to decode request payload: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Unable to decode request payload.")
	}
	//validating product
	if err := c.Validate(product); err != nil {
		log.Errorf("Unable to validate request payload: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Unable to validate request payload.")
	}
	//updating database
	if _, err := h.Col.UpdateOne(context.Background(), bson.M{"_id": c.Param("id")}, bson.M{"$set": product}); err != nil {
		log.Errorf("Unable to update database: %v", err)
		return err
	}

	return c.JSON(http.StatusOK, product)
}
