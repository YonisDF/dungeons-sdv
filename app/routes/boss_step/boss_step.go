package boss_step

import (
	controller "dungeons/app/controllers/boss_step"
	service "dungeons/app/services/boss_step"

	"github.com/gin-gonic/gin"
)

func SetupRouter(g *gin.Engine) {

	servicesBossStep := service.New()
	bossStepController := controller.New(servicesBossStep)

	v1 := g.Group("/v1")
	{
		bossSteps := v1.Group("/bossSteps")
		{
			bossSteps.POST("", bossStepController.Create)
			bossSteps.GET("", bossStepController.Get)
			bossSteps.GET("/:id", bossStepController.GetByID)
			bossSteps.DELETE("/:id", bossStepController.Suspend)
			bossSteps.PATCH("/:id", bossStepController.Update)
			bossSteps.GET("/IDS/:ids", bossStepController.GetByIDs)
		}

		mjDungeons := v1.Group("/mj/dungeons")
		{
			mjDungeons.POST("/:id/steps", bossStepController.CreateForDungeon)
			mjDungeons.PUT("/:id/steps/:stepId", bossStepController.UpdateForDungeon)
			mjDungeons.PUT("/:id/steps/reorder", bossStepController.ReorderForDungeon)
		}
	}
}
