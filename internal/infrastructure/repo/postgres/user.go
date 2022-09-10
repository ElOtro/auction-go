package repo

import (
	"context"
	"errors"
	"time"

	"github.com/ElOtro/auction-go/internal/entity"
	"github.com/ElOtro/auction-go/pkg/postgres"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
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

// Get the User
func (r *UserRepo) Get(userID int64) (*entity.User, error) {
	query := `
		SELECT id, active, role, name, email, password_hash, created_at, updated_at FROM users
		WHERE id = $1`

	var user entity.User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.Pool.QueryRow(ctx, query, userID).Scan(
		&user.ID,
		&user.Active,
		&user.Role,
		&user.Name,
		&user.Email,
		&user.Password.Hash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, entity.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil

}

// GetByEmail the User details from the database based on the user's email address.
// Because we have a UNIQUE constraint on the email column, this SQL query will only
// return one record (or none at all, in which case we return a ErrRecordNotFound error).
func (r *UserRepo) GetByEmail(email string) (*entity.User, error) {
	query := `
		SELECT id, active, role, name, email, password_hash, created_at, updated_at FROM users
		WHERE email = $1`

	var user entity.User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.Pool.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Active,
		&user.Role,
		&user.Name,
		&user.Email,
		&user.Password.Hash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, entity.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil

}

// Insert the User
func (r *UserRepo) Insert(user *entity.User) error {
	query := `
		INSERT INTO users (name, email, password_hash, active) VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at`

	args := []interface{}{user.Name, user.Email, user.Password.Hash, user.Active}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// If the table already contains a record with this email address, then when we try
	// to perform the insert there will be a violation of the UNIQUE "users_email_key"
	// constraint that we set up in the previous chapter. We check for this error
	// specifically, and return custom ErrDuplicateEmail error instead.
	var e *pgconn.PgError
	err := r.Pool.QueryRow(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		switch {
		case errors.As(err, &e) && e.Code == pgerrcode.UniqueViolation:
			return entity.ErrDuplicateEmail
		default:
			return err
		}
	}

	return nil
}
