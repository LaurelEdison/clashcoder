package handlers

import (
	"database/sql"
	"time"

	"github.com/LaurelEdison/clashcoder/backend/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID    `json:"id"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
	Name         string       `json:"name"`
	Email        string       `json:"email"`
	PasswordHash string       `json:"password_hash"`
	LastLoginAt  sql.NullTime `json:"last_login_at"`
}

func databaseUserToUser(dbUser database.User) User {
	return User{
		ID:           dbUser.ID,
		CreatedAt:    dbUser.CreatedAt,
		UpdatedAt:    dbUser.UpdatedAt,
		Name:         dbUser.Name,
		Email:        dbUser.Email,
		PasswordHash: dbUser.PasswordHash,
		LastLoginAt:  dbUser.LastLoginAt,
	}
}
