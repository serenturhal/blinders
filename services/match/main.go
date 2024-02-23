package main

import (
	"context"
	"log"
	"os"
	"time"

	"blinders/packages/auth"
	"blinders/packages/match"
	"blinders/packages/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	matchapi "blinders/services/match/api"
)

var (
	service matchapi.Service
	client  *mongo.Client
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}
	app := fiber.New()

	adminJSON, _ := utils.GetFile("firebase.admin.development.json")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	opts := options.Client().ApplyURI(os.Getenv("MONGO_URL")).SetAuth(options.Credential{
		Username: os.Getenv("MONGO_DB_USERNAME"),
		Password: os.Getenv("MONGO_DB_PASSWORD"),
	})

	mongoClient, err := mongo.Connect(ctx, opts)
	if err != nil {
		panic(err)
	}
	client = mongoClient
	if err := client.Ping(ctx, &readpref.ReadPref{}); err != nil {
		panic(err)
	}

	db := client.Database(os.Getenv("MONGO_DB"))
	userCol := db.Collection("user")
	matchCol := db.Collection("match")

	auth, _ := auth.NewFirebaseManager(adminJSON)
	service = matchapi.Service{
		Auth: auth,
		App:  app,
		Core: &match.MongoMatcher{
			UserCol:  userCol,
			MatchCol: matchCol,
			Embedder: match.MockEmbedder{},
		},
	}
	service.InitRoute()
}

func main() {
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	port := os.Getenv("MATCH_SERVICE_PORT")
	log.Panic(service.App.Listen(":" + port))
}
