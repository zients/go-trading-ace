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
	HSet(key string, field string, value interface{}) error
	HGet(key string, field string) (string, error)
	HIncrFloat(key string, field string, value float64) error
	ZAdd(key string, members ...*redis.Z) error
	ZRange(key string, start, stop int64) ([]string, error)
	ZRangeWithScores(key string, start, stop int64) ([]string, []float64, error)
	ZRevRange(key string, start, stop int64) ([]string, error)
	ZRevRangeWithScores(key string, start, stop int64) ([]string, []float64, error)
	SetTTL(key string, expiration time.Duration) error
}

type RedisHelper struct {
	redisClient *redis.Client
	prefix      string
}

func NewRedisHelper(redisClient *redis.Client, config *config.Config) IRedisHelper {
	return &RedisHelper{
		redisClient: redisClient,
		prefix:      config.Redis.Prefix,
	}
}

func (r *RedisHelper) Set(key string, value string, expiration time.Duration) error {
	err := r.redisClient.Set(context.Background(), r.prefix+key, value, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set key %s: %w", key, err)
	}

	return nil
}

func (r *RedisHelper) Get(key string) (string, error) {
	val, err := r.redisClient.Get(context.Background(), r.prefix+key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", fmt.Errorf("key %s does not exist", key)
		}

		return "", fmt.Errorf("failed to get key %s: %w", key, err)
	}

	return val, nil
}

func (r *RedisHelper) Delete(key string) error {
	err := r.redisClient.Del(context.Background(), r.prefix+key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete key %s: %w", key, err)
	}

	return nil
}

func (r *RedisHelper) HSet(key string, field string, value interface{}) error {
	err := r.redisClient.HSet(context.Background(), r.prefix+key, field, value).Err()
	if err != nil {
		return fmt.Errorf("failed to HSET field %s in key %s: %w", field, key, err)
	}

	return nil
}

func (r *RedisHelper) HGet(key string, field string) (string, error) {
	val, err := r.redisClient.HGet(context.Background(), r.prefix+key, field).Result()
	if err != nil {
		if err == redis.Nil {
			return "", fmt.Errorf("field %s does not exist in key %s", field, key)
		}

		return "", fmt.Errorf("failed to HGET field %s in key %s: %w", field, key, err)
	}

	return val, nil
}

func (r *RedisHelper) HIncrFloat(key string, field string, value float64) error {
	err := r.redisClient.HIncrByFloat(context.Background(), r.prefix+key, field, value).Err()
	if err != nil {
		return fmt.Errorf("failed to HIncrFloat field %s in key %s: %w", field, key, err)
	}

	return nil
}

func (r *RedisHelper) ZAdd(key string, members ...*redis.Z) error {
	err := r.redisClient.ZAdd(context.Background(), r.prefix+key, members...).Err()
	if err != nil {
		return fmt.Errorf("failed to ZADD to key %s: %w", key, err)
	}

	return nil
}

func (r *RedisHelper) ZRange(key string, start, stop int64) ([]string, error) {
	vals, err := r.redisClient.ZRange(context.Background(), r.prefix+key, start, stop).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to ZRANGE key %s: %w", key, err)
	}

	return vals, nil
}

func (r *RedisHelper) ZRangeWithScores(key string, start, stop int64) ([]string, []float64, error) {
	result, err := r.redisClient.ZRangeWithScores(context.Background(), r.prefix+key, start, stop).Result()
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
	vals, err := r.redisClient.ZRevRange(context.Background(), r.prefix+key, start, stop).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to ZRANGE key %s: %w", key, err)
	}

	return vals, nil
}

func (r *RedisHelper) ZRevRangeWithScores(key string, start, stop int64) ([]string, []float64, error) {
	result, err := r.redisClient.ZRevRangeWithScores(context.Background(), r.prefix+key, start, stop).Result()
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
	err := r.redisClient.Expire(context.Background(), r.prefix+key, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set key %s: %w", key, err)
	}

	return nil
}
