package mocks

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

type HomeControllerMock struct {
	mock.Mock
}

func (m *HomeControllerMock) Home(c *gin.Context) {
	m.Called(c)
}
