package insert

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type insertService struct{}

func (s *insertService) Insert(c echo.Context, db *mongo.Collection, user User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user.CreatedAt = time.Now()
	HashPassword, err := HashPassword(user.Password)
	if err != nil {
		log.Fatal(err)
	}
	user.Password = HashPassword

	_, err = db.InsertOne(ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return echo.NewHTTPError(http.StatusInternalServerError, "email already exists")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Insert Error: %s", err.Error()))
	}

	return err
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
