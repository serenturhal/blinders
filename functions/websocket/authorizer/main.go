package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"blinders/packages/auth"
	"blinders/packages/db"
	"blinders/packages/utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type APIGatewayWebsocketProxyRequest struct {
	events.APIGatewayWebsocketProxyRequest `       json:",inline"`
	MethodArn                              string `json:"methodArn"` // ??? refs: https://gist.github.com/praveen001/1b045d1c31cd9c72e4e6638e9f883f83
}

var (
	authManager auth.Manager
	database    *db.MongoManager
)

func init() {
	adminConfig, err := utils.GetFile("firebase.admin.json")
	if err != nil {
		log.Fatal(err)
	}

	authManager, err = auth.NewFirebaseManager(adminConfig)
	if err != nil {
		log.Fatal(err)
	}

	url := fmt.Sprintf(
		db.MongoURLTemplate,
		os.Getenv("MONGO_USERNAME"),
		os.Getenv("MONGO_PASSWORD"),
		os.Getenv("MONGO_HOST"),
		os.Getenv("MONGO_PORT"),
		os.Getenv("MONGO_DATABASE"),
	)

	database = db.NewMongoManager(url, os.Getenv("MONGO_DATABASE"))
	if database == nil {
		log.Fatal("cannot create database manager")
	}
}

func HandleRequest(
	_ context.Context,
	request APIGatewayWebsocketProxyRequest,
) (events.APIGatewayCustomAuthorizerResponse, error) {
	authorization := request.Headers["Authorization"]
	if !strings.Contains(authorization, "Bearer ") {
		// TODO: need to response with unauthorized response code
		// (it's currently 500 as default)
		return events.APIGatewayCustomAuthorizerResponse{}, fmt.Errorf("[authorizer] invalid token")
	}
	jwt := strings.Split(authorization, " ")[1]
	authUser, err := authManager.Verify(jwt)
	if err != nil {
		return events.APIGatewayCustomAuthorizerResponse{}, err
	}

	user, err := database.Users.GetUserByFirebaseUID(authUser.AuthID)
	if err != nil {
		return events.APIGatewayCustomAuthorizerResponse{}, err
	}

	// Is it secure to log the id out to cloudwatch?
	// how to log the request tracking efficient and secure
	log.Println("[authorizer] issued user's policy of", user.ID.Hex())

	userBytes, _ := json.Marshal(authUser)
	return events.APIGatewayCustomAuthorizerResponse{
		PrincipalID: user.ID.Hex(),
		PolicyDocument: events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Action:   []string{"execute-api:Invoke"},
					Effect:   "Allow",
					Resource: []string{request.MethodArn},
				},
			},
		},
		Context: map[string]interface{}{
			"user": string(userBytes),
		},
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
