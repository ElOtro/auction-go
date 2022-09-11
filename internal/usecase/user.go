package usecase

import (
	"fmt"

	"github.com/ElOtro/auction-go/internal/entity"
)

type UserRepository interface {
	GetAll() ([]*entity.User, error)
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

// List - getting translate list from store.
func (uc *UserUseCase) List() ([]*entity.User, error) {
	Users, err := uc.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("UserUseCase - History - s.repo.GetAll: %w", err)
	}

	return Users, nil
}
