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
	h.campaignService.StartCampaign()

	ctx.JSON(200, nil)
}
