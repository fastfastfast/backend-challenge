package fetchuserbyid

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

type Servicer interface {
	FetchUserByID(c echo.Context, db *mongo.Collection, idStr string) (*User, error)
}

type UserHandler struct {
	Service Servicer
}

func NewUserService() Servicer {
	return &userService{}
}

func (h *UserHandler) HandlerFetchUserByID(c echo.Context, db *mongo.Collection) error {
	idStr := c.Param("id")
	if idStr == "" {
		return echo.NewHTTPError(http.StatusInternalServerError, "id is null")
	}

	createdUser, err := h.Service.FetchUserByID(c, db, idStr)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, createdUser)
}
