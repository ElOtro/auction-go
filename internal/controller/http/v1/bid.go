package v1

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ElOtro/auction-go/internal/entity"
	"github.com/ElOtro/auction-go/internal/validator"
)

type BidUseCase interface {
	List(lotID int64) ([]*entity.Bid, error)
	Create(bid *entity.Bid) error
}

type BidController struct {
	uc  BidUseCase
	ucl LotUseCase
}

func NewBidController(uc BidUseCase, ucl LotUseCase) *BidController {
	return &BidController{uc: uc, ucl: ucl}
}

type listBidResponse struct {
	Bid []*entity.Bid `json:"bids"`
}

// @Summary     Show bid list
// @Description Show all bid list
// @ID          bidList
// @Tags        bids
// @Accept      json
// @Produce     json
// @Param       id            path     int    true "Lot ID"                   Format(int64)
// @Param       Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success     200           {object} listBidResponse
// @Failure     500
// @Router      /lots/{id}/bids [get]
func (c *BidController) List(w http.ResponseWriter, r *http.Request) {
	lotID, err := readIDParam("ID", r)
	if err != nil {
		notFoundResponse(w, r)
		return
	}

	_, err = c.ucl.Show(lotID)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrRecordNotFound):
			notFoundResponse(w, r)
		default:
			serverErrorResponse(w, r, err)
		}
		return
	}

	bids, err := c.uc.List(lotID)
	if err != nil {
		serverErrorResponse(w, r, err)
		return
	}

	err = writeJSON(w, http.StatusOK, listBidResponse{bids}, nil)
	if err != nil {
		serverErrorResponse(w, r, err)
	}
}

// Get          godoc
// @Summary     Create bid
// @Description create bid
// @ID          create-bid
// @Tags        bids
// @Accept      json
// @Produce     json
// @Param       id            path   int    true "Lot ID"                   Format(int64)
// @Param       Authorization header   string true "Insert your access token" default(Bearer <Add access token here>)
// @Success     201
// @Failure     400
// @Failure     422
// @Failure     500
// @Router      /lots/{id}/bids [post]
func (c *BidController) Create(w http.ResponseWriter, r *http.Request) {
	lotID, err := readIDParam("ID", r)
	if err != nil {
		notFoundResponse(w, r)
		return
	}

	lot, err := c.ucl.Show(lotID)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrRecordNotFound):
			notFoundResponse(w, r)
		default:
			serverErrorResponse(w, r, err)
		}
		return
	}

	bid := &entity.Bid{
		Amount:   lot.StepPrice,
		LotID:    lotID,
		BidderID: nil,
	}

	// Validate the record, sending the client a 422 Unprocessable Entity
	// response if any checks fail.
	v := validator.New()

	// Call the validate function and return a response containing the errors if
	// any of the checks fail.
	if entity.ValidateBid(v, bid); !v.Valid() {
		failedValidationResponse(w, r, v.Errors)
		return
	}

	err = c.uc.Create(bid)
	if err != nil {
		serverErrorResponse(w, r, err)
		return
	}

	// When sending a HTTP response, we want to include a Location header to let the
	// client know which URL they can find the newly-created resource at.
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/lots/%d/bids/%d", lotID, bid.ID))

	// Write a JSON response with a 201 Created status code, the lot data in the
	// response body, and the Location header.
	err = writeJSON(w, http.StatusCreated, envelope{"bid": bid}, headers)
	if err != nil {
		serverErrorResponse(w, r, err)
	}

}
