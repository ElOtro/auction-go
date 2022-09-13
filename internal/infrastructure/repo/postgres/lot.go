package repo

import (
	"context"
	"errors"
	"time"

	"github.com/ElOtro/auction-go/internal/entity"
	"github.com/ElOtro/auction-go/pkg/postgres"
	"github.com/jackc/pgx/v4"
)

// LotRepo -.
type LotRepo struct {
	*postgres.Postgres
}

// NewLotRepo -.
func NewLotRepo(pg *postgres.Postgres) *LotRepo {
	return &LotRepo{pg}
}

// GetAll -.
func (r LotRepo) GetAll() ([]*entity.Lot, error) {
	// Construct the SQL query to retrieve all records.
	query := `SELECT id, status, title, description, start_price, end_price, step_price, creator_id, winner_id, 
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

	defer rows.Close()

	lots := []*entity.Lot{}

	// Use rows.Next to iterate through the rows in the resultset.
	for rows.Next() {
		// Initialize an empty struct to hold the data for an individual.
		var lot entity.Lot

		err := rows.Scan(
			&lot.ID,
			&lot.Status,
			&lot.Title,
			&lot.Description,
			&lot.StartPrice,
			&lot.EndPrice,
			&lot.StepPrice,
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

// Get method for fetching a specific record from the lots table.
func (r LotRepo) Get(id int64) (*entity.Lot, error) {
	if id < 1 {
		return nil, entity.ErrRecordNotFound
	}

	// Define the SQL query for retrieving data.
	query := `SELECT id, status, title, description, start_price, end_price, step_price, creator_id, winner_id, 
	          start_at, end_at, notify, created_at, updated_at
		      FROM lots
			  WHERE id = $1`

	var lot entity.Lot

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	// Execute the query using the QueryRow() method, passing in the provided id value
	err := r.Pool.QueryRow(ctx, query, id).Scan(
		&lot.ID,
		&lot.Status,
		&lot.Title,
		&lot.Description,
		&lot.StartPrice,
		&lot.EndPrice,
		&lot.StepPrice,
		&lot.CreatorID,
		&lot.WinnerID,
		&lot.StartAt,
		&lot.EndAt,
		&lot.Notify,
		&lot.CreatedAt,
		&lot.UpdatedAt,
	)

	// Handle any errors. If there was no matching found, Scan() will return
	// a pgx.ErrNoRows error. We check for this and return our custom ErrRecordNotFound
	// error instead.
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, entity.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &lot, nil
}

// Insert method for inserting a new record in the table.
func (r LotRepo) Insert(lot *entity.Lot) error {
	// Define the SQL query for inserting a new record
	query := `
		INSERT INTO lots (status, title, description, start_price, end_price, step_price, creator_id, start_at, end_at, notify) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at`

	args := []interface{}{
		&lot.Status,
		&lot.Title,
		&lot.Description,
		&lot.StartPrice,
		&lot.EndPrice,
		&lot.StepPrice,
		&lot.CreatorID,
		&lot.StartAt,
		&lot.EndAt,
		&lot.Notify,
	}

	// Use the QueryRow() method to execute the SQL query on our connection pool
	return r.Pool.QueryRow(context.Background(), query, args...).Scan(
		&lot.ID,
		&lot.CreatedAt,
		&lot.UpdatedAt,
	)
}

// Update method for updating a specific record.
func (r LotRepo) Update(lot *entity.Lot) error {
	query := `
		UPDATE lots
		SET status = $1, title = $2, description = $3, start_price = $4, end_price = $5, step_price = $6, 
		creator_id = $7, winner_id = $8, start_at = $9, end_at = $10, notify = $11, destroyed_at = $12, updated_at = NOW() 
		WHERE id = $13
		RETURNING updated_at`

	// Create an args slice containing the values for the placeholder parameters.
	args := []interface{}{
		&lot.Status,
		&lot.Title,
		&lot.Description,
		&lot.StartPrice,
		&lot.EndPrice,
		&lot.StepPrice,
		&lot.CreatorID,
		&lot.WinnerID,
		&lot.StartAt,
		&lot.EndAt,
		&lot.Notify,
		&lot.DestroyedAt,
		&lot.ID,
	}

	// Use the QueryRow() method to execute the query, passing in the args slice as a
	// variadic parameter and scanning the new version value into the movie struct.
	return r.Pool.QueryRow(context.Background(), query, args...).Scan(
		&lot.UpdatedAt,
	)
}

// Delete method for deleting a specific record.
func (r LotRepo) Delete(id int64) error {
	if id < 1 {
		return entity.ErrRecordNotFound
	}

	// Construct the SQL query to delete the record.
	query := `
		DELETE FROM lots WHERE id = $1`

	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the SQL query using the Exec() method, passing in the id variable as
	// the value for the placeholder parameter. The Exec() method returns a sql.Result
	// object.
	result, err := r.Pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	// Call the RowsAffected() method on the sql.Result object to get the number of rows
	// affected by the query.
	rowsAffected := result.RowsAffected()

	// If no rows were affected, we know that the products table didn't contain a record
	// with the provided ID at the moment we tried to delete it. In that case we
	// return an ErrRecordNotFound error.
	if rowsAffected == 0 {
		return entity.ErrRecordNotFound
	}

	return nil
}
