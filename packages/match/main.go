package match

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Matcher interface {
	Match(ctx context.Context, fromID string, toID string) error
	// Suggest returns list of users that maybe match with given user
	Suggest(ctx context.Context, id string) ([]UserMatch, error)
	// TODO: Temporarily expose this method. Users should be automatically added after a new user event is fired from the user service.
	AddUser(ctx context.Context, user UserMatch) error
}

type UserMatch struct {
	UserID   string `json:"user_id" bson:"user_id"`
	Name     string `json:"name" bson:"name"`
	Gender   string `json:"gender" bson:"gender"`
	Native   string `json:"native" bson:"native"`
	Learning string `json:"learning" bson:"learning"`
	Age      int    `json:"age" bson:"age"`
}

type UserStore struct {
	UserMatch `bson:"data"`
	Vector    []float32 `bson:"vector"`
}

type MongoMatcher struct {
	UserCol  *mongo.Collection
	MatchCol *mongo.Collection
	Embedder Embedder
}

func (m *MongoMatcher) Match(ctx context.Context, fromID, toID string) error {
	// Here, we assume that users with ID fromID and toID already exist.
	if exists, _ := m._findMatchEntry(ctx, fromID, toID); exists {
		fmt.Println("already request")
		return nil
	}
	// check if toID user already match fromID user
	exists, err := m._findMatchEntry(ctx, toID, fromID)
	if err != nil {
		return err
	}
	if exists {
		return m._fulfillMatchRequest(ctx, fromID, toID)
	}

	return m._addMatchEntry(ctx, fromID, toID)
}

func (m *MongoMatcher) Suggest(ctx context.Context, fromID string) ([]UserMatch, error) {
	// TODO: Temporarily get all users
	return m.Get(ctx, bson.M{"data.user_id": bson.M{"$ne": fromID}})
}

func (m *MongoMatcher) AddUser(ctx context.Context, user UserMatch) error {
	embedding, err := m.Embedder.Embed(user)
	if err != nil {
		return err
	}

	userStore := UserStore{
		UserMatch: user,
		Vector:    embedding,
	}
	if _, err := m.UserCol.InsertOne(ctx, userStore); err != nil {
		return err
	}

	return nil
}

func (m *MongoMatcher) Get(ctx context.Context, filter bson.M) ([]UserMatch, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	cur, err := m.UserCol.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	res := []UserMatch{}
	for cur.Next(ctx) {
		user := new(UserStore)
		if err := cur.Decode(user); err != nil {
			return nil, err
		}
		res = append(res, user.UserMatch)
	}
	return res, nil
}

type matchEntry struct {
	FromID string `bson:"fromID"`
	ToID   string `bson:"toID"`
}

func (m *MongoMatcher) _findMatchEntry(ctx context.Context, fromID string, toID string) (bool, error) {
	filter := bson.M{"fromID": fromID, "toID": toID}
	res := m.MatchCol.FindOne(ctx, filter)
	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, fmt.Errorf("match: cannot find match entry, err: %v", res.Err())
	}
	entry := new(matchEntry)
	if err := res.Decode(entry); err != nil {
		return false, fmt.Errorf("match: cannot find match entry, err: %v", err)
	}

	return true, nil
}

func (m *MongoMatcher) _addMatchEntry(ctx context.Context, fromID string, toID string) error {
	entry := matchEntry{
		FromID: fromID,
		ToID:   toID,
	}
	if _, err := m.MatchCol.InsertOne(ctx, entry); err != nil {
		return fmt.Errorf("match: cannot add match entry, err: %v", err)
	}

	return nil
}

// _fulfillMatchRequest runs necessary processes after 2 users are matched.
func (m *MongoMatcher) _fulfillMatchRequest(_ context.Context, fromID string, toID string) error {
	fmt.Println("matched" + fromID + toID)
	return nil
}
