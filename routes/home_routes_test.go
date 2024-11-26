package routes_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"trading-ace/routes"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockController is a mock for the IHomeController
type MockHomeController struct {
	mock.Mock
}

func (m *MockHomeController) Home(c *gin.Context) {
	m.Called(c)
}

func TestHomeRoutes(t *testing.T) {
	// Set up the Gin engine
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// Create a mock home controller
	mockHomeController := new(MockHomeController)

	// Create the HomeRoutes object
	homeRoutes := routes.NewHomeRoutes(r, mockHomeController)

	// Register the routes
	homeRoutes.RegisterHomeRoutes()

	// Test the / route (Home route)
	t.Run("GET /", func(t *testing.T) {
		// Set up the mock expectation
		mockHomeController.On("Home", mock.Anything).Return()

		// Perform the request
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Assert the response
		assert.Equal(t, http.StatusOK, w.Code)
		mockHomeController.AssertExpectations(t)
	})
}
