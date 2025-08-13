package routes

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/LaurelEdison/clashcoder/backend/handlers"
	"github.com/go-chi/chi/v5"
)

func TestRoutes(t *testing.T) {
	h := &handlers.Handlers{}
	router := chi.NewRouter()
	SetupRoutes(router, h)

	tests := []struct {
		name            string
		method          string
		path            string
		wantStatus      int
		wantBody        string
		wantContentType string
	}{
		{
			name:            "GETHealth_200",
			method:          "GET",
			path:            "/healthz",
			wantStatus:      http.StatusOK,
			wantBody:        "{}\n",
			wantContentType: "application/json",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assertStatus(t, w, tt.wantStatus)
			assertContentTypeJSON(t, w)
			assertBody(t, w, tt.wantBody)
		})
	}
}

// HELPER FUNCTIONS
func assertStatus(t *testing.T, r *httptest.ResponseRecorder, want int) {
	t.Helper()
	if got := r.Result().StatusCode; got != want {
		t.Errorf("Unexpected status code, expected %d, actual %d", want, got)
	}
}

func assertContentTypeJSON(t *testing.T, r *httptest.ResponseRecorder) {
	t.Helper()
	if ct := r.Result().Header.Get("Content-Type"); !strings.Contains(ct, "application/json") {
		t.Errorf("Unexpected Content-Type , expected application/json , actual %q", ct)
	}
}

func assertBody(t *testing.T, r *httptest.ResponseRecorder, wantBody string) {
	t.Helper()
	if body := r.Body.String(); body != wantBody {
		t.Errorf("Unexpected body, expected %q, got %q", wantBody, body)
	}
}
