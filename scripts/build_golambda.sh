# !/bin/bash

rm -rf dist
echo "cleaned dist directory"

GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GOFLAGS=-trimpath go build -mod=readonly -ldflags='-s -w' -o ./dist/ ./functions/translate
echo "build translate lambda function completed"

GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GOFLAGS=-trimpath go build -mod=readonly -ldflags='-s -w' -o ./dist/ ./functions/websocket/connect
echo "build connect lambda function completed"

GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GOFLAGS=-trimpath go build -mod=readonly -ldflags='-s -w' -o ./dist/authorizer/handler ./functions/websocket/authorizer
echo "build authorizer lambda function completed"
cp ./firebase.admin.json ./dist/authorizer
echo "copied firebase.admin.json to authorizer"

GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GOFLAGS=-trimpath go build -mod=readonly -ldflags='-s -w' -o ./dist/ ./functions/websocket/disconnect
echo "build disconnect lambda function completed"
