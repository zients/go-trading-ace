package helpers

import (
	"errors"
	"testing"
	"time"

	"trading-ace/config"

	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
)

func setupRedisHelper() (IRedisHelper, redismock.ClientMock) {
	redisClient, mock := redismock.NewClientMock()
	return NewRedisHelper(redisClient, &config.Config{
		Redis: config.RedisConfig{Prefix: "test:"},
	}), mock
}

func TestRedisHelper_Set(t *testing.T) {
	r, mock := setupRedisHelper()

	key := "key"
	value := "value"
	expiration := time.Minute

	mock.ExpectSet("test:"+key, value, expiration).SetVal("OK")

	err := r.Set(key, value, expiration)
	assert.NoError(t, err)

	// Simulate failure
	mock.ExpectSet("test:"+key, value, expiration).SetErr(errors.New("redis error"))

	err = r.Set(key, value, expiration)
	assert.Error(t, err)
}

func TestRedisHelper_Get(t *testing.T) {
	r, mock := setupRedisHelper()

	key := "key"
	expected := "value"

	mock.ExpectGet("test:" + key).SetVal(expected)

	val, err := r.Get(key)
	assert.NoError(t, err)
	assert.Equal(t, expected, val)

	// Simulate key does not exist
	mock.ExpectGet("test:" + key).RedisNil()

	_, err = r.Get(key)
	assert.Error(t, err)

	// Simulate redis error
	mock.ExpectGet("test:" + key).SetErr(errors.New("redis error"))

	_, err = r.Get(key)
	assert.Error(t, err)
}

func TestRedisHelper_Delete(t *testing.T) {
	r, mock := setupRedisHelper()

	key := "key"

	mock.ExpectDel("test:" + key).SetVal(1)

	err := r.Delete(key)
	assert.NoError(t, err)

	// Simulate redis error
	mock.ExpectDel("test:" + key).SetErr(errors.New("redis error"))

	err = r.Delete(key)
	assert.Error(t, err)
}

func TestRedisHelper_IncrFloat(t *testing.T) {
	r, mock := setupRedisHelper()

	key := "key"
	increment := 1.5

	mock.ExpectIncrByFloat("test:"+key, increment).SetVal(2.5)

	err := r.IncrFloat(key, increment)
	assert.NoError(t, err)

	// Simulate redis error
	mock.ExpectIncrByFloat("test:"+key, increment).SetErr(errors.New("redis error"))

	err = r.IncrFloat(key, increment)
	assert.Error(t, err)
}

func TestRedisHelper_HSet(t *testing.T) {
	r, mock := setupRedisHelper()

	key := "key"
	field := "field"
	value := "value"

	mock.ExpectHSet("test:"+key, field, value).SetVal(1)

	err := r.HSet(key, field, value)
	assert.NoError(t, err)

	// Simulate redis error
	mock.ExpectHSet("test:"+key, field, value).SetErr(errors.New("redis error"))

	err = r.HSet(key, field, value)
	assert.Error(t, err)
}

func TestRedisHelper_HGet(t *testing.T) {
	r, mock := setupRedisHelper()

	key := "key"
	field := "field"
	expected := "value"

	mock.ExpectHGet("test:"+key, field).SetVal(expected)

	val, err := r.HGet(key, field)
	assert.NoError(t, err)
	assert.Equal(t, expected, val)

	// Simulate field does not exist
	mock.ExpectHGet("test:"+key, field).RedisNil()

	_, err = r.HGet(key, field)
	assert.Error(t, err)

	// Simulate redis error
	mock.ExpectHGet("test:"+key, field).SetErr(errors.New("redis error"))

	_, err = r.HGet(key, field)
	assert.Error(t, err)
}

func TestRedisHelper_HIncrFloat(t *testing.T) {
	r, mock := setupRedisHelper()

	key := "key"
	field := "field"
	increment := 1.5

	mock.ExpectHIncrByFloat("test:"+key, field, increment).SetVal(2.5)

	err := r.HIncrFloat(key, field, increment)
	assert.NoError(t, err)

	// Simulate redis error
	mock.ExpectHIncrByFloat("test:"+key, field, increment).SetErr(errors.New("redis error"))

	err = r.HIncrFloat(key, field, increment)
	assert.Error(t, err)
}

func TestRedisHelper_ZAdd(t *testing.T) {
	r, mock := setupRedisHelper()

	key := "key"
	member := redis.Z{Score: 1.0, Member: "value"}

	mock.ExpectZAdd("test:"+key, &member).SetVal(1)

	err := r.ZAdd(key, &member)
	assert.NoError(t, err)

	// Simulate redis error
	mock.ExpectZAdd("test:"+key, &member).SetErr(errors.New("redis error"))

	err = r.ZAdd(key, &member)
	assert.Error(t, err)
}

func TestRedisHelper_ZRange(t *testing.T) {
	r, mock := setupRedisHelper()

	key := "key"
	expected := []string{"value1", "value2"}

	mock.ExpectZRange("test:"+key, 0, -1).SetVal(expected)

	vals, err := r.ZRange(key, 0, -1)
	assert.NoError(t, err)
	assert.Equal(t, expected, vals)

	// Simulate redis error
	mock.ExpectZRange("test:"+key, 0, -1).SetErr(errors.New("redis error"))

	_, err = r.ZRange(key, 0, -1)
	assert.Error(t, err)
}

func TestRedisHelper_SetTTL(t *testing.T) {
	r, mock := setupRedisHelper()

	key := "key"
	expiration := time.Minute

	mock.ExpectExpire("test:"+key, expiration).SetVal(true)

	err := r.SetTTL(key, expiration)
	assert.NoError(t, err)

	// Simulate redis error
	mock.ExpectExpire("test:"+key, expiration).SetErr(errors.New("redis error"))

	err = r.SetTTL(key, expiration)
	assert.Error(t, err)
}
