package main

import (
	"fmt"
	"log"
	"os"

	"blinders/packages/auth"
	"blinders/packages/utils"
	chatapi "blinders/services/chat/api"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

var service chatapi.Service

func init() {
	app := fiber.New()
	adminJSON, _ := utils.GetFile("firebase.admin.development.json")

	authManager, _ := auth.NewFirebaseManager(adminJSON)
	service = chatapi.Service{App: app, Auth: authManager}
	service.InitRoute()
	if err := godotenv.Load(".env.development"); err != nil {
		log.Fatal("failed to load env", err)
	}
}

func main() {
	port := os.Getenv("CHAT_SERVICE_PORT")
	err := service.App.Listen(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Println("launch chat service error", err)
	}
}
