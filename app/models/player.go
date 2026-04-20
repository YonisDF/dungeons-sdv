package models

import "time"

type PlayerID string

type Player struct {
	CustomID    string    `bson:"customID" json:"id"`
	DisplayName string    `bson:"displayname" json:"display_name"`
	Gold        int64     `bson:"gold" json:"gold"`
	CreatedAt   time.Time `bson:"createdat" json:"created_at"`
	UpdatedAt   time.Time `bson:"updatedat" json:"updated_at"`
	Suspended   bool      `json:"suspended"`
}

type PlayerResponse struct {
	ID          string    `json:"id"`
	DisplayName string    `json:"display_name"`
	Wallet      Wallet    `json:"wallet"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Suspended   bool      `json:"suspended"`
}

type UpdatePlayerInput struct {
	DisplayName *string `json:"display_name"`
	Suspended   bool    `json:"suspended"`
	Gold        *int64  `json:"gold"`
}

type Wallet struct {
	Gold int64 `json:"gold"` // int64 pour éviter les soucis quand ça grossit
}

// Collection Mongodb collection
func (p *Player) Collection() string {
	return "player"
}
