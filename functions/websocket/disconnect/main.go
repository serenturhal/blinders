package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest(
	ctx context.Context,
	request events.APIGatewayWebsocketProxyRequest,
) (events.APIGatewayProxyResponse, error) {
	// Handle the connect request here
	fmt.Println(ctx, request)
	return events.APIGatewayProxyResponse{StatusCode: 200, Body: "Disconnected."}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
