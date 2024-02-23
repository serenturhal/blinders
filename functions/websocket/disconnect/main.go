package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"blinders/packages/auth"
	"blinders/packages/session"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/redis/go-redis/v9"
)

var sessionManager *session.Manager

func init() {
	// TODO: need to store these secrets to aws secret manager instead of pass in env
	sessionManager = session.NewManager(redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Username: os.Getenv("REDIS_USERNAME"),
		Password: os.Getenv("REDIS_PASSWORD"),
	}))
}

func HandleRequest(
	_ context.Context,
	request events.APIGatewayWebsocketProxyRequest,
) (events.APIGatewayProxyResponse, error) {
	connectionID := request.RequestContext.ConnectionID
	userStr := request.RequestContext.Authorizer.(map[string]interface{})["user"].(string)

	var user auth.UserAuth
	err := json.Unmarshal([]byte(userStr), &user)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 404, Body: "user not found"}, nil
	}

	err = sessionManager.RemoveSession(user.AuthID, connectionID)
	if err != nil {
		log.Println(err)
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "failed to remove session"}, nil
	}

	return events.APIGatewayProxyResponse{StatusCode: 200, Body: "Connected."}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
