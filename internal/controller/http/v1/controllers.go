package v1

import "github.com/ElOtro/auction-go/internal/usecase"

// Create a Controllers struct which wraps all controllers.
type Controllers struct {
	Lot     LotController
	Bid     BidController
	User    UserController
	Session SessionController
}

// For ease of use, we also add a NewControllers() method which returns a Controllers struct
func NewControllers(usecases *usecase.UseCases, jwtSecret string) Controllers {
	return Controllers{
		Lot:     *NewLotController(&usecases.Lot),
		Bid:     *NewBidController(&usecases.Bid, &usecases.Lot),
		User:    *NewUserController(&usecases.User),
		Session: *NewSessionController(&usecases.User, jwtSecret),
	}
}
