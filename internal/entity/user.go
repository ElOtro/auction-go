package entity

import "time"

// User type
type User struct {
	ID          int64      `json:"id"`
	Active      bool       `json:"active"`
	Name        string     `json:"name"`
	Email       string     `json:"email"`
	Password    password   `json:"-"`
	DestroyedAt *time.Time `json:"destroyed_at,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

// Create a custom password type
type password struct {
	plaintext *string
	hash      []byte
}
