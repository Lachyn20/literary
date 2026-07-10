package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hemra-siirow/literary/internal/domain/repository"
)

func TestRequireRoleEnforcesRole(t *testing.T) {
	handler := RequireRoles("admin")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := context.WithValue(req.Context(), userKey, &repository.TokenClaims{Role: "editor"})
	req = req.WithContext(ctx)

	rrw := httptest.NewRecorder()
	handler.ServeHTTP(rrw, req)

	if rrw.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", rrw.Code)
	}
}
