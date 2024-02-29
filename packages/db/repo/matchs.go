package repo

import (
	"context"
	"fmt"
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

	if _, err := col.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.M{"userID": 1},
		Options: options.Index().SetUnique(true),
	}); err != nil {
		log.Println("can not create index for userID:", err)
		return nil
	}

	if _, err := col.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.M{"firebaseUID": 1},
		Options: options.Index().SetUnique(true),
	}); err != nil {
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

// GetMatchCandidates returns `numReturn` ID of users that speak one language of `learnings` and are currently learning `native`.
func (r *MatchsRepo) GetUsersByLanguage(userID string) ([]string, error) {
	user, err := r.GetUserByFirebaseUID(userID)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	stages := []bson.M{
		{"$match": bson.M{
			"native":      bson.M{"$in": user.Learnings},        // user must speak at least one language of `learnings`
			"learnings":   bson.M{"$in": []string{user.Native}}, // user should learning `native`.
			"firebaseUID": bson.M{"$ne": userID},
		}},
		// at here we may sort users based on any rank mark from the system.
		// currently, we random choose 1000 user.
		{
			"$sample": bson.M{"size": 1000},
		},
		{"$project": bson.M{"_id": 0, "firebaseUID": 1}},
	}

	cur, err := r.Col.Aggregate(ctx, stages)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	type ReturnType struct {
		FirebaseUID string `bson:"firebaseUID"`
	}

	ids := []string{}
	for cur.Next(ctx) {
		doc := new(ReturnType)
		if err := cur.Decode(doc); err != nil {
			return nil, err
		}
		fmt.Println(doc)
		ids = append(ids, doc.FirebaseUID)
	}
	return ids, nil
}
