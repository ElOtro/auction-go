package v1

import "github.com/ElOtro/auction-go/internal/usecase"

// Create a Controllers struct which wraps all controllers.
type Controllers struct {
	Lot  LotController
	User UserController
}

// For ease of use, we also add a NewControllers() method which returns a Controllers struct
func NewControllers(usecases *usecase.UseCases, jwtSecret string) Controllers {
	return Controllers{
		Lot:  *NewLotController(&usecases.Lot),
		User: *NewUserController(&usecases.User, jwtSecret),
	}
}
