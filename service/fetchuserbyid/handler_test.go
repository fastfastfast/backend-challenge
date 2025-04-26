package fetchuserbyid

import (
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

func (m *MockService) FetchUserByID(c echo.Context, db *mongo.Collection, idStr string) (*User, error) {
	args := m.Called(c, db, idStr)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func TestHandlerFetchUserByID(t *testing.T) {
	// Create valid ObjectID for testing
	validID := primitive.NewObjectID()
	validIDStr := validID.Hex()

	// Test cases
	testCases := []struct {
		name           string
		idParam        string
		setupMock      func(*MockService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:    "Success",
			idParam: validIDStr,
			setupMock: func(ms *MockService) {
				user := &User{ID: validID, Name: "Test User"}
				ms.On("FetchUserByID", mock.Anything, mock.Anything, validIDStr).Return(user, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:    "Error from Service",
			idParam: "invalidID",
			setupMock: func(ms *MockService) {
				ms.On("FetchUserByID", mock.Anything, mock.Anything, "invalidID").Return(nil, errors.New("invalid ObjectID"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
		},
		{
			name:           "Empty ID Parameter",
			idParam:        "",
			setupMock:      func(ms *MockService) {},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Set the ID parameter
			if tc.idParam != "" {
				c.SetParamNames("id")
				c.SetParamValues(tc.idParam)
			}

			// Create mock service
			mockService := new(MockService)
			tc.setupMock(mockService)

			// Create handler with mock service
			handler := &UserHandler{
				Service: mockService,
			}

			// Call the handler
			err := handler.HandlerFetchUserByID(c, nil) // Pass nil for DB since it's mocked

			// Assertions
			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedStatus, rec.Code)
			}

			// Verify that all expectations were met
			mockService.AssertExpectations(t)
		})
	}
}
