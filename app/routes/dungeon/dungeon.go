package dungeon

import (
	controller "dungeons/app/controllers/player"
	service "dungeons/app/services/player"

	"github.com/gin-gonic/gin"
)

func SetupRouter(g *gin.Engine) {

	servicesDungeon := service.New()
	dungeonController := controller.New(servicesDungeon)

	v1 := g.Group("/v1")
	{
		dungeons := v1.Group("/dungeons")
		{
			dungeons.POST("", dungeonController.Create)
			dungeons.GET("", dungeonController.Get)
			dungeons.GET("/:id", dungeonController.GetByID)
			dungeons.POST("/:id", dungeonController.Update)
			dungeons.GET("/IDS/:ids", dungeonController.GetByIDs)
		}
	}
}
