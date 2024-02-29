package match

import (
	"context"
	"encoding/json"
	"fmt"

	"blinders/packages/db"
	"blinders/packages/db/models"

	"github.com/redis/go-redis/v9"
)

type Matcher interface {
	// Suggest returns list of users that maybe match with given user
	Suggest(ctx context.Context, id string) ([]models.MatchInfo, error)
}

type MongoMatcher struct {
	Db          *db.MongoManager
	RedisClient *redis.Client
}

func NewMongoMatcher(Db *db.MongoManager, RedisClient *redis.Client) *MongoMatcher {
	return &MongoMatcher{
		Db:          Db,
		RedisClient: RedisClient,
	}
}

// currently, suggest suggests 5 users that are not friend of current user.
func (m *MongoMatcher) Suggest(ctx context.Context, fromID string) ([]models.MatchInfo, error) {
	// Get 1000 users that maybe match with current user (user that speak and learn language)
	user, err := m.Db.Users.GetUserByFirebaseUID(fromID)
	if err != nil {
		return nil, err
	}

	jsonStr, _ := m.RedisClient.JSONGet(ctx, CreateMatchKeyWithUserID(fromID), "$.embed").Result()
	embedArr := []EmbeddingVector{}
	if err := json.Unmarshal([]byte(jsonStr), &embedArr); err != nil {
		return nil, err
	}
	embed := embedArr[0]

	// exclude friends of current user
	// TODO: Need to optimize; currently, excluding 700 users takes 230ms on M1 chip with 16GB RAM.
	excludeFilter := fromID
	for _, friendID := range user.FriendsFirebaseUID {
		excludeFilter += " | " + friendID
	}
	excludeFilter = fmt.Sprintf("-@id:(%s)", excludeFilter)

	candidates, err := m.Db.Matchs.GetUsersByLanguage(user.FirebaseUID)
	if err != nil {
		return nil, err
	}

	includeFilter := ""
	if len(candidates) != 0 {
		includeFilter = candidates[0]
		for idx := 1; idx < len(candidates); idx++ {
			includeFilter += " | " + candidates[idx]
		}
		includeFilter = fmt.Sprintf("@id:(%s)", includeFilter)
	}

	prefilter := fmt.Sprintf("(%s %s)", excludeFilter, includeFilter)

	cmd := m.RedisClient.Do(ctx,
		"FT.SEARCH",
		"idx:match_vss",
		fmt.Sprintf("%s=>[KNN 5 @embed $query_vector as vector_score]", prefilter),
		"SORTBY", "vector_score",
		"PARAMS", "2",
		"query_vector", &embed,
		"DIALECT", "2",
		"RETURN", "1", "id",
	)
	if err := cmd.Err(); err != nil {
		fmt.Println(err)
		return nil, err
	}

	res := []models.MatchInfo{}
	for _, doc := range cmd.Val().(map[any]any)["results"].([]any) {
		userID := doc.(map[any]any)["extra_attributes"].(map[any]any)["id"].(string)
		user, err := m.Db.Matchs.GetUserByFirebaseUID(userID)
		if err != nil {
			return nil, err
		}
		res = append(res, user)
	}
	return res, nil
}
