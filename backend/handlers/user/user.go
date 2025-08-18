package users

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/LaurelEdison/clashcoder/backend/handlers"
	"github.com/LaurelEdison/clashcoder/backend/internal/database"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(h *handlers.Handlers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		decoder := json.NewDecoder(r.Body)
		params := parameters{}
		if err := decoder.Decode(&params); err != nil {
			h.RespondWithError(w, 400, "Error parsing json: %v")
			return
		}

		passhash, err := HashPassword(params.Password)
		if err != nil {
			h.ZapLogger.Error("Error hashing password", zap.Error(err))
			return
		}

		user, err := h.DB.CreateUser(r.Context(), database.CreateUserParams{
			ID:           uuid.New(),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			Name:         params.Name,
			Email:        params.Email,
			PasswordHash: passhash,
		})
		if err != nil {
			h.RespondWithError(w, 500, fmt.Sprintf("Error creating user: %v", err))
			return
		}

		h.RespondWithJSON(w, 200, handlers.DatabaseUserToUser(user))
	}
}

func GetUserByEmail(h *handlers.Handlers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			Email string `json:"email"`
		}
		decoder := json.NewDecoder(r.Body)
		params := parameters{}
		if err := decoder.Decode(&params); err != nil {
			h.RespondWithError(w, 400, "Error parsing json: %v")
			return
		}
		user, err := h.DB.GetUserByEmail(r.Context(), params.Email)
		if err != nil {
			h.RespondWithError(w, http.StatusInternalServerError, "Could not fetch user")
			return
		}

		h.RespondWithJSON(w, http.StatusOK, handlers.DatabaseUserToUser(user))
	}
}

func FetchProfileSelf(h *handlers.Handlers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userID, ok := r.Context().Value(UserIDKey).(uuid.UUID)
		if !ok {
			h.RespondWithError(w, http.StatusUnauthorized, "No user in context")
			return
		}
		user, err := h.DB.GetUserByID(r.Context(), userID)
		if err != nil {
			h.RespondWithError(w, http.StatusInternalServerError, "Could not fetch user")
			return
		}

		h.RespondWithJSON(w, http.StatusOK, handlers.DatabaseUserToUser(user))

	}
}

// HELPERS
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), err
}

func CheckPasswordHash(HashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(HashedPassword), []byte(password))
	return err == nil
}

type ContextKey string

const (
	UserIDKey ContextKey = "user_id"
)

func GetUserId(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(UserIDKey).(uuid.UUID)
	return id, ok
}
