package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/hemra-siirow/literary/internal/domain/repository"
)

type key int

const (
	userKey key = iota
)

func UserFromContext(ctx context.Context) *repository.TokenClaims {
	v := ctx.Value(userKey)
	if v == nil {
		return nil
	}
	if tc, ok := v.(*repository.TokenClaims); ok {
		return tc
	}
	return nil
}

func JWTAuth(tokenGen repository.TokenGenerator) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if auth == "" {
				next.ServeHTTP(w, r.WithContext(r.Context()))
				return
			}
			parts := strings.SplitN(auth, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				next.ServeHTTP(w, r.WithContext(r.Context()))
				return
			}
			token := parts[1]
			claims, err := tokenGen.ValidateToken(token)
			if err != nil {
				next.ServeHTTP(w, r.WithContext(r.Context()))
				return
			}
			ctx := context.WithValue(r.Context(), userKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireRoles(roles ...string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := UserFromContext(r.Context())
			if claims == nil {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			for _, role := range roles {
				if claims.Role == role {
					next.ServeHTTP(w, r)
					return
				}
			}
			w.WriteHeader(http.StatusForbidden)
		})
	}
}
