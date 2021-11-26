#!/bin/sh
export PATH="$PATH:$(go env GOPATH)/bin:/var/opt/protoc/bin"
protoc -I ./api/v1 -I ./api \
--go_out=./api/v1 --go_opt=paths=source_relative \
--go-grpc_out=./api/v1 --go-grpc_opt=paths=source_relative \
--grpc-gateway_out=./api/v1 --grpc-gateway_opt logtostderr=true --grpc-gateway_opt=paths=source_relative \
./api/v1/capturedimage.proto 

# protoc -I ././api/v1 --openapiv2_out ./api/v1 \
#     --openapiv2_opt logtostderr=true \
#     ./api/v1/capturedimage.proto 

# for JS gRPC client
#protoc -I ./api/v1 -I ./api \
#--js_out=import_style=commonjs,binary:client ./api/v1/capturedimage.proto 