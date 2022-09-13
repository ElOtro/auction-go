package repo

import (
	"context"
	"time"

	"github.com/ElOtro/auction-go/internal/entity"
	"github.com/ElOtro/auction-go/pkg/postgres"
)

// BidRepo -.
type BidRepo struct {
	*postgres.Postgres
}

// New -.
func NewBidRepo(pg *postgres.Postgres) *BidRepo {
	return &BidRepo{pg}
}

// GetAll method for fetching all records from the bids table for given lot.
func (r *BidRepo) GetAll(lotID int64) ([]*entity.Bid, error) {
	// Construct the SQL query to retrieve all records.
	query := "SELECT id, amount, bidder_id, created_at, updated_at FROM bids WHERE lot_id = $1"

	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Use QueryContext() to execute the query. This returns a sql.Rows resultset
	// containing the result.
	rows, err := r.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	bids := []*entity.Bid{}

	// Use rows.Next to iterate through the rows in the resultset.
	for rows.Next() {
		// Initialize an empty struct to hold the data for an individual record.
		var bid entity.Bid

		// Scan the values from the row into the struct. Again, note that we're
		// using the pq.Array() adapter on the genres field here.
		err := rows.Scan(
			&bid.ID,
			&bid.Amount,
			&bid.BidderID,
			&bid.CreatedAt,
			&bid.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Add the Bid struct to the slice.
		bids = append(bids, &bid)
	}

	// When the rows.Next() loop has finished, call rows.Err() to retrieve any error
	// that was encountered during the iteration.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return bids, nil
}
