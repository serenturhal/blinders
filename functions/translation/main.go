package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"blinders/packages/translation"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type TranslatePayload struct {
	Text string `json:"text"`
}

type TranslateResponse struct {
	Text       string `json:"text"`
	Translated string `json:"translated"`
	Languages  string `json:"languages"`
}

var translator translation.Translator

func init() {
	translator = translation.YandexTranslator{APIKey: os.Getenv("YANDEX_API_KEY")}
}

func HandleRequest(_ context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	text, ok := event.QueryStringParameters["text"]
	if !ok {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 400,
			Body:       "required text param",
		}, nil
	}

	langs, ok := event.QueryStringParameters["languages"]
	if !ok {
		langs = "en-vi"
	}

	translated, err := translator.Translate(text, translation.Languages(langs))
	if err != nil {
		log.Println("error translating: ", err)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 400,
			Body:       fmt.Sprintf("cannot translate \"%s\"", text),
		}, nil
	}

	res := TranslateResponse{
		Text:       text,
		Translated: translated,
		Languages:  langs,
	}

	resInBytes, _ := json.Marshal(res)
	return events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
		Body:       string(resInBytes),
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
