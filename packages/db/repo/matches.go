package repo

import (
	"context"
	"log"
	"time"

	"blinders/packages/db/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MatchesRepo struct {
	Col *mongo.Collection
}

func NewMatchesRepo(col *mongo.Collection) *MatchesRepo {
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

	return &MatchesRepo{
		Col: col,
	}
}

func (r *MatchesRepo) InsertNewRawMatchInfo(doc models.MatchInfo) (models.MatchInfo, error) {
	ctx, cal := context.WithTimeout(context.Background(), 5*time.Second)
	defer cal()

	_, err := r.Col.InsertOne(ctx, doc)

	return doc, err
}

func (r *MatchesRepo) GetMatchInfoByUserID(id primitive.ObjectID) (models.MatchInfo, error) {
	ctx, cal := context.WithTimeout(context.Background(), 5*time.Second)
	defer cal()

	var doc models.MatchInfo
	err := r.Col.FindOne(ctx, bson.M{"userID": id}).Decode(&doc)

	return doc, err
}

func (r *MatchesRepo) GetMatchInfoByFirebaseUID(uid string) (models.MatchInfo, error) {
	ctx, cal := context.WithTimeout(context.Background(), 5*time.Second)
	defer cal()

	var doc models.MatchInfo
	err := r.Col.FindOne(ctx, bson.M{"firebaseUID": uid}).Decode(&doc)

	return doc, err
}

// GetUsersByLanguage returns `numReturn` ID of users that speak one language of `learnings` and are currently learning `native` or are currently learning same language as user.
func (r *MatchesRepo) GetUsersByLanguage(firebaseUID string, numReturn uint32) ([]string, error) {
	user, err := r.GetMatchInfoByFirebaseUID(firebaseUID)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	stages := []bson.M{
		{"$match": bson.M{
			"firebaseUID": bson.M{"$ne": firebaseUID},
			"$or": []bson.M{
				{
					"native":    bson.M{"$in": user.Learnings},        // Users must speak at least one language of `learnings`.
					"learnings": bson.M{"$in": []string{user.Native}}, // Users should be learning their `native`.
				},
				{
					"learnings": bson.M{"$in": user.Learnings}, // Users who learn the same language as the current user.
				},
			},
		}},
		// at here we may sort users based on any ranking mark from the system.
		// currently, we random choose 1000 user.
		{
			"$sample": bson.M{"size": numReturn},
		},
		{"$project": bson.M{"_id": 0, "firebaseUID": 1}},
	}

	cur, err := r.Col.Aggregate(ctx, stages)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = cur.Close(ctx); err != nil {
			log.Panicf("hepo: cannot close cursor, err: %v", err)
		}
	}()

	type ReturnType struct {
		FirebaseUID string `bson:"firebaseUID"`
	}

	var ids []string
	for cur.Next(ctx) {
		doc := new(ReturnType)
		if err := cur.Decode(doc); err != nil {
			return nil, err
		}
		ids = append(ids, doc.FirebaseUID)
	}
	return ids, nil
}

func (r *MatchesRepo) DropUserWithFirebaseUID(firebaseUID string) (models.MatchInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	filter := bson.M{"firebaseUID": firebaseUID}
	res := r.Col.FindOneAndDelete(ctx, filter)
	if err := res.Err(); err != nil {
		return models.MatchInfo{}, err
	}
	var deletedUser models.MatchInfo
	if err := res.Decode(&deletedUser); err != nil {
		return models.MatchInfo{}, err
	}
	return deletedUser, nil
}
