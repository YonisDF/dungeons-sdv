package boss_step

import (
	"dungeons/app/controllers/common"
	"dungeons/app/models"
	bossStep "dungeons/app/services/boss_step"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type BossStep struct {
	BossStepService *bossStep.BossStep
}

func New(bossStepService *bossStep.BossStep) *BossStep {
	return &BossStep{
		BossStepService: bossStepService,
	}
}

// Get controller to get list of dungeon
func (s *BossStep) Get(ctx *gin.Context) {
	var params models.QueryParams

	params.Parse(ctx)
	messageTypes := &models.MessageTypes{
		OK:                  "bossStep.Search.Found",
		BadRequest:          "bossStep.Search.BadRequest",
		NotFound:            "bossStep.Search.NotFound",
		InternalServerError: "bossStep.Search.Error",
	}

	bossSteps, err := s.BossStepService.Get(params)
	if err != nil {
		common.SendResponse(ctx, http.StatusInternalServerError, models.KnownError(http.StatusInternalServerError, messageTypes.InternalServerError, err))
		return
	}
	totalCount := len(bossSteps)
	if totalCount == 0 {
		status := http.StatusNotFound
		common.SendResponse(ctx, status, models.KnownError(status, messageTypes.NotFound, errors.New(" Data not found. ")))
		return
	}

	low := params.Offset - 1
	if low == -1 {
		low = 0
	}

	// Available CountMax calculation
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
		common.SendResponse(ctx, status, models.KnownError(status, messageTypes.NotFound, errors.New(" Offset cannot be higher than count. ")))
		return
	}

	sendingBossSteps := bossSteps[low:high]

	meta := models.MetaResponse{
		ObjectName: "BossStep",
		TotalCount: totalCount,
		Count:      len(sendingBossSteps),
		Offset:     low + 1,
	}

	response := &models.WSResponse{
		Meta: meta,
		Data: sendingBossSteps,
	}

	common.SendResponse(ctx, http.StatusOK, response)
}

// Create controller to create new boss step
func (s *BossStep) Create(ctx *gin.Context) {
	var in models.BossStep
	messageTypes := &models.MessageTypes{
		Created:             "bossStep.Create.Created",
		BadRequest:          "bossStep.Create.BadRequest",
		InternalServerError: "bossStep.Create.Error",
	}

	if err := ctx.BindJSON(&in); err != nil {
		common.SendResponse(ctx, http.StatusBadRequest, models.KnownError(http.StatusBadRequest, messageTypes.BadRequest, err))
		return
	}

	bossStep, err := s.BossStepService.Create(&in)
	if err != nil {
		common.SendResponse(ctx, http.StatusBadRequest, models.KnownError(http.StatusInternalServerError, messageTypes.InternalServerError, err))
		return
	}

	meta := models.MetaResponse{
		ObjectName: "BossStep",
		TotalCount: 1,
		Count:      1,
		Offset:     0,
	}
	response := &models.WSResponse{
		Meta: meta,
		Data: bossStep,
	}

	common.SendResponse(ctx, http.StatusCreated, response)
}

// GetByID controller to get one boss step by id
func (s *BossStep) GetByID(ctx *gin.Context) {
	messageTypes := &models.MessageTypes{
		OK:                  "bossStep.get.founded",
		NotFound:            "bossStep.get.NotFound",
		BadRequest:          "bossStep.get.BadRequest",
		InternalServerError: "bossStep.get.Error",
	}

	id := ctx.Param("id")

	bossStep, err := s.BossStepService.GetByID(id)
	if err != nil {
		common.SendResponse(ctx, http.StatusInternalServerError, models.KnownError(http.StatusInternalServerError, messageTypes.InternalServerError, err))
		return
	}

	// Créer la réponse avec le client
	meta := models.MetaResponse{
		ObjectName: "BossStep",
		TotalCount: 1,
		Count:      1,
		Offset:     0,
	}
	response := &models.WSResponse{
		Meta: meta,
		Data: bossStep,
	}

	// Envoyer la réponse
	common.SendResponse(ctx, http.StatusOK, response)
}

// Update controller to update boss step
func (s *BossStep) Update(ctx *gin.Context) {
	var in models.BossStep
	messageTypes := &models.MessageTypes{
		OK:                  "bossStep.Update.Updated",
		BadRequest:          "bossStep.Update.BadRequest",
		InternalServerError: "bossStep.Update.Error",
	}

	if err := ctx.BindJSON(&in); err != nil {
		common.SendResponse(ctx, http.StatusBadRequest, models.KnownError(http.StatusBadRequest, messageTypes.BadRequest, err))
		return
	}

	id := ctx.Param("id")
	err := s.BossStepService.Update(id, &in)
	if err != nil {
		common.SendResponse(ctx, http.StatusBadRequest, models.KnownError(http.StatusInternalServerError, messageTypes.InternalServerError, err))
		return
	}
	common.SendResponse(ctx, http.StatusOK, models.Success(http.StatusOK, messageTypes.OK, "boss step updated"))
}

func (s *BossStep) Suspend(ctx *gin.Context) {
	messageTypes := &models.MessageTypes{
		OK:                  "bossStep.Suspend.Updated",
		InternalServerError: "bossStep.Suspend.Error",
	}
	id := ctx.Param("id")
	err := s.BossStepService.Suspend(id)
	if err != nil {
		common.SendResponse(ctx, http.StatusBadRequest, models.KnownError(http.StatusInternalServerError, messageTypes.InternalServerError, err))
		return
	}
	common.SendResponse(ctx, http.StatusOK, models.Success(http.StatusOK, messageTypes.OK, "bossStep suspended"))
}

// GetByIDs To get all boss step customID
func (s *BossStep) GetByIDs(ctx *gin.Context) {
	var params models.QueryParams
	params.Parse(ctx)
	messageTypes := &models.MessageTypes{
		OK:                  "bossStep.Search.Updated",
		NotFound:            "bossStep.Search.NotFound",
		BadRequest:          "bossStep.Search.BadRequest",
		InternalServerError: "bossStep.Search.Error",
	}

	// Extraction des identifiants de la requête URL
	ids := ctx.Param("ids")
	idList := strings.Split(ids, "&")

	bossSteps, err := s.BossStepService.GetByIds(idList)
	if err != nil {
		common.SendResponse(ctx, http.StatusInternalServerError, models.KnownError(http.StatusInternalServerError, messageTypes.InternalServerError, err))
		return
	}

	totalCount := len(bossSteps)
	if totalCount == 0 {
		status := http.StatusNotFound
		common.SendResponse(ctx, status, models.KnownError(status, messageTypes.NotFound, errors.New("Data not found.")))
		return
	}

	low := params.Offset - 1
	if low == -1 {
		low = 0
	}

	// Calcul du CountMax disponible
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
		common.SendResponse(ctx, status, models.KnownError(status, messageTypes.NotFound, errors.New("Offset cannot be higher than count.")))
		return
	}

	sendingBossSteps := bossSteps[low:high]

	meta := models.MetaResponse{
		ObjectName: "BossStep",
		TotalCount: totalCount,
		Count:      len(sendingBossSteps),
		Offset:     low + 1,
	}

	response := &models.WSResponse{
		Meta: meta,
		Data: sendingBossSteps,
	}

	common.SendResponse(ctx, http.StatusOK, response)
}
