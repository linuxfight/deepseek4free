package application

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/linuxfight/deepseek4free/internal/dto"
)

func (i *Instance) models(ctx echo.Context) error {
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
	return ctx.JSON(http.StatusOK, models)
}
