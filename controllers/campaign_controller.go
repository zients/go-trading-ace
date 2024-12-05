package controllers

import (
	"strconv"
	"trading-ace/config"
	"trading-ace/dtos"
	"trading-ace/services"

	"github.com/gin-gonic/gin"
)

type ICampaignController interface {
	StartCampaign(ctx *gin.Context)
	GetPointHistories(ctx *gin.Context)
	GetTaskStatus(ctx *gin.Context)
	GetLeaderboard(ctx *gin.Context)
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

	ctx.JSON(200, gin.H{"status": "ok"})
}

func (h *CampaignController) GetPointHistories(ctx *gin.Context) {
	address := ctx.Param("address")

	pointHistories, err := h.campaignService.GetPointHistories(address)
	if err != nil {
		ctx.JSON(500, gin.H{"status": "error", "message": err.Error()})
		return
	}

	results := []*dtos.GetPointHistoryDTO{}
	for _, data := range pointHistories {
		dto := &dtos.GetPointHistoryDTO{
			Task:        dtos.ConvertTaskToDTO(data.Task),
			TaskHistory: dtos.ConvertTaskHistoryToDTO(data.TaskHistory),
		}

		results = append(results, dto)
	}

	ctx.JSON(200, gin.H{"status": "ok", "result": results})
}

func (h *CampaignController) GetTaskStatus(ctx *gin.Context) {
	address := ctx.Param("address")

	taskStatus, err := h.campaignService.GetTaskStatus(address)
	if err != nil {
		ctx.JSON(500, gin.H{"status": "error", "message": err.Error()})
		return
	}

	results := []*dtos.TaskWithTaskHistoryDTO{}
	for _, v := range taskStatus {
		results = append(results, dtos.CovertTaskWithTaskHistoryToDTO(v))
	}

	ctx.JSON(200, gin.H{"status": "ok", "result": results})
}

func (h *CampaignController) GetLeaderboard(ctx *gin.Context) {
	taskName := ctx.Param("taskName")
	periodStr := ctx.Param("period")
	period, err := strconv.ParseInt(periodStr, 10, 32)
	if err != nil {
		ctx.JSON(500, gin.H{"status": "error", "message": err.Error()})
		return
	}

	leaderboardEntries, err := h.campaignService.GetLeaderboard(taskName, int(period))
	if err != nil {
		ctx.JSON(500, gin.H{"status": "error", "message": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"status": "ok", "result": leaderboardEntries})
}
