package v1

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ElOtro/auction-go/internal/entity"
	"github.com/ElOtro/auction-go/internal/validator"
)

type LotUseCase interface {
	List() ([]*entity.Lot, error)
	Show(id int64) (*entity.Lot, error)
	Create(lot *entity.Lot) error
	Update(lot *entity.Lot) error
	Delete(id int64) error
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

type lotResponse struct {
	Lot *entity.Lot `json:"lot"`
}

type lotRequest struct {
	Lot *entity.BaseLot `json:"lot"`
}

type lotUpdate struct {
	Status *entity.LotStatus `json:"status" example:"1"`
	*entity.BaseLot
}

type lotUpdateRequest struct {
	Lot lotUpdate `json:"lot"`
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
		serverErrorResponse(w, r, err)
		return
	}

	err = writeJSON(w, http.StatusOK, listLotResponse{lots}, nil)
	if err != nil {
		serverErrorResponse(w, r, err)
	}
}

// Get          godoc
// @Summary     Show lot
// @Description show lot
// @ID          lot
// @Tags        lots
// @Accept      json
// @Produce     json
// @Param       id path int true "Lot ID" Format(int64)
// @Success     200 {object} lotResponse
// @Failure     404
// @Failure     500
// @Router      /lots/{id} [get]
func (c *LotController) Show(w http.ResponseWriter, r *http.Request) {
	id, err := readIDParam("ID", r)
	if err != nil {
		notFoundResponse(w, r)
		return
	}
	lot, err := c.uc.Show(id)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrRecordNotFound):
			notFoundResponse(w, r)
		default:
			serverErrorResponse(w, r, err)
		}
		return
	}

	err = writeJSON(w, http.StatusOK, lotResponse{lot}, nil)
	if err != nil {
		serverErrorResponse(w, r, err)
	}
}

// Get          godoc
// @Summary     Create lot
// @Description create lot
// @ID          create-lot
// @Tags        lots
// @Accept      json
// @Produce     json
// @Param       lot body lotRequest true "Create Lot"
// @Success     201
// @Failure     400
// @Failure     500
// @Router      /lots [post]
func (c *LotController) Create(w http.ResponseWriter, r *http.Request) {
	var input lotRequest

	err := readJSON(w, r, &input)
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	var fields = input.Lot
	lot := &entity.Lot{
		Status:      entity.LotPending,
		Title:       fields.Title,
		Description: fields.Description,
		StartPrice:  *fields.StartPrice,
		StepPrice:   *fields.StepPrice,
		StartAt:     *fields.StartAt,
		EndAt:       *fields.EndAt,
		Notify:      fields.Notify,
	}

	// Initialize a new Validator instance.
	v := validator.New()

	// Call the validate function and return a response containing the errors if
	// any of the checks fail.
	if entity.ValidateLot(v, lot); !v.Valid() {
		failedValidationResponse(w, r, v.Errors)
		return
	}

	err = c.uc.Create(lot)
	if err != nil {
		serverErrorResponse(w, r, err)
		return
	}

	// When sending a HTTP response, we want to include a Location header to let the
	// client know which URL they can find the newly-created resource at.
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/lots/%d", lot.ID))

	// Write a JSON response with a 201 Created status code, the lot data in the
	// response body, and the Location header.
	err = writeJSON(w, http.StatusCreated, envelope{"data": lot}, headers)
	if err != nil {
		serverErrorResponse(w, r, err)
	}

}

// Get          godoc
// @Summary     Update lot
// @Description update lot
// @ID          update-lot
// @Tags        lots
// @Accept      json
// @Produce     json
// @Param       id  path     int        true "Lot ID" Format(int64)
// @Param       lot body     lotUpdateRequest true "Update Lot"
// @Success     200 {object} lotResponse
// @Failure     400
// @Failure     500
// @Router      /lots/{id} [patch]
func (c *LotController) Update(w http.ResponseWriter, r *http.Request) {
	// Extract the ID from the URL.
	id, err := readIDParam("ID", r)
	if err != nil {
		notFoundResponse(w, r)
		return
	}

	// Fetch the existing record from the database, sending a 404 Not Found
	// response to the client if we couldn't find a matching record.
	lot, err := c.uc.Show(id)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrRecordNotFound):
			notFoundResponse(w, r)
		default:
			serverErrorResponse(w, r, err)
		}
		return
	}

	// Declare an input struct to hold the expected data from the client.
	var input lotUpdateRequest

	err = readJSON(w, r, &input)
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	var fields = input.Lot

	if fields.Status != nil {
		lot.Status = *fields.Status
	}

	if fields.Title != "" {
		lot.Title = fields.Title
	}

	if fields.Description != "" {
		lot.Description = fields.Description
	}

	if fields.StartPrice != nil {
		lot.StartPrice = *fields.StartPrice
	}

	if fields.StepPrice != nil {
		lot.StepPrice = *fields.StepPrice
	}

	if fields.StartAt != nil {
		lot.StartAt = *fields.StartAt
	}

	if fields.EndAt != nil {
		lot.EndAt = *fields.EndAt
	}

	lot.Notify = fields.Notify

	// Validate the updated lot record, sending the client a 422 Unprocessable Entity
	// response if any checks fail.
	v := validator.New()

	if entity.ValidateLot(v, lot); !v.Valid() {
		failedValidationResponse(w, r, v.Errors)
		return
	}

	err = c.uc.Update(lot)
	if err != nil {
		serverErrorResponse(w, r, err)
		return
	}

	responseLot := entity.Lot{
		ID:          lot.ID,
		Status:      lot.Status,
		Title:       lot.Title,
		Description: lot.Description,
		StartPrice:  lot.StartPrice,
		EndPrice:    lot.EndPrice,
		StepPrice:   lot.StepPrice,
		CreatorID:   lot.CreatorID,
		WinnerID:    lot.WinnerID,
		StartAt:     lot.StartAt,
		EndAt:       lot.EndAt,
		Notify:      lot.Notify,
		CreatedAt:   lot.CreatedAt,
		UpdatedAt:   lot.UpdatedAt,
	}

	// Write the updated lot record in a JSON response.
	err = writeJSON(w, http.StatusOK, envelope{"lot": responseLot}, nil)
	if err != nil {
		serverErrorResponse(w, r, err)
	}

}

// Get          godoc
// @Summary     Delete lot
// @Description delete lot
// @ID          delete-lot
// @Tags        lots
// @Accept      json
// @Produce     json
// @Param       id  path     int true "Lot ID" Format(int64)
// @Success     200
// @Failure     404
// @Failure     500
// @Router      /lots/{id} [delete]
func (c *LotController) Delete(w http.ResponseWriter, r *http.Request) {
	// Extract the ID from the URL.
	id, err := readIDParam("ID", r)
	if err != nil {
		notFoundResponse(w, r)
		return
	}

	// Delete the record from the database, sending a 404 Not Found response to the
	// client if there isn't a matching record.
	err = c.uc.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrRecordNotFound):
			notFoundResponse(w, r)
		default:
			serverErrorResponse(w, r, err)
		}
		return
	}

	// Return a 200 OK status code along with a success message.
	err = writeJSON(w, http.StatusOK, envelope{"message": "lot successfully deleted"}, nil)
	if err != nil {
		serverErrorResponse(w, r, err)
	}
}
