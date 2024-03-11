package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type MatchInfo struct {
	UserID    primitive.ObjectID `json:"userId"    bson:"userId,omitempty"`
	Name      string             `json:"name"      bson:"name,omitempty"`
	Gender    string             `json:"gender"    bson:"gender,omitempty"`
	Major     string             `json:"major"     bson:"major,omitempty"`
	Native    string             `json:"native"    bson:"native,omitempty"`
	Country   string             `json:"country"   bson:"country,omitempty"` // ISO-3166 format
	Learnings []string           `json:"learnings" bson:"learnings,omitempty"`
	Interests []string           `json:"interests" bson:"interests,omitempty"`
	Age       int                `json:"age"       bson:"age,omitempty"`
}
