package inventory

import (
	itemController "dungeons/app/controllers/item"
	itemService "dungeons/app/services/item"

	inventoryController "dungeons/app/controllers/inventory"
	inventoryService "dungeons/app/services/inventory"

	"github.com/gin-gonic/gin"
)

func SetupRouter(g *gin.Engine) {

	servicesItem := itemService.New()
	itemController := itemController.New(servicesItem)

	servicesInventory := inventoryService.New()
	inventoryController := inventoryController.New(servicesInventory)

	v1 := g.Group("/v1")
	{
		inventory := v1.Group("/inventory")
		{
			inventory.GET("/:playerID", inventoryController.GetByPlayerID)
			inventory.POST("/:playerID/add", inventoryController.AddItem)
			inventory.POST("/:playerID/remove", inventoryController.RemoveItem)
		}

		items := v1.Group("/items")
		{
			items.GET("", itemController.Get)
			items.GET("/:id", itemController.GetByID)
			items.POST("", itemController.Create)
			items.PUT("/:id", itemController.Update)
		}
	}
}
