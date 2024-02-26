package repo

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"blinders/packages/db/models"
)

type MatchsRepo struct {
	Col *mongo.Collection
}

func NewMatchsRepo(col *mongo.Collection) *MatchsRepo {
	ctx, cal := context.WithTimeout(context.Background(), time.Second*5)
	defer cal()

	_, err := col.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.M{"userID": 1},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		log.Println("can not create index for firebaseUID:", err)
		return nil
	}

	return &MatchsRepo{
		Col: col,
	}
}

func (r *MatchsRepo) InsertNewRawMatchInfo(doc models.MatchInfo) (models.MatchInfo, error) {
	ctx, cal := context.WithTimeout(context.Background(), 5*time.Second)
	defer cal()

	_, err := r.Col.InsertOne(ctx, doc)

	return doc, err
}

func (r *MatchsRepo) GetMatchInfoByUserID(id primitive.ObjectID) (models.MatchInfo, error) {
	ctx, cal := context.WithTimeout(context.Background(), 5*time.Second)
	defer cal()

	var doc models.MatchInfo
	err := r.Col.FindOne(ctx, bson.M{"userID": id}).Decode(&doc)

	return doc, err
}

func (r *MatchsRepo) GetUserByFirebaseUID(uid string) (models.MatchInfo, error) {
	ctx, cal := context.WithTimeout(context.Background(), 5*time.Second)
	defer cal()

	var doc models.MatchInfo
	err := r.Col.FindOne(ctx, bson.M{"firebaseUID": uid}).Decode(&doc)

	return doc, err
}
