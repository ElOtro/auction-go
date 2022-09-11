package usecase

import (
	"fmt"

	"github.com/ElOtro/auction-go/internal/entity"
)

type LotRepository interface {
	GetAll() ([]*entity.Lot, error)
}

// LotUseCase -.
type LotUseCase struct {
	repo LotRepository
}

// New -.
func NewLotUseCase(r LotRepository) *LotUseCase {
	return &LotUseCase{
		repo: r,
	}
}

// History - getting translate history from store.
func (uc *LotUseCase) List() ([]*entity.Lot, error) {
	companies, err := uc.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("LotUseCase - History - s.repo.GetAll: %w", err)
	}

	return companies, nil
}
