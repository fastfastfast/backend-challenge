package listalluser

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

type Servicer interface {
	ListAllUser(c echo.Context, db *mongo.Collection) (*[]User, error)
}

type ListAllUserHandler struct {
	Service Servicer
}

func NewListAllUserService() Servicer {
	return &listAllUserService{}
}

func (h *ListAllUserHandler) HandlerListAllUser(c echo.Context, db *mongo.Collection) error {

	createdUser, err := h.Service.ListAllUser(c, db)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, createdUser)
}
