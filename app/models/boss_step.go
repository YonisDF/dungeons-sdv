package models

type BossStepID string

type Difficulty string

const (
	DifficultyEasy   Difficulty = "easy"
	DifficultyMedium Difficulty = "medium"
	DifficultyHard   Difficulty = "hard"
)

type BossStep struct {
	CustomID        string     `bson:"customID" json:"id"`
	DungeonID       string     `bson:"dungeonID" json:"dungeonID"`
	Name            string     `json:"name"`
	Order           int        `json:"order"`
	Location        string     `json:"location"`
	Latitude        float64    `json:"latitude"`
	Longitude       float64    `json:"longitude"`
	RadiusMeter     float64    `json:"radiusMeter"`
	ZoneDescription string     `json:"zoneDescription"`
	Difficulty      Difficulty `json:"difficulty"`
	Rewards         []string   `json:"rewards"`
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
}

// Collection Mongodb collection
func (bs *BossStep) Collection() string {
	return "boss_step"
}
