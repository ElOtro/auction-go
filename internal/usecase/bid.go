package usecase

import (
	"fmt"

	"github.com/ElOtro/auction-go/internal/entity"
)

type BidRepository interface {
	GetAll() ([]*entity.Bid, error)
}

// BidUseCase -.
type BidUseCase struct {
	repo BidRepository
}

// New -.
func NewBidUseCase(r BidRepository) *BidUseCase {
	return &BidUseCase{
		repo: r,
	}
}

// History - getting translate history from store.
func (uc *BidUseCase) List() ([]*entity.Bid, error) {
	companies, err := uc.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("BidUseCase - List - s.repo.List: %w", err)
	}

	return companies, nil
}
