package session

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Sessions []string

type Manager struct {
	RedisClient *redis.Client
}

func NewManager(redisClient *redis.Client) *Manager {
	return &Manager{
		RedisClient: redisClient,
	}
}

func (m *Manager) AddSession(userID string, connectionID string) error {
	key := ConstructUserKey(userID)
	value := ConstructConnectionKey(connectionID)
	return m.RedisClient.Set(context.Background(), key, value, 0).Err()
}

func (m *Manager) GetSession(userID string) (string, error) {
	key := ConstructUserKey(userID)
	return m.RedisClient.Get(context.Background(), key).Result()
}
