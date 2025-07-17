package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

func CacheID(ctx context.Context, key string, ids []int) error {
	rdb := redis.NewClient(&redis.Options{
		Addr:     GetRedisConfig().Addr,
		Password: GetRedisConfig().Password,
		DB:       GetRedisConfig().DB,
	})

	idsJSON, err := json.Marshal(ids)
	if err != nil {
		return fmt.Errorf("failed to marshal IDs: %w", err)
	}

	err = rdb.Set(ctx, key, string(idsJSON), 0).Err()
	if err != nil {
		return fmt.Errorf("failed to set user IDs in Redis: %w", err)
	}

	log.Printf("Published %d user IDs to Redis", len(ids))
	return nil
}

func CacheUserIDs(ctx context.Context, key string, ids []string) error {
	rdb := redis.NewClient(&redis.Options{
		Addr:     GetRedisConfig().Addr,
		Password: GetRedisConfig().Password,
		DB:       GetRedisConfig().DB,
	})

	idsJSON, err := json.Marshal(ids)
	if err != nil {
		return fmt.Errorf("failed to marshal user IDs: %w", err)
	}

	err = rdb.Set(ctx, key, string(idsJSON), 0).Err()
	if err != nil {
		return fmt.Errorf("failed to set user IDs in Redis: %w", err)
	}

	log.Printf("Published %d user IDs to Redis", len(ids))
	return nil
}
