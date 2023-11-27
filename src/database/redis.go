package database

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

func Cache(client *redis.Client) *RedisCache {
	return &RedisCache{
		client: client,
	}
}

func (c *RedisCache) GetCache(key string) (string, error) {

	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	return val, nil
}

func (c *RedisCache) SetCache(key string, value interface{}, expiration time.Duration) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	// Set the JSON string value in the cache
	err = c.client.Set(ctx, key, jsonValue, expiration).Err()
	if err != nil {
		return err
	}

	return nil
}

// Clear a cache by key
func (c *RedisCache) ClearCache(key string) error {
	err := c.client.Del(ctx, key).Err()
	if err != nil {
		return err
	}

	return nil
}

var ctx = context.Background()

func RedisClient() *redis.Client {
	url := os.Getenv("REDIS_URL")
	opts, err := redis.ParseURL(url)
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(opts)

	// Ping the Redis server to check the connection
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}

	if pong == "PONG" {
		log.Printf("✅ Sucessfully connected to the Redis Database!")
	} else {
		log.Printf("❌ Failed to connect to the Redis Database!")
	}

	return rdb
}
