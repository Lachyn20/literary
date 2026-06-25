package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hemra-siirow/literary/internal/domain/repository"
	m "github.com/hemra-siirow/literary/internal/presentation/http/middleware"
)

type RouteRegistrar interface {
	RegisterRoutes(r chi.Router)
}

func NewRouter(registrars []RouteRegistrar, allowedOrigin string, tokenGen interface{}) http.Handler {
	r := chi.NewRouter()

	// CORS
	r.Use(m.CORS(allowedOrigin))

	// JWT middleware if provided
	if tg, ok := tokenGen.(mJWTProvider); ok {
		r.Use(m.JWTAuth(tg.TokenGenerator()))
	}

	// register routes
	for _, reg := range registrars {
		reg.RegisterRoutes(r)
	}

	return r
}

// mJWTProvider is a tiny adapter to avoid import cycles — router accepts anything implementing this.
type mJWTProvider interface {
	TokenGenerator() repository.TokenGenerator
}
