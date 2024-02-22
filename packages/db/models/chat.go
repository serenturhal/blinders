package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Conversation struct {
	ID        primitive.ObjectID `bson:"_id" json:"Id"`
	Members   []Member
	Messages  []Message
	CreatedAt primitive.DateTime
	UpdatedAt primitive.DateTime
	Metadata  ConversationMetadata
}

type ConversationMetadata struct {
	Name  string `omitempty:"true"`
	Image string `omitempty:"true"`
}

type Member struct {
	ID                    primitive.ObjectID `bson:"_id" json:"Id"`
	UserID                primitive.ObjectID
	Nickname              string
	LatestViewedMessageID primitive.ObjectID
	CreatedAt             primitive.DateTime
	UpdatedAt             primitive.DateTime
	JoinedAt              primitive.DateTime
}

type Message struct {
	ID        primitive.ObjectID `bson:"_id" json:"Id"`
	SenderID  primitive.ObjectID
	Content   string
	Status    string
	CreatedAt primitive.DateTime
	UpdatedAt primitive.DateTime
	Emotions  []MessageEmotion
}

type MessageEmotion struct {
	SenderID  primitive.ObjectID
	Content   string
	CreatedAt primitive.DateTime
	UpdatedAt primitive.DateTime
}
