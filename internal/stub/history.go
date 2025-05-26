package stub

import (
	"context"
	"fmt"
	"github.com/linuxfight/deepseek4free/internal/stub/gen"
)

func (stub *Stub) GetChatHistory(_ context.Context, chat *gen.Chat) (*gen.ChatHistoryResponse, error) {
	apiHistory, err := stub.api.GetMessageHistory(chat.Id)
	if err != nil {
		return nil, err
	}

	var messages []*gen.ChatMessage
	for _, apiMessage := range apiHistory.ChatMessages {
		content := ""

		if apiMessage.SearchEnabled {
			content += "<searching>"
			for _, result := range apiMessage.SearchResults {
				link := ""

				if result.Title != "" && result.Url != "" {
					if result.Snippet == "" {
						link = fmt.Sprintf("[%s](%s)\n",
							result.Title,
							result.Url)
					} else {
						link = fmt.Sprintf("[%s](%s) - %s\n",
							result.Title,
							result.Url,
							result.Snippet)
					}
				} else {
					continue
				}

				content += link
			}
			content += "</searching>" + "\n"
		}

		if apiMessage.ThinkingEnabled {
			content += "<thinking>"
			content += *apiMessage.ThinkingContent
			content += "</thinking>" + "\n"
		}

		content += "<answer>"
		content += apiMessage.Content
		content += "</answer>"

		messages = append(messages, &gen.ChatMessage{
			Id:   int32(apiMessage.MessageId),
			Role: apiMessage.Role,
		})
	}

	chatInfo := &gen.Chat{
		Id:               apiHistory.ChatSession.Id,
		CurrentMessageId: apiHistory.ChatSession.CurrentMessageId,
		Title:            apiHistory.ChatSession.Title,
	}

	response := &gen.ChatHistoryResponse{
		Chat:     chatInfo,
		Messages: messages,
	}

	return response, nil
}
