package repo

import (
	"context"
	"time"

	"github.com/ElOtro/auction-go/internal/entity"
	"github.com/ElOtro/auction-go/pkg/postgres"
)

// LotRepo -.
type LotRepo struct {
	*postgres.Postgres
}

// New -.
func NewLotRepo(pg *postgres.Postgres) *LotRepo {
	return &LotRepo{pg}
}

// GetAll method for fetching all records from the lots table.
func (r LotRepo) GetAll() ([]*entity.Lot, error) {
	// Construct the SQL query to retrieve all records.
	query := `SELECT id, status, title, description, start_price, end_price, creator_id, winner_id, 
	          start_at, end_at, notify, created_at, updated_at
		      FROM lots`

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

	lots := []*entity.Lot{}

	// Use rows.Next to iterate through the rows in the resultset.
	for rows.Next() {
		// Initialize an empty struct to hold the data for an individual.
		var lot entity.Lot

		// Scan the values from the row into the struct. Again, note that we're
		// using the pq.Array() adapter on the genres field here.
		err := rows.Scan(
			&lot.ID,
			&lot.Status,
			&lot.Title,
			&lot.Description,
			&lot.StartPrice,
			&lot.EndPrice,
			&lot.CreatorID,
			&lot.WinnerID,
			&lot.StartAt,
			&lot.EndAt,
			&lot.Notify,
			&lot.CreatedAt,
			&lot.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Add the Lot struct to the slice.
		lots = append(lots, &lot)
	}

	// When the rows.Next() loop has finished, call rows.Err() to retrieve any error
	// that was encountered during the iteration.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return lots, nil
}
