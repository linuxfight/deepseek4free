package kv

import (
	"context"
	"errors"
	"fmt"
	"hash/fnv"
	"strconv"
	"strings"

	"github.com/valkey-io/valkey-go"
)

type ChatData struct {
	CurrentMessageId string
	ChatId           string
}

func (c *ChatData) Serialize() string {
	return fmt.Sprintf("%s;%s", c.ChatId, c.CurrentMessageId)
}

func (c *ChatData) Deserialize(text string) error {
	data := strings.Split(text, ";")
	if len(data) > 2 {
		return fmt.Errorf("invalid cache data")
	}

	c.ChatId = data[0]
	if len(data) != 2 {
		c.CurrentMessageId = ""
	}

	return nil
}

func (i *Instance) GetChatData(ctx context.Context, token, title string) (*ChatData, error) {
	h := fnv.New64a()
	_, err := h.Write([]byte(token + ";" + title))
	if err != nil {
		return nil, err
	}
	hs := strconv.Itoa(int(h.Sum64()))

	stringData, err := i.cache.Do(ctx, i.cache.B().Get().Key(hs).Build()).ToString()
	if err != nil {
		if errors.Is(err, valkey.Nil) {
			return &ChatData{
				CurrentMessageId: "",
				ChatId:           "",
			}, nil
		}

		return nil, err
	}

	chatData := &ChatData{}
	err = chatData.Deserialize(stringData)
	if err != nil {
		return nil, err
	}

	return chatData, nil
}

func (i *Instance) SetChatData(ctx context.Context, token, title string, data *ChatData) error {
	h := fnv.New64a()
	_, err := h.Write([]byte(token + ";" + title))
	if err != nil {
		return err
	}
	hs := strconv.Itoa(int(h.Sum64()))

	err = i.cache.Do(ctx, i.cache.B().Set().Key(hs).Value(data.Serialize()).Build()).Error()
	if err != nil {
		return err
	}

	return nil
}
