package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"trading-ace/config"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx/fxtest"
)

func TestNewGinServerUsesConfiguredRequestTimeout(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := NewGinServer(&config.Config{
		Server: config.ServerConfig{
			RequestTimeoutSeconds: 2,
		},
	})

	router.GET("/deadline", func(c *gin.Context) {
		deadline, ok := c.Request.Context().Deadline()
		require.True(t, ok)
		assert.LessOrEqual(t, time.Until(deadline), 2*time.Second)
		assert.Greater(t, time.Until(deadline), time.Second)
		c.Status(http.StatusNoContent)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/deadline", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestNewAppContextCancelsOnLifecycleStop(t *testing.T) {
	lifecycle := fxtest.NewLifecycle(t)

	ctx := NewAppContext(lifecycle)

	require.NoError(t, lifecycle.Start(contextWithTimeout(t)))
	assert.NoError(t, ctx.Err())

	require.NoError(t, lifecycle.Stop(contextWithTimeout(t)))
	assert.ErrorIs(t, ctx.Err(), context.Canceled)
}

func contextWithTimeout(t *testing.T) context.Context {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	t.Cleanup(cancel)
	return ctx
}
