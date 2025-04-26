package authenticate

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
)

// Mock Service
type mockService struct{}

func (m *mockService) Authenticate(c echo.Context, db *mongo.Collection, user, password string) (*Response, error) {
	if user == "validuser" && password == "validpassword" {
		return &Response{Token: "valid-token"}, nil
	}
	return nil, errors.New("authentication failed")
}

func TestHandlerAuthenticate_Success(t *testing.T) {
	e := echo.New()
	mockSvc := &mockService{}
	handler := &AuthenticateHandler{Service: mockSvc}

	// สร้าง request payload
	user := User{
		Name:     "validuser",
		Password: "validpassword",
	}
	body, _ := json.Marshal(user)
	req := httptest.NewRequest(http.MethodPost, "/authenticate", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// เรียก function ที่จะ test
	err := handler.HandlerAuthenticate(c, nil)

	// ตรวจสอบผลลัพธ์
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response Response
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, "valid-token", response.Token)
}

func TestHandlerAuthenticate_BindError(t *testing.T) {
	e := echo.New()
	mockSvc := &mockService{}
	handler := &AuthenticateHandler{Service: mockSvc}

	req := httptest.NewRequest(http.MethodPost, "/authenticate", bytes.NewReader([]byte(`invalid-json`)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.HandlerAuthenticate(c, nil)

	assert.Error(t, err)
	httpError, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusInternalServerError, httpError.Code)
}

// func TestHandlerAuthenticate_ValidationError(t *testing.T) {
// 	e := echo.New()
// 	mockSvc := &mockService{}
// 	handler := &AuthenticateHandler{Service: mockSvc}

// 	// missing email (required)
// 	user := User{
// 		Name:     "validuser",
// 		Password: "validpassword",
// 	}
// 	body, _ := json.Marshal(user)
// 	req := httptest.NewRequest(http.MethodPost, "/authenticate", bytes.NewReader(body))
// 	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

// 	rec := httptest.NewRecorder()
// 	c := e.NewContext(req, rec)

// 	err := handler.HandlerAuthenticate(c, nil)

// 	assert.Error(t, err)
// 	httpError, ok := err.(*echo.HTTPError)
// 	assert.True(t, ok)
// 	assert.Equal(t, http.StatusInternalServerError, httpError.Code)
// }

func TestHandlerAuthenticate_AuthFail(t *testing.T) {
	e := echo.New()
	mockSvc := &mockService{}
	handler := &AuthenticateHandler{Service: mockSvc}

	// user or password invalid
	user := User{
		Name:     "wronguser",
		Password: "wrongpassword",
	}
	body, _ := json.Marshal(user)
	req := httptest.NewRequest(http.MethodPost, "/authenticate", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.HandlerAuthenticate(c, nil)

	assert.Error(t, err)
	assert.Equal(t, "authentication failed", err.Error())
}
