package dungeon

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
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Dungeon struct {
	validate *validator.Validate
}

func New() *Dungeon {
	return &Dungeon{
		validate: validator.New(),
	}
}

// Get services to get list of Dungeon on db
func (d *Dungeon) Get(queryParams models.QueryParams) ([]models.Dungeon, error) {
	var (
		err      error
		dungeons []models.Dungeon
		dungeon  models.Dungeon
		cursor   *mongo.Cursor
	)

	srv := server.GetServer()
	collection := srv.Database.Collection(dungeon.Collection())

	filter := mongodb.SelectConstructeur(queryParams)
	cursor, err = collection.Find(context.TODO(), filter)
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		// A new result variable should be declared for each document.
		var dungeon models.Dungeon
		err = cursor.Decode(&dungeon)
		if err != nil {
			log.Error().Err(err).Msg("")
			return nil, err
		}
		dungeons = append(dungeons, dungeon)
	}

	err = cursor.Err()
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	return dungeons, err
}

// Create new dungeon on db
func (d *Dungeon) Create(in *models.Dungeon) (*models.Dungeon, error) {
	var dungeon models.Dungeon

	srv := server.GetServer()
	collection := srv.Database.Collection(dungeon.Collection())

	// Check input fields
	err := d.validate.Struct(in)
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	err = functions.ConvertInputStructToDataStruct(in, &dungeon)
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	dungeon.CustomID = functions.NewUUID()
	dungeon.CreatedAt = time.Now()

	_, err = collection.InsertOne(context.TODO(), dungeon)
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	return &dungeon, nil
}

// GetByID controller to get one Dungeon by ID
func (d *Dungeon) GetByID(id string) (models.Dungeon, error) {
	var dungeon models.Dungeon
	var queryParams models.QueryParams

	srv := server.GetServer()
	collection := srv.Database.Collection(dungeon.Collection())

	queryParams.FilterClause = append(queryParams.FilterClause, "customID,"+id)
	filter := mongodb.SelectConstructeur(queryParams)

	err := collection.FindOne(context.TODO(), filter).Decode(&dungeon)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return dungeon, errors.New("dungeon not found")
		}
		log.Error().Err(err).Msg("")
		return dungeon, err
	}

	return dungeon, nil
}

// Update controller to update a Dungeon
func (d *Dungeon) Update(id string, in *models.UpdateDungeonInput) error {
	var (
		result      *mongo.UpdateResult
		err         error
		queryParams models.QueryParams
	)

	srv := server.GetServer()

	// Check input fields
	err = d.validate.Struct(in)
	if err != nil {
		log.Error().Err(err).Msg("")
		return err
	}

	updateFields := bson.M{}

	if in.Title != nil {
		updateFields["title"] = *in.Title
	}

	if in.Description != nil {
		updateFields["description"] = *in.Description
	}

	if in.CreatedBy != nil {
		updateFields["createdby"] = *in.CreatedBy
	}

	if in.Area != nil {
		updateFields["area"] = *in.Area
	}

	if in.Status != nil {
		updateFields["status"] = *in.Status
	}

	updateFields["updatedAt"] = time.Now()

	if len(updateFields) == 1 {
		return errors.New("no fields to update")
	}

	queryParams.FilterClause = append(queryParams.FilterClause, "customID,"+id)
	filter := mongodb.SelectConstructeur(queryParams)

	collection := srv.Database.Collection((&models.Dungeon{}).Collection())

	update := bson.M{
		"$set": updateFields,
	}

	result, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Error().Err(err).Msg("")
		return err
	}

	if result.MatchedCount == 0 {
		err = errors.New("dungeon to be modified was not found")
	}

	if err == nil && result.ModifiedCount == 0 {
		err = errors.New("dungeon could not be updated")
	}

	if err != nil {
		log.Error().Err(err).Msg("")
	}

	return err
}

// Suspend or Delete controller to suspend a Dungeon
func (s *Dungeon) Suspend(id string) error {
	var (
		err         error
		queryParams models.QueryParams
		Dungeon     models.Dungeon
	)

	srv := server.GetServer()
	collection := srv.Database.Collection(Dungeon.Collection())

	queryParams.FilterClause = append(queryParams.FilterClause, "customID,"+id)
	filter := mongodb.SelectConstructeur(queryParams)
	update := bson.M{
		"$set": bson.M{
			"suspended": true,
			"updatedAt": time.Now(),
		},
	}
	_, err = collection.UpdateOne(context.TODO(), filter, update)

	return err
}

// GetByIds controller to get list of Dungeon by Ids
func (s *Dungeon) GetByIds(ids []string) ([]models.Dungeon, error) {
	var Dungeons []models.Dungeon
	for _, id := range ids {
		Dungeon, err := s.GetByID(id)
		if err != nil {
			return nil, err
		}
		Dungeons = append(Dungeons, Dungeon)
	}
	return Dungeons, nil
}

// Publish
func (d *Dungeon) Publish(id string) error {
	var (
		result      *mongo.UpdateResult
		err         error
		queryParams models.QueryParams
	)

	srv := server.GetServer()
	collection := srv.Database.Collection((&models.Dungeon{}).Collection())

	now := time.Now()

	queryParams.FilterClause = append(queryParams.FilterClause, "customID,"+id)
	filter := mongodb.SelectConstructeur(queryParams)

	update := bson.M{
		"$set": bson.M{
			"status":      models.DungeonStatusPublished,
			"publishedAt": now,
			"suspended":   false,
		},
	}

	result, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Error().Err(err).Msg("")
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("Dungeon to be published was not found")
	}

	if result.ModifiedCount == 0 {
		return errors.New("Dungeon could not be published")
	}

	return nil
}
