package routes

import (
	"trading-ace/controllers"

	"github.com/gin-gonic/gin"
)

type ICampaignRoutes interface {
	RegisterCampaignRoutes()
}

type CampaignRoutes struct {
	r                  *gin.Engine
	campaignController controllers.ICampaignController
}

func NewCampaignRoutes(r *gin.Engine, campaignController controllers.ICampaignController) ICampaignRoutes {
	return &CampaignRoutes{
		r:                  r,
		campaignController: campaignController,
	}
}

func (h *CampaignRoutes) RegisterCampaignRoutes() {
	group := h.r.Group("/campaign")

	group.GET("/start", h.campaignController.StartCampaign)
	group.GET("/histories/:address", h.campaignController.GetPointHistories)
}
