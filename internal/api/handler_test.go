package api

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/clevertechware/demo-renovate-basic-service-plan/internal/plan"
)

type stubFinder struct {
	plans map[string]plan.Plan
	err   error
}

func (s stubFinder) FindByCode(_ context.Context, code string) (plan.Plan, error) {
	if s.err != nil {
		return plan.Plan{}, s.err
	}
	p, ok := s.plans[code]
	if !ok {
		return plan.Plan{}, plan.ErrNotFound
	}
	return p, nil
}

func get(t *testing.T, finder PlanFinder, path string) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	NewRouter(finder).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
	return rec
}

func TestPlanEndpointReturnsSVGForKnownWarehouse(t *testing.T) {
	t.Parallel()

	finder := stubFinder{plans: map[string]plan.Plan{
		"LYS-01": {WarehouseCode: "LYS-01", Width: 600, Height: 300, Zones: []plan.Zone{{Label: "Stockage"}}},
	}}

	rec := get(t, finder, "/plans/LYS-01.svg")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if got := rec.Header().Get("Content-Type"); got != "image/svg+xml" {
		t.Errorf("Content-Type = %q, want image/svg+xml", got)
	}
	if !strings.Contains(rec.Body.String(), "Stockage") {
		t.Errorf("body does not contain the zone label:\n%s", rec.Body.String())
	}
}

func TestPlanEndpointStatusForUnsuccessfulLookups(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		finder stubFinder
		path   string
		want   int
	}{
		{"unknown warehouse is not found", stubFinder{}, "/plans/XXX-99.svg", http.StatusNotFound},
		{"missing svg extension is not found", stubFinder{}, "/plans/LYS-01", http.StatusNotFound},
		{"empty code is not found", stubFinder{}, "/plans/.svg", http.StatusNotFound},
		{"repository failure is an internal error", stubFinder{err: errors.New("connection refused")},
			"/plans/LYS-01.svg", http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := get(t, tt.finder, tt.path).Code; got != tt.want {
				t.Errorf("GET %s status = %d, want %d", tt.path, got, tt.want)
			}
		})
	}
}

func TestPlanEndpointDoesNotLeakRepositoryErrorToClient(t *testing.T) {
	t.Parallel()

	rec := get(t, stubFinder{err: errors.New("password authentication failed")}, "/plans/LYS-01.svg")

	if strings.Contains(rec.Body.String(), "password") {
		t.Errorf("body leaks the internal error: %q", rec.Body.String())
	}
}

func TestHealthzReportsOk(t *testing.T) {
	t.Parallel()

	rec := get(t, stubFinder{}, "/healthz")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if !strings.Contains(rec.Body.String(), `"status":"ok"`) {
		t.Errorf("body = %q, want status ok", rec.Body.String())
	}
}
