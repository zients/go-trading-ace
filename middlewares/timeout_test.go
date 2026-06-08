package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTimeoutAddsDeadlineToRequestContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(Timeout(10 * time.Second))

	router.GET("/deadline", func(c *gin.Context) {
		deadline, ok := c.Request.Context().Deadline()
		require.True(t, ok)
		assert.LessOrEqual(t, time.Until(deadline), 10*time.Second)
		assert.Greater(t, time.Until(deadline), 9*time.Second)
		c.Status(http.StatusNoContent)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/deadline", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}
