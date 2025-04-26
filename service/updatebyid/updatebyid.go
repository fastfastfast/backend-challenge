package updatebyid

import (
	"context"
	"fmt"
	"time"

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
		return fmt.Errorf("invalid id format: %w", err)
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
		return err
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("no user found with the given id")
	}

	return nil
}
