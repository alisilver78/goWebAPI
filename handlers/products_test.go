package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// integration test
func TestProduct(t *testing.T) {
	var docID string

	t.Run("test create product", func(t *testing.T) {
		var IDs []string
		body := `
		[{
			"product_name":"GoogleTalk",
			"price":79,
			"currency":"USD",
			"vendor":"Google",
			"quantity":150,
			"accessories":["charger", "subscription" ]
		}]`
		req := httptest.NewRequest(http.MethodPost, "/producs", strings.NewReader(body))
		res := httptest.NewRecorder()
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		e := echo.New()
		c := e.NewContext(req, res)
		h.Col = col
		err := h.CreateProducts(c)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, res.Code)

		err = json.Unmarshal(res.Body.Bytes(), &IDs)
		assert.Nil(t, err)
		docID = IDs[0]
		t.Logf("IDs: %#+v", IDs)
		for _, ID := range IDs {
			assert.NotNil(t, ID)
		}
	})

	t.Run("get products", func(t *testing.T) {
		var products []Product
		req := httptest.NewRequest(http.MethodGet, "/producs", nil)
		res := httptest.NewRecorder()
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		e := echo.New()
		c := e.NewContext(req, res)
		h.Col = col
		err := h.GetProducts(c)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.Code)

		err = json.Unmarshal(res.Body.Bytes(), &products)
		assert.Nil(t, err)
		for _, product := range products {
			assert.Equal(t, "GoogleTalk", product.Name)
		}
	})

	t.Run("get products with query parameters", func(t *testing.T) {
		var products []Product
		req := httptest.NewRequest(http.MethodGet, "/producs/?currency=USD&price=79", nil)
		res := httptest.NewRecorder()
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		e := echo.New()
		c := e.NewContext(req, res)
		h.Col = col
		err := h.GetProducts(c)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.Code)

		err = json.Unmarshal(res.Body.Bytes(), &products)
		assert.Nil(t, err)
		for _, product := range products {
			assert.Equal(t, "GoogleTalk", product.Name)
		}
	})

	t.Run("get a product", func(t *testing.T) {
		var product Product
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/products/%s", docID), nil)
		res := httptest.NewRecorder()
		//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		e := echo.New()
		c := e.NewContext(req, res)
		c.SetParamNames("id")
		c.SetParamValues(docID)
		h.Col = col
		err := h.GetProduct(c)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.Code)

		err = json.Unmarshal(res.Body.Bytes(), &product)
		assert.Nil(t, err)
		assert.Equal(t, "USD", product.Currency)
	})

	t.Run("PUT product", func(t *testing.T) {
		var product Product
		body := `
		{
			"product_name":"GoogleTalk",
			"price":109,
			"currency":"USD",
			"vendor":"Google",
			"quantity":150,
			"accessories":["charger", "subscription" ]
		}`
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/producs/%s", docID), strings.NewReader(body))
		res := httptest.NewRecorder()
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		e := echo.New()
		c := e.NewContext(req, res)
		c.SetParamNames("id")
		c.SetParamValues(docID)
		h.Col = col
		err := h.UpdateProduct(c)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.Code)

		err = json.Unmarshal(res.Body.Bytes(), &product)
		assert.Nil(t, err)
		assert.Equal(t, "USD", product.Currency)
	})

	t.Run("delete a product", func(t *testing.T) {
		var deCount int64
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/products/%s", docID), nil)
		res := httptest.NewRecorder()
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		e := echo.New()
		c := e.NewContext(req, res)
		c.SetParamNames("id")
		c.SetParamValues(docID)
		h.Col = col
		err := h.DeleteProduct(c)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.Code)

		err = json.Unmarshal(res.Body.Bytes(), &deCount)
		assert.Nil(t, err)
		assert.Equal(t, int64(1), deCount)
	})
}
