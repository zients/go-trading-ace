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

// MockController is a mock for the ICampaignController
type MockController struct {
	mock.Mock
}

func (m *MockController) StartCampaign(c *gin.Context) {
	m.Called(c)
}

func (m *MockController) GetPointHistories(c *gin.Context) {
	m.Called(c)
}

func TestCampaignRoutes(t *testing.T) {
	// Set up the Gin engine
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// Create a mock controller
	mockController := new(MockController)

	// Create the CampaignRoutes object
	campaignRoutes := routes.NewCampaignRoutes(r, mockController)

	// Register the routes
	campaignRoutes.RegisterCampaignRoutes()

	// Test the /campaign/start route
	t.Run("GET /campaign/start", func(t *testing.T) {
		// Set up the mock expectation
		mockController.On("StartCampaign", mock.Anything).Return()

		// Perform the request
		req, _ := http.NewRequest(http.MethodGet, "/campaign/start", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Assert the response
		assert.Equal(t, http.StatusOK, w.Code)
		mockController.AssertExpectations(t)
	})

	// Test the /campaign/histories/:address route
	t.Run("GET /campaign/histories/:address", func(t *testing.T) {
		// Set up the mock expectation
		mockController.On("GetPointHistories", mock.Anything).Return()

		// Perform the request
		req, _ := http.NewRequest(http.MethodGet, "/campaign/histories/123", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Assert the response
		assert.Equal(t, http.StatusOK, w.Code)
		mockController.AssertExpectations(t)
	})
}
