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

// GetHistory -.
func (r *BidRepo) GetAll() ([]*entity.Bid, error) {
	// Construct the SQL query to retrieve all records.
	query := `SELECT id, amount, lot_id, bidder_id, created_at, updated_at
			  FROM bids`

	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Use QueryContext() to execute the query. This returns a sql.Rows resultset
	// containing the result.
	rows, err := r.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	// Importantly, defer a call to rows.Close() to ensure that the resultset is closed
	// before GetAll() returns.
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
			&bid.LotID,
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
