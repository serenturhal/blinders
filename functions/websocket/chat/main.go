package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"blinders/packages/db"
	"blinders/packages/session"
	"blinders/packages/utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/redis/go-redis/v9"
)

var (
	sessionManager *session.Manager
	database       *db.MongoManager
)

func init() {
	fmt.Println("init chat websocket lambda function")
	// TODO: need to store these secrets to aws secret manager instead of pass in env
	sessionManager = session.NewManager(redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Username: os.Getenv("REDIS_USERNAME"),
		Password: os.Getenv("REDIS_PASSWORD"),
	}))

	url := fmt.Sprintf(
		db.MongoURLTemplate,
		os.Getenv("MONGO_USERNAME"),
		os.Getenv("MONGO_PASSWORD"),
		os.Getenv("MONGO_HOST"),
		os.Getenv("MONGO_PORT"),
		os.Getenv("MONGO_DATABASE"),
	)

	fmt.Println(url)
	database = db.NewMongoManager(url, os.Getenv("MONGO_DATABASE"))
	if database == nil {
		log.Fatal("cannot create database manager")
	}
}

func HandleRequest(
	_ context.Context,
	request events.APIGatewayWebsocketProxyRequest,
) (any, error) {
	connectionID := request.RequestContext.ConnectionID
	userID := request.RequestContext.Authorizer.(map[string]interface{})["principalId"].(string)

	genericEvent, err := utils.JSONConvert[ChatEvent](request.Body)
	if err != nil {
		return nil, err
	}

	switch genericEvent.Type {
	case UserSendMessage:
		payload, err := utils.JSONConvert[UserSendMessagePayload](request.Body)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusBadRequest,
				Body:       "invalid send message event",
			}, nil
		}
		handleSendMessage(userID, connectionID, *payload)
	default:
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}, nil
	}

	return events.APIGatewayProxyResponse{StatusCode: http.StatusOK}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
