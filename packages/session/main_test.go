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
		redisClient.Close()
	}
}

func TestRedisConn(t *testing.T) {
	manager, teardown := setup()
	defer teardown()

	err := manager.RedisClient.Ping(context.Background()).Err()
	assert.Nil(t, err)
}
