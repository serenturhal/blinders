package main

import (
	"fmt"
	"log"
	"os"

	"blinders/packages/auth"
	"blinders/packages/db"
	"blinders/packages/utils"
	chatapi "blinders/services/chat/api"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

var service chatapi.Service

func init() {
	app := fiber.New()
	adminJSON, _ := utils.GetFile("firebase.admin.development.json")
	url := fmt.Sprintf(
		db.MongoURLTemplate,
		os.Getenv("MONGO_USERNAME"),
		os.Getenv("MONGO_PASSWORD"),
		os.Getenv("MONGO_HOST"),
		os.Getenv("MONGO_PORT"),
		os.Getenv("MONGO_DATABASE"),
	)

	mongoManager := db.NewMongoManager(url, os.Getenv("MONGO_DATABASE"))

	fmt.Println("Connect to mongo url", url)
	auth, _ := auth.NewFirebaseManager(mongoManager.Users, adminJSON)
	service = chatapi.Service{App: app, Auth: auth}
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
