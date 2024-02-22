package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID            primitive.ObjectID `bson:"_id" json:"Id"`
	Conversations []EmbeddedConversation
	CreatedAt     primitive.DateTime
	UpdatedAt     primitive.DateTime
}

type EmbeddedConversation struct {
	ID             primitive.ObjectID `bson:"_id" json:"Id"`
	ConversationID primitive.ObjectID
	CreatedAt      primitive.DateTime
	UpdatedAt      primitive.DateTime
	Settings       struct {
		Notification bool
	}
}
