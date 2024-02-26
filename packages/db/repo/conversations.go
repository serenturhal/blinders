package repo

import (
	"context"
	"time"

	"blinders/packages/db/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ConversationsRepo struct {
	Col *mongo.Collection
}

func NewConversationsRepo(col *mongo.Collection) *ConversationsRepo {
	return &ConversationsRepo{
		Col: col,
	}
}

func (r *ConversationsRepo) GetConversationByID(id primitive.ObjectID) (models.Conversation, error) {
	ctx, cal := context.WithTimeout(context.Background(), time.Second)
	defer cal()

	var conversation models.Conversation
	err := r.Col.FindOne(ctx, bson.M{"_id": id}).Decode(&conversation)

	return conversation, err
}

func (r *ConversationsRepo) InsertNewConversation(c models.Conversation) (models.Conversation, error) {
	ctx, cal := context.WithTimeout(context.Background(), time.Second)
	defer cal()

	_, err := r.Col.InsertOne(ctx, c)

	return c, err
}

// this function creates new ID and time and insert the document to database
func (r *ConversationsRepo) InsertNewRawConversation(conversation models.Conversation) (models.Conversation, error) {
	conversation.ID = primitive.NewObjectID()
	now := primitive.NewDateTimeFromTime(time.Now())
	conversation.CreatedAt = now
	conversation.UpdatedAt = now

	return r.InsertNewConversation(conversation)
}
