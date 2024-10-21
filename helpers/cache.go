package helpers

import (
	"encoding/json"
	"time"

	"products-api/config"

	"github.com/go-redis/redis/v8"
)

func SetCache(key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return config.RedisClient.Set(config.Ctx, key, data, ttl).Err()
}

func GetCache(key string, dest interface{}) (bool, error) {
	data, err := config.RedisClient.Get(config.Ctx, key).Result()
	if err == redis.Nil {
		return false, nil // Cache miss
	} else if err != nil {
		return false, err // Redis error
	}

	err = json.Unmarshal([]byte(data), dest)
	if err != nil {
		return false, err // Unmarshal error
	}
	return true, nil // Cache hit
}
