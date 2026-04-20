package inventory

import (
	"dungeons/app/controllers/common"
	"dungeons/app/models"
	inventory "dungeons/app/services/inventory"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Inventory struct {
	InventoryService *inventory.Inventory
}

func New(inventoryService *inventory.Inventory) *Inventory {
	return &Inventory{
		InventoryService: inventoryService,
	}
}

// Get Inventory of a specific Player
func (s *Inventory) GetByPlayerID(ctx *gin.Context) {
	messageTypes := &models.MessageTypes{
		OK:                  "inventory.Search.Found",
		NotFound:            "inventory.Search.NotFound",
		BadRequest:          "inventory.Search.BadRequest",
		InternalServerError: "inventory.Search.Error",
	}

	playerID := ctx.Param("playerID")

	entries, err := s.InventoryService.GetByPlayerID(playerID)
	if err != nil {
		common.SendResponse(
			ctx,
			http.StatusInternalServerError,
			models.KnownError(http.StatusInternalServerError, messageTypes.InternalServerError, err),
		)
		return
	}

	if len(entries) == 0 {
		status := http.StatusNotFound
		common.SendResponse(
			ctx,
			status,
			models.KnownError(status, messageTypes.NotFound, errors.New("Data not found.")),
		)
		return
	}

	items := make([]models.InventoryItemDTO, 0, len(entries))
	for _, entry := range entries {
		items = append(items, models.InventoryItemDTO{
			ItemID: entry.ItemID,
			Qty:    entry.Qty,
		})
	}

	inventoryResponse := &models.InventoryResponse{
		PlayerID: playerID,
		Items:    items,
	}

	meta := models.MetaResponse{
		ObjectName: "Inventory",
		TotalCount: len(items),
		Count:      len(items),
		Offset:     1,
	}

	response := &models.WSResponse{
		Meta: meta,
		Data: inventoryResponse,
	}

	common.SendResponse(ctx, http.StatusOK, response)
}

// Add item to inventory
func (s *Inventory) AddItem(ctx *gin.Context) {
	var in models.InventoryChangeInput

	messageTypes := &models.MessageTypes{
		Created:             "inventory.Add.Updated",
		BadRequest:          "inventory.Add.BadRequest",
		InternalServerError: "inventory.Add.Error",
	}

	if err := ctx.BindJSON(&in); err != nil {
		common.SendResponse(
			ctx,
			http.StatusBadRequest,
			models.KnownError(http.StatusBadRequest, messageTypes.BadRequest, err),
		)
		return
	}

	playerID := ctx.Param("playerID")

	err := s.InventoryService.AddItem(playerID, in.ItemID, in.Qty)
	if err != nil {
		common.SendResponse(
			ctx,
			http.StatusInternalServerError,
			models.KnownError(http.StatusInternalServerError, messageTypes.InternalServerError, err),
		)
		return
	}

	common.SendResponse(
		ctx,
		http.StatusOK,
		models.Success(http.StatusOK, messageTypes.Created, "item added to inventory"),
	)
}

// Remove Item from Inventory
func (s *Inventory) RemoveItem(ctx *gin.Context) {
	var in models.InventoryChangeInput

	messageTypes := &models.MessageTypes{
		OK:                  "inventory.Remove.Updated",
		BadRequest:          "inventory.Remove.BadRequest",
		InternalServerError: "inventory.Remove.Error",
	}

	if err := ctx.BindJSON(&in); err != nil {
		common.SendResponse(
			ctx,
			http.StatusBadRequest,
			models.KnownError(http.StatusBadRequest, messageTypes.BadRequest, err),
		)
		return
	}

	playerID := ctx.Param("playerID")

	err := s.InventoryService.RemoveItem(playerID, in.ItemID, in.Qty)
	if err != nil {
		common.SendResponse(
			ctx,
			http.StatusInternalServerError,
			models.KnownError(http.StatusInternalServerError, messageTypes.InternalServerError, err),
		)
		return
	}

	common.SendResponse(
		ctx,
		http.StatusOK,
		models.Success(http.StatusOK, messageTypes.OK, "item removed from inventory"),
	)
}
