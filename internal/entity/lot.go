package entity

import "time"

// Lot type
type Lot struct {
	ID          int64      `json:"id"`
	Status      int        `json:"status"`
	Title       string     `json:"title"`
	Description string     `json:"description,omitempty"`
	StartPrice  int64      `json:"start_price,omitempty"`
	EndPrice    int64      `json:"end_price,omitempty"`
	CreatorID   int64      `json:"creator_id,omitempty"`
	WinnerID    int64      `json:"winner_id,omitempty"`
	StartAt     *time.Time `json:"start_at,omitempty"`
	DestroyedAt *time.Time `json:"destroyed_at,omitempty"`
	EndAt       *time.Time `json:"end_at,omitempty"`
	Notify      bool       `json:"notify"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

// LotSearch  type
type LotSearch struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type LotFilters struct {
	Name string
}
