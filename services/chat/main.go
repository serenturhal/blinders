package main

import (
	"fmt"
	"log"
	"os"

	"blinders/packages/auth"
	"blinders/packages/utils"
	chatcore "blinders/services/chat/core"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

var service chatcore.Service

func init() {
	app := fiber.New()
	adminJSON, _ := utils.GetFile("firebase.admin.development.json")
	auth, _ := auth.NewFirebaseManager(adminJSON)
	service = chatcore.Service{App: app, Auth: auth}
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
