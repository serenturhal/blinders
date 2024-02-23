package db

import (
	"context"
	"log"
	"time"

	"blinders/packages/db/repo"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// username:password@host:port
const MongoURLTemplate = "mongodb://%s:%s@%s/%s"

const (
	UserCollection         = "users"
	ConversationCollection = "conversations"
)

type MongoManager struct {
	Client        *mongo.Client
	Database      string
	Users         *repo.Users
	Conversations *repo.Conversations
}

func NewMongoManager(url string, name string) *MongoManager {
	ctx, cal := context.WithTimeout(context.Background(), time.Second*10)
	defer cal()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	if err != nil {
		log.Fatal("cannot connect to mongo")
	}

	return &MongoManager{
		Client:        client,
		Database:      name,
		Users:         repo.NewUsers(client.Database(name).Collection(UserCollection)),
		Conversations: repo.NewConversations(client.Database(name).Collection(ConversationCollection)),
	}
}
