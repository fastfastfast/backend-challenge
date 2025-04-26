package listalluser

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// MockService is a mock implementation of the Servicer interface
type MockService struct {
	mock.Mock
}

// ListAllUser mocks the ListAllUser method
func (m *MockService) ListAllUser(c echo.Context, db *mongo.Collection) (*[]User, error) {
	args := m.Called(c, db)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*[]User), args.Error(1)
}

func TestHandlerListAllUser_Success(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService := new(MockService)
	handler := &ListAllUserHandler{Service: mockService}

	// Mock database collection - we're not actually using it in the test
	var db *mongo.Collection

	// Mock data
	objID := primitive.NewObjectID()
	objIDs := primitive.NewObjectID()

	mockUsers := &[]User{
		{ID: objID, Name: "user1", Email: "user1@example.com"},
		{ID: objIDs, Name: "user2", Email: "user2@example.com"},
	}

	// Expectations
	mockService.On("ListAllUser", c, db).Return(mockUsers, nil)

	// Test
	err := handler.HandlerListAllUser(c, db)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Verify response body
	var responseUsers []User
	err = json.Unmarshal(rec.Body.Bytes(), &responseUsers)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(responseUsers))
	assert.Equal(t, "user1", responseUsers[0].Name)
	assert.Equal(t, "user2", responseUsers[1].Name)

	// Verify expectations
	mockService.AssertExpectations(t)
}

func TestHandlerListAllUser_Error(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService := new(MockService)
	handler := &ListAllUserHandler{Service: mockService}

	// Mock database collection - we're not actually using it in the test
	var db *mongo.Collection

	// Expectations - service returns an error
	expectedError := errors.New("database error")
	mockService.On("ListAllUser", c, db).Return(nil, expectedError)

	// Test
	err := handler.HandlerListAllUser(c, db)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, expectedError, err) // The error should be returned unchanged

	// Verify expectations
	mockService.AssertExpectations(t)
}

func TestNewListAllUserService(t *testing.T) {
	// Test that NewListAllUserService returns an implementation of Servicer
	service := NewListAllUserService()
	assert.NotNil(t, service)
	assert.IsType(t, &listAllUserService{}, service)
}
