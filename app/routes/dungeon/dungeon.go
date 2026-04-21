package dungeon

import (
	controller "dungeons/app/controllers/dungeon"
	service "dungeons/app/services/dungeon"

	"github.com/gin-gonic/gin"
)

func SetupRouter(g *gin.Engine) {
	servicesDungeon := service.New()
	dungeonController := controller.New(servicesDungeon)

	v1 := g.Group("/v1")
	{
		dungeons := v1.Group("/dungeons")
		{
			dungeons.GET("", dungeonController.Get)
			dungeons.GET("/:id", dungeonController.GetByID)
			dungeons.GET("/ids/:ids", dungeonController.GetByIDs)
		}

		mjDungeons := v1.Group("/mj/dungeons")
		{
			mjDungeons.POST("", dungeonController.Create)
			mjDungeons.PATCH("/:id", dungeonController.Update)
			mjDungeons.POST("/:id/publish", dungeonController.Publish)
			mjDungeons.DELETE("/:id", dungeonController.Suspend)
		}
	}
}
