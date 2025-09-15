package dto

import (
	"fmt"
	"strings"

	"github.com/bytedance/sonic"
)

// UnmarshalJSON handles both string and array formats for Content
func (m *Message) UnmarshalJSON(data []byte) error {
	// Temporary struct to handle dynamic content
	type tempMessage struct {
		Role    string      `json:"role"`
		Content interface{} `json:"content"`
	}

	var tmp tempMessage
	// Use Sonic for fast unmarshaling
	if err := sonic.Unmarshal(data, &tmp); err != nil {
		return err
	}

	m.Role = tmp.Role

	// Handle different content types
	switch content := tmp.Content.(type) {
	case string:
		m.Content = content
	case []interface{}:
		// Process array of content objects
		var texts []string
		for _, item := range content {
			switch v := item.(type) {
			case string:
				texts = append(texts, v)
			case map[string]interface{}:
				if text, ok := v["text"].(string); ok {
					texts = append(texts, text)
				}
			}
		}
		m.Content = strings.Join(texts, " ")
	default:
		return fmt.Errorf("unexpected content type: %T", content)
	}

	return nil
}
