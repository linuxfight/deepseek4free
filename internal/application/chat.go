package application

import (
	"github.com/labstack/echo/v4"
	"github.com/linuxfight/deepseek4free/internal/dto"
	"github.com/linuxfight/deepseek4free/pkg/api"
	"net/http"
	"strings"
)

func (i *Instance) chat(c echo.Context) error {
	apiKey := c.Request().Header.Get("Authorization")
	if apiKey == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": "Authorization header required"})
	}

	apiKey = strings.Replace(apiKey, "Bearer ", "", 1)

	var req dto.ChatCompletionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	apiClient := api.New(i.solver, apiKey)

	chat, err := apiClient.CreateChat()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	err = apiClient.ChangeTitle(chat, "ApiChat"+chat)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	defer func(apiClient *api.Client, chatSessionId string) {
		err := apiClient.DeleteChatSession(chatSessionId)
		if err != nil {
			panic(err)
		}
	}(apiClient, chat)

	resp := make(chan string)

	err = apiClient.Completion(chat, "", req.Messages[len(req.Messages)-1].Content, true, true, resp)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if req.Stream {
		return handleStreaming(c, req, resp)
	} else {
		return handleNonStreaming(c, req, resp)
	}
}
