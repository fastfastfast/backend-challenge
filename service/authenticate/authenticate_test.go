package authenticate

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"

	"github.com/labstack/echo/v4"
)

type mockCollection struct {
	mock.Mock
}

func (m *mockCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	args := m.Called(ctx, filter)
	return args.Get(0).(*mongo.SingleResult)
}

// func TestAuthenticate_Success(t *testing.T) {
// 	e := echo.New()
// 	c := e.NewContext(nil, nil)

// 	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
// 	service := &AuthenticateService{}

// 	// hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
// 	// user := User{
// 	// 	Name:     "testuser",
// 	// 	Password: string(hashedPassword),
// 	// }

// 	collection := mt.Coll

// 	resp, err := service.Authenticate(c, collection, "testuser", "password")

//		assert.NoError(t, err)
//		assert.NotNil(t, resp)
//		assert.NotEmpty(t, resp.Token)
//	}
func TestAuthenticate_WrongPassword(t *testing.T) {
	e := echo.New()
	c := e.NewContext(nil, nil)

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	// defer mt.Close()

	service := &AuthenticateService{}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
	userDoc := bson.D{
		{Key: "name", Value: "testuser"},
		{Key: "password", Value: string(hashedPassword)},
	}

	mt.Run("wrong password", func(mt *mtest.T) {
		firstBatch := []bson.D{userDoc}
		cursor := mtest.CreateCursorResponse(1, "db.users", mtest.FirstBatch, firstBatch...)
		mt.AddMockResponses(cursor)

		collection := mt.Coll

		resp, err := service.Authenticate(c, collection, "testuser", "wrongpassword")

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.Equal(t, echo.ErrUnauthorized, err)
	})
}

func TestAuthenticate_UserNotFound(t *testing.T) {
	e := echo.New()
	c := e.NewContext(nil, nil)

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	service := &AuthenticateService{}

	mt.Run("user not found", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "db.users", mtest.FirstBatch)) // No document
		collection := mt.Coll

		resp, err := service.Authenticate(c, collection, "notfound", "password")

		assert.Nil(t, resp)
		assert.Error(t, err)
	})
}
