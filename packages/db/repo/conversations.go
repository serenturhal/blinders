package repo

import "go.mongodb.org/mongo-driver/mongo"

type Conversations struct {
	Col *mongo.Collection
}

func NewConversations(col *mongo.Collection) *Conversations {
	return &Conversations{
		Col: col,
	}
}
