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
	HomeController controllers.IHomeController
}

func NewHomeRoutes(r *gin.Engine, homeController controllers.IHomeController) IHomeRoutes {
	return &HomeRoutes{
		r:              r,
		HomeController: homeController,
	}
}

func (h *HomeRoutes) RegisterHomeRoutes() {
	h.r.GET("/", h.HomeController.Home)
}
