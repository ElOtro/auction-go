package usecase

import repo "github.com/ElOtro/auction-go/internal/infrastructure/repo/postgres"

// Create a UseCases struct which wraps all repos.
type UseCases struct {
	User UserUseCase
	Lot  LotUseCase
	Bid  BidUseCase
}

// For ease of use, we also add a NewUseCases() method which returns a UseCases struct containing
// the initialized UseCases.
func NewUseCases(repos *repo.Repo) UseCases {
	return UseCases{
		User: *NewUserUseCase(&repos.Users),
		Lot:  *NewLotUseCase(&repos.Lots),
		Bid:  *NewBidUseCase(&repos.Bids),
	}
}
