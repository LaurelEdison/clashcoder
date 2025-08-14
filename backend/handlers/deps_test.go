package handlers

import (
	"go.uber.org/zap"
	"testing"
)

func TestInitHandlers(t *testing.T) {
	zapLogger, err := zap.NewProduction()
	if err != nil {
		t.Errorf("Failed to initialize zapLogger: Error: %v", err)
	}
	h := New(zapLogger)
	if h.zapLogger == nil {
		t.Errorf("Expected zapLogger to be initialized, got nil")
	}
}
