package match

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

type UserStore struct {
	Vector    []float32 `bson:"vector"`
	UserMatch `bson:",inline,omiempty"`
}
