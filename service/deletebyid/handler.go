package deletebyid

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

type Servicer interface {
	deleteUserByID(collection *mongo.Collection, id string) error
}

type DeleteByIDHandler struct {
	Service Servicer
}

func NewUpdateUserByIDService() Servicer {
	return &updateUserByIDService{}
}

func (h *DeleteByIDHandler) HandlerDeleteByID(c echo.Context, db *mongo.Collection) error {
	idStr := c.Param("id")
	if idStr == "" {
		return echo.NewHTTPError(http.StatusInternalServerError, "id is null")
	}

	err := h.Service.deleteUserByID(db, idStr)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, "")
}
