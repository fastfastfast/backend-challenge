package deletebyid

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestDeleteUserByID(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	service := &updateUserByIDService{}
	mt.Run("delete success", func(mt *mtest.T) {
		// Arrange
		collection := mt.Coll
		objID := primitive.NewObjectID()

		// Mock success deleteOne response, บอกลบได้ 1 อัน
		mt.AddMockResponses(
			mtest.CreateSuccessResponse(
				bson.E{Key: "n", Value: 1},
			),
		)

		err := service.deleteUserByID(collection, objID.Hex())

		assert.NoError(t, err)
	})

	mt.Run("invalid id format", func(mt *mtest.T) {
		collection := mt.Coll

		// Act
		err := service.deleteUserByID(collection, "invalid-id")

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid id format")
	})

	mt.Run("user not found", func(mt *mtest.T) {
		collection := mt.Coll
		objID := primitive.NewObjectID()

		// Mock the delete response (DeletedCount = 0)
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		// Act
		err := service.deleteUserByID(collection, objID.Hex())

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no user found with the given id")
	})
}
