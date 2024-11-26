package mocks

import (
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/mock"
)

// MockRedisHelper 用來模擬 IRedisHelper 接口
type MockRedisHelper struct {
	mock.Mock
}

func (m *MockRedisHelper) Set(key string, value string, expiration time.Duration) error {
	args := m.Called(key, value, expiration)
	return args.Error(0)
}

func (m *MockRedisHelper) Get(key string) (string, error) {
	args := m.Called(key)
	return args.String(0), args.Error(1)
}

func (m *MockRedisHelper) Delete(key string) error {
	args := m.Called(key)
	return args.Error(0)
}

func (m *MockRedisHelper) IncrFloat(key string, value float64) error {
	args := m.Called(key, value)
	return args.Error(0)
}

func (m *MockRedisHelper) HSet(key string, field string, value interface{}) error {
	args := m.Called(key, field, value)
	return args.Error(0)
}

func (m *MockRedisHelper) HGet(key string, field string) (string, error) {
	args := m.Called(key, field)
	return args.String(0), args.Error(1)
}

func (m *MockRedisHelper) HGetAll(key string) (map[string]string, error) {
	args := m.Called(key)
	return args.Get(0).(map[string]string), args.Error(1)
}

func (m *MockRedisHelper) HIncrFloat(key string, field string, value float64) error {
	args := m.Called(key, field, value)
	return args.Error(0)
}

func (m *MockRedisHelper) ZAdd(key string, members ...*redis.Z) error {
	args := m.Called(key, members)
	return args.Error(0)
}

func (m *MockRedisHelper) ZRange(key string, start, stop int64) ([]string, error) {
	args := m.Called(key, start, stop)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockRedisHelper) ZRangeWithScores(key string, start, stop int64) ([]string, []float64, error) {
	args := m.Called(key, start, stop)
	return args.Get(0).([]string), args.Get(1).([]float64), args.Error(2)
}

func (m *MockRedisHelper) ZRevRange(key string, start, stop int64) ([]string, error) {
	args := m.Called(key, start, stop)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockRedisHelper) ZRevRangeWithScores(key string, start, stop int64) ([]string, []float64, error) {
	args := m.Called(key, start, stop)
	return args.Get(0).([]string), args.Get(1).([]float64), args.Error(2)
}

func (m *MockRedisHelper) SetTTL(key string, expiration time.Duration) error {
	args := m.Called(key, expiration)
	return args.Error(0)
}
