package routes

import (
	"fmt"
	"trading-ace/config"
	"trading-ace/controllers"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

type IHomeRoutes interface {
	RegisterHomeRoutes()
}

type HomeRoutes struct {
	r              *gin.Engine
	homeController controllers.IHomeController
	config         *config.Config
}

func NewHomeRoutes(r *gin.Engine, homeController controllers.IHomeController, config *config.Config) IHomeRoutes {
	return &HomeRoutes{
		r:              r,
		homeController: homeController,
		config:         config,
	}
}

func (h *HomeRoutes) RegisterHomeRoutes() {
	h.r.GET("/", h.homeController.Home)

	h.r.GET("/swagger.json", func(c *gin.Context) {
		c.File("./docs/swagger.json")
	})

	url := ginSwagger.URL(fmt.Sprintf("http://localhost:%d/swagger.json", h.config.Server.Port))
	h.r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
}
