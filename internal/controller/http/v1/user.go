package v1

import (
	"net/http"

	"github.com/ElOtro/auction-go/internal/entity"
)

type UserUseCase interface {
	List() ([]*entity.User, error)
	Get(userID int64) (*entity.User, error)
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

type showUserResponse struct {
	User *entity.User `json:"user"`
}

type registerUser struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
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

// Get          godoc
// @Summary     Show user
// @Description Show user
// @ID          user
// @Tags        users
// @Accept      json
// @Produce     json
// @Success     200  {object} showUserResponse
// @Router      /users/{id} [get]
func (c *UserController) Show(w http.ResponseWriter, r *http.Request) {
	id, err := readIDParam("ID", r)
	if err != nil {
		errorResponse(w, r, http.StatusBadRequest, err)
		return
	}
	user, err := c.uc.Get(id)
	if err != nil {
		errorResponse(w, r, http.StatusInternalServerError, err)

		return
	}

	err = writeJSON(w, http.StatusOK, showUserResponse{user}, nil)
	if err != nil {
		serverErrorResponse(w, r, err)
	}
}
