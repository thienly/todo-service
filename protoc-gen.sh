#!/bin/bash
protoc --proto_path=./api/proto/v1 --go_out=./pkg/api/v1 --go_opt=paths=source_relative \
    --go-grpc_out=./pkg/api/v1 --go-grpc_opt=paths=source_relative todo.proto
