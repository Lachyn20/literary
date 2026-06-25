package middleware

import (
	"net/http"

	"golang.org/x/time/rate"
)

func RateLimit(rps float64, burst int) func(next http.Handler) http.Handler {
	lim := rate.NewLimiter(rate.Limit(rps), burst)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !lim.Allow() {
				w.WriteHeader(http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
