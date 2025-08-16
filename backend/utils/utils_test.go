package utils

import (
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestGetPort(t *testing.T) {
	zapLogger := zap.NewNop()

	t.Run("Port not set", func(t *testing.T) {
		original := os.Getenv("PORT")
		if err := os.Setenv("PORT", original); err != nil {
			t.Errorf("Error could reset env: %v", err)
		}

		if err := os.Unsetenv("PORT"); err != nil {
			t.Errorf("Error could not unset env: %v", err)
		}
		got := GetPort(zapLogger)
		assert.Equal(t, got, "8080")

	})

	t.Run("Port set", func(t *testing.T) {
		original := os.Getenv("PORT")
		if err := os.Setenv("PORT", original); err != nil {
			t.Errorf("Error could reset env: %v", err)
		}

		if err := os.Setenv("PORT", "1234"); err != nil {
			t.Errorf("Error, could not set env: %v", err)
		}
		got := GetPort(zapLogger)
		assert.Equal(t, got, "1234")
	})
}

func TestGetDbUrl(t *testing.T) {
	zapLogger := zap.NewNop()

	t.Run("DB url not set", func(t *testing.T) {
		original := os.Getenv("DB_URL")
		if err := os.Setenv("DB_URL", original); err != nil {
			t.Errorf("Error resetting env: %v", err)
		}

		got := GetDBUrl(zapLogger)

		assert.Empty(t, got, "Expected empty string when db is not set")

	})

	t.Run("DB Url set", func(t *testing.T) {
		original := os.Getenv("DB_URL")
		if err := os.Setenv("DB_URL", original); err != nil {
			t.Errorf("Error resetting env: %v", err)
		}
		if err := os.Setenv("DB_URL", "postgres://postgres:postgres@localhost:5433/clashcoder?sslmode=disable"); err != nil {
			t.Errorf("Error, could not set env: %v", err)
		}
		got := GetDBUrl(zapLogger)
		assert.Equal(t, got, "postgres://postgres:postgres@localhost:5433/clashcoder?sslmode=disable")

	})
}

func TestSetupCors(t *testing.T) {
	zapLogger, err := zap.NewProduction()
	if err != nil {
		t.Errorf("Error, could not create zapLogger: %v", err)
	}
	router := chi.NewRouter()
	SetupCors(zapLogger, router)
	assert.NotNil(t, router)
}
