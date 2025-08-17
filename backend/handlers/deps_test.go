package handlers

import (
	"github.com/LaurelEdison/clashcoder/backend/internal/database"
	"go.uber.org/zap"
	"testing"
)

type MockDB struct {
	DB database.Queries
}

func TestInitHandlers(t *testing.T) {
	zapLogger := zap.NewNop()
	mockDb := &MockDB{}
	h := New(zapLogger, &mockDb.DB)
	if h.ZapLogger == nil {
		t.Errorf("Expected zapLogger to be initialized, got nil")
	}
}
