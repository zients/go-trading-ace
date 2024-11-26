package controllers

import (
	"trading-ace/config"

	"github.com/gin-gonic/gin"
)

type IHomeController interface {
	Home(ctx *gin.Context)
}

type HomeController struct {
	config *config.Config
}

func NewHomeController(config *config.Config) IHomeController {
	return &HomeController{
		config: config,
	}
}

func (h *HomeController) Home(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"hello": "world"})
}
