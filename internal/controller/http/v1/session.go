package v1

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/ElOtro/auction-go/internal/entity"
	"github.com/ElOtro/auction-go/internal/validator"
	"github.com/pascaldekloe/jwt"
)

// List         godoc
// @Summary     Login user
// @Description Login user
// @ID          email
// @Tags        users
// @Accept      json
// @Produce     json
// @Success     200 {object} listUserResponse
// @Router      /auth [post]
func (c *UserController) login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := readJSON(w, r, &input)
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	// Validate the email and password provided by the client.
	v := validator.New()

	entity.ValidateEmail(v, input.Email)
	entity.ValidatePasswordPlaintext(v, input.Password)

	if !v.Valid() {
		failedValidationResponse(w, r, v.Errors)
		return
	}

	// Lookup the user record based on the email address. If no matching user was
	// found, then we call the app.invalidCredentialsResponse() helper to send a 401
	// Unauthorized response to the client (we will create this helper in a moment).
	user, err := c.uc.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrRecordNotFound):
			invalidCredentialsResponse(w, r)
		default:
			serverErrorResponse(w, r, err)
		}
		return
	}

	// Check if the provided password matches the actual password for the user.
	match, err := user.Password.Matches(input.Password)
	if err != nil {
		serverErrorResponse(w, r, err)
		return
	}

	// If the passwords don't match, then we call the app.invalidCredentialsResponse()
	// helper again and return.
	if !match {
		invalidCredentialsResponse(w, r)
		return
	}

	// Create a JWT claims struct containing the user ID as the subject, with an issued
	// time of now and validity window of the next 24 hours. We also set the issuer and
	// audience to a unique identifier for our application.
	var claims jwt.Claims
	claims.Subject = strconv.FormatInt(user.ID, 10)
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.NotBefore = jwt.NewNumericTime(time.Now())
	claims.Expires = jwt.NewNumericTime(time.Now().Add(24 * time.Hour))
	claims.Issuer = "auction-go"
	claims.Audiences = []string{"auction-go"}

	// Sign the JWT claims using the HMAC-SHA256 algorithm and the secret key from the // application config. This returns a []byte slice containing the JWT as a base64- // encoded string.
	jwtBytes, err := claims.HMACSign(jwt.HS256, []byte(c.jwtSecret))
	if err != nil {
		serverErrorResponse(w, r, err)
		return
	}

	//Encode the token to JSON and send it in the response along with a 201 Created
	//status code.
	err = writeJSON(w, http.StatusCreated, envelope{"token": string(jwtBytes)}, nil)
	if err != nil {
		serverErrorResponse(w, r, err)
	}
}
