package handlers

import (
	"github.com/LaurelEdison/clashcoder/backend/internal/database"
	"go.uber.org/zap"
	"testing"
)

type mockDB struct {
	DB database.Queries
}

func TestInitHandlers(t *testing.T) {
	zapLogger := zap.NewNop()
	mockDb := &mockDB{}
	h := New(zapLogger, &mockDb.DB)
	if h.zapLogger == nil {
		t.Errorf("Expected zapLogger to be initialized, got nil")
	}
}
