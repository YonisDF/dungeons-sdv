package item

import (
	"dungeons/app/controllers/common"
	"dungeons/app/models"
	itemService "dungeons/app/services/item"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Item struct {
	ItemService *itemService.Item
}

func New(itemService *itemService.Item) *Item {
	return &Item{
		ItemService: itemService,
	}
}

// Get All Items
func (s *Item) Get(ctx *gin.Context) {
	var params models.QueryParams
	params.Parse(ctx)

	messageTypes := &models.MessageTypes{
		OK:                  "item.Search.Found",
		NotFound:            "item.Search.NotFound",
		BadRequest:          "item.Search.BadRequest",
		InternalServerError: "item.Search.Error",
	}

	items, err := s.ItemService.Get(params)
	if err != nil {
		common.SendResponse(
			ctx,
			http.StatusInternalServerError,
			models.KnownError(http.StatusInternalServerError, messageTypes.InternalServerError, err),
		)
		return
	}

	totalCount := len(items)
	if totalCount == 0 {
		status := http.StatusNotFound
		common.SendResponse(
			ctx,
			status,
			models.KnownError(status, messageTypes.NotFound, errors.New("Data not found.")),
		)
		return
	}

	low := params.Offset - 1
	if low == -1 {
		low = 0
	}

	maxCount := params.Count
	if maxCount == 0 {
		maxCount = 100
	}

	high := maxCount + low
	if high > totalCount {
		high = totalCount
	}

	if low > high {
		status := http.StatusBadRequest
		common.SendResponse(
			ctx,
			status,
			models.KnownError(status, messageTypes.BadRequest, errors.New("Offset cannot be higher than count.")),
		)
		return
	}

	sendingItems := items[low:high]

	meta := models.MetaResponse{
		ObjectName: "Item",
		TotalCount: totalCount,
		Count:      len(sendingItems),
		Offset:     low + 1,
	}

	response := &models.WSResponse{
		Meta: meta,
		Data: sendingItems,
	}

	common.SendResponse(ctx, http.StatusOK, response)
}

// Get Items by ID
func (s *Item) GetByID(ctx *gin.Context) {
	messageTypes := &models.MessageTypes{
		OK:                  "item.Get.Found",
		NotFound:            "item.Get.NotFound",
		BadRequest:          "item.Get.BadRequest",
		InternalServerError: "item.Get.Error",
	}

	id := ctx.Param("id")

	item, err := s.ItemService.GetByID(id)
	if err != nil {
		common.SendResponse(
			ctx,
			http.StatusInternalServerError,
			models.KnownError(http.StatusInternalServerError, messageTypes.InternalServerError, err),
		)
		return
	}

	meta := models.MetaResponse{
		ObjectName: "Item",
		TotalCount: 1,
		Count:      1,
		Offset:     0,
	}

	response := &models.WSResponse{
		Meta: meta,
		Data: item,
	}

	common.SendResponse(ctx, http.StatusOK, response)
}

// Create Item
func (s *Item) Create(ctx *gin.Context) {
	messageTypes := &models.MessageTypes{
		Created:             "item.Create.Created",
		BadRequest:          "item.Create.BadRequest",
		InternalServerError: "item.Create.Error",
	}

	var in models.ItemDef
	if err := ctx.ShouldBindJSON(&in); err != nil {
		common.SendResponse(
			ctx,
			http.StatusBadRequest,
			models.KnownError(http.StatusBadRequest, messageTypes.BadRequest, err),
		)
		return
	}

	now := time.Now()
	in.CreatedAt = now
	in.UpdatedAt = now

	item, err := s.ItemService.Create(&in)
	if err != nil {
		common.SendResponse(
			ctx,
			http.StatusInternalServerError,
			models.KnownError(http.StatusInternalServerError, messageTypes.InternalServerError, err),
		)
		return
	}

	meta := models.MetaResponse{
		ObjectName: "Item",
		TotalCount: 1,
		Count:      1,
		Offset:     0,
	}

	response := &models.WSResponse{
		Meta: meta,
		Data: item,
	}

	common.SendResponse(ctx, http.StatusCreated, response)
}

// Update Item
func (s *Item) Update(ctx *gin.Context) {
	messageTypes := &models.MessageTypes{
		OK:                  "item.Update.Ok",
		BadRequest:          "item.Update.BadRequest",
		NotFound:            "item.Update.NotFound",
		InternalServerError: "item.Update.Error",
	}

	id := ctx.Param("id")

	var in models.ItemDef
	if err := ctx.ShouldBindJSON(&in); err != nil {
		common.SendResponse(
			ctx,
			http.StatusBadRequest,
			models.KnownError(http.StatusBadRequest, messageTypes.BadRequest, err),
		)
		return
	}

	err := s.ItemService.Update(id, &in)
	if err != nil {
		if err.Error() == "item not found" {
			common.SendResponse(
				ctx,
				http.StatusNotFound,
				models.KnownError(http.StatusNotFound, messageTypes.NotFound, err),
			)
			return
		}

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
		models.Success(http.StatusOK, messageTypes.OK, "item updated"),
	)
}
