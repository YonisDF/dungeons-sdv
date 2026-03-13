package run

import (
	"dungeons/app/controllers/common"
	"dungeons/app/models"
	run "dungeons/app/services/run"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Run struct {
	RunService *run.Run
}

func New(runService *run.Run) *Run {
	return &Run{
		RunService: runService,
	}
}

// Get controller to get list of run
func (s *Run) Get(ctx *gin.Context) {
	var params models.QueryParams

	params.Parse(ctx)
	messageTypes := &models.MessageTypes{
		OK:                  "run.Search.Found",
		BadRequest:          "run.Search.BadRequest",
		NotFound:            "run.Search.NotFound",
		InternalServerError: "run.Search.Error",
	}

	runs, err := s.RunService.Get(params)
	if err != nil {
		common.SendResponse(ctx, http.StatusInternalServerError, models.KnownError(http.StatusInternalServerError, messageTypes.InternalServerError, err))
		return
	}
	totalCount := len(runs)
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

	sendingRuns := runs[low:high]

	meta := models.MetaResponse{
		ObjectName: "Run",
		TotalCount: totalCount,
		Count:      len(sendingRuns),
		Offset:     low + 1,
	}

	response := &models.WSResponse{
		Meta: meta,
		Data: sendingRuns,
	}

	common.SendResponse(ctx, http.StatusOK, response)
}

// Create controller to create new run
func (s *Run) Create(ctx *gin.Context) {
	var in models.Run
	messageTypes := &models.MessageTypes{
		Created:             "run.Create.Created",
		BadRequest:          "run.Create.BadRequest",
		InternalServerError: "run.Create.Error",
	}

	if err := ctx.BindJSON(&in); err != nil {
		common.SendResponse(ctx, http.StatusBadRequest, models.KnownError(http.StatusBadRequest, messageTypes.BadRequest, err))
		return
	}

	run, err := s.RunService.Create(&in)
	if err != nil {
		common.SendResponse(ctx, http.StatusBadRequest, models.KnownError(http.StatusInternalServerError, messageTypes.InternalServerError, err))
		return
	}

	meta := models.MetaResponse{
		ObjectName: "Run",
		TotalCount: 1,
		Count:      1,
		Offset:     0,
	}
	response := &models.WSResponse{
		Meta: meta,
		Data: run,
	}

	common.SendResponse(ctx, http.StatusCreated, response)
}

// GetByID controller to get one run by id
func (s *run) GetByID(ctx *gin.Context) {
	messageTypes := &models.MessageTypes{
		OK:                  "run.get.founded",
		NotFound:            "run.get.NotFound",
		BadRequest:          "run.get.BadRequest",
		InternalServerError: "run.get.Error",
	}

	id := ctx.Param("id")

	run, err := s.RunService.GetByID(id)
	if err != nil {
		common.SendResponse(ctx, http.StatusInternalServerError, models.KnownError(http.StatusInternalServerError, messageTypes.InternalServerError, err))
		return
	}

	// Créer la réponse avec le client
	meta := models.MetaResponse{
		ObjectName: "Run",
		TotalCount: 1,
		Count:      1,
		Offset:     0,
	}
	response := &models.WSResponse{
		Meta: meta,
		Data: run,
	}

	// Envoyer la réponse
	common.SendResponse(ctx, http.StatusOK, response)
}

// Update controller to update run
func (s *Run) Update(ctx *gin.Context) {
	var in models.Run
	messageTypes := &models.MessageTypes{
		OK:                  "run.Update.Updated",
		BadRequest:          "run.Update.BadRequest",
		InternalServerError: "run.Update.Error",
	}

	if err := ctx.BindJSON(&in); err != nil {
		common.SendResponse(ctx, http.StatusBadRequest, models.KnownError(http.StatusBadRequest, messageTypes.BadRequest, err))
		return
	}

	id := ctx.Param("id")
	err := s.RunService.Update(id, &in)
	if err != nil {
		common.SendResponse(ctx, http.StatusBadRequest, models.KnownError(http.StatusInternalServerError, messageTypes.InternalServerError, err))
		return
	}
	common.SendResponse(ctx, http.StatusOK, models.Success(http.StatusOK, messageTypes.OK, "run updated"))
}

func (s *Run) Suspend(ctx *gin.Context) {
	messageTypes := &models.MessageTypes{
		OK:                  "run.Suspend.Updated",
		InternalServerError: "run.Suspend.Error",
	}
	id := ctx.Param("id")
	err := s.RunService.Suspend(id)
	if err != nil {
		common.SendResponse(ctx, http.StatusBadRequest, models.KnownError(http.StatusInternalServerError, messageTypes.InternalServerError, err))
		return
	}
	common.SendResponse(ctx, http.StatusOK, models.Success(http.StatusOK, messageTypes.OK, "run suspended"))
}

// GetByIDs To get all run customID
func (s *Run) GetByIDs(ctx *gin.Context) {
	var params models.QueryParams
	params.Parse(ctx)
	messageTypes := &models.MessageTypes{
		OK:                  "run.Search.Updated",
		NotFound:            "run.Search.NotFound",
		BadRequest:          "run.Search.BadRequest",
		InternalServerError: "run.Search.Error",
	}

	// Extraction des identifiants de la requête URL
	ids := ctx.Param("ids")
	idList := strings.Split(ids, "&")

	runs, err := s.RunService.GetByIds(idList)
	if err != nil {
		common.SendResponse(ctx, http.StatusInternalServerError, models.KnownError(http.StatusInternalServerError, messageTypes.InternalServerError, err))
		return
	}

	totalCount := len(runs)
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

	sendingRun := runs[low:high]

	meta := models.MetaResponse{
		ObjectName: "Run",
		TotalCount: totalCount,
		Count:      len(sendingRun),
		Offset:     low + 1,
	}

	response := &models.WSResponse{
		Meta: meta,
		Data: sendingRun,
	}

	common.SendResponse(ctx, http.StatusOK, response)
}
