package models

import "time"

type DungeonID string

type DungeonStatus string

const (
	DungeonStatusDraft     DungeonStatus = "draft"
	DungeonStatusPublished DungeonStatus = "published"
	DungeonStatusArchived  DungeonStatus = "archived"
)

type Dungeon struct {
	CustomID    string        `bson:"customID" json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	CreatedBy   string        `json:"createdBy"`
	Area        string        `json:"area"`
	Status      DungeonStatus `json:"status"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	Suspended   bool          `json:"suspended"`
}

type DungeonResponse struct {
	ID          string        `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	CreatedBy   string        `json:"createdBy"`
	Area        string        `json:"area"`
	Status      DungeonStatus `json:"status"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	Suspended   bool          `json:"suspended"`
}

type UpdateDungeonInput struct {
	Title       *string        `json:"title"`
	Description *string        `json:"description"`
	CreatedBy   *string        `json:"createdBy"`
	Area        *string        `json:"area"`
	Status      *DungeonStatus `json:"status"`
	Suspended   bool           `json:"suspended"`
}

// Collection Mongodb collection
func (d *Dungeon) Collection() string {
	return "dungeon"
}
