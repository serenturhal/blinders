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
	return m.RedisClient.SAdd(context.Background(), key, value).Err()
}

func (m *Manager) RemoveSession(userID string, connectionID string) error {
	key := ConstructUserKey(userID)
	value := ConstructConnectionKey(connectionID)
	return m.RedisClient.SRem(context.Background(), key, value).Err()
}

func (m *Manager) GetSessions(userID string) ([]string, error) {
	key := ConstructUserKey(userID)
	return m.RedisClient.SMembers(context.Background(), key).Result()
}
