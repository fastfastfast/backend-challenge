package authenticate

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

type Servicer interface {
	Authenticate(c echo.Context, db *mongo.Collection, user, password string) (*Response, error)
}

type AuthenticateHandler struct {
	Service Servicer
}

func NewAuthenticateService() Servicer {
	return &AuthenticateService{}
}

func (h *AuthenticateHandler) HandlerAuthenticate(c echo.Context, db *mongo.Collection) error {
	var user User
	if err := c.Bind(&user); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Bind json Error: %s", err.Error()))
	}

	validate := validator.New()
	err := validate.Struct(user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("validate Error: %s", err.Error()))
	}

	token, err := h.Service.Authenticate(c, db, user.Name, user.Password)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, token)
}
