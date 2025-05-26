package stub

import (
	"context"
	"github.com/linuxfight/deepseek4free/internal/stub/gen"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (stub *Stub) GetQuota(_ context.Context, _ *emptypb.Empty) (*gen.QuotaResponse, error) {
	apiQuota, err := stub.api.GetQuota()
	if err != nil {
		return nil, err
	}
	return &gen.QuotaResponse{
		Quota: int32(apiQuota.Quota),
		Used:  int32(apiQuota.Used),
	}, nil
}

func (stub *Stub) Completion(req *gen.CompletionRequest, stream grpc.ServerStreamingServer[gen.CompletionResponse]) error {
	history, err := stub.api.GetMessageHistory(req.Chat.Id)
	if err != nil {
		return err
	}

	parentMessageId := ""
	if history.ChatSession.CurrentMessageId != nil {
		parentMessageId = string(*history.ChatSession.CurrentMessageId)
	}

	tokensCh := make(chan string)

	err = stub.api.Completion(req.Chat.Id, parentMessageId, req.Prompt, req.Think, req.Search, tokensCh)
	if err != nil {
		return err
	}

	for msg := range tokensCh {
		err := stream.Send(&gen.CompletionResponse{
			Token: msg,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
