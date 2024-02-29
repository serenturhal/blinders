package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MatchInfo struct {
	FirebaseUID string             `json:"firebaseUID" bson:"firebaseUID,omiempty"`
	Name        string             `json:"name" bson:"name,omiempty"`
	Gender      string             `json:"gender" bson:"gender,omiempty"`
	Major       string             `json:"major" bson:"major,omiempty"`
	Native      string             `json:"native" bson:"native,omiempty"`
	Country     string             `json:"country" bson:"country,omiempty"` // ISO-3166 format
	Learnings   []string           `json:"learnings" bson:"learnings,omiempty"`
	Interests   []string           `json:"interests" bson:"interests,omiempty"`
	UserID      primitive.ObjectID `json:"userID" bson:"userID,omiempty"`
	Age         int                `json:"age" bson:"age,omiempty"`
}
