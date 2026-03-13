package models

import "time"

type RunID string

type State string

const (
	StateActive     State = "active"
	StateCompleted  State = "completed"
	StateAbandonned State = "abandonned"
)

type Run struct {
	CustomID    string    `bson:"customID" json:"id"`
	DungeonID   string    `bson:"dungeonID" json:"dungeonID"`
	PlayerID    string    `bson:"playerID" json:"playerID"`
	State       State     `json:"state"`
	CurrentStep string    `json:"currentStep"`
	KilledSteps []string  `json:"killed_steps"`
	BossStepID  string    `json:"boss_step_id"`
	KilledAt    time.Time `json:"killed_at"`
	Proof       string    `json:"proof"`
	StartedAt   time.Time `json:"started_at"`
	EndedAt     time.Time `json:"ended_at"`
}

type RunResponse struct {
	ID          string    `json:"id"`
	DungeonID   string    `json:"dungeonID"`
	PlayerID    string    `json:"playerID"`
	State       State     `json:"state"`
	CurrentStep string    `json:"currentStep"`
	KilledSteps []string  `json:"killed_steps"`
	BossStepID  string    `json:"boss_step_id"`
	KilledAt    time.Time `json:"killed_at"`
	Proof       string    `json:"proof"`
	StartedAt   time.Time `json:"started_at"`
	EndedAt     time.Time `json:"ended_at"`
}

// Collection Mongodb collection
func (r *Run) Collection() string {
	return "run"
}
