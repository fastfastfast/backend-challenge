package deletebyid

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
)

// Mock service
type MockService struct {
	mock.Mock
}

func (m *MockService) deleteUserByID(collection *mongo.Collection, id string) error {
	args := m.Called(collection, id)
	return args.Error(0)
}

func TestHandlerDeleteByID_Success(t *testing.T) {
	e := echo.New()

	var db *mongo.Collection

	mockService := new(MockService)
	handler := DeleteByIDHandler{Service: mockService}

	req := httptest.NewRequest(http.MethodDelete, "/users/12345", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("12345")

	mockService.On("deleteUserByID", db, "12345").Return(nil)

	err := handler.HandlerDeleteByID(c, db)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "\"\"\n", rec.Body.String())

	mockService.AssertExpectations(t)
}

func TestHandlerDeleteByID_MissingID(t *testing.T) {
	e := echo.New()

	var db *mongo.Collection

	mockService := new(MockService)
	handler := DeleteByIDHandler{Service: mockService}

	req := httptest.NewRequest(http.MethodDelete, "/users/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.HandlerDeleteByID(c, db)

	if httpError, ok := err.(*echo.HTTPError); ok {
		assert.Equal(t, http.StatusInternalServerError, httpError.Code)
		assert.Equal(t, "id is null", httpError.Message)
	} else {
		t.Errorf("expected HTTPError, got %v", err)
	}
}

func TestHandlerDeleteByID_DeleteError(t *testing.T) {
	e := echo.New()

	var db *mongo.Collection

	mockService := new(MockService)
	handler := DeleteByIDHandler{Service: mockService}

	req := httptest.NewRequest(http.MethodDelete, "/users/6789", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("6789")

	mockService.On("deleteUserByID", db, "6789").Return(errors.New("delete failed"))

	err := handler.HandlerDeleteByID(c, db)

	assert.EqualError(t, err, "delete failed")

	mockService.AssertExpectations(t)
}
