package v1

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/pascaldekloe/jwt"
)

func (c *UserController) authenticate(next http.Handler) http.Handler {
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
			case errors.Is(err, ErrRecordNotFound):
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
