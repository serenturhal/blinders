package explore

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"blinders/packages/db"
	"blinders/packages/db/models"

	"github.com/redis/go-redis/v9"
)

type Explorer interface {
	// Suggest returns list of users that maybe match with given user
	Suggest(id string) ([]models.MatchInfo, error)
	// AddUserMatchInformation adds user match information to the database. A new user-created event must be fired for the embed worker to embed the recently added user.
	AddUserMatchInformation(info models.MatchInfo) (models.MatchInfo, error)
}

type MongoExplorer struct {
	Db          *db.MongoManager
	RedisClient *redis.Client
}

func NewMongoExplorer(Db *db.MongoManager, RedisClient *redis.Client) *MongoExplorer {
	return &MongoExplorer{
		Db:          Db,
		RedisClient: RedisClient,
	}
}

// currently, suggest suggests 5 users that are not friend of current user.
func (m *MongoExplorer) Suggest(fromID string) ([]models.MatchInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

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

func (m *MongoExplorer) AddUserMatchInformation(info models.MatchInfo) (models.MatchInfo, error) {
	user, err := m.Db.Users.GetUserByFirebaseUID(info.FirebaseUID)
	if err != nil {
		return models.MatchInfo{}, err
	}
	info.UserID = user.ID

	// duplicated match information will be handled by the repository since we have already indexed the collection with firebaseUID.
	matchInfo, err := m.Db.Matchs.InsertNewRawMatchInfo(info)
	if err != nil {
		return models.MatchInfo{}, err
	}
	return matchInfo, nil
}
