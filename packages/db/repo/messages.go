package repo

import (
	"context"
	"time"

	"blinders/packages/db/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MessagesRepo struct {
	Col *mongo.Collection
}

func NewMessagesRepo(col *mongo.Collection) *MessagesRepo {
	return &MessagesRepo{
		Col: col,
	}
}

func (r MessagesRepo) ConstructNewMessage(
	senderID primitive.ObjectID,
	conversationID primitive.ObjectID,
	replyTo primitive.ObjectID,
	content string,
) models.Message {
	now := primitive.NewDateTimeFromTime(time.Now())
	return models.Message{
		ID:             primitive.NewObjectID(),
		Status:         "delivered",
		Emotions:       make([]models.MessageEmotion, 0),
		SenderID:       senderID,
		ConversationID: conversationID,
		ReplyTo:        replyTo,
		Content:        content,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

func (r *MessagesRepo) GetMessageByID(id primitive.ObjectID) (models.Message, error) {
	ctx, cal := context.WithTimeout(context.Background(), time.Second)
	defer cal()

	var message models.Message
	err := r.Col.FindOne(ctx, bson.M{"_id": id}).Decode(&message)

	return message, err
}

func (r *MessagesRepo) InsertNewMessage(m models.Message) (models.Message, error) {
	ctx, cal := context.WithTimeout(context.Background(), time.Second)
	defer cal()

	_, err := r.Col.InsertOne(ctx, m)

	return m, err
}

// this function creates new ID and time and insert the document to database
func (r *MessagesRepo) InsertNewRawMessage(m models.Message) (models.Message, error) {
	m.ID = primitive.NewObjectID()
	now := primitive.NewDateTimeFromTime(time.Now())
	m.CreatedAt = now
	m.UpdatedAt = now

	return r.InsertNewMessage(m)
}
