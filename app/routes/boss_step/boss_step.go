package boss_step

import (
	controller "dungeons/app/controllers/player"
	service "dungeons/app/services/player"

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
			bossSteps.POST("/:id", bossStepController.Update)
			bossSteps.GET("/IDS/:ids", bossStepController.GetByIDs)
		}
	}
}
