package v1

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ElOtro/auction-go/internal/entity"
	"github.com/ElOtro/auction-go/internal/validator"
	"github.com/pascaldekloe/jwt"
)

type SessionUseCase interface {
	Get(userID int64) (*entity.User, error)
	GetByEmail(email string) (*entity.User, error)
	Insert(*entity.User) error
}

type SessionController struct {
	uc        SessionUseCase
	jwtSecret string
}

type registerUser struct {
	Name     string `json:"name" example:"Test User"`
	Email    string `json:"email" example:"test@example.com"`
	Password string `json:"password" example:"12345678"`
}

type authUser struct {
	Email    string `json:"email" example:"test@example.com"`
	Password string `json:"password" example:"12345678"`
}

type tokenResponse struct {
	Token string `json:"token"`
}

func NewSessionController(uc SessionUseCase, jwtSecret string) *SessionController {
	return &SessionController{uc: uc, jwtSecret: jwtSecret}
}

func (c *SessionController) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add the "Vary: Authorization" header to the response. This indicates to any
		// caches that the response may vary based on the value of the Authorization
		// header in the request.
		w.Header().Add("Vary", "Authorization")

		// Retrieve the value of the Authorization header from the request. This will
		// return the empty string "" if there is no such header found.
		authorizationHeader := r.Header.Get("Authorization")

		if authorizationHeader == "" {
			invalidAuthenticationTokenResponse(w, r)
			// next.ServeHTTP(w, r)
			return
		}
		// Otherwise, we expect the value of the Authorization header to be in the format
		// "Bearer <token>". We try to split this into its constituent parts, and if the
		// header isn't in the expected format we return a 401 Unauthorized response
		// using the invalidAuthenticationTokenResponse() helper (which we will create
		// in a moment).
		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			invalidAuthenticationTokenResponse(w, r)
			return
		}

		// Extract the actual authentication token from the header parts.
		token := headerParts[1]

		// Parse the JWT and extract the claims. This will return an error if the JWT
		// contents doesn't match the signature (i.e. the token has been tampered with)
		// or the algorithm isn't valid.
		claims, err := jwt.HMACCheck([]byte(token), []byte(c.jwtSecret))
		if err != nil {
			invalidAuthenticationTokenResponse(w, r)
			return
		}

		// Check if the JWT is still valid at this moment in time.
		if !claims.Valid(time.Now()) {
			invalidAuthenticationTokenResponse(w, r)
			return
		}

		// Check that the issuer is our application.
		if claims.Issuer != "auction-go" {
			invalidAuthenticationTokenResponse(w, r)
			return
		}

		// Check that our application is in the expected audiences for the JWT.
		if !claims.AcceptAudience("auction-go") {
			invalidAuthenticationTokenResponse(w, r)
			return
		}

		// At this point, we know that the JWT is all OK and we can trust the data in
		// it. We extract the user ID from the claims subject and convert it from a
		// string into an int64.
		userID, err := strconv.ParseInt(claims.Subject, 10, 64)
		if err != nil {
			serverErrorResponse(w, r, err)
			return
		}

		/// Lookup the user record from the database.
		user, err := c.uc.Get(userID)
		if err != nil {
			switch {
			case errors.Is(err, entity.ErrRecordNotFound):
				invalidAuthenticationTokenResponse(w, r)
			default:
				serverErrorResponse(w, r, err)
			}
			return
		}

		// Call the contextSetUser() helper to add the user information to the request // context.
		r = contextSetUser(r, user)
		// Call the next handler in the chain.
		next.ServeHTTP(w, r)

	})
}

// List         godoc
// @Summary     Login user
// @Description login user
// @ID          email
// @Tags        sessions
// @Accept      json
// @Produce     json
// @Param       login body     authUser true "Login"
// @Success     201   {object} tokenResponse
// @Router      /auth [post]
func (c *SessionController) login(w http.ResponseWriter, r *http.Request) {
	var input authUser

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
	err = writeJSON(w, http.StatusCreated, tokenResponse{Token: string(jwtBytes)}, nil)
	if err != nil {
		serverErrorResponse(w, r, err)
	}
}

// Get          godoc
// @Summary     Register user
// @Description add by json user
// @Tags        sessions
// @Accept      json
// @Produce     json
// @Param       user body     registerUser true "Register user"
// @Success     201  {object} showUserResponse
// @Failure     422
// @Router      /register [post]
func (c *SessionController) Register(w http.ResponseWriter, r *http.Request) {
	// Create an anonymous struct to hold the expected data from the request body.
	var input registerUser

	// Parse the request body into the anonymous struct.
	err := readJSON(w, r, &input)
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

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
		return
	}

	// Write a JSON response containing the user data along with a 201 Created status
	// code.
	err = writeJSON(w, http.StatusCreated, envelope{"data": user}, nil)
	if err != nil {
		serverErrorResponse(w, r, err)
	}
}
