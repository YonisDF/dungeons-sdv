package models

type BossStepID string

type Difficulty string

const (
	DifficultyEasy   Difficulty = "easy"
	DifficultyMedium Difficulty = "medium"
	DifficultyHard   Difficulty = "hard"
)

type ReorderBossStepsInput struct {
	Steps []struct {
		ID    string `json:"id"`
		Order int    `json:"order"`
	} `json:"steps"`
}

type BossStep struct {
	CustomID        string     `bson:"customID" json:"id"`
	DungeonID       string     `bson:"dungeonID" json:"dungeonID"`
	Name            string     `bson:"name" json:"name"`
	Order           int        `bson:"order" json:"order"`
	Location        string     `bson:"location" json:"location"`
	Latitude        float64    `bson:"latitude" json:"latitude"`
	Longitude       float64    `bson:"longitude" json:"longitude"`
	RadiusMeter     float64    `bson:"radiusMeter" json:"radiusMeter"`
	ZoneDescription string     `bson:"zoneDescription" json:"zoneDescription"`
	Difficulty      Difficulty `bson:"difficulty" json:"difficulty"`
	Rewards         []string   `bson:"rewards" json:"rewards"`
	Suspended       bool       `bson:"suspended" json:"suspended"`
}

type BossStepResponse struct {
	ID              string     `json:"id"`
	DungeonID       string     `json:"dungeonID"`
	Name            string     `json:"name"`
	Order           int        `json:"order"`
	Location        string     `json:"location"`
	Latitude        float64    `json:"latitude"`
	Longitude       float64    `json:"longitude"`
	RadiusMeter     float64    `json:"radiusMeter"`
	ZoneDescription string     `json:"zoneDescription"`
	Difficulty      Difficulty `json:"difficulty"`
	Rewards         []string   `json:"rewards"`
	Suspended       bool       `json:"suspended"`
}

type UpdateBossStepInput struct {
	DungeonID       *string     `json:"dungeonID"`
	Name            *string     `json:"name"`
	Order           *int        `json:"order"`
	Location        *string     `json:"location"`
	Latitude        *float64    `json:"latitude"`
	Longitude       *float64    `json:"longitude"`
	RadiusMeter     *float64    `json:"radiusMeter"`
	ZoneDescription *string     `json:"zoneDescription"`
	Difficulty      *Difficulty `json:"difficulty"`
	Rewards         *[]string   `json:"rewards"`
	Suspended       *bool       `json:"suspended"`
}

// Collection Mongodb collection
func (bs *BossStep) Collection() string {
	return "boss_step"
}
