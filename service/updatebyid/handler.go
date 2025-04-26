package updatebyid

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

type Servicer interface {
	updateUserByID(collection *mongo.Collection, id, newName, newEmail string) error
}

type UpdateUserByIDHandler struct {
	Service Servicer
}

func NewUpdateUserByIDService() Servicer {
	return &updateUserByIDService{}
}

func (h *UpdateUserByIDHandler) HandlerUpdateUserByID(c echo.Context, db *mongo.Collection) error {
	var user User
	if err := c.Bind(&user); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Bind json Error: %s", err.Error()))
	}

	validate := validator.New()
	err := validate.Struct(user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("validate Error: %s", err.Error()))
	}

	err = h.Service.updateUserByID(db, user.ID, user.Name, user.Email)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, "")
}
