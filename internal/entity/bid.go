package entity

import "time"

// Bid type
type Bid struct {
	ID        int64      `json:"id"`
	Price     int64      `json:"price,omitempty"`
	LotID     int64      `json:"lot_id,omitempty"`
	BidderID  int64      `json:"bidder_id,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}
