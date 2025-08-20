package kv

import (
	"context"
	"github.com/valkey-io/valkey-go"
)

type Instance struct {
	cache valkey.Client
}

func (i *Instance) Close() {
	i.cache.Close()
}

func New() (*Instance, error) {
	cache, err := valkey.NewClient(valkey.ClientOption{InitAddress: []string{"127.0.0.1:6379"}})
	if err != nil {
		return nil, err
	}

	err = cache.Do(context.Background(), cache.B().Ping().Build()).Error()
	if err != nil {
		return nil, err
	}

	return &Instance{
		cache: cache,
	}, nil
}
