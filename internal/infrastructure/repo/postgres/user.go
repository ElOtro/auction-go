package repo

import (
	"context"
	"time"

	"github.com/ElOtro/auction-go/internal/entity"
	"github.com/ElOtro/auction-go/pkg/postgres"
)

// UserRepo -.
type UserRepo struct {
	*postgres.Postgres
}

// New -.
func NewUserRepo(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{pg}
}

// GetHistory -.
func (r *UserRepo) GetAll() ([]*entity.User, error) {
	// Construct the SQL query to retrieve all records.
	query := `SELECT id, created_at, name, email, password_hash, is_active, updated_at
			  FROM users 
			  WHERE destroyed_at IS NULL`

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

	users := []*entity.User{}

	// Use rows.Next to iterate through the rows in the resultset.
	for rows.Next() {
		// Initialize an empty struct to hold the data for an individual record.
		var user entity.User

		// Scan the values from the row into the struct. Again, note that we're
		// using the pq.Array() adapter on the genres field here.
		err := rows.Scan(
			&user.ID,
			&user.Active,
			&user.Name,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Add the User struct to the slice.
		users = append(users, &user)
	}

	// When the rows.Next() loop has finished, call rows.Err() to retrieve any error
	// that was encountered during the iteration.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
