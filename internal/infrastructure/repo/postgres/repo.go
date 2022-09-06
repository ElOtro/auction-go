package repo

import (
	"errors"

	"github.com/ElOtro/auction-go/pkg/postgres"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

// Create a Repo struct which wraps all repo.
type Repo struct {
	Users UserRepo
	Lots  LotRepo
	Bids  BidRepo
}

// For ease of use, we also add a NewRepo() method which returns a Repo struct
func NewRepo(pg *postgres.Postgres) Repo {
	return Repo{
		Users: UserRepo{pg},
		Lots:  LotRepo{pg},
		Bids:  BidRepo{pg},
	}
}
