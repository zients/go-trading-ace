package controllers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"trading-ace/config"
	"trading-ace/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestStartCampaignPassesRequestContextToService(t *testing.T) {
	gin.SetMode(gin.TestMode)

	campaignService := new(mocks.MockCampaignService)
	controller := NewCampaignController(&config.Config{}, campaignService)
	requestKey := struct{}{}
	requestCtx := context.WithValue(context.Background(), requestKey, "request-context")

	campaignService.On("StartCampaign", mock.MatchedBy(func(actual context.Context) bool {
		return actual == requestCtx && actual.Value(requestKey) == "request-context"
	})).Return(nil).Once()

	w := httptest.NewRecorder()
	ginCtx, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodGet, "/campaign/start", nil)
	ginCtx.Request = req.WithContext(requestCtx)

	controller.StartCampaign(ginCtx)

	assert.Equal(t, http.StatusOK, w.Code)
	campaignService.AssertExpectations(t)
}
