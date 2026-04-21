package player

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

type Player struct {
	validate *validator.Validate
}

func New() *Player {
	return &Player{
		validate: validator.New(),
	}
}

// Get services to get list of Player on db
func (p *Player) Get(queryParams models.QueryParams) ([]models.Player, error) {
	var (
		err     error
		players []models.Player
		player  models.Player
		cursor  *mongo.Cursor
	)

	srv := server.GetServer()
	collection := srv.Database.Collection(player.Collection())

	filter := mongodb.SelectConstructeur(queryParams)
	cursor, err = collection.Find(context.TODO(), filter)
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		// A new result variable should be declared for each document.
		var player models.Player
		err = cursor.Decode(&player)
		if err != nil {
			log.Error().Err(err).Msg("")
			return nil, err
		}
		players = append(players, player)
	}

	err = cursor.Err()
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	return players, err
}

// Create new player on db
func (p *Player) Create(in *models.Player) (*models.Player, error) {
	var player models.Player

	srv := server.GetServer()
	collection := srv.Database.Collection(player.Collection())

	// Check input fields
	err := p.validate.Struct(in)
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	err = functions.ConvertInputStructToDataStruct(in, &player)
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	player.CustomID = functions.NewUUID()
	player.CreatedAt = time.Now()

	_, err = collection.InsertOne(context.TODO(), player)
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	return &player, nil
}

// GetByID controller to get one Player by ID
func (p *Player) GetByID(id string) (models.Player, error) {
	var player models.Player
	var queryParams models.QueryParams

	srv := server.GetServer()
	collection := srv.Database.Collection(player.Collection())

	queryParams.FilterClause = append(queryParams.FilterClause, "customID,"+id)
	filter := mongodb.SelectConstructeur(queryParams)

	err := collection.FindOne(context.TODO(), filter).Decode(&player)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return player, errors.New("player not found")
		}
		log.Error().Err(err).Msg("")
		return player, err
	}

	return player, nil
}

// Update controller to update a Player
func (p *Player) Update(id string, in *models.UpdatePlayerInput) error {
	var (
		result      *mongo.UpdateResult
		err         error
		queryParams models.QueryParams
	)

	srv := server.GetServer()

	// Check input fields
	err = p.validate.Struct(in)
	if err != nil {
		log.Error().Err(err).Msg("")
		return err
	}

	updateFields := bson.M{
		"updatedAt": time.Now(),
	}

	if in.DisplayName != nil {
		updateFields["displayName"] = *in.DisplayName
	}

	if in.Gold != nil {
		updateFields["gold"] = *in.Gold
	}

	if len(updateFields) == 1 {
		return errors.New("no fields to update")
	}

	queryParams.FilterClause = append(queryParams.FilterClause, "customID,"+id)
	filter := mongodb.SelectConstructeur(queryParams)

	collection := srv.Database.Collection((&models.Player{}).Collection())

	update := bson.M{
		"$set": updateFields,
	}

	result, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Error().Err(err).Msg("")
		return err
	}

	if result.MatchedCount == 0 {
		err = errors.New("player to be modified was not found")
	}

	if err == nil && result.ModifiedCount == 0 {
		err = errors.New("player could not be updated")
	}

	if err != nil {
		log.Error().Err(err).Msg("")
	}

	return err
}

// Suspend or Delete controller to suspend a Player
func (s *Player) Suspend(id string) error {
	var (
		err         error
		queryParams models.QueryParams
		player      models.Player
	)

	srv := server.GetServer()
	collection := srv.Database.Collection(player.Collection())

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

// GetByIds controller to get list of Player by Ids
func (s *Player) GetByIds(ids []string) ([]models.Player, error) {
	var Players []models.Player
	for _, id := range ids {
		Player, err := s.GetByID(id)
		if err != nil {
			return nil, err
		}
		Players = append(Players, Player)
	}
	return Players, nil
}
