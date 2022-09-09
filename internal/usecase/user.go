package usecase

import (
	"fmt"

	"github.com/ElOtro/auction-go/internal/entity"
)

type UserRepository interface {
	GetAll() ([]*entity.User, error)
	Get(userID int64) (*entity.User, error)
}

// UserUseCase -.
type UserUseCase struct {
	repo UserRepository
}

// New -.
func NewUserUseCase(r UserRepository) *UserUseCase {
	return &UserUseCase{
		repo: r,
	}
}

// List - getting user list from store.
func (uc *UserUseCase) List() ([]*entity.User, error) {
	users, err := uc.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("UserUseCase - List - s.repo.GetAll: %w", err)
	}

	return users, nil
}

// Get - getting user from store.
func (uc *UserUseCase) Get(userID int64) (*entity.User, error) {
	user, err := uc.repo.Get(userID)
	if err != nil {
		return nil, fmt.Errorf("UserUseCase - Get - s.repo.Get: %w", err)
	}

	return user, nil
}
