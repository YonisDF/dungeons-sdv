package models

import "time"

type ItemID string

type InventoryEntry struct {
	PlayerID  string    `bson:"playerID" json:"playerID"`
	ItemID    string    `bson:"itemID" json:"itemID"`
	Qty       int64     `bson:"qty" json:"qty"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}

type InventoryChangeInput struct {
	ItemID string `json:"itemId"`
	Qty    int64  `json:"qty"`
}

// Collection Mongodb collection
func (ie *InventoryEntry) Collection() string {
	return "inventory"
}

type ItemDef struct {
	ID          string    `bson:"id" json:"id"`
	Type        string    `bson:"type" json:"type"`
	Rarity      string    `bson:"rarity" json:"rarity"`
	Name        string    `bson:"name" json:"name"`
	Description string    `bson:"description" json:"description"`
	StatsJSON   []byte    `bson:"statsJSON" json:"-"`
	Tradable    bool      `bson:"tradable" json:"tradable"`
	BaseValue   int64     `bson:"baseValue" json:"baseValue"`
	CreatedAt   time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time `bson:"updatedAt" json:"updatedAt"`
}

// Collection Mongodb collection
func (i *ItemDef) Collection() string {
	return "item"
}

type InventoryResponse struct {
	PlayerID string             `json:"playerId"`
	Items    []InventoryItemDTO `json:"items"`
}

type InventoryItemDTO struct {
	ItemID string `json:"itemId"`
	Qty    int64  `json:"qty"`
}

type ItemDefResponse struct {
	ID          string         `json:"id"`
	Type        string         `json:"type"`
	Rarity      string         `json:"rarity"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Stats       map[string]any `json:"stats,omitempty"`
	Tradable    bool           `json:"tradable"`
	BaseValue   int64          `json:"baseValue,omitempty"`
}
