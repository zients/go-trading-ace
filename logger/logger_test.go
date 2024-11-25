package logger

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func TestLogrusLogger(t *testing.T) {
	// Create a new Logrus logger
	logger := NewLogrusLogger()

	// Create a hook to capture logs
	hook := test.NewLocal(logger.(*LogrusLogger).logger)

	logger.Info("This is an info message")
	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, logrus.InfoLevel, hook.LastEntry().Level)
	assert.Equal(t, "This is an info message", hook.LastEntry().Message)

	logger.Warn("This is a warning message")
	assert.Equal(t, 2, len(hook.Entries))
	assert.Equal(t, logrus.WarnLevel, hook.Entries[1].Level)
	assert.Equal(t, "This is a warning message", hook.Entries[1].Message)

	logger.Error("This is an error message")
	assert.Equal(t, 3, len(hook.Entries))
	assert.Equal(t, logrus.ErrorLevel, hook.Entries[2].Level)
	assert.Equal(t, "This is an error message", hook.Entries[2].Message)

	logger.Debug("This is a debug message")
	assert.Equal(t, 3, len(hook.Entries))
}
