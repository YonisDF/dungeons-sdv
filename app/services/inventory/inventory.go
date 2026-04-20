package inventory

import (
	"context"
	"errors"
	"time"

	"dungeons/app/models"
	"dungeons/app/server"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Inventory struct{}

func New() *Inventory {
	return &Inventory{}
}

// Get inventory to get a list of items of a specific player
func (s *Inventory) GetByPlayerID(playerID string) ([]models.InventoryEntry, error) {
	var entries = make([]models.InventoryEntry, 0)

	srv := server.GetServer()
	collection := srv.Database.Collection((&models.InventoryEntry{}).Collection())

	cursor, err := collection.Find(context.TODO(), bson.M{"playerID": playerID})
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var entry models.InventoryEntry
		if err := cursor.Decode(&entry); err != nil {
			log.Error().Err(err).Msg("")
			return nil, err
		}
		entries = append(entries, entry)
	}

	if err := cursor.Err(); err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	return entries, nil
}

// Add item to a player's inventory
func (s *Inventory) AddItem(playerID, itemID string, qty int64) error {
	if qty <= 0 {
		return errors.New("qty must be greater than 0")
	}

	srv := server.GetServer()
	collection := srv.Database.Collection((&models.InventoryEntry{}).Collection())

	filter := bson.M{
		"playerID": playerID,
		"itemID":   itemID,
	}

	update := bson.M{
		"$inc": bson.M{
			"qty": qty,
		},
		"$set": bson.M{
			"updatedAt": time.Now(),
		},
		"$setOnInsert": bson.M{
			"playerID": playerID,
			"itemID":   itemID,
		},
	}

	opts := options.UpdateOne().SetUpsert(true)

	_, err := collection.UpdateOne(
		context.TODO(),
		filter,
		update,
		opts,
	)
	if err != nil {
		log.Error().Err(err).Msg("")
		return err
	}

	return nil
}

// Remove an item from a player's inventory
func (s *Inventory) RemoveItem(playerID, itemID string, qty int64) error {
	if qty <= 0 {
		return errors.New("qty must be greater than 0")
	}

	srv := server.GetServer()
	collection := srv.Database.Collection((&models.InventoryEntry{}).Collection())

	var entry models.InventoryEntry
	err := collection.FindOne(context.TODO(), bson.M{
		"playerID": playerID,
		"itemID":   itemID,
	}).Decode(&entry)
	if err != nil {
		log.Error().Err(err).Msg("")
		return err
	}

	if entry.Qty < qty {
		return errors.New("not enough quantity in inventory")
	}

	newQty := entry.Qty - qty
	if newQty == 0 {
		_, err = collection.DeleteOne(context.TODO(), bson.M{
			"playerID": playerID,
			"itemID":   itemID,
		})
	} else {
		_, err = collection.UpdateOne(
			context.TODO(),
			bson.M{
				"playerID": playerID,
				"itemID":   itemID,
			},
			bson.M{
				"$set": bson.M{
					"qty":       newQty,
					"updatedAt": time.Now(),
				},
			},
		)
	}

	if err != nil {
		log.Error().Err(err).Msg("")
		return err
	}

	return nil
}
