package session

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/test-go/testify/assert"
)

func setup() (*Manager, func()) {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	manager := NewManager(redisClient)
	return manager, func() {
		_ = manager.RedisClient.Del(context.Background(), ConstructUserKey("1"))
		redisClient.Close()
	}
}

func TestRedisConn(t *testing.T) {
	manager, teardown := setup()
	defer teardown()

	err := manager.RedisClient.Ping(context.Background()).Err()
	assert.Nil(t, err)
}

func TestAddSession(t *testing.T) {
	manager, teardown := setup()
	defer teardown()

	err := manager.AddSession("1", "1")
	assert.Nil(t, err)

	value, err := manager.GetSessions("1")
	assert.Nil(t, err)
	assert.Contains(t, value, ConstructConnectionKey("1"))
}

func TestRemoveSession(t *testing.T) {
	manager, teardown := setup()
	defer teardown()

	err := manager.RemoveSession("1", "1")
	assert.Nil(t, err)
	value, err := manager.GetSessions("1")
	assert.Nil(t, err)
	assert.NotContains(t, value, ConstructConnectionKey("1"))
}

func TestGetSessions(t *testing.T) {
	manager, teardown := setup()
	defer teardown()

	_ = manager.AddSession("1", "1")
	_ = manager.AddSession("1", "2")
	value, err := manager.GetSessions("1")
	assert.Nil(t, err)
	assert.Len(t, value, 2)
	assert.Contains(t, value, ConstructConnectionKey("1"))
	assert.Contains(t, value, ConstructConnectionKey("2"))
}
