package updatebyid

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

func (s *updateUserByIDService) updateUserByID(collection *mongo.Collection, id string, newName, newEmail string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("invalid id format: %s", err.Error()))
	}

	update := bson.M{
		"$set": bson.M{},
	}
	if newName != "" {
		update["$set"].(bson.M)["name"] = newName
	}
	if newEmail != "" {
		update["$set"].(bson.M)["email"] = newEmail
	}

	filter := bson.M{"_id": objID}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("err when update: %s", err.Error()))
	}
	if result.MatchedCount == 0 {
		return echo.NewHTTPError(http.StatusInternalServerError, "no user found with the given id")
	}

	return nil
}
