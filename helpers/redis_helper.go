package helpers

import (
	"context"
	"fmt"
	"time"
	"trading-ace/config"

	"github.com/go-redis/redis/v8"
)

type IRedisHelper interface {
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
	IncrFloat(ctx context.Context, key string, value float64) error
	HSet(ctx context.Context, key string, field string, value interface{}) error
	HGet(ctx context.Context, key string, field string) (string, error)
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	HIncrFloat(ctx context.Context, key string, field string, value float64) error
	ZAdd(ctx context.Context, key string, members ...*redis.Z) error
	ZRange(ctx context.Context, key string, start, stop int64) ([]string, error)
	ZRangeWithScores(ctx context.Context, key string, start, stop int64) ([]string, []float64, error)
	ZRevRange(ctx context.Context, key string, start, stop int64) ([]string, error)
	ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) ([]string, []float64, error)
	SetTTL(ctx context.Context, key string, expiration time.Duration) error
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

func (r *RedisHelper) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	err := r.redisClient.Set(ctx, r.prefix+key, value, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set key %s: %w", key, err)
	}

	return nil
}

func (r *RedisHelper) Get(ctx context.Context, key string) (string, error) {
	val, err := r.redisClient.Get(ctx, r.prefix+key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", fmt.Errorf("key %s does not exist", key)
		}

		return "", fmt.Errorf("failed to get key %s: %w", key, err)
	}

	return val, nil
}

func (r *RedisHelper) Delete(ctx context.Context, key string) error {
	err := r.redisClient.Del(ctx, r.prefix+key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete key %s: %w", key, err)
	}

	return nil
}

func (r *RedisHelper) IncrFloat(ctx context.Context, key string, value float64) error {
	err := r.redisClient.IncrByFloat(ctx, r.prefix+key, value).Err()
	if err != nil {
		return fmt.Errorf("failed to IncrFloat key %s: %w", key, err)
	}

	return nil
}

func (r *RedisHelper) HSet(ctx context.Context, key string, field string, value interface{}) error {
	err := r.redisClient.HSet(ctx, r.prefix+key, field, value).Err()
	if err != nil {
		return fmt.Errorf("failed to HSET field %s in key %s: %w", field, key, err)
	}

	return nil
}

func (r *RedisHelper) HGet(ctx context.Context, key string, field string) (string, error) {
	val, err := r.redisClient.HGet(ctx, r.prefix+key, field).Result()
	if err != nil {
		if err == redis.Nil {
			return "", fmt.Errorf("field %s does not exist in key %s", field, key)
		}

		return "", fmt.Errorf("failed to HGET field %s in key %s: %w", field, key, err)
	}

	return val, nil
}

func (r *RedisHelper) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	result, err := r.redisClient.HGetAll(ctx, r.prefix+key).Result()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *RedisHelper) HIncrFloat(ctx context.Context, key string, field string, value float64) error {
	err := r.redisClient.HIncrByFloat(ctx, r.prefix+key, field, value).Err()
	if err != nil {
		return fmt.Errorf("failed to HIncrFloat field %s in key %s: %w", field, key, err)
	}

	return nil
}

func (r *RedisHelper) ZAdd(ctx context.Context, key string, members ...*redis.Z) error {
	err := r.redisClient.ZAdd(ctx, r.prefix+key, members...).Err()
	if err != nil {
		return fmt.Errorf("failed to ZADD to key %s: %w", key, err)
	}

	return nil
}

func (r *RedisHelper) ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	vals, err := r.redisClient.ZRange(ctx, r.prefix+key, start, stop).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to ZRANGE key %s: %w", key, err)
	}

	return vals, nil
}

func (r *RedisHelper) ZRangeWithScores(ctx context.Context, key string, start, stop int64) ([]string, []float64, error) {
	result, err := r.redisClient.ZRangeWithScores(ctx, r.prefix+key, start, stop).Result()
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

func (r *RedisHelper) ZRevRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	vals, err := r.redisClient.ZRevRange(ctx, r.prefix+key, start, stop).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to ZRANGE key %s: %w", key, err)
	}

	return vals, nil
}

func (r *RedisHelper) ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) ([]string, []float64, error) {
	result, err := r.redisClient.ZRevRangeWithScores(ctx, r.prefix+key, start, stop).Result()
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

func (r *RedisHelper) SetTTL(ctx context.Context, key string, expiration time.Duration) error {
	err := r.redisClient.Expire(ctx, r.prefix+key, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set key %s: %w", key, err)
	}

	return nil
}
