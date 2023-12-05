package main

import (
	"blinders/packages/common"
	blinderContext "blinders/packages/context"
	"blinders/packages/suggestion"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"time"

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
	UserID string          `json:"userId"`
	Msgs   []ClientMessage `json:"messages"`
}

type ClientMessage struct {
	Timestamp any    `json:"time"`
	ID        string `json:"id"`
	Content   string `json:"content"`
	FromID    string `json:"senderId"`
	ChatID    string `json:"roomId"`
}

func (m ClientMessage) ToCommonMessage() common.Message {
	var Timestamp int64
	switch timestamp := m.Timestamp.(type) {
	case int:
		Timestamp = int64(timestamp)
	case string:
		// expect date time as string type, "Tue Dec 05 2023 12:35:04 GMT+0700"
		layout := "Mon Jan 02 2006 15:04:05 GMT-0700"
		t, err := time.Parse(layout, timestamp)
		if err != nil {
			panic(fmt.Sprintf("clienmessage: given time (%s) cannot parse with layout (%s)", timestamp, layout))
		}
		Timestamp = t.Unix()
	default:
		panic(fmt.Sprintf("clienmessage: unknow timestamp type (%s)", reflect.TypeOf(m.Timestamp).String()))
	}

	return common.Message{
		FromID:    m.FromID,
		ToID:      m.ChatID,
		Content:   m.Content,
		Timestamp: Timestamp,
	}
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

	msgs := []common.Message{}
	for _, m := range suggestionRequest.Msgs {
		msgs = append(msgs, m.ToCommonMessage())
	}
	userContext, err := blinderContext.GetUserContext(suggestionRequest.UserID)
	if err != nil {
		return APIGatewayProxyResponseWithJSON(
			400,
			map[string]any{
				"error": fmt.Sprintf("function: cannot get context of user, userid: (%s) err: (%s)", suggestionRequest.UserID, err.Error()),
			})
	}

	suggestions, err := suggester.ChatCompletion(ctx, userContext, msgs)
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
