package router

import (
	"net/http"
	"pokedex/pkg/http/middleware"

	"github.com/go-chi/chi"
)

// Handlers holds http.HandlerFuncs for API endpoints. This allows for loose coupling between the
// route and the handler. Handlers can not be constructed outside of the router allowing for easier
// dependency injection for each handler
type Handlers struct {
	GetPokemonByName           http.HandlerFunc
	GetPokemonByNameTranslated http.HandlerFunc

	HealthCheck http.HandlerFunc
}

// New constructs a new router
func New(handlers Handlers) *chi.Mux {
	router := chi.NewRouter()

	middleware.Common(router)

	router.Get("/_healthcheck", handlers.HealthCheck)

	router.Route("/v1/pokemon", func(r chi.Router) {
		r.Get("/{name}", handlers.GetPokemonByName)
		r.Get("/translated/{name}", handlers.GetPokemonByNameTranslated)
	})

	return router
}
