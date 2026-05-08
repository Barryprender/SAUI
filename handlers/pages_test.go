package handlers_test

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"saui/handlers"
	"saui/internal/testutil"
	"saui/middleware"
)

func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

// newTestServer wires up a handler with a real in-memory gateway.
func newTestServer(t *testing.T) http.Handler {
	t.Helper()
	_, gw := testutil.NewStore(t)
	h := handlers.New(gw, discardLogger())

	mux := http.NewServeMux()
	page := func(hf http.HandlerFunc) http.Handler {
		return middleware.Chain(hf, middleware.Session, middleware.CSRF)
	}
	mux.Handle("GET /{$}", page(h.Home))
	mux.Handle("GET /why", page(h.Why))
	mux.Handle("GET /architecture", page(h.Architecture))
	mux.Handle("GET /stack", page(h.Stack))
	mux.Handle("GET /cases", page(h.Cases))
	mux.Handle("GET /cases/food-ordering", page(h.CaseFoodOrdering))
	mux.Handle("GET /cases/banking", page(h.CaseBanking))
	mux.Handle("GET /cases/healthcare", page(h.CaseHealthcare))
	mux.Handle("GET /cases/saas-dashboard", page(h.CaseSaaS))
	mux.Handle("GET /cases/distributed-systems", page(h.CaseDistributed))
	mux.Handle("GET /cases/micro-frontends", page(h.CaseMFE))
	mux.Handle("GET /testing", page(h.Testing))
	mux.Handle("GET /limits", page(h.Limits))
	mux.Handle("GET /blog", page(h.Blog))
	mux.Handle("GET /code", page(h.Code))
	return mux
}

func get(t *testing.T, srv http.Handler, path string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	return w
}

func assertStatus(t *testing.T, w *httptest.ResponseRecorder, want int) {
	t.Helper()
	if w.Code != want {
		t.Errorf("status %d, want %d", w.Code, want)
	}
}

func assertContains(t *testing.T, w *httptest.ResponseRecorder, substr string) {
	t.Helper()
	if !strings.Contains(w.Body.String(), substr) {
		t.Errorf("response body does not contain %q", substr)
	}
}

// --- Status code tests ---

func TestAllRoutes_Return200(t *testing.T) {
	srv := newTestServer(t)
	routes := []string{
		"/", "/why", "/architecture", "/stack",
		"/cases", "/cases/food-ordering", "/cases/banking",
		"/cases/healthcare", "/cases/saas-dashboard",
		"/cases/distributed-systems", "/cases/micro-frontends",
		"/testing", "/limits", "/blog", "/code",
	}
	for _, route := range routes {
		t.Run(route, func(t *testing.T) {
			w := get(t, srv, route)
			assertStatus(t, w, http.StatusOK)
		})
	}
}

func TestHome_Returns404ForUnknownPath(t *testing.T) {
	srv := newTestServer(t)
	w := get(t, srv, "/nonexistent-page")
	assertStatus(t, w, http.StatusNotFound)
}

// --- Content tests ---

func TestHome_ContainsCoreProposition(t *testing.T) {
	srv := newTestServer(t)
	w := get(t, srv, "/")
	assertStatus(t, w, http.StatusOK)
	assertContains(t, w, "The server owns state")
	assertContains(t, w, "The browser renders truth")
}

func TestHome_ContainsNavLinks(t *testing.T) {
	srv := newTestServer(t)
	w := get(t, srv, "/")
	assertStatus(t, w, http.StatusOK)
	for _, href := range []string{"/why", "/architecture", "/stack", "/limits"} {
		assertContains(t, w, href)
	}
}

func TestLimits_ContainsHonestLimits(t *testing.T) {
	srv := newTestServer(t)
	w := get(t, srv, "/limits")
	assertStatus(t, w, http.StatusOK)
	assertContains(t, w, "Real-time collaborative editing")
	assertContains(t, w, "Offline-first")
	assertContains(t, w, "Games and simulations")
}

func TestCases_ContainsAllCaseLinks(t *testing.T) {
	srv := newTestServer(t)
	w := get(t, srv, "/cases")
	assertStatus(t, w, http.StatusOK)
	for _, href := range []string{
		"/cases/food-ordering", "/cases/banking", "/cases/healthcare",
		"/cases/saas-dashboard", "/cases/distributed-systems", "/cases/micro-frontends",
	} {
		assertContains(t, w, href)
	}
}

func TestAllPages_ContainDoctype(t *testing.T) {
	srv := newTestServer(t)
	for _, route := range []string{"/", "/why", "/limits", "/cases"} {
		t.Run(route, func(t *testing.T) {
			w := get(t, srv, route)
			assertContains(t, w, "<!doctype html>")
		})
	}
}

func TestAllPages_ContainCSRFMeta(t *testing.T) {
	srv := newTestServer(t)
	for _, route := range []string{"/", "/why", "/architecture"} {
		t.Run(route, func(t *testing.T) {
			w := get(t, srv, route)
			assertContains(t, w, `name="csrf-token"`)
		})
	}
}

