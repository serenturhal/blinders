package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Conversation struct {
	ID        primitive.ObjectID   `bson:"_id"       json:"id"`
	Members   []Member             `bson:"members"   json:"members"`
	CreatedAt primitive.DateTime   `bson:"createdAt" json:"createdAt"`
	UpdatedAt primitive.DateTime   `bson:"updatedAt" json:"updatedAt"`
	Metadata  ConversationMetadata `bson:"metadata"  json:"metadata"`
}

type ConversationMetadata struct {
	Name  string `omitempty:"true"`
	Image string `omitempty:"true"`
}

type Member struct {
	ID                    primitive.ObjectID `bson:"_id"                   json:"id"`
	UserID                primitive.ObjectID `bson:"userId"                json:"userId"`
	Nickname              string             `bson:"nickname"              json:"nickname"`
	LatestViewedMessageID primitive.ObjectID `bson:"latestViewedMessageId" json:"latestViewedMessageId"`
	CreatedAt             primitive.DateTime `bson:"createdAt"             json:"createdAt"`
	UpdatedAt             primitive.DateTime `bson:"updatedAt"             json:"updatedAt"`
	JoinedAt              primitive.DateTime `bson:"joinedAt"              json:"joinedAt"`
}

type MessageStatus string

const (
	DeliveredStatus MessageStatus = "delivered"
	ReceivedStatus  MessageStatus = "received"
	SeenStatus      MessageStatus = "seen"
)

type Message struct {
	ID             primitive.ObjectID `bson:"_id"            json:"id"`
	SenderID       primitive.ObjectID `bson:"senderId"       json:"senderId"`
	ConversationID primitive.ObjectID `bson:"conversationId" json:"conversationId"`
	ReplyTo        primitive.ObjectID `bson:"replyTo"        json:"replyTo"`
	Content        string             `bson:"content"        json:"content"`
	Status         MessageStatus      `bson:"status"         json:"status"`
	CreatedAt      primitive.DateTime `bson:"createdAt"      json:"createdAt"`
	UpdatedAt      primitive.DateTime `bson:"updatedAt"      json:"updatedAt"`
	Emotions       []MessageEmotion   `bson:"emotions"       json:"emotions"`
}

type MessageEmotion struct {
	SenderID  primitive.ObjectID `bson:"senderId"  json:"senderId"`
	Content   string             `bson:"content"   json:"content"`
	CreatedAt primitive.DateTime `bson:"createdAt" json:"createdAt"`
	UpdatedAt primitive.DateTime `bson:"updatedAt" json:"updatedAt"`
}
