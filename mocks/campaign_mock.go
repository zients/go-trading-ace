package mocks

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

type CampaignControllerMock struct {
	mock.Mock
}

func (m *CampaignControllerMock) StartCampaign(c *gin.Context) {
	m.Called(c)
}

func (m *CampaignControllerMock) GetPointHistories(c *gin.Context) {
	m.Called(c)
}
