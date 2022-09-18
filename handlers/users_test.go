package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestUsers(t *testing.T) {

	//invalid password(less than eight character) (unhappy_senario)
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
		con := e.NewContext(req, res)
		uh.Col = usersCol
		err := uh.CreateUser(con)
		// t.Logf("res: %#+v\n", res.Body.String())
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, res.Code)

	})

	//valid user
	t.Run("Test create user", func(t *testing.T) {
		var user User
		body := `
		{
			"username": "test@example.com",
			"password": "abc12345"
		}`
		req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(body))
		res := httptest.NewRecorder()
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		e := echo.New()
		con := e.NewContext(req, res)
		uh.Col = usersCol
		err := uh.CreateUser(con)
		//t.Logf("res: %#+v\n", res.Body.String())
		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, res.Code)
		token := res.Header().Get("X-Auth-Token")
		//token should not be empty
		assert.NotEmpty(t, token)
		err = json.Unmarshal(res.Body.Bytes(), &user)
		assert.Nil(t, err)
		assert.Equal(t, "test@example.com", user.Email)
		//Password field must be empty in response
		assert.Empty(t, user.Password)
	})

	//recreating a user (unhappy_senario)
	t.Run("recreating a user unhappy_senario", func(t *testing.T) {
		body := `{
			"username": "test@example.com",
			"password": "abc12345"
		}`
		req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(body))
		res := httptest.NewRecorder()
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		e := echo.New()
		con := e.NewContext(req, res)
		uh.Col = usersCol
		err := uh.CreateUser(con)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusBadRequest, err)
	})

	//testing authenticate
	t.Run("test authentication", func(t *testing.T) {
		var user User
		body := `
		{
			"username": "test@example.com",
			"password": "abc12345"
		}`
		req := httptest.NewRequest(http.MethodPost, "/auth", strings.NewReader(body))
		res := httptest.NewRecorder()
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		e := echo.New()
		con := e.NewContext(req, res)
		uh.Col = usersCol
		err := uh.AthnUser(con)
		// t.Logf("res: %#+v\n", res.Body.String())
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.Code)
		token := res.Header().Get("X-Auth-Token")
		//token should not be empty
		assert.NotEmpty(t, token)
		err = json.Unmarshal(res.Body.Bytes(), &user)
		assert.Nil(t, err)
		assert.Equal(t, "test@example.com", user.Email)
		//Password field must be empty in response
		assert.Empty(t, user.Password)
	})

}
