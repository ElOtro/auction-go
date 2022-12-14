// Package v1 implements routing paths. Each services in own file.
package v1

import (
	_ "github.com/ElOtro/auction-go/docs" // docs is generated by Swag CLI, you have to import it.
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	// Swagger docs.
)

// Create a Handlers struct which wraps all models.
type Handlers struct {
	controllers Controllers
}

// For ease of use, we also add a NewHandlers() method which
// returns a Handlers struct
func NewHandlers(controllers Controllers) *Handlers {
	return &Handlers{controllers: controllers}
}

// NewHandlers -.
// Swagger spec:
// @title       Auction API
// @description Using an api service
// @version     1.0
// @host        localhost:8080
// @BasePath    /v1
func (h *Handlers) Routes() *chi.Mux {
	mux := chi.NewRouter()
	// A good base middleware stack
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	// Swagger
	mux.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), // The url pointing to API definition
	))

	// Routers
	mux.Route("/v1", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Post("/register", h.controllers.Session.Register)
			r.Post("/auth", h.controllers.Session.login)
		})

		r.Route("/users", func(r chi.Router) {
			r.Use(h.controllers.Session.authenticate)
			{
				r.Get("/", h.controllers.User.List)
				r.Get("/{ID}", h.controllers.User.Show)
			}
		})

		r.Route("/lots", func(r chi.Router) {
			r.Use(h.controllers.Session.authenticate)
			{
				r.Get("/", h.controllers.Lot.List)
				r.Get("/{ID}", h.controllers.Lot.Show)
				r.Post("/", h.controllers.Lot.Create)
				r.Patch("/{ID}", h.controllers.Lot.Update)
				r.Delete("/{ID}", h.controllers.Lot.Delete)
				// bids
				r.Get("/{ID}/bids", h.controllers.Bid.List)
				r.Post("/{ID}/bids", h.controllers.Bid.Create)
			}
		})
	})

	return mux
}
