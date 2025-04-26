package deletebyid

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type updateUserByIDService struct{}

func (s *updateUserByIDService) deleteUserByID(collection *mongo.Collection, id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("invalid id format: %s", err.Error()))
	}

	filter := bson.M{"_id": objID}
	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return echo.NewHTTPError(http.StatusInternalServerError, "no user found with the given id")
	}

	return nil
}
