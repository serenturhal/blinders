package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"blinders/packages/auth"
	"blinders/packages/db"
	"blinders/packages/explore"
	"blinders/packages/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"

	exploreapi "blinders/services/explore/api"
)

var service *exploreapi.Service

func init() {
	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}
	app := fiber.New(fiber.Config{
		Prefork:      true,
		WriteTimeout: time.Second * 5,
		ReadTimeout:  time.Second * 5,

		DisableStartupMessage: true,
	})

	adminJSON, _ := utils.GetFile("firebase.admin.json")
	authManager, err := auth.NewFirebaseManager(adminJSON)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

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

	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	if err := redisClient.Ping(ctx).Err(); err != nil {
		panic(err)
	}

	core := explore.NewMongoExplorer(mongoManager, redisClient)

	service = exploreapi.NewService(app, authManager, core, redisClient)
	service.InitRoute()
}

func main() {
	port := os.Getenv("EXPLORE_SERVICE_PORT")
	go service.Loop()
	fmt.Println("listening on: ", port)
	log.Panic(service.App.Listen(":" + port))
}
