package match

import (
	"bytes"
	"encoding/binary"
)

type UserMatch struct {
	UserID    string   `json:"id" bson:"userID,omiempty"`
	Name      string   `json:"name" bson:"name,omiempty"`
	Gender    string   `json:"gender" bson:"gender,omiempty"`
	Major     string   `json:"major" bson:"major,omiempty"`
	Native    string   `json:"native" bson:"native,omiempty"`
	Country   string   `json:"country" bson:"country,omiempty"` // ISO-3166 format
	Learnings []string `json:"learnings" bson:"learning,omiempty"`
	Interests []string `json:"interests" bson:"interests,omiempty"`
	Age       int      `json:"age" bson:"age,omiempty"`
}

type (
	EmbeddingVector [128]float32
	UserStore       struct {
		UserMatch `bson:",inline,omiempty"`
		Vector    EmbeddingVector `bson:"vector"`
	}
)

func (v EmbeddingVector) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, v)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
