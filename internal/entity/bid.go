package entity

import (
	"time"

	"github.com/ElOtro/auction-go/internal/validator"
)

// Bid type
type Bid struct {
	ID        int64      `json:"id"`
	Amount    int64      `json:"amount"`
	Price     int64      `json:"price"`
	LotID     int64      `json:"lot_id,omitempty"`
	BidderID  *int64     `json:"bidder_id,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

func ValidateBid(v *validator.Validator, bid *Bid) {
	v.Check(bid.Amount > 0, "start_price", "must be greater than zero")
	v.Check(*bid.BidderID != 0, "bidder_id", "must be provided")
}
