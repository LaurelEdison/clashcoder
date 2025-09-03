package lobby

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/LaurelEdison/clashcoder/backend/handlers"
	users "github.com/LaurelEdison/clashcoder/backend/handlers/user"
	"github.com/LaurelEdison/clashcoder/backend/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func GenerateLobbyCode(length int) (string, error) {
	randomBytes := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, randomBytes); err != nil {
		return "", err
	}
	hash := sha256.Sum256(randomBytes)
	hexStr := hex.EncodeToString(hash[:])
	return hexStr[:length], nil
}

