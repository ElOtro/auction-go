package v1

import (
	"errors"
	"net/http"

	"github.com/ElOtro/auction-go/internal/entity"
	"github.com/ElOtro/auction-go/internal/validator"
)

type UserUseCase interface {
	List() ([]*entity.User, error)
	Get(userID int64) (*entity.User, error)
	Insert(*entity.User) error
}

type UserController struct {
	uc        UserUseCase
	jwtSecret string
}

func NewUserController(uc UserUseCase, jwtSecret string) *UserController {
	return &UserController{uc: uc, jwtSecret: jwtSecret}
}

type listUserResponse struct {
	User []*entity.User `json:"users"`
}

type showUserResponse struct {
	User *entity.User `json:"user"`
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
// @Success     200 {object} showUserResponse
// @Router      /users/id [get]
func (c *UserController) Show(w http.ResponseWriter, r *http.Request) {
	id, err := readIDParam("userID", r)
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

func (c *UserController) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	// Create an anonymous struct to hold the expected data from the request body.
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Parse the request body into the anonymous struct.
	err := readJSON(w, r, &input)
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	// Copy the data from the request body into a new User struct. Notice also that we
	// set the Activated field to false, which isn't strictly necessary because the
	// Activated field will have the zero-value of false by default. But setting this
	// explicitly helps to make our intentions clear to anyone reading the code.
	user := &entity.User{
		Name:   input.Name,
		Email:  input.Email,
		Active: true,
	}

	// Use the Password.Set() method to generate and store the hashed and plaintext
	// passwords.
	err = user.Password.Set(input.Password)
	if err != nil {
		serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()

	// Validate the user struct and return the error messages to the client if any of
	// the checks fail.
	if entity.ValidateUser(v, user); !v.Valid() {
		failedValidationResponse(w, r, v.Errors)
		return
	}

	// Insert the user data into the database.
	err = c.uc.Insert(user)
	if err != nil {
		switch {
		// If we get a ErrDuplicateEmail error, use the v.AddError() method to manually
		// add a message to the validator instance, and then call our
		// failedValidationResponse() helper.
		case errors.Is(err, entity.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			failedValidationResponse(w, r, v.Errors)
		default:
			serverErrorResponse(w, r, err)
		}
	}

	// Write a JSON response containing the user data along with a 201 Created status
	// code.
	err = writeJSON(w, http.StatusCreated, envelope{"data": user}, nil)
	if err != nil {
		serverErrorResponse(w, r, err)
	}
}
