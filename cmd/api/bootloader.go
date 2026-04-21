package main

import (
	"context"
	"errors"
	"time"

	"dungeons/app/functions"
	"dungeons/app/models"
	"dungeons/app/server"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func LoadDemoData() error {
	srv := server.GetServer()
	ctx := context.TODO()
	now := time.Now()

	metaCollection := srv.Database.Collection("bootloader_meta")

	var existing bson.M
	err := metaCollection.FindOne(ctx, bson.M{"key": "demo-data-v1"}).Decode(&existing)
	if err == nil {
		log.Info().Msg("demo bootloader already loaded, skipping")
		return nil
	}
	if !errors.Is(err, mongo.ErrNoDocuments) {
		return err
	}

	playerCollection := srv.Database.Collection((&models.Player{}).Collection())
	dungeonCollection := srv.Database.Collection((&models.Dungeon{}).Collection())
	bossStepCollection := srv.Database.Collection((&models.BossStep{}).Collection())
	itemCollection := srv.Database.Collection((&models.ItemDef{}).Collection())
	inventoryCollection := srv.Database.Collection((&models.InventoryEntry{}).Collection())
	runCollection := srv.Database.Collection((&models.Run{}).Collection())

	// Players
	player1 := models.Player{
		CustomID:    functions.NewUUID(),
		DisplayName: "Pedro",
		Gold:        1200,
		CreatedAt:   now,
		UpdatedAt:   now,
		Suspended:   false,
	}
	player2 := models.Player{
		CustomID:    functions.NewUUID(),
		DisplayName: "Bob",
		Gold:        800,
		CreatedAt:   now,
		UpdatedAt:   now,
		Suspended:   false,
	}

	// Dungeons
	dungeon1 := models.Dungeon{
		CustomID:    functions.NewUUID(),
		Title:       "Catacombes de Paris",
		Description: "Un ancien réseau souterrain infesté de créatures hostiles.",
		CreatedBy:   player1.CustomID,
		Area:        "Paris Centre",
		Status:      models.DungeonStatusPublished,
		CreatedAt:   now,
		UpdatedAt:   now,
		Suspended:   false,
	}
	dungeon2 := models.Dungeon{
		CustomID:    functions.NewUUID(),
		Title:       "Forteresse Oubliée",
		Description: "Un bastion en ruine gardant un artéfact douteusement sécurisé.",
		CreatedBy:   player2.CustomID,
		Area:        "Lot Nord",
		Status:      models.DungeonStatusDraft,
		CreatedAt:   now,
		UpdatedAt:   now,
		Suspended:   false,
	}

	// Items
	item1 := models.ItemDef{
		ID:          "item-sword-bronze",
		Type:        "weapon",
		Rarity:      "common",
		Name:        "Bronze Sword",
		Description: "Une épée basique pour ceux qui aiment survivre un minimum.",
		StatsJSON:   []byte(`{"attack":12,"speed":3}`),
		Tradable:    true,
		BaseValue:   100,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	item2 := models.ItemDef{
		ID:          "item-potion-heal",
		Type:        "consumable",
		Rarity:      "common",
		Name:        "Healing Potion",
		Description: "Remet un peu de vie, ce qui reste pratique quand on tient à ses organes.",
		StatsJSON:   []byte(`{"heal":50}`),
		Tradable:    true,
		BaseValue:   35,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	item3 := models.ItemDef{
		ID:          "item-orb-void",
		Type:        "artifact",
		Rarity:      "epic",
		Name:        "Void Orb",
		Description: "Un artefact dont personne ne devrait vouloir, donc évidemment tout le monde le veut.",
		StatsJSON:   []byte(`{"power":80,"curse":15}`),
		Tradable:    false,
		BaseValue:   900,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Boss steps
	step1 := models.BossStep{
		CustomID:        functions.NewUUID(),
		DungeonID:       dungeon1.CustomID,
		Name:            "Manticore",
		Order:           1,
		Location:        "Central Plaza",
		Latitude:        48.8529,
		Longitude:       2.2990,
		RadiusMeter:     50,
		ZoneDescription: "Zone centrale très dangereuse",
		Difficulty:      models.DifficultyMedium,
		Rewards:         []string{"item-sword-bronze", "item-potion-heal"},
		Suspended:       false,
	}
	step2 := models.BossStep{
		CustomID:        functions.NewUUID(),
		DungeonID:       dungeon1.CustomID,
		Name:            "Liche du Puits",
		Order:           2,
		Location:        "Ancient Well",
		Latitude:        48.8535,
		Longitude:       2.3001,
		RadiusMeter:     35,
		ZoneDescription: "Le puits pulse d'une magie franchement peu rassurante.",
		Difficulty:      models.DifficultyHard,
		Rewards:         []string{"item-orb-void"},
		Suspended:       false,
	}
	step3 := models.BossStep{
		CustomID:        functions.NewUUID(),
		DungeonID:       dungeon2.CustomID,
		Name:            "Gardien de Fer",
		Order:           1,
		Location:        "Broken Gate",
		Latitude:        44.4491,
		Longitude:       1.4366,
		RadiusMeter:     40,
		ZoneDescription: "L'entrée est gardée par une machine rouillée mais vindicative.",
		Difficulty:      models.DifficultyEasy,
		Rewards:         []string{"item-sword-bronze"},
		Suspended:       false,
	}

	// Inventory
	inventory1 := models.InventoryEntry{
		PlayerID:  player1.CustomID,
		ItemID:    item1.ID,
		Qty:       1,
		UpdatedAt: now,
	}
	inventory2 := models.InventoryEntry{
		PlayerID:  player1.CustomID,
		ItemID:    item2.ID,
		Qty:       5,
		UpdatedAt: now,
	}
	inventory3 := models.InventoryEntry{
		PlayerID:  player2.CustomID,
		ItemID:    item3.ID,
		Qty:       1,
		UpdatedAt: now,
	}

	// Runs
	run1 := models.Run{
		CustomID:    functions.NewUUID(),
		DungeonID:   dungeon1.CustomID,
		PlayerID:    player1.CustomID,
		State:       models.StateActive,
		CurrentStep: step1.CustomID,
		KilledSteps: []string{},
		BossStepID:  step1.CustomID,
		Proof:       "photo://demo-proof-1",
		StartedAt:   now.Add(-45 * time.Minute),
		EndedAt:     time.Time{},
	}
	run2 := models.Run{
		CustomID:    functions.NewUUID(),
		DungeonID:   dungeon1.CustomID,
		PlayerID:    player2.CustomID,
		State:       models.StateCompleted,
		CurrentStep: step2.CustomID,
		KilledSteps: []string{step1.CustomID},
		BossStepID:  step2.CustomID,
		KilledAt:    now.Add(-2 * time.Hour),
		Proof:       "photo://demo-proof-2",
		StartedAt:   now.Add(-3 * time.Hour),
		EndedAt:     now.Add(-90 * time.Minute),
	}

	// Upserts

	// Players
	if err := upsertByCustomID(ctx, playerCollection, player1.CustomID, player1); err != nil {
		return err
	}
	if err := upsertByCustomID(ctx, playerCollection, player2.CustomID, player2); err != nil {
		return err
	}

	// Dungeons
	if err := upsertByCustomID(ctx, dungeonCollection, dungeon1.CustomID, dungeon1); err != nil {
		return err
	}
	if err := upsertByCustomID(ctx, dungeonCollection, dungeon2.CustomID, dungeon2); err != nil {
		return err
	}

	// Boss steps
	if err := upsertByCustomID(ctx, bossStepCollection, step1.CustomID, step1); err != nil {
		return err
	}
	if err := upsertByCustomID(ctx, bossStepCollection, step2.CustomID, step2); err != nil {
		return err
	}
	if err := upsertByCustomID(ctx, bossStepCollection, step3.CustomID, step3); err != nil {
		return err
	}

	// Items
	if err := upsertByID(ctx, itemCollection, item1.ID, item1); err != nil {
		return err
	}
	if err := upsertByID(ctx, itemCollection, item2.ID, item2); err != nil {
		return err
	}
	if err := upsertByID(ctx, itemCollection, item3.ID, item3); err != nil {
		return err
	}

	// Inventory
	if err := upsertInventory(ctx, inventoryCollection, inventory1); err != nil {
		return err
	}
	if err := upsertInventory(ctx, inventoryCollection, inventory2); err != nil {
		return err
	}
	if err := upsertInventory(ctx, inventoryCollection, inventory3); err != nil {
		return err
	}

	// Runs
	if err := upsertByCustomID(ctx, runCollection, run1.CustomID, run1); err != nil {
		return err
	}
	if err := upsertByCustomID(ctx, runCollection, run2.CustomID, run2); err != nil {
		return err
	}

	// Bootloader meta
	_, err = metaCollection.InsertOne(ctx, bson.M{
		"key":      "demo-data-v1",
		"loadedAt": now,
		"version":  1,
		"comment":  "initial demo seed",
	})
	if err != nil {
		return err
	}

	log.Info().Msg("demo bootloader data loaded")
	return nil
}

func upsertByCustomID(ctx context.Context, collection *mongo.Collection, customID string, doc interface{}) error {
	_, err := collection.UpdateOne(
		ctx,
		bson.M{"customID": customID},
		bson.M{"$set": doc},
		options.UpdateOne().SetUpsert(true),
	)
	return err
}

func upsertByID(ctx context.Context, collection *mongo.Collection, id string, doc interface{}) error {
	_, err := collection.UpdateOne(
		ctx,
		bson.M{"id": id},
		bson.M{"$set": doc},
		options.UpdateOne().SetUpsert(true),
	)
	return err
}

func upsertInventory(ctx context.Context, collection *mongo.Collection, entry models.InventoryEntry) error {
	_, err := collection.UpdateOne(
		ctx,
		bson.M{
			"playerID": entry.PlayerID,
			"itemID":   entry.ItemID,
		},
		bson.M{"$set": entry},
		options.UpdateOne().SetUpsert(true),
	)
	return err
}
