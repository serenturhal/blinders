package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"blinders/packages/auth"
	"blinders/packages/utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type APIGatewayWebsocketProxyRequest struct {
	events.APIGatewayWebsocketProxyRequest `       json:",inline"`
	MethodArn                              string `json:"methodArn"` // ??? refs: https://gist.github.com/praveen001/1b045d1c31cd9c72e4e6638e9f883f83
}

var authManager auth.Manager

func init() {
	adminConfig, err := utils.GetFile("firebase.admin.json")
	if err != nil {
		log.Fatal(err)
	}

	authManager, err = auth.NewFirebaseManager(adminConfig)
	if err != nil {
		log.Fatal(err)
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
	user, err := authManager.Verify(jwt)
	if err != nil {
		return events.APIGatewayCustomAuthorizerResponse{}, err
	}

	// Is it secure to log the uid out to cloudwatch?
	// how to log the request tracking efficient and secure
	log.Println("[authorizer] issued user's policy of", user.AuthID)

	userBytes, _ := json.Marshal(user)
	return events.APIGatewayCustomAuthorizerResponse{
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
