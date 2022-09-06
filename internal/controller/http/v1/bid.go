package v1

import (
	"net/http"

	"github.com/ElOtro/auction-go/internal/entity"
)

type BidUseCase interface {
	List() ([]*entity.Bid, error)
}

type BidController struct {
	uc BidUseCase
}

func NewBidController(uc BidUseCase) *BidController {
	return &BidController{uc: uc}
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
// @Success     200 {object} listBidResponse
// @Router      /bids [get]
func (c *BidController) List(w http.ResponseWriter, r *http.Request) {
	bids, err := c.uc.List()
	if err != nil {
		errorResponse(w, r, http.StatusInternalServerError, "get bids")

		return
	}

	err = writeJSON(w, http.StatusOK, listBidResponse{bids}, nil)
	if err != nil {
		serverErrorResponse(w, r, err)
	}
}
