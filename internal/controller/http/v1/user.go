package v1

import (
	"net/http"

	"github.com/ElOtro/auction-go/internal/entity"
)

type UserUseCase interface {
	List() ([]*entity.User, error)
}

type UserController struct {
	uc UserUseCase
}

func NewUserController(uc UserUseCase) *UserController {
	return &UserController{uc: uc}
}

type listUserResponse struct {
	User []*entity.User `json:"users"`
}

// List         godoc
// @Summary     Show user list
// @Description Show all user list
// @ID          userList
// @Tags        users
// @Accept      json
// @Produce     json
// @Success     200 {object} listUserResponse
// @Router      /users [get]
func (c *UserController) List(w http.ResponseWriter, r *http.Request) {
	users, err := c.uc.List()
	if err != nil {
		errorResponse(w, r, http.StatusInternalServerError, "get users")

		return
	}

	err = writeJSON(w, http.StatusOK, listUserResponse{users}, nil)
	if err != nil {
		serverErrorResponse(w, r, err)
	}
}
