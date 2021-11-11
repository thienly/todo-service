#!/bin/bash
#protoc --proto_path=./api/proto/v1  --go_out=./pkg/api/v1 --go_opt=paths=source_relative \
#    --go-grpc_out=./pkg/api/v1 --go-grpc_opt=paths=source_relative todo.proto

#protoc -I . --go_out=./pkg/api/v1 --go_opt=paths=source_relative \
#    --go-grpc_out=./pkg/api/v1 --go-grpc_opt=paths=source_relative api/proto/todo.proto

protoc --proto_path=api/proto v1/todo.proto  --go_out=:pb \
  --go-grpc_out=:pb  --grpc-gateway_out=:pb  --openapiv2_out=:swagger
