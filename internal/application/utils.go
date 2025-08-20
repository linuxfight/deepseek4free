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

func handleNonStreaming(ctx echo.Context, req dto.ChatCompletionRequest, data chan string) error {
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
	return ctx.JSON(http.StatusOK, response)
}

func handleStreaming(ctx echo.Context, req dto.ChatCompletionRequest, data chan string) error {
	ctx.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	ctx.Response().Header().Set("Cache-Control", "no-cache")
	ctx.Response().Header().Set("Connection", "keep-alive")
	ctx.Response().Header().Set("Access-Control-Allow-Origin", "*")
	ctx.Response().WriteHeader(http.StatusOK)

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

		if _, err := fmt.Fprintf(ctx.Response(), "data: %s\n\n", data); err != nil {
			return err
		}
		ctx.Response().Flush()
		time.Sleep(time.Duration(100+rand.Intn(100)) * time.Millisecond)
	}

	if _, err := ctx.Response().Write([]byte("data: [DONE]\n\n")); err != nil {
		return err
	}
	ctx.Response().Flush()

	return nil
}
