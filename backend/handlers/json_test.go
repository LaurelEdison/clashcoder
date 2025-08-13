package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go.uber.org/zap"
)

func TestRespondWithJSON(t *testing.T) {
	tests := []struct {
		name     string
		status   int
		payload  any
		wantBody string
	}{
		{
			name:     "emptyStruct_200",
			status:   http.StatusOK,
			payload:  struct{}{},
			wantBody: "{}\n",
		},
		{
			name:     "mapPayload_200",
			status:   http.StatusOK,
			payload:  map[string]string{"msg": "ok"},
			wantBody: `{"msg":"ok"}` + "\n",
		},
		{
			name:     "slicePayload_200",
			status:   http.StatusOK,
			payload:  []string{"a", "b"},
			wantBody: `["a","b"]` + "\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handlers{
				zapLogger: zap.NewNop(),
			}
			w := httptest.NewRecorder()

			h.RespondWithJSON(w, tt.status, tt.payload)

			assertStatus(t, w, tt.status)
			assertContentTypeJSON(t, w)
			assertBody(t, w, tt.wantBody)

		})
	}

}

//HELPER FUNCTIONS

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
