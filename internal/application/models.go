package application

import (
	"github.com/labstack/echo/v4"
	"github.com/linuxfight/deepseek4free/internal/dto"
	"net/http"
)

func (i *Instance) models(c echo.Context) error {
	models := map[string]interface{}{
		"object": "list",
		"data": []dto.Model{
			{
				ID:      "r1",
				Object:  "model",
				OwnedBy: "deepseek",
			},
		},
	}
	return c.JSON(http.StatusOK, models)
}
