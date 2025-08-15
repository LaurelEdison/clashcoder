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
	zapLogger, err := zap.NewProduction()
	mockDb := &mockDB{}
	if err != nil {
		t.Errorf("Failed to initialize zapLogger: Error: %v", err)
	}
	h := New(zapLogger, &mockDb.DB)

	if h.zapLogger == nil {
		t.Errorf("Expected zapLogger to be initialized, got nil")
	}
}
