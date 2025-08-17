package handlers

import (
	"testing"
	"time"

	"github.com/LaurelEdison/clashcoder/backend/internal/database"
	"github.com/google/uuid"
)

func TestDataBaseUserToUser(t *testing.T) {
	now := time.Now()
	userId := uuid.New()

	dbUser := database.User{
		ID:           userId,
		CreatedAt:    now,
		UpdatedAt:    now,
		Name:         "Alice",
		Email:        "Alice@gmail.com",
		PasswordHash: "hashedpassword",
	}

	want := User{
		ID:        userId,
		CreatedAt: now,
		UpdatedAt: now,
		Name:      "Alice",
		Email:     "Alice@gmail.com",
	}
	got := databaseUserToUser(dbUser)

	if got != want {
		t.Errorf("Database user to user got %v, expected %v", got, want)
	}

}
