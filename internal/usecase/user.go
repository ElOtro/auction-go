package usecase

import (
	"github.com/ElOtro/auction-go/internal/entity"
)

type UserRepository interface {
	GetAll() ([]*entity.User, error)
	Get(userID int64) (*entity.User, error)
	GetByEmail(email string) (*entity.User, error)
	Insert(user *entity.User) error
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
		return nil, err
	}

	return users, nil
}

// Get - getting user from store.
func (uc *UserUseCase) Get(userID int64) (*entity.User, error) {
	user, err := uc.repo.Get(userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetByEmail - getting user from store.
func (uc *UserUseCase) GetByEmail(email string) (*entity.User, error) {
	user, err := uc.repo.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Insert - inserting user to store.
func (uc *UserUseCase) Insert(user *entity.User) error {
	err := uc.repo.Insert(user)
	if err != nil {
		return err
	}

	return nil
}
