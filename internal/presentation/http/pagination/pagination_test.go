package pagination

import (
	"net/http/httptest"
	"testing"
)

func TestParseUsesDefaultValuesWhenQueryIsEmpty(t *testing.T) {
	req := httptest.NewRequest("GET", "/items", nil)
	params := Parse(req)

	if params.Limit != DefaultLimit {
		t.Fatalf("expected default limit %d, got %d", DefaultLimit, params.Limit)
	}
	if params.Offset != 0 {
		t.Fatalf("expected offset 0, got %d", params.Offset)
	}
}

func TestNewInfoHandlesZeroLimitWithoutPanicking(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("NewInfo panicked with zero limit: %v", r)
		}
	}()

	info := NewInfo(0, 15, 45)
	if info.Limit != DefaultLimit {
		t.Fatalf("expected fallback limit %d, got %d", DefaultLimit, info.Limit)
	}
	if info.Page != 1 {
		t.Fatalf("expected page 1, got %d", info.Page)
	}
}
