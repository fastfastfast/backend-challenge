package insert

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

type Servicer interface {
	Insert(c echo.Context, db *mongo.Collection, user User) error
}

type InsertHandler struct {
	Service Servicer
}

func NewInsertService() Servicer {
	return &insertService{}
}

func (h *InsertHandler) HandlerInsert(c echo.Context, db *mongo.Collection) error {
	var user User
	if err := c.Bind(&user); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Bind json Error: %s", err.Error()))
	}

	validate := validator.New()
	err := validate.Struct(user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("validate Error: %s", err.Error()))
	}

	err = h.Service.Insert(c, db, user)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, "")
}
