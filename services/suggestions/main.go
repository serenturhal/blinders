package main

import (
	"os"

	pkg "blinders/packages/suggestion"
	"blinders/services/suggestion/core"

	"github.com/sashabaranov/go-openai"
)

func initService() *core.Service {
	port := os.Getenv("PORT")
	config := &core.ServiceConfig{
		Port: port,
	}

	openaiKey := os.Getenv("OPENAI_API_KEY")

	client := openai.NewClient(openaiKey)

	suggester, err := pkg.NewGPTSuggester(client)
	if err != nil {
		panic(err)
	}

	sv, err := core.NewTransporter(suggester, config)
	if err != nil {
		panic(err)
	}

	return sv
}

func main() {
	service := initService()
	panic(service.Listen())
}
