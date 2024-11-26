package logger

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func TestLogrusLogger_Info(t *testing.T) {
	// 設置測試鉤子來捕獲日誌
	logger := NewLogrusLogger()
	hook := test.NewGlobal()
	logger.(*LogrusLogger).logger.AddHook(hook)

	// 呼叫 Info 方法
	logger.Info("This is an info message")

	// 檢查 Info 日誌是否存在
	assert.Len(t, hook.Entries, 1)
	assert.Equal(t, logrus.InfoLevel, hook.Entries[0].Level)
	assert.Equal(t, "This is an info message", hook.Entries[0].Message)
}

func TestLogrusLogger_Warn(t *testing.T) {
	// 設置測試鉤子來捕獲日誌
	logger := NewLogrusLogger()
	hook := test.NewGlobal()
	logger.(*LogrusLogger).logger.AddHook(hook)

	// 呼叫 Warn 方法
	logger.Warn("This is a warn message")

	// 檢查 Warn 日誌是否存在
	assert.Len(t, hook.Entries, 1)
	assert.Equal(t, logrus.WarnLevel, hook.Entries[0].Level)
	assert.Equal(t, "This is a warn message", hook.Entries[0].Message)
}

func TestLogrusLogger_Error(t *testing.T) {
	// 設置測試鉤子來捕獲日誌
	logger := NewLogrusLogger()
	hook := test.NewGlobal()
	logger.(*LogrusLogger).logger.AddHook(hook)

	// 呼叫 Error 方法
	logger.Error("This is an error message")

	// 檢查 Error 日誌是否存在
	assert.Len(t, hook.Entries, 1)
	assert.Equal(t, logrus.ErrorLevel, hook.Entries[0].Level)
	assert.Equal(t, "This is an error message", hook.Entries[0].Message)
}
