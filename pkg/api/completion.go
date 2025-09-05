package api

import (
	"bytes"
	"io"
	"net/http"
	"strconv"

	"github.com/bytedance/sonic"
	"github.com/linuxfight/deepseek4free/pkg/api/models"
)

func (c *Client) Completion(chatSessionId, parentMessage, prompt string, think, search bool, response chan string) error {
	pow, err := c.getPow("/api/v0/chat/completion")
	if err != nil {
		return err
	}

	answer, err := c.powSolver.CalculateHash(pow.Challenge, pow.Salt, pow.Difficulty, int(pow.ExpireAt))
	if err != nil {
		return err
	}

	/*
		if parentMessage == "" {
			parentMessage = "null"
		} else {
			parentMessage = fmt.Sprintf(`"%s"`, parentMessage)
		}
	*/

	// body := fmt.Sprintf(`{"chat_session_id":"%s","parent_message_id":%s,"prompt":"%s","ref_file_ids":[],"thinking_enabled":%v,"search_enabled":%v}`, chatSessionId, parentMessage, prompt, think, search)
	completionData := models.CompletionData{
		ChatSessionId:   chatSessionId,
		ParentMessageId: nil,
		Prompt:          prompt,
		RefFileIds:      []string{},
		ThinkingEnabled: think,
		SearchEnabled:   search,
	}

	if parentMessage != "" {
		parentMsgId, err := strconv.ParseInt(parentMessage, 10, 64)
		if err != nil {
			return err
		}
		parentMsgIdInt := int(parentMsgId)
		completionData.ParentMessageId = &parentMsgIdInt
	}

	body, err := sonic.Marshal(&completionData)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", completionUrl, bytes.NewReader(body))
	if err != nil {
		return err
	}

	c.applyHeaders(request, len(body))
	c.applyPowHeader(request, int(answer), pow)

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return err
	}

	var parseError error

	if resp != nil {
		go func() {
			parseError = parseEvents(resp.Body, response)
		}()
	}

	if resp.StatusCode != http.StatusOK {
		respBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		println(string(respBytes))
	}

	return parseError
}
