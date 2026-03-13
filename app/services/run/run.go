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
func (r *Run) Update(id string, in *models.Run) error {
	var (
		doc         interface{}
		result      *mongo.UpdateResult
		err         error
		queryParams models.QueryParams
		Run         models.Run
	)

	srv := server.GetServer()

	// Check input fields
	err = r.validate.Struct(in)
	if err != nil {
		log.Error().Err(err).Msg("")
		return err
	}

	Run, err = r.GetByID(id)
	if err != nil {
		log.Error().Err(err).Msg("")
		return err
	}

	err = functions.ConvertInputStructToDataStruct(in, &Run)
	if err != nil {
		log.Error().Err(err).Msg("")
		return err
	}

	Run.EndedAt = time.Now()
	collection := srv.Database.Collection(Run.Collection())

	queryParams.FilterClause = append(queryParams.FilterClause, "customID,"+id)
	filter := mongodb.SelectConstructeur(queryParams)
	if doc, err = mongodb.ToDoc(Run); err != nil {
		log.Error().Err(err).Msg("")
		return err
	}

	update := bson.M{"$set": doc}
	result, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Error().Err(err).Msg("")
		return err
	}

	if result.MatchedCount == 0 {
		err = errors.New("Run to be modified was not found")
	}

	if err == nil && result.ModifiedCount == 0 {
		err = errors.New("Run could not be updated")
	}
	if err != nil {
		log.Error().Err(err).Msg("")
	}

	return err
}

// Suspend or Delete controller to suspend a Run
func (s *Run) Suspend(id string) error {
	var (
		err         error
		queryParams models.QueryParams
		Run         models.Run
	)

	srv := server.GetServer()
	collection := srv.Database.Collection(Run.Collection())

	queryParams.FilterClause = append(queryParams.FilterClause, "customID,"+id)
	filter := mongodb.SelectConstructeur(queryParams)
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "suspended", Value: true},
		}},
	}
	_, err = collection.UpdateOne(context.TODO(), filter, update)

	return err
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
