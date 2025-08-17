package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/LaurelEdison/clashcoder/backend/internal/database"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handlers) SignUp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		h.RespondWithError(w, 400, "Error parsing json: %v")
	}

	passhash, err := h.HashPassword(params.Password)
	if err != nil {
		h.zapLogger.Error("Error hashing password", zap.Error(err))
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
	}
	h.RespondWithJSON(w, 200, databaseUserToUser(user))

}

func (h *Handlers) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		h.RespondWithError(w, 400, "Error parsing json: %v")
	}

}

func (h *Handlers) FetchProfileSelf(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		h.RespondWithError(w, http.StatusUnauthorized, "No user in context")
		return
	}
	user, err := h.DB.GetUserByID(r.Context(), userID)
	if err != nil {
		h.RespondWithError(w, http.StatusInternalServerError, "Could not fetch user")
		return
	}

	h.RespondWithJSON(w, http.StatusOK, databaseUserToUser(user))

}

// HELPERS
func (h *Handlers) HashPassword(password string) (string, error) {
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
