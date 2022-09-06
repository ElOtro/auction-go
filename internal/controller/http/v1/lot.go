package v1

import (
	"net/http"

	"github.com/ElOtro/auction-go/internal/entity"
)

type LotUseCase interface {
	List() ([]*entity.Lot, error)
}

type LotController struct {
	uc LotUseCase
}

func NewLotController(uc LotUseCase) *LotController {
	return &LotController{uc: uc}
}

type listLotResponse struct {
	Lot []*entity.Lot `json:"lots"`
}

// @Summary     Show lot list
// @Description Show all lot list
// @ID          lotList
// @Tags        lots
// @Accept      json
// @Produce     json
// @Success     200 {object} listLotResponse
// @Router      /lots [get]
func (c *LotController) List(w http.ResponseWriter, r *http.Request) {
	lots, err := c.uc.List()
	if err != nil {
		errorResponse(w, r, http.StatusInternalServerError, "get lots")

		return
	}

	err = writeJSON(w, http.StatusOK, listLotResponse{lots}, nil)
	if err != nil {
		serverErrorResponse(w, r, err)
	}
}
