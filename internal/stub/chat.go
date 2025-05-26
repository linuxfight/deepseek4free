package stub

import (
	"context"
	"github.com/linuxfight/deepseek4free/internal/stub/gen"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (stub *Stub) CreateChat(_ context.Context, _ *emptypb.Empty) (*gen.Chat, error) {
	chatId, err := stub.api.CreateChat()
	if err != nil {
		return nil, err
	}
	return &gen.Chat{
		Id:               chatId,
		Title:            nil,
		CurrentMessageId: nil,
	}, nil
}

func (stub *Stub) EditChat(_ context.Context, chat *gen.Chat) (*emptypb.Empty, error) {
	err := stub.api.ChangeTitle(chat.Id, *chat.Title)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (stub *Stub) DeleteChat(_ context.Context, chat *gen.Chat) (*emptypb.Empty, error) {
	err := stub.api.DeleteChatSession(chat.Id)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (stub *Stub) GetAllChats(_ context.Context, _ *emptypb.Empty) (*gen.ChatListResponse, error) {
	apiChats, err := stub.api.GetAllChats()
	if err != nil {
		return nil, err
	}
	response := &gen.ChatListResponse{}
	for _, apiChat := range apiChats {
		chat := gen.Chat{
			Id:               apiChat.Id,
			CurrentMessageId: apiChat.CurrentMessageId,
			Title:            apiChat.Title,
		}
		response.Chats = append(response.Chats, &chat)
	}
	return response, nil
}
