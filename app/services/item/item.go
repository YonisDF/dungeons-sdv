package item

import (
	"context"
	"errors"

	"dungeons/app/models"
	"dungeons/app/server"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Item struct{}

func New() *Item {
	return &Item{}
}

// Get Item list
func (s *Item) Get(params models.QueryParams) ([]models.ItemDef, error) {
	var items = make([]models.ItemDef, 0)

	srv := server.GetServer()
	collection := srv.Database.Collection((&models.ItemDef{}).Collection())

	filter := bson.M{}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var item models.ItemDef
		if err := cursor.Decode(&item); err != nil {
			log.Error().Err(err).Msg("")
			return nil, err
		}
		items = append(items, item)
	}

	if err := cursor.Err(); err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	return items, nil
}

// Get Item by id
func (s *Item) GetByID(id string) (*models.ItemDef, error) {
	srv := server.GetServer()
	collection := srv.Database.Collection((&models.ItemDef{}).Collection())

	filter := bson.M{"id": id}

	var item models.ItemDef
	err := collection.FindOne(context.TODO(), filter).Decode(&item)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("item not found")
		}
		log.Error().Err(err).Msg("")
		return nil, err
	}

	return &item, nil
}

// Create Item
func (s *Item) Create(in *models.ItemDef) (*models.ItemDef, error) {
	srv := server.GetServer()
	collection := srv.Database.Collection((&models.ItemDef{}).Collection())

	_, err := collection.InsertOne(context.TODO(), in)
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	return in, nil
}
