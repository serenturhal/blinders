package models

import (
	"bytes"
	"encoding/binary"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MatchInfo struct {
	UserID      primitive.ObjectID `json:"userID" bson:"userID,omiempty"`
	FirebaseUID string             `json:"firebaseUID" bson:"firebaseUID,omiempty"`
	Name        string             `json:"name" bson:"name,omiempty"`
	Gender      string             `json:"gender" bson:"gender,omiempty"`
	Major       string             `json:"major" bson:"major,omiempty"`
	Native      string             `json:"native" bson:"native,omiempty"`
	Country     string             `json:"country" bson:"country,omiempty"` // ISO-3166 format
	Learnings   []string           `json:"learnings" bson:"learning,omiempty"`
	Interests   []string           `json:"interests" bson:"interests,omiempty"`
	Age         int                `json:"age" bson:"age,omiempty"`
}

type (
	MatchEmbedding [128]float32
	MatchStore     struct {
		MatchInfo `bson:",inline,omiempty"`
		Vector    MatchEmbedding `bson:"vector"`
	}
)

func (v MatchEmbedding) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, v)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
