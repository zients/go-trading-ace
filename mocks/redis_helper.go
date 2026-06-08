package mocks

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/mock"
)

// MockRedisHelper 用來模擬 IRedisHelper 接口
type MockRedisHelper struct {
	mock.Mock
}

func (m *MockRedisHelper) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}

func (m *MockRedisHelper) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *MockRedisHelper) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *MockRedisHelper) IncrFloat(ctx context.Context, key string, value float64) error {
	args := m.Called(ctx, key, value)
	return args.Error(0)
}

func (m *MockRedisHelper) HSet(ctx context.Context, key string, field string, value interface{}) error {
	args := m.Called(ctx, key, field, value)
	return args.Error(0)
}

func (m *MockRedisHelper) HGet(ctx context.Context, key string, field string) (string, error) {
	args := m.Called(ctx, key, field)
	return args.String(0), args.Error(1)
}

func (m *MockRedisHelper) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	args := m.Called(ctx, key)
	return args.Get(0).(map[string]string), args.Error(1)
}

func (m *MockRedisHelper) HIncrFloat(ctx context.Context, key string, field string, value float64) error {
	args := m.Called(ctx, key, field, value)
	return args.Error(0)
}

func (m *MockRedisHelper) RecordSwapVolumeOnce(ctx context.Context, eventKey string, volumeKey string, totalKey string, address string, amount float64, expiration time.Duration) (float64, bool, error) {
	args := m.Called(ctx, eventKey, volumeKey, totalKey, address, amount, expiration)
	return args.Get(0).(float64), args.Bool(1), args.Error(2)
}

func (m *MockRedisHelper) ZAdd(ctx context.Context, key string, members ...*redis.Z) error {
	args := m.Called(ctx, key, members)
	return args.Error(0)
}

func (m *MockRedisHelper) ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	args := m.Called(ctx, key, start, stop)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockRedisHelper) ZRangeWithScores(ctx context.Context, key string, start, stop int64) ([]string, []float64, error) {
	args := m.Called(ctx, key, start, stop)
	return args.Get(0).([]string), args.Get(1).([]float64), args.Error(2)
}

func (m *MockRedisHelper) ZRevRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	args := m.Called(ctx, key, start, stop)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockRedisHelper) ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) ([]string, []float64, error) {
	args := m.Called(ctx, key, start, stop)
	return args.Get(0).([]string), args.Get(1).([]float64), args.Error(2)
}

func (m *MockRedisHelper) SetTTL(ctx context.Context, key string, expiration time.Duration) error {
	args := m.Called(ctx, key, expiration)
	return args.Error(0)
}
