package run

import (
	"context"
	"dungeons/app/functions"
	"dungeons/app/models"
	"dungeons/app/mongodb"
	"dungeons/app/server"
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Run struct {
	validate *validator.Validate
}

func New() *Run {
	return &Run{
		validate: validator.New(),
	}
}

// Get services to get list of Run on db
func (r *Run) Get(queryParams models.QueryParams) ([]models.Run, error) {
	var (
		err    error
		runs   []models.Run
		run    models.Run
		cursor *mongo.Cursor
	)

	srv := server.GetServer()
	collection := srv.Database.Collection(run.Collection())

	filter := mongodb.SelectConstructeur(queryParams)
	cursor, err = collection.Find(context.TODO(), filter)
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		// A new result variable should be declared for each document.
		var run models.Run
		err = cursor.Decode(&run)
		if err != nil {
			log.Error().Err(err).Msg("")
			return nil, err
		}
		runs = append(runs, run)
	}

	err = cursor.Err()
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	return runs, err
}

// Create new run on db
func (r *Run) Create(in *models.Run) (*models.Run, error) {
	var run models.Run

	srv := server.GetServer()
	collection := srv.Database.Collection(run.Collection())

	// Check input fields
	err := r.validate.Struct(in)
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	err = functions.ConvertInputStructToDataStruct(in, &run)
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	run.CustomID = functions.NewUUID()
	run.StartedAt = time.Now()

	_, err = collection.InsertOne(context.TODO(), run)
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	return &run, nil
}

// GetByID controller to get one Run by ID
func (r *Run) GetByID(id string) (models.Run, error) {
	var (
		err         error
		run         models.Run
		queryParams models.QueryParams
	)

	srv := server.GetServer()
	collection := srv.Database.Collection(run.Collection())

	queryParams.FilterClause = append(queryParams.FilterClause, "customID,"+id)
	filter := mongodb.SelectConstructeur(queryParams)
	err = collection.FindOne(context.TODO(), filter).Decode(&run)
	if err == nil {
		if err == mongo.ErrNoDocuments {
			log.Error().Err(err).Msg("")
			return run, err
		}

	}
	return run, err
}

// Update controller to update a Run
func (r *Run) Update(id string, in *models.UpdateRunInput) error {
	var (
		result      *mongo.UpdateResult
		err         error
		queryParams models.QueryParams
	)

	srv := server.GetServer()

	// Check input fields
	if err = r.validate.Struct(in); err != nil {
		log.Error().Err(err).Msg("")
		return err
	}

	updateFields := bson.M{}

	if in.DungeonID != nil {
		updateFields["dungeonID"] = *in.DungeonID
	}
	if in.PlayerID != nil {
		updateFields["playerID"] = *in.PlayerID
	}
	if in.State != nil {
		updateFields["state"] = *in.State
	}
	if in.CurrentStep != nil {
		updateFields["currentStep"] = *in.CurrentStep
	}
	if in.KilledSteps != nil {
		updateFields["killedSteps"] = *in.KilledSteps
	}
	if in.BossStepID != nil {
		updateFields["bossStepID"] = *in.BossStepID
	}
	if in.KilledAt != nil {
		updateFields["killedAt"] = *in.KilledAt
	}
	if in.Proof != nil {
		updateFields["proof"] = *in.Proof
	}
	if in.StartedAt != nil {
		updateFields["startedAt"] = *in.StartedAt
	}
	if in.EndedAt != nil {
		updateFields["endedAt"] = *in.EndedAt
	}

	if len(updateFields) == 0 {
		return errors.New("no fields to update")
	}

	queryParams.FilterClause = append(queryParams.FilterClause, "customID,"+id)
	filter := mongodb.SelectConstructeur(queryParams)

	collection := srv.Database.Collection((&models.Run{}).Collection())

	update := bson.M{
		"$set": updateFields,
	}

	result, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Error().Err(err).Msg("")
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("Run to be modified was not found")
	}

	if result.ModifiedCount == 0 {
		return errors.New("Run could not be updated")
	}

	return nil
}

// GetByIds controller to get list of Run by Ids
func (s *Run) GetByIds(ids []string) ([]models.Run, error) {
	var Runs []models.Run
	for _, id := range ids {
		Run, err := s.GetByID(id)
		if err != nil {
			return nil, err
		}
		Runs = append(Runs, Run)
	}
	return Runs, nil
}
