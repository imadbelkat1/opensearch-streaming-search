package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func IsItemInCache(ctx context.Context, key string, targetID int) (bool, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     GetRedisConfig().Addr,
		Password: GetRedisConfig().Password,
		DB:       GetRedisConfig().DB,
	})

	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil // Key does not exist
		}
		return false, fmt.Errorf("failed to get value from Redis: %w", err)
	}

	var ids []int
	if err := json.Unmarshal([]byte(val), &ids); err != nil {
		return false, fmt.Errorf("failed to unmarshal IDs: %w", err)
	}

	for _, id := range ids {
		if id == targetID {
			return true, nil // ID found in the list
		}
	}

	return false, nil // ID not found in the list
}

func IsUserIDInCache(ctx context.Context, key string, targetID string) (bool, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     GetRedisConfig().Addr,
		Password: GetRedisConfig().Password,
		DB:       GetRedisConfig().DB,
	})

	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil // Key does not exist
		}
		return false, fmt.Errorf("failed to get value from Redis: %w", err)
	}

	var ids []string
	if err := json.Unmarshal([]byte(val), &ids); err != nil {
		return false, fmt.Errorf("failed to unmarshal IDs: %w", err)
	}

	for _, id := range ids {
		if id == targetID {
			return true, nil // ID found in the list
		}
	}

	return false, nil // ID not found in the list
}
