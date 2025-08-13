package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandlerHealth(t *testing.T) {
	//Arrange
	h := &Handlers{}
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	//Act
	h.HandlerHealth(w, req)

	//Assert
	resp := w.Result()

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d got %d",
			http.StatusOK, resp.StatusCode)
	}

	if ct := resp.Header.Get("content-type"); !strings.Contains(ct, "application/json") {
		t.Errorf("Expected Content-Type header to contain application/json got %s", ct)
	}

	bodyBytes := w.Body.Bytes()
	expected := "{}\n"
	if string(bodyBytes) != expected {
		t.Errorf("Expected body %s, got %s", string(bodyBytes), expected)
	}

}
