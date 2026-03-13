package main

import (
	"dungeons/app/models"
	playerService "dungeons/app/services/player"
	"fmt"
	"time"
)

func CreatePlayer() {
	p := models.Player{
		CustomID:    "Jean_Val_Jean",
		DisplayName: "Jean Val Jean",
		Gold:        123,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	player, err := playerService.New().Create(&p)
	if err != nil {
		fmt.Println("Il y a une erreur lors de la création d'un player", err)
	} else {
		fmt.Println("Bienvenue au nouveau joueur", player)
	}
}

func getPlayer() {

}

func updatePlayer() {

}

/*d := Dungeon{
CustomID:    "d-123",
Title:       "Dark Cave",
Description: "Desc...",
CreatedBy:   "gm-42",
Area:        "Cave",
Bosses:      []int64{1, 2, 3},
Status:      DungeonStatusDraft,
CreatedAt:   time.Now(),
UpdatedAt:   time.Now(),
}*/
