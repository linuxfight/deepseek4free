package api

import (
	"fmt"
	"hashTest/api/models"
	"net/http"
)

// CreateChat is a method to create a new chat session. Returns UUID of a new chat session.
func (c *Client) CreateChat() (string, error) {
	var data models.ChatCreateResponse
	if err := c.execute(chatCreateUrl, chatCreateBody, http.MethodPost, &data); err != nil {
		return "", err
	}

	return data.Data.BizData.Id, nil
}

func (c *Client) GetAllChats() ([]models.ChatSession, error) {
	var data models.ChatListResponse
	if err := c.execute(chatListUrl, "", http.MethodGet, &data); err != nil {
		return []models.ChatSession{}, err
	}
	return data.Data.BizData.ChatSessions, nil
}

func (c *Client) ChangeTitle(chatSessionId string, title string) error {
	body := fmt.Sprintf(`{"chat_session_id":"%s","title":"%s"}`, chatSessionId, title)
	var data models.ChatEditResponse
	if err := c.execute(chatEditUrl, body, http.MethodPost, &data); err != nil {
		return err
	}
	return nil
}

func (c *Client) DeleteChatSession(chatSessionId string) error {
	body := fmt.Sprintf(`{"chat_session_id":"%s"}`, chatSessionId)
	var data models.NullResponse
	if err := c.execute(chatEditUrl, body, http.MethodPost, &data); err != nil {
		return err
	}
	return nil
}

func (c *Client) GetMessageHistory(chatSessionId string) (models.ChatHistory, error) {
	var data models.ChatHistoryResponse
	if err := c.execute(historyBaseUrl+chatSessionId, "", http.MethodGet, &data); err != nil {
		return models.ChatHistory{}, err
	}

	return data.Data.BizData, nil
}
