package run

import (
	controller "dungeons/app/controllers/player"
	service "dungeons/app/services/player"

	"github.com/gin-gonic/gin"
)

func SetupRouter(g *gin.Engine) {

	servicesRun := service.New()
	runController := controller.New(servicesRun)

	v1 := g.Group("/v1")
	{
		runs := v1.Group("/runs")
		{
			runs.POST("", runController.Create)
			runs.GET("", runController.Get)
			runs.GET("/:id", runController.GetByID)
			runs.POST("/:id", runController.Update)
			runs.GET("/IDS/:ids", runController.GetByIDs)
		}
	}
}
