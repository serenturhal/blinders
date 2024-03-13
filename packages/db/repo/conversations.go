package repo

import (
	"context"
	"fmt"
	"log"
	"time"

	"blinders/packages/db/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ConversationsRepo struct {
	Col *mongo.Collection
}

func NewConversationsRepo(col *mongo.Collection) *ConversationsRepo {
	return &ConversationsRepo{
		Col: col,
	}
}

func (r *ConversationsRepo) GetConversationByID(
	id primitive.ObjectID,
) (*models.Conversation, error) {
	ctx, cal := context.WithTimeout(context.Background(), time.Second)
	defer cal()

	var conversation models.Conversation
	err := r.Col.FindOne(ctx, bson.M{"_id": id}).Decode(&conversation)

	return &conversation, err
}

func (r *ConversationsRepo) InsertNewConversation(
	c models.Conversation,
) (*models.Conversation, error) {
	ctx, cal := context.WithTimeout(context.Background(), time.Second)
	defer cal()

	_, err := r.Col.InsertOne(ctx, c)

	return &c, err
}

// this function creates new ID and time and insert the document to database
func (r *ConversationsRepo) InsertNewRawConversation(
	conversation models.Conversation,
) (*models.Conversation, error) {
	conversation.ID = primitive.NewObjectID()
	now := primitive.NewDateTimeFromTime(time.Now())
	conversation.CreatedAt = now
	conversation.UpdatedAt = now

	return r.InsertNewConversation(conversation)
}

func (r *ConversationsRepo) InsertIndividualConversation(
	userID, friendID primitive.ObjectID,
) (*models.Conversation, error) {
	ctx, cal := context.WithTimeout(context.Background(), time.Second)
	defer cal()

	upsert := true
	now := primitive.NewDateTimeFromTime(time.Now())
	result, err := r.Col.UpdateOne(ctx,
		bson.M{
			"type": models.IndividualConversation,
			"members": bson.M{
				"$all": []bson.M{
					{"$elemMatch": bson.M{"userId": userID}},
					{"$elemMatch": bson.M{"userId": friendID}},
				},
				"$size": 2,
			},
		},
		bson.M{
			"$setOnInsert": models.Conversation{
				ID:   primitive.NewObjectID(),
				Type: models.IndividualConversation,
				Members: []models.Member{{
					ID:        primitive.NewObjectID(),
					UserID:    userID,
					CreatedAt: now,
					UpdatedAt: now,
					JoinedAt:  now,
				}, {
					ID:        primitive.NewObjectID(),
					UserID:    friendID,
					CreatedAt: now,
					UpdatedAt: now,
					JoinedAt:  now,
				}},
				CreatedBy: userID,
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
		&options.UpdateOptions{Upsert: &upsert},
	)
	if err != nil {
		log.Println("can not insert conversation:", err)
		return nil, fmt.Errorf("something went wrong when inserting conversation")
	}
	if result.UpsertedCount == 0 {
		log.Println("conversation already existed")
		return nil, fmt.Errorf("conversation already existed")
	}

	conv, err := r.GetConversationByID(result.UpsertedID.(primitive.ObjectID))

	return conv, err
}
