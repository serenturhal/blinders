package main

import (
	"fmt"
	"log"
	"os"

	"blinders/packages/auth"
	"blinders/packages/suggest"
	"blinders/packages/utils"
	suggestcore "blinders/services/suggestion/core"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
)

var service suggestcore.Service

func init() {
	if err := godotenv.Load(".env.development"); err != nil {
		log.Fatal("failed to load env", err)
	}
	app := fiber.New()
	adminJSON, _ := utils.GetFile("firebase.admin.development.json")
	auth, _ := auth.NewFirebaseManager(adminJSON)

	openaiKey := os.Getenv("OPENAI_API_KEY")
	client := openai.NewClient(openaiKey)
	suggester, err := suggest.NewGPTSuggester(client)
	if err != nil {
		log.Fatal("failed to init openai client", err)
	}

	service = suggestcore.Service{App: app, Auth: auth, Suggester: suggester}
	service.InitRoute()
}

func main() {
	port := os.Getenv("SUGGEST_SERVICE_PORT")
	err := service.App.Listen(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Println("launch suggest service error", err)
	}
}
