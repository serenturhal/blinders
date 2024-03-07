package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID            primitive.ObjectID     `bson:"_id"           json:"id"`
	Conversations []EmbeddedConversation `bson:"conversations" json:"conversations"`
	FirebaseUID   string                 `bson:"firebaseUID"   json:"firebaseUID"`
	ImageURL      string                 `bson:"imageURL"      json:"imageURL"`
	FriendsUserID []string               `bson:"friends"       json:"friends"`
	CreatedAt     primitive.DateTime     `bson:"createdAt"     json:"createdAt"`
	UpdatedAt     primitive.DateTime     `bson:"updatedAt"     json:"updatedAt"`
}

type EmbeddedConversation struct {
	ID             primitive.ObjectID `bson:"_id"            json:"id"`
	ConversationID primitive.ObjectID `bson:"conversationId" json:"conversationId"`
	CreatedAt      primitive.DateTime `bson:"createdAt"      json:"createdAt"`
	UpdatedAt      primitive.DateTime `bson:"updatedAt"      json:"updatedAt"`
	Settings       struct {
		Notification bool `bson:"notification"   json:"notification"`
	} `bson:"settings"       json:"settings"`
}
