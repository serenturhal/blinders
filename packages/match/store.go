package match

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Store interface {
	Add(ctx context.Context, id string, info any) (primitive.ObjectID, error)
	GetAll(ctx context.Context) ([]UserStore, error)
}

type MongoStore struct {
	UserCol *mongo.Collection
}

type UserStore struct {
	Info   any                `bson:"info"`
	UserID string             `bson:"userId"`
	ID     primitive.ObjectID `bson:"_id,omiempty"`
}

func (s *MongoStore) Add(ctx context.Context, id string, info any) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	userStore := UserStore{
		UserID: id,
		Info:   info,
	}
	res, err := s.UserCol.InsertOne(ctx, userStore)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return res.InsertedID.(primitive.ObjectID), nil
}

func (s *MongoStore) GetAll(ctx context.Context) ([]UserStore, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	cur, err := s.UserCol.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	res := []UserStore{}
	for cur.Next(ctx) {
		user := new(UserStore)
		if err := cur.Decode(user); err != nil {
			return nil, err
		}
		res = append(res, *user)
	}
	return res, nil
}
