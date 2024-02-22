package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type APIGatewayWebsocketProxyRequest struct {
	events.APIGatewayWebsocketProxyRequest `       json:",inline"`
	MethodArn                              string `json:"methodArn"` // ??? refs: https://gist.github.com/praveen001/1b045d1c31cd9c72e4e6638e9f883f83
}

func HandleRequest(
	_ context.Context,
	request APIGatewayWebsocketProxyRequest,
) (events.APIGatewayCustomAuthorizerResponse, error) {
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
			"user": "hello user from authorizer",
		},
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
