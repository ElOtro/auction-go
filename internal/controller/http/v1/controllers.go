package v1

import "github.com/ElOtro/auction-go/internal/usecase"

// Create a Models struct which wraps all models.
type Controllers struct {
	Lot  LotController
	User UserController
}

// For ease of use, we also add a NewControllers() method which returns a Controllers struct
func NewControllers(usecases *usecase.UseCases) Controllers {
	return Controllers{
		Lot:  *NewLotController(&usecases.Lot),
		User: *NewUserController(&usecases.User),
	}
}
