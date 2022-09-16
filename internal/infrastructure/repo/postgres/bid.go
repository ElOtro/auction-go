package repo

import (
	"context"
	"time"

	"github.com/ElOtro/auction-go/internal/entity"
	"github.com/ElOtro/auction-go/pkg/postgres"
	"github.com/jackc/pgx/v4"
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
	rows, err := r.Pool.Query(ctx, query, lotID)
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

// Insert method for inserting a new record in the table.
func (r *BidRepo) Insert(bid *entity.Bid) error {

	tx, err := r.Pool.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback(context.Background())
		} else {
			tx.Commit(context.Background())
		}
	}()

	// Get sum of bids for certain lot
	query := "SELECT COALESCE(SUM(amount), 0) FROM bids WHERE lot_id = $1"
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	var sum int64
	err = tx.QueryRow(ctx, query, bid.LotID).Scan(&sum)
	defer cancel()

	// Define the SQL query for inserting a new record
	sum += bid.Amount
	query = `
		INSERT INTO bids (amount, price, lot_id, bidder_id) VALUES ($1, $2, $3, $4)
		RETURNING id, price, created_at, updated_at`

	args := []interface{}{
		&bid.Amount,
		sum,
		&bid.LotID,
		&bid.BidderID,
	}

	// Use the QueryRow() method to execute the SQL query on our connection pool
	return tx.QueryRow(context.Background(), query, args...).Scan(
		&bid.ID,
		&bid.Price,
		&bid.CreatedAt,
		&bid.UpdatedAt,
	)
}
