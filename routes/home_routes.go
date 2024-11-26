package routes

import (
	"trading-ace/controllers"

	"github.com/gin-gonic/gin"
)

type IHomeRoutes interface {
	RegisterHomeRoutes()
}

type HomeRoutes struct {
	r              *gin.Engine
	homeController controllers.IHomeController
}

func NewHomeRoutes(r *gin.Engine, homeController controllers.IHomeController) IHomeRoutes {
	return &HomeRoutes{
		r:              r,
		homeController: homeController,
	}
}

func (h *HomeRoutes) RegisterHomeRoutes() {
	h.r.GET("/", h.homeController.Home)
}
