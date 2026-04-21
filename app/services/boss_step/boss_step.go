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
	"go.mongodb.org/mongo-driver/v2/bson"
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
	var bossStep models.BossStep
	var queryParams models.QueryParams

	srv := server.GetServer()
	collection := srv.Database.Collection(bossStep.Collection())

	queryParams.FilterClause = append(queryParams.FilterClause, "customID,"+id)
	filter := mongodb.SelectConstructeur(queryParams)

	err := collection.FindOne(context.TODO(), filter).Decode(&bossStep)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return bossStep, errors.New("not found")
		}
		log.Error().Err(err).Msg("")
		return bossStep, err
	}

	return bossStep, nil
}

// Update controller to update a BossStep
func (b *BossStep) Update(id string, in *models.UpdateBossStepInput) error {
	var (
		result      *mongo.UpdateResult
		err         error
		queryParams models.QueryParams
	)

	srv := server.GetServer()

	// Check input fields
	if err = b.validate.Struct(in); err != nil {
		log.Error().Err(err).Msg("")
		return err
	}

	updateFields := bson.M{}

	if in.Name != nil {
		updateFields["name"] = *in.Name
	}
	if in.DungeonID != nil {
		updateFields["dungeonID"] = *in.DungeonID
	}
	if in.Order != nil {
		updateFields["order"] = *in.Order
	}
	if in.Location != nil {
		updateFields["location"] = *in.Location
	}
	if in.Latitude != nil {
		updateFields["latitude"] = *in.Latitude
	}
	if in.Longitude != nil {
		updateFields["longitude"] = *in.Longitude
	}
	if in.RadiusMeter != nil {
		updateFields["radiusMeter"] = *in.RadiusMeter
	}
	if in.ZoneDescription != nil {
		updateFields["zoneDescription"] = *in.ZoneDescription
	}
	if in.Difficulty != nil {
		updateFields["difficulty"] = *in.Difficulty
	}
	if in.Rewards != nil {
		updateFields["rewards"] = *in.Rewards
	}
	if in.Suspended != nil {
		updateFields["suspended"] = *in.Suspended
	}

	if len(updateFields) == 0 {
		return errors.New("no fields to update")
	}

	queryParams.FilterClause = append(queryParams.FilterClause, "customID,"+id)
	filter := mongodb.SelectConstructeur(queryParams)

	collection := srv.Database.Collection((&models.BossStep{}).Collection())

	update := bson.M{
		"$set": updateFields,
	}

	result, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Error().Err(err).Msg("")
		return err
	}

	if result.MatchedCount == 0 {
		err = errors.New("BossStep to be modified was not found")
	} else if result.ModifiedCount == 0 {
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

// Create BossStep linked to a dungeon through url
func (b *BossStep) CreateForDungeon(dungeonID string, in *models.BossStep) (*models.BossStep, error) {
	srv := server.GetServer()
	var bossStep models.BossStep

	collection := srv.Database.Collection(bossStep.Collection())

	in.DungeonID = dungeonID

	if err := b.validate.Struct(in); err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	if err := functions.ConvertInputStructToDataStruct(in, &bossStep); err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	bossStep.CustomID = functions.NewUUID()

	if _, err := collection.InsertOne(context.TODO(), bossStep); err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	return &bossStep, nil
}

// Update BossStep
func (b *BossStep) UpdateForDungeon(dungeonID, stepID string, in *models.UpdateBossStepInput) error {
	srv := server.GetServer()
	collection := srv.Database.Collection((&models.BossStep{}).Collection())

	if err := b.validate.Struct(in); err != nil {
		log.Error().Err(err).Msg("")
		return err
	}

	in.DungeonID = &dungeonID

	updateFields := bson.M{}

	if in.Name != nil {
		updateFields["name"] = *in.Name
	}
	if in.DungeonID != nil {
		updateFields["dungeonID"] = *in.DungeonID
	}
	if in.Order != nil {
		updateFields["order"] = *in.Order
	}
	if in.Location != nil {
		updateFields["location"] = *in.Location
	}
	if in.Latitude != nil {
		updateFields["latitude"] = *in.Latitude
	}
	if in.Longitude != nil {
		updateFields["longitude"] = *in.Longitude
	}
	if in.RadiusMeter != nil {
		updateFields["radiusMeter"] = *in.RadiusMeter
	}
	if in.ZoneDescription != nil {
		updateFields["zoneDescription"] = *in.ZoneDescription
	}
	if in.Difficulty != nil {
		updateFields["difficulty"] = *in.Difficulty
	}
	if in.Rewards != nil {
		updateFields["rewards"] = *in.Rewards
	}
	if in.Suspended != nil {
		updateFields["suspended"] = *in.Suspended
	}

	if len(updateFields) == 0 {
		return errors.New("no fields to update")
	}

	filter := bson.M{
		"customID":  stepID,
		"dungeonID": dungeonID,
	}

	update := bson.M{
		"$set": updateFields,
	}

	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Error().Err(err).Msg("")
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("BossStep to be modified was not found")
	}
	if result.ModifiedCount == 0 {
		return errors.New("BossStep could not be updated")
	}

	return nil
}

// Reorder BossSteps for a dungeon
func (b *BossStep) Reorder(dungeonID string, in *models.ReorderBossStepsInput) error {
	if len(in.Steps) == 0 {
		return errors.New("no steps to reorder")
	}

	srv := server.GetServer()
	collection := srv.Database.Collection((&models.BossStep{}).Collection())

	for _, s := range in.Steps {
		filter := bson.M{
			"customID":  s.ID,
			"dungeonID": dungeonID,
		}

		update := bson.M{
			"$set": bson.M{
				"order": s.Order,
			},
		}

		result, err := collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			log.Error().Err(err).Msg("")
			return err
		}

		if result.MatchedCount == 0 {
			return errors.New("BossStep " + s.ID + " not found in this dungeon")
		}
	}

	return nil
}
