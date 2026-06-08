package helpers

import (
	"context"
	_ "embed"
	"fmt"
	"strconv"
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
	RecordSwapVolumeOnce(ctx context.Context, eventKey string, volumeKey string, totalKey string, address string, amount float64, expiration time.Duration) (float64, bool, error)
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

//go:embed lua/record_swap_volume_once.lua
var recordSwapVolumeOnceLua string

var recordSwapVolumeOnceScript = redis.NewScript(recordSwapVolumeOnceLua)

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

func (r *RedisHelper) RecordSwapVolumeOnce(ctx context.Context, eventKey string, volumeKey string, totalKey string, address string, amount float64, expiration time.Duration) (float64, bool, error) {
	ttlSeconds := int64(expiration.Seconds())
	if ttlSeconds < 1 {
		ttlSeconds = 1
	}

	result, err := recordSwapVolumeOnceScript.Run(
		ctx,
		r.redisClient,
		[]string{r.prefix + eventKey, r.prefix + volumeKey, r.prefix + totalKey},
		address,
		amount,
		ttlSeconds,
	).Result()
	if err != nil {
		return 0, false, fmt.Errorf("failed to record swap volume once: %w", err)
	}

	values, ok := result.([]interface{})
	if !ok || len(values) != 2 {
		return 0, false, fmt.Errorf("unexpected swap volume record result: %v", result)
	}

	recordedFlag, err := redisResultInt64(values[0])
	if err != nil {
		return 0, false, fmt.Errorf("unexpected swap volume recorded flag: %w", err)
	}
	if recordedFlag != 0 && recordedFlag != 1 {
		return 0, false, fmt.Errorf("unexpected swap volume recorded flag: %d", recordedFlag)
	}

	totalAmount, err := redisResultFloat64(values[1])
	if err != nil {
		return 0, false, fmt.Errorf("unexpected swap volume total amount: %w", err)
	}

	return totalAmount, recordedFlag == 1, nil
}

func redisResultInt64(value interface{}) (int64, error) {
	switch v := value.(type) {
	case int:
		return int64(v), nil
	case int64:
		return v, nil
	case string:
		return strconv.ParseInt(v, 10, 64)
	case []byte:
		return strconv.ParseInt(string(v), 10, 64)
	default:
		return 0, fmt.Errorf("unsupported type %T", value)
	}
}

func redisResultFloat64(value interface{}) (float64, error) {
	switch v := value.(type) {
	case float32:
		return float64(v), nil
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case string:
		return strconv.ParseFloat(v, 64)
	case []byte:
		return strconv.ParseFloat(string(v), 64)
	default:
		return 0, fmt.Errorf("unsupported type %T", value)
	}
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
