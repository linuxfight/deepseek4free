package application

import (
	"github.com/labstack/echo/v4"
	"github.com/linuxfight/deepseek4free/internal/dto"
	"github.com/linuxfight/deepseek4free/pkg/api"
	"net/http"
	"strconv"
	"strings"
)

func (i *Instance) chat(ctx echo.Context) error {
	apiKey := ctx.Request().Header.Get("Authorization")
	if apiKey == "" {
		return ctx.JSON(http.StatusUnauthorized, map[string]interface{}{"error": "Authorization header required"})
	}

	apiKey = strings.Replace(apiKey, "Bearer ", "", 1)

	var req dto.ChatCompletionRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	apiClient := api.New(i.solver, apiKey)

	chatData, err := i.cache.GetChatData(ctx.Request().Context(), apiKey, req.Messages[0].Content)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
	}

	if chatData.ChatId == "" {
		chatData.ChatId, err = apiClient.CreateChat()
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
	} else {
		if chatData.CurrentMessageId == "" || chatData.CurrentMessageId == "null" {
			chatData.CurrentMessageId = "0"
		}
		msgId, err := strconv.Atoi(chatData.CurrentMessageId)
		if err != nil {
			panic(err)
		}
		msgId += 2
		chatData.CurrentMessageId = strconv.Itoa(msgId)
	}

	defer func(apiClient *api.Client, chatSessionId string) {
		text := req.Messages[0].Content
		// This stuff is here, because chat title requests cannot be handled with markdown, fuck this
		/*
			if strings.HasPrefix(text, "### Task:\nGenerate a concise, 3-5 word") {
					text = "title_req_" + strconv.FormatInt(time.Now().Unix(), 10)
				} else if strings.HasPrefix(text, "### Task:\nSuggest 3-5") {
					text = "follow_req_" + strconv.FormatInt(time.Now().Unix(), 10)
				} else if strings.HasPrefix(text, "### Task:\nGenerate 1-3 broad") {
					text = "tags_req_" + strconv.FormatInt(time.Now().Unix(), 10)
				}
		*/

		err := apiClient.ChangeTitle(chatSessionId, text)
		if err != nil {
			panic(err)
		}

		err = i.cache.SetChatData(ctx.Request().Context(), apiKey, req.Messages[0].Content, chatData)
		if err != nil {
			panic(err)
		}
	}(apiClient, chatData.ChatId)

	resp := make(chan string)

	err = apiClient.Completion(chatData.ChatId, chatData.CurrentMessageId, req.Messages[len(req.Messages)-1].Content, true, true, resp)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if req.Stream {
		return handleStreaming(ctx, req, resp)
	} else {
		return handleNonStreaming(ctx, req, resp)
	}
}
