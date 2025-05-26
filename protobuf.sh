#!/bin/zsh

protoc --go_out=./internal/stub/gen \
        --go_opt=paths=source_relative \
          --go-grpc_out=./internal/stub/gen \
            --go-grpc_opt=paths=source_relative \
              api.proto