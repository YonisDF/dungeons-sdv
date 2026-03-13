package boss_step

import (
	"context"
	"dungeons/app/functions"
	"dungeons/app/models"
	"dungeons/app/mongodb"
	"dungeons/app/server"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type BossStep struct {
	validate *validator.Validate
}

func New() *BossStep {
	return &BossStep{
		validate: validator.New(),
	}
}

// Get services to get list of BossStep on db
func (b *BossStep) Get(queryParams models.QueryParams) ([]models.BossStep, error) {
	var (
		err       error
		bossSteps []models.BossStep
		bossStep  models.BossStep
		cursor    *mongo.Cursor
	)

	srv := server.GetServer()
	collection := srv.Database.Collection(bossStep.Collection())

	filter := mongodb.SelectConstructeur(queryParams)
	cursor, err = collection.Find(context.TODO(), filter)
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		// A new result variable should be declared for each document.
		var bossStep models.BossStep
		err = cursor.Decode(&bossStep)
		if err != nil {
			log.Error().Err(err).Msg("")
			return nil, err
		}
		bossSteps = append(bossSteps, bossStep)
	}

	err = cursor.Err()
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	return bossSteps, err
}

// Create new bossStep on db
func (b *BossStep) Create(in *models.BossStep) (*models.BossStep, error) {
	var bossStep models.BossStep

	srv := server.GetServer()
	collection := srv.Database.Collection(bossStep.Collection())

	// Check input fields
	err := b.validate.Struct(in)
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	err = functions.ConvertInputStructToDataStruct(in, &bossStep)
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	bossStep.CustomID = functions.NewUUID()

	_, err = collection.InsertOne(context.TODO(), bossStep)
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	return &bossStep, nil
}

// GetByID controller to get one BossStep by ID
func (b *BossStep) GetByID(id string) (models.BossStep, error) {
	var (
		err         error
		bossStep    models.BossStep
		queryParams models.QueryParams
	)

	srv := server.GetServer()
	collection := srv.Database.Collection(bossStep.Collection())

	queryParams.FilterClause = append(queryParams.FilterClause, "customID,"+id)
	filter := mongodb.SelectConstructeur(queryParams)
	err = collection.FindOne(context.TODO(), filter).Decode(&bossStep)
	if err == nil {
		if err == mongo.ErrNoDocuments {
			log.Error().Err(err).Msg("")
			return bossStep, err
		}

	}
	return bossStep, err
}

// Update controller to update a BossStep
func (b *BossStep) Update(id string, in *models.BossStep) error {
	var (
		doc         interface{}
		result      *mongo.UpdateResult
		err         error
		queryParams models.QueryParams
		BossStep    models.BossStep
	)

	srv := server.GetServer()

	// Check input fields
	err = b.validate.Struct(in)
	if err != nil {
		log.Error().Err(err).Msg("")
		return err
	}

	BossStep, err = b.GetByID(id)
	if err != nil {
		log.Error().Err(err).Msg("")
		return err
	}

	err = functions.ConvertInputStructToDataStruct(in, &BossStep)
	if err != nil {
		log.Error().Err(err).Msg("")
		return err
	}

	collection := srv.Database.Collection(BossStep.Collection())

	queryParams.FilterClause = append(queryParams.FilterClause, "customID,"+id)
	filter := mongodb.SelectConstructeur(queryParams)
	if doc, err = mongodb.ToDoc(BossStep); err != nil {
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
		err = errors.New("BossStep to be modified was not found")
	}

	if err == nil && result.ModifiedCount == 0 {
		err = errors.New("BossStep could not be updated")
	}
	if err != nil {
		log.Error().Err(err).Msg("")
	}

	return err
}

// Suspend or Delete controller to suspend a BossStep
func (s *BossStep) Suspend(id string) error {
	var (
		err         error
		queryParams models.QueryParams
		BossStep    models.BossStep
	)

	srv := server.GetServer()
	collection := srv.Database.Collection(BossStep.Collection())

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

// GetByIds controller to get list of BossStep by Ids
func (s *BossStep) GetByIds(ids []string) ([]models.BossStep, error) {
	var BossSteps []models.BossStep
	for _, id := range ids {
		BossStep, err := s.GetByID(id)
		if err != nil {
			return nil, err
		}
		BossSteps = append(BossSteps, BossStep)
	}
	return BossSteps, nil
}
