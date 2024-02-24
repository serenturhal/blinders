package match

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Matcher interface {
	Match(ctx context.Context, fromID string, toID string) error
	// Suggest returns list of users that maybe match with given user
	Suggest(ctx context.Context, id string) ([]UserMatch, error)
	// TODO: Temporarily expose this method. User should be automatically added to match db after a new user event is fired from the user service.
	AddUserMatch(ctx context.Context, user UserMatch) error
}

type MongoMatcher struct {
	UserCol     *mongo.Collection
	MatchCol    *mongo.Collection
	Embedder    Embedder
	RedisClient *redis.Client
}

func (m *MongoMatcher) InitIndex(ctx context.Context) {
	err := m.RedisClient.Do(
		ctx,
		"FT.CREATE",
		"my_idx",
		"ON", "HASH",
		"PREFIX", "1", "match:",
		"SCHEMA", "embed",
		"VECTOR", "HNSW", "6", "TYPE", "FLOAT32", "DIM", "128", "DISTANCE_METRIC", "L2",
	).Err()

	fmt.Println(err)
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
	// TODO: Temporarily  get 5 random users
	user, err := m._getUserMatchWithID(ctx, fromID)
	if err != nil {
		return nil, err
	}
	embed, err := m.Embedder.Embed(*user)
	if err != nil {
		return nil, err
	}

	slice := m.RedisClient.Do(ctx,
		"FT.SEARCH",
		"my_idx",
		"(*)=>[KNN 6 @embed $vec]", "PARAMS", "2", "vec", embed, "DIALECT", "2",
		"RETURN", "1", "id",
	).Val().(map[any]any)

	res := []UserMatch{}
	for _, doc := range slice["results"].([]any) {
		userID := doc.(map[any]any)["extra_attributes"].(map[any]any)["id"].(string)
		user, err := m._getUserMatchWithID(ctx, userID)
		if err != nil {
			return nil, err
		}
		res = append(res, *user)
	}
	return res, nil
}

func (m *MongoMatcher) AddUserMatch(ctx context.Context, user UserMatch) error {
	embedding, err := m.Embedder.Embed(user)
	if err != nil {
		return err
	}

	userStore := UserStore{
		UserMatch: user,
		Vector:    embedding,
	}
	if err := m.RedisClient.HSet(ctx, fmt.Sprintf("match:%v", user.UserID), "embed", embedding, "id", user.UserID).Err(); err != nil {
		return err
	}

	if _, err := m.UserCol.InsertOne(ctx, userStore); err != nil {
		return err
	}

	return nil
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

func (m *MongoMatcher) _getUserMatchWithID(ctx context.Context, userID string) (*UserMatch, error) {
	filter := bson.M{"userID": userID}
	cur := m.UserCol.FindOne(ctx, filter)
	if err := cur.Err(); err != nil {
		return nil, err
	}

	user := new(UserMatch)
	if err := cur.Decode(user); err != nil {
		return nil, err
	}
	return user, nil
}
