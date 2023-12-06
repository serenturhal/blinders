package main

import (
	"blinders/functions/utils"
	"blinders/packages/common"
	"blinders/packages/suggestion"
	"blinders/packages/user"
	commonUtils "blinders/utils"
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
		suggester, err = suggestion.NewGPTSuggester(client)
		if err != nil {
			panic(err)
		}
	}
}

type ChatSuggestRequest struct {
	Msgs []ClientMessage `json:"messages"`
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
	token, ok := event.Headers["Authorization"]
	if !ok {
		return utils.APIGatewayProxyResponseWithJSON(
			400,
			map[string]any{
				"error": "function: Token not found",
			})
	}

	usr, err := commonUtils.VerifyFirestoreToken(token)
	if err != nil {
		return utils.APIGatewayProxyResponseWithJSON(
			400,
			map[string]any{
				"error": fmt.Sprintf("function: cannot verify user, err : (%s)", err.Error()),
			})
	}

	suggestionRequest := new(ChatSuggestRequest)
	if err := json.Unmarshal([]byte(event.Body), suggestionRequest); err != nil {
		return utils.APIGatewayProxyResponseWithJSON(
			400,
			map[string]any{
				"error": fmt.Sprintf("function: cannot unmarshal request body, err : (%s)", err.Error()),
			})
	}

	msgs := []common.Message{}
	for _, m := range suggestionRequest.Msgs {
		msgs = append(msgs, m.ToCommonMessage())
	}

	userData, err := user.GetUserData(usr.ID)
	if err != nil {
		return utils.APIGatewayProxyResponseWithJSON(
			400,
			map[string]any{
				"error": fmt.Sprintf("function: cannot get data of user, userid: (%s) err: (%s)", usr.ID, err.Error()),
			})
	}

	suggestions, err := suggester.ChatCompletion(ctx, userData, msgs)
	if err != nil {
		return utils.APIGatewayProxyResponseWithJSON(
			400,
			map[string]any{
				"error": fmt.Sprintf("function: cannot get suggestions, err: (%s)", err.Error()),
			})
	}

	return utils.APIGatewayProxyResponseWithJSON(
		200,
		map[string]any{
			"suggestions": suggestions,
		},
	)
}

func main() {
	lambda.Start(HandleRequest)
}
