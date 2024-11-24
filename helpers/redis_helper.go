package helpers

import (
	"context"
	"fmt"
	"time"
	"trading-ace/config"

	"github.com/go-redis/redis/v8"
)

type IRedisHelper interface {
	Set(key string, value string, expiration time.Duration) error
	Get(key string) (string, error)
	Delete(key string) error
	HSet(hash string, key string, value interface{}) error
	HGet(hash string, key string) (string, error)
	ZAdd(key string, members ...*redis.Z) error
	ZRange(key string, start, stop int64) ([]string, error)
	ZRangeWithScores(key string, start, stop int64) ([]string, []float64, error)
	ZRevRange(key string, start, stop int64) ([]string, error)
	ZRevRangeWithScores(key string, start, stop int64) ([]string, []float64, error)
	SetTTL(key string, expiration time.Duration) error
}

type RedisHelper struct {
	RedisClient *redis.Client
	Prefix      string
}

func NewRedisHelper(redisClient *redis.Client, config *config.Config) IRedisHelper {
	return &RedisHelper{
		RedisClient: redisClient,
		Prefix:      config.Redis.Prefix,
	}
}

func (r *RedisHelper) Set(key string, value string, expiration time.Duration) error {
	err := r.RedisClient.Set(context.Background(), r.Prefix+key, value, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set key %s: %w", key, err)
	}

	return nil
}

func (r *RedisHelper) Get(key string) (string, error) {
	val, err := r.RedisClient.Get(context.Background(), r.Prefix+key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", fmt.Errorf("key %s does not exist", key)
		}

		return "", fmt.Errorf("failed to get key %s: %w", key, err)
	}

	return val, nil
}

func (r *RedisHelper) Delete(key string) error {
	err := r.RedisClient.Del(context.Background(), r.Prefix+key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete key %s: %w", key, err)
	}

	return nil
}

func (r *RedisHelper) HSet(hash string, key string, value interface{}) error {
	err := r.RedisClient.HSet(context.Background(), r.Prefix+hash, key, value).Err()
	if err != nil {
		return fmt.Errorf("failed to HSET field %s in hash %s: %w", key, hash, err)
	}
	return nil
}

func (r *RedisHelper) HGet(hash string, key string) (string, error) {
	val, err := r.RedisClient.HGet(context.Background(), r.Prefix+hash, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", fmt.Errorf("field %s does not exist in hash %s", key, hash)
		}
		return "", fmt.Errorf("failed to HGET field %s in hash %s: %w", key, hash, err)
	}
	return val, nil
}

func (r *RedisHelper) ZAdd(key string, members ...*redis.Z) error {
	err := r.RedisClient.ZAdd(context.Background(), r.Prefix+key, members...).Err()
	if err != nil {
		return fmt.Errorf("failed to ZADD to key %s: %w", key, err)
	}
	return nil
}

func (r *RedisHelper) ZRange(key string, start, stop int64) ([]string, error) {
	vals, err := r.RedisClient.ZRange(context.Background(), r.Prefix+key, start, stop).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to ZRANGE key %s: %w", key, err)
	}
	return vals, nil
}

func (r *RedisHelper) ZRangeWithScores(key string, start, stop int64) ([]string, []float64, error) {
	result, err := r.RedisClient.ZRangeWithScores(context.Background(), r.Prefix+key, start, stop).Result()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get range with scores from ZSET %s: %w", key, err)
	}

	var values []string
	var scores []float64

	for _, z := range result {
		values = append(values, z.Member.(string))
		scores = append(scores, z.Score)
	}

	return values, scores, nil
}

func (r *RedisHelper) ZRevRange(key string, start, stop int64) ([]string, error) {
	vals, err := r.RedisClient.ZRevRange(context.Background(), r.Prefix+key, start, stop).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to ZRANGE key %s: %w", key, err)
	}
	return vals, nil
}

func (r *RedisHelper) ZRevRangeWithScores(key string, start, stop int64) ([]string, []float64, error) {
	result, err := r.RedisClient.ZRevRangeWithScores(context.Background(), r.Prefix+key, start, stop).Result()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get range with scores from ZSET %s: %w", key, err)
	}

	var values []string
	var scores []float64

	for _, z := range result {
		values = append(values, z.Member.(string))
		scores = append(scores, z.Score)
	}

	return values, scores, nil
}

func (r *RedisHelper) SetTTL(key string, expiration time.Duration) error {
	err := r.RedisClient.Expire(context.Background(), r.Prefix+key, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set key %s: %w", key, err)
	}

	return nil
}
