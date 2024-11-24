package controllers

import (
	"trading-ace/config"
	"trading-ace/services"

	"github.com/gin-gonic/gin"
)

type ICampaignController interface {
	StartCampaign(ctx *gin.Context)
}

type CampaignController struct {
	config          *config.Config
	campaignService services.ICampaignService
}

func NewCampaignController(config *config.Config, campaignService services.ICampaignService) ICampaignController {
	return &CampaignController{
		config:          config,
		campaignService: campaignService,
	}
}

func (h *CampaignController) StartCampaign(ctx *gin.Context) {
	if err := h.campaignService.StartCampaign(); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, nil)
}
