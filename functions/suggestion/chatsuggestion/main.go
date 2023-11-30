package main

import (
	"blinders/packages/common"
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

type ChatSuggestRequest struct {
	UserContext common.UserContext `json:"userContext"`
	Msgs        []common.Message   `json:"messages"`
}

func HandleRequest(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	suggestionRequest := new(ChatSuggestRequest)
	if err := json.Unmarshal([]byte(event.Body), suggestionRequest); err != nil {
		return APIGatewayProxyResponseWithJSON(
			400,
			map[string]any{
				"error": fmt.Sprintf("function: cannot unmarshal request body, err : (%s)", err.Error()),
			})
	}

	suggestions, err := suggester.ChatCompletion(ctx, suggestionRequest.UserContext, suggestionRequest.Msgs)
	if err != nil {
		return APIGatewayProxyResponseWithJSON(
			400,
			map[string]any{
				"error": fmt.Sprintf("function: cannot get suggestions, err: (%s)", err.Error()),
			})
	}

	return APIGatewayProxyResponseWithJSON(
		200,
		map[string]any{
			"suggestions": suggestions,
		},
	)
}

func APIGatewayProxyResponseWithJSON(code int, v any) (events.APIGatewayProxyResponse, error) {
	defaultResponse := events.APIGatewayProxyResponse{
		StatusCode: 400,
	}
	bodyBytes, err := json.Marshal(v)
	if err != nil {
		defaultResponse.Body = fmt.Sprintf("function: cannot marshall struct, err : (%s)", err.Error())
		return defaultResponse, err
	}

	defaultResponse.StatusCode = code
	defaultResponse.Body = string(bodyBytes)
	return defaultResponse, nil
}

func main() {
	lambda.Start(HandleRequest)
}
