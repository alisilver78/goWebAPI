package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestUsers(t *testing.T) {
	//sending invalid password (less than eight character)
	t.Run("Test create user invalid data unhappy_senario", func(t *testing.T) {
		body := `
		{
			"username": "test@example.com",
			"password": "abc123"
		}`
		req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(body))
		res := httptest.NewRecorder()
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		e := echo.New()
		c := e.NewContext(req, res)
		uh.Col = usersCol
		err := uh.CreateUser(c)
		t.Logf("res: %#+v\n", res.Body.String())
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, res.Code)

	})
}
