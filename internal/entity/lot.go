package entity

import (
	"time"

	"github.com/ElOtro/auction-go/internal/validator"
)

type LotStatus int

const (
	LotPending LotStatus = 1 << iota
	LotPublished
	LotProcessing
	LotFinished
)

type BaseLot struct {
	Title       string     `json:"title" example:"Lot #1"`
	Description string     `json:"description,omitempty" example:"Some Precious Items"`
	StartPrice  *int64     `json:"start_price,omitempty" example:"100000"`
	StepPrice   *int64     `json:"step_price,omitempty" example:"15000"`
	StartAt     *time.Time `json:"start_at,omitempty" example:"2022-09-09T12:45:00+03:00"`
	EndAt       *time.Time `json:"end_at,omitempty" example:"2022-09-09T13:45:00+03:00"`
	Notify      bool       `json:"notify" example:"true"`
}

// Lot type
// @Description Lot
type Lot struct {
	ID          int64      `json:"id"`
	Status      LotStatus  `json:"status"`
	Title       string     `json:"title"`
	Description string     `json:"description,omitempty"`
	StartPrice  int64      `json:"start_price,omitempty"`
	EndPrice    int64      `json:"end_price,omitempty"`
	StepPrice   int64      `json:"step_price,omitempty"`
	CreatorID   *int64     `json:"creator_id,omitempty"`
	WinnerID    *int64     `json:"winner_id,omitempty"`
	StartAt     time.Time  `json:"start_at,omitempty"`
	EndAt       time.Time  `json:"end_at,omitempty"`
	Notify      bool       `json:"notify"`
	DestroyedAt *time.Time `json:"-"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

// LotSearch  type
type LotSearch struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
}

type LotFilters struct {
	Title string
}

func ValidateLot(v *validator.Validator, lot *Lot) {
	v.Check(lot.Title != "", "title", "must be provided")
	v.Check(lot.Description != "", "description", "must be provided")
	v.Check(lot.StartPrice > 0, "start_price", "must be greater than zero")
}
