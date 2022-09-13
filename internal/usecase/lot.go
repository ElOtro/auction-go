package usecase

import (
	"github.com/ElOtro/auction-go/internal/entity"
)

type LotRepository interface {
	GetAll() ([]*entity.Lot, error)
	Get(id int64) (*entity.Lot, error)
	Insert(lot *entity.Lot) error
	Update(lot *entity.Lot) error
	Delete(id int64) error
}

// LotUseCase -.
type LotUseCase struct {
	repo LotRepository
}

// NewLotUseCase -.
func NewLotUseCase(r LotRepository) *LotUseCase {
	return &LotUseCase{
		repo: r,
	}
}

// List - getting all lots from store.
func (uc *LotUseCase) List() ([]*entity.Lot, error) {
	lots, err := uc.repo.GetAll()
	if err != nil {
		return nil, err
	}

	return lots, nil
}

// Show - getting a lot from store.
func (uc *LotUseCase) Show(id int64) (*entity.Lot, error) {
	lot, err := uc.repo.Get(id)
	if err != nil {
		return nil, err
	}

	return lot, nil
}

// Create - creating a lot in store.
func (uc *LotUseCase) Create(lot *entity.Lot) error {
	err := uc.repo.Insert(lot)
	if err != nil {
		return err
	}

	return nil
}

// Update - updating a lot to store.
func (uc *LotUseCase) Update(lot *entity.Lot) error {
	err := uc.repo.Insert(lot)
	if err != nil {
		return err
	}

	return nil
}

// Delete - deleting a lot from store.
func (uc *LotUseCase) Delete(id int64) error {
	err := uc.repo.Delete(id)
	if err != nil {
		return err
	}

	return nil
}
