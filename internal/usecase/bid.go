package usecase

import (
	"github.com/ElOtro/auction-go/internal/entity"
)

type BidRepository interface {
	GetAll(lotID int64) ([]*entity.Bid, error)
	Insert(bid *entity.Bid) error
}

// BidUseCase -.
type BidUseCase struct {
	repo    BidRepository
	lotRepo LotRepository
}

// New -.
func NewBidUseCase(r BidRepository, lr LotRepository) *BidUseCase {
	return &BidUseCase{
		repo:    r,
		lotRepo: lr,
	}
}

// History - getting translate history from store.
func (uc *BidUseCase) List(lotID int64) ([]*entity.Bid, error) {
	companies, err := uc.repo.GetAll(lotID)
	if err != nil {
		return nil, err
	}

	return companies, nil
}

// Create - creating a lot in store.
func (uc *BidUseCase) Create(bid *entity.Bid) error {
	err := uc.repo.Insert(bid)
	if err != nil {
		return err
	}

	return nil
}
