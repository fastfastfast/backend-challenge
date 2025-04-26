package fetchuserbyid

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type userService struct{}

func (s *userService) FetchUserByID(c echo.Context, db *mongo.Collection, idStr string) (*User, error) {

	objID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("invalid id format: %s", err.Error()))

	}

	var user User
	err = db.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("err when find by id: %s", err.Error()))
		}
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "server error")

	}

	return &user, nil

}
