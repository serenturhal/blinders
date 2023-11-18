package main

import (
	"blinders/packages/translation"
	"context"
	"encoding/json"
	"fmt"
	"os"

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
	translator = translation.YandexTranslator{ApiKey: os.Getenv("YANDEX_API_KEY")}
	fmt.Println(translator, "<-- translator")
}

func HandleRequest(_ context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	fmt.Println(event, "<-- event")
	translated, err := translator.TranslateEnToVi("hello")
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 400,
			Body:       err.Error(),
		}, err
	}

	res := TranslateResponse{
		Text:       "hello",
		Translated: translated,
		Languages:  "en-vi",
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
