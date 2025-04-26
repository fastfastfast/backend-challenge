package authenticate

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthenticateService struct{}

func (s *AuthenticateService) Authenticate(c echo.Context, db *mongo.Collection, name, password string) (*Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var userDB User
	err := db.FindOne(ctx, bson.M{"name": name}).Decode(&userDB)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Find Error: %s", err.Error()))
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userDB.Password), []byte(password)); err != nil {
		return nil, echo.ErrUnauthorized
	}

	claims := &jwtCustomClaims{
		name,
		true,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return nil, err
	}

	return &Response{
		Token: t,
	}, nil
}
