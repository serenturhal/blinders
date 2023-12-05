package main

import (
	"blinders/packages/suggestion"
	"context"
	"encoding/json"
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

	suggestionRequest := new(SuggestionPayload)

	if err := json.Unmarshal([]byte(event.Body), suggestionRequest); err != nil {
		return APIGatewayProxyResponseWithJSON(400, map[string]any{
			"error": fmt.Sprintf("functions: cannot unmarshal struct from json, err: (%s)", err.Error()),
		})
	}

	suggestion, err := suggester.TextCompletion(ctx, suggestionRequest.Text)
	if err != nil {
		return APIGatewayProxyResponseWithJSON(400, map[string]any{
			"error": fmt.Sprintf("functions: cannot get suggestions, err: (%s)", err.Error()),
		})
	}
	return APIGatewayProxyResponseWithJSON(200, map[string]any{
		"suggestions": suggestion,
	})
}

func APIGatewayProxyResponseWithJSON(code int, v any) (events.APIGatewayProxyResponse, error) {
	bodyByte, err := json.Marshal(v)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       fmt.Sprintf("functions: cannot marshal struct to json, err: (%s)", err.Error()),
		}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: code,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(bodyByte),
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
