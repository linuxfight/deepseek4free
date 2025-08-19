package application

import (
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/labstack/echo/v4"
	"github.com/linuxfight/deepseek4free/internal/dto"
	"math/rand"
	"net/http"
	"time"
)

func handleNonStreaming(c echo.Context, req dto.ChatCompletionRequest, data chan string) error {
	answer := ""

	for msg := range data {
		answer = answer + msg
	}

	response := dto.ChatCompletionResponse{
		ID:      fmt.Sprintf("chatcmpl-%d", rand.Int()),
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   req.Model,
		Choices: []dto.Choice{
			{
				Index: 0,
				Message: dto.Message{
					Role:    "assistant",
					Content: answer,
				},
				FinishReason: "stop",
			},
		},
	}
	return c.JSON(http.StatusOK, response)
}

func handleStreaming(c echo.Context, req dto.ChatCompletionRequest, data chan string) error {
	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().WriteHeader(http.StatusOK)

	for msg := range data {
		chunkResp := dto.ChunkResponse{
			ID:      fmt.Sprintf("chatcmpl-%d", rand.Int()),
			Object:  "chat.completion.chunk",
			Created: time.Now().Unix(),
			Model:   req.Model,
			Choices: []dto.ChunkChoice{
				{
					Index: 0,
					Delta: dto.Delta{
						Content: msg,
					},
					FinishReason: nil,
				},
			},
		}

		data, err := sonic.Marshal(chunkResp)
		if err != nil {
			return err
		}

		if _, err := fmt.Fprintf(c.Response(), "data: %s\n\n", data); err != nil {
			return err
		}
		c.Response().Flush()
		time.Sleep(time.Duration(100+rand.Intn(100)) * time.Millisecond)
	}

	if _, err := c.Response().Write([]byte("data: [DONE]\n\n")); err != nil {
		return err
	}
	c.Response().Flush()

	return nil
}
