package match

import (
	"context"
	"log"

	"blinders/packages/db/models"
	"blinders/packages/db/repo"

	"github.com/redis/go-redis/v9"
)

type Matcher interface {
	// Suggest returns list of users that maybe match with given user
	Suggest(ctx context.Context, id string) ([]models.MatchInfo, error)
	// TODO: Temporarily expose this method. User should be automatically added to match db after a new user event is fired from the user service.
	AddUserMatch(ctx context.Context, user models.MatchInfo) error
}

type MongoMatcher struct {
	Embedder    Embedder
	Repo        *repo.MatchsRepo
	RedisClient *redis.Client
}

func (m *MongoMatcher) InitIndex(ctx context.Context) {
	err := m.RedisClient.Do(
		ctx,
		"FT.CREATE",
		"idx:match_vss",
		"ON", "HASH",
		"PREFIX", "1", "match:",
		"SCHEMA", "embed", "VECTOR",
		"HNSW", "6",
		"TYPE", "FLOAT32",
		"DIM", "128",
		"DISTANCE_METRIC", "L2",
	).Err()
	if err != nil {
		log.Println(err)
	}
}

func (m *MongoMatcher) Suggest(ctx context.Context, fromID string) ([]models.MatchInfo, error) {
	// TODO: Temporarily  get 5 random users
	user, err := m.Repo.GetUserByFirebaseUID(fromID)
	if err != nil {
		return nil, err
	}
	embed, err := m.Embedder.Embed(user)
	if err != nil {
		return nil, err
	}

	slice := m.RedisClient.Do(ctx,
		"FT.SEARCH",
		"idx:match_vss",
		"(*)=>[KNN 5 @embed $vec]",
		"PARAMS", "2",
		"vec", embed,
		"DIALECT", "2",
		"RETURN", "1", "id",
	).Val().(map[any]any)

	res := []models.MatchInfo{}
	for _, doc := range slice["results"].([]any) {
		userID := doc.(map[any]any)["extra_attributes"].(map[any]any)["id"].(string)
		user, err := m.Repo.GetUserByFirebaseUID(userID)
		if err != nil {
			return nil, err
		}
		res = append(res, user)
	}
	return res, nil
}

func (m *MongoMatcher) AddUserMatch(ctx context.Context, user models.MatchInfo) error {
	embedding, err := m.Embedder.Embed(user)
	if err != nil {
		return err
	}

	if err := m.RedisClient.HSet(
		ctx,
		CreateMatchKeyWithUserID(user.FirebaseUID),
		"embed", embedding,
		"id", user.FirebaseUID,
	).Err(); err != nil {
		return err
	}

	if _, err := m.Repo.InsertNewRawMatchInfo(user); err != nil {
		return err
	}

	return nil
}
