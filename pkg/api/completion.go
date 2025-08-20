package api

import (
	"fmt"
	"net/http"
	"strings"
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

	if parentMessage == "" {
		parentMessage = "null"
	} else {
		parentMessage = fmt.Sprintf(`"%s"`, parentMessage)
	}

	body := fmt.Sprintf(`{"chat_session_id":"%s","parent_message_id":%s,"prompt":"%s","ref_file_ids":[],"thinking_enabled":%v,"search_enabled":%v}`, chatSessionId, parentMessage, prompt, think, search)

	request, err := http.NewRequest("POST", completionUrl, strings.NewReader(body))
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

	return parseError
}
