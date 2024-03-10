package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"blinders/packages/auth"
	"blinders/packages/db"
	"blinders/packages/utils"
	restapi "blinders/services/rest/api"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	fiberadapter "github.com/awslabs/aws-lambda-go-api-proxy/fiber"
	"github.com/gofiber/fiber/v2"
)

var fiberLambda *fiberadapter.FiberLambda

func init() {
	log.Println("rest api running on environment:", os.Getenv("ENVIRONMENT"))
	app := fiber.New()

	url := fmt.Sprintf(
		db.MongoURLTemplate,
		os.Getenv("MONGO_USERNAME"),
		os.Getenv("MONGO_PASSWORD"),
		os.Getenv("MONGO_HOST"),
		os.Getenv("MONGO_PORT"),
		os.Getenv("MONGO_DATABASE"),
	)

	database := db.NewMongoManager(url, os.Getenv("MONGO_DATABASE"))
	if database == nil {
		log.Fatal("cannot create database manager")
	}

	adminConfig, err := utils.GetFile("firebase.admin.json")
	if err != nil {
		log.Fatal(err)
	}
	authManager, err := auth.NewFirebaseManager(adminConfig)
	if err != nil {
		log.Fatal(err)
	}

	api := restapi.NewManager(app, authManager, database)
	err = api.InitRoute(restapi.InitOptions{})
	if err != nil {
		panic(err)
	}

	fiberLambda = fiberadapter.New(api.App)
}

func Handler(
	ctx context.Context,
	req events.APIGatewayV2HTTPRequest,
) (events.APIGatewayV2HTTPResponse, error) {
	log.Println(req.RawPath)
	log.Println(req)
	return fiberLambda.ProxyWithContextV2(ctx, req)
}

func main() {
	lambda.Start(Handler)
}
