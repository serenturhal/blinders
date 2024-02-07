package main

import (
	"log"
	"os"

	"blinders/packages/auth"
	"blinders/packages/match"
	"blinders/packages/utils"

	"github.com/gofiber/fiber/v2"

	matchcore "blinders/services/match/core"
)

var service matchcore.Service

func init() {
	app := fiber.New()

	adminJSON, _ := utils.GetFile("firebase.admin.development.json")
	auth, _ := auth.NewFirebaseManager(adminJSON)
	service = matchcore.Service{
		Auth: auth,
		App:  app,
		Core: &match.MockMatcher{},
	}
}

func main() {
	port := os.Getenv("MATCH_SERVICE_PORT")
	log.Panic(service.App.Listen(port))
}
