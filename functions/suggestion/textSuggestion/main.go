package main

import (
	"blinders/packages/suggestion"
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sashabaranov/go-openai"
)

var (
	suggester suggestion.Suggester = nil
	apiKey                         = os.Getenv("OPENAI_API_KEY")
)

func init() {
	if suggester == nil {
		var err error
		client := openai.NewClient(apiKey)
		suggester, err = suggestion.NewGPTSuggestor(client)
		if err != nil {
			panic(err)
		}
	}
}

type SuggestionPayload struct {
	Text string `json:"text"`
}

func HandleRequest(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Printf("processing requests data for event %s.\n", event.RequestContext.RequestID)
	fmt.Printf("Body size = %d.\n", len(event.Body))
	fmt.Println("Headers:")
	for key, value := range event.Headers {
		fmt.Printf("\t%s: %s\n", key, value)
	}
	event.Body
}

func main() {
	lambda.Start(HandleRequest)
}
