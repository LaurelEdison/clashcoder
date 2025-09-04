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

func CreateLobby(h *handlers.Handlers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			MaxUsers int `json:"max_users"`
		}

		UserID, ok := users.GetUserId(r.Context())
		if !ok {
			h.RespondWithError(w, http.StatusUnauthorized, "Failed getting userid")
			return
		}
		decoder := json.NewDecoder(r.Body)
		params := &parameters{}
		err := decoder.Decode(params)
		if err != nil {
			h.RespondWithError(w, http.StatusBadRequest, "Could not decode json")
			return
		}
		if params.MaxUsers < 2 {
			h.RespondWithError(w, http.StatusBadRequest, "Invalid max users must be atleast 2")
			return
		}
		if params.MaxUsers > 8 {
			h.RespondWithError(w, http.StatusBadRequest, "Invalid max users must be atleast 2")
			return
		}

		LobbyCode, err := GenerateLobbyCode(8)
		if err != nil {
			h.RespondWithError(w, http.StatusInternalServerError, "Could not generate lobby code")
			return
		}
		Lobby, err := h.DB.CreateLobby(r.Context(), database.CreateLobbyParams{
			ID:         uuid.New(),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
			InviteCode: LobbyCode,
			MaxUsers:   int32(params.MaxUsers),
			StartedAt:  sql.NullTime{Time: time.Now()},
			ReadyState: false,
			Status:     "waiting",
		})
		if err != nil {
			h.RespondWithError(w, http.StatusInternalServerError, "Failed to create lobby")
			return
		}

		_, err = CreateLobbyHost(h, r, Lobby.ID, UserID)
		if err != nil {
			h.RespondWithError(w, http.StatusInternalServerError, "Failed to create lobby host")
			return
		}

		h.RespondWithJSON(w, http.StatusOK, handlers.DatabaseLobbyToLobby(Lobby))
	}
}

func GetLobbyById(h *handlers.Handlers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		LobbyIDstr := chi.URLParam(r, "lobby_id")
		LobbyID, err := uuid.Parse(LobbyIDstr)
		if err != nil {
			h.RespondWithError(w, http.StatusBadRequest, "Failed decoding lobby id")
			return
		}
		Lobby, err := h.DB.GetLobbyById(r.Context(), LobbyID)
		if err != nil {
			h.RespondWithError(w, http.StatusBadRequest, "Could not fetch lobby")
			return
		}
		h.RespondWithJSON(w, http.StatusOK, handlers.DatabaseLobbyToLobby(Lobby))

	}
}

func StartMatch(h *handlers.Handlers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			ProblemIDstr string `json:"problem_id"`
		}
		decoder := json.NewDecoder(r.Body)
		params := parameters{}
		if err := decoder.Decode(&params); err != nil {
			h.RespondWithError(w, http.StatusBadRequest, "Failed decoding json")
			return
		}

		LobbyIDstr := chi.URLParam(r, "lobby_id")
		LobbyID, err := uuid.Parse(LobbyIDstr)
		if err != nil {
			h.RespondWithError(w, http.StatusBadRequest, "Failed decoding lobby id")
			return
		}
		LobbyHost, err := h.DB.GetHostFromLobbyID(r.Context(), LobbyID)
		if err != nil {
			h.RespondWithError(w, http.StatusInternalServerError, "Could not get lobby host")
			return
		}

		pUserID, ok := users.GetUserId(r.Context())
		if !ok {
			h.RespondWithError(w, http.StatusUnauthorized, "Failed getting userid")
			return
		}

		isHost := LobbyHost.UserID == pUserID
		if !isHost {
			h.RespondWithError(w, http.StatusUnauthorized, "User is not host")
			return
		}

		var ProblemID uuid.UUID
		if params.ProblemIDstr == "" {
			Problem, err := h.DB.GetRandomProblem(r.Context())
			if err != nil {
				h.RespondWithError(w, http.StatusBadRequest, "Could not get problem")
				return
			}
			ProblemID = Problem.ID
		} else {
			ProblemID, err = uuid.Parse(params.ProblemIDstr)
			if err != nil {
				h.RespondWithError(w, http.StatusBadRequest, "Could not parse problem id")
				return
			}
			_, err = h.DB.GetProblemByID(r.Context(), ProblemID)
			if err != nil {
				h.RespondWithError(w, http.StatusBadRequest, "Could not get problem")
				return
			}
		}
		err = h.DB.SelectProblem(r.Context(), database.SelectProblemParams{
			ID:        LobbyID,
			ProblemID: uuid.NullUUID{UUID: ProblemID},
		})
		if err != nil {
			h.RespondWithError(w, http.StatusBadRequest, "Could update with new problem")
			return
		}
		err = h.DB.UpdateLobbyStartEnd(r.Context(), database.UpdateLobbyStartEndParams{
			ID:        LobbyID,
			StartedAt: sql.NullTime{Time: time.Now()},
			EndedAt:   sql.NullTime{},
		})
		if err != nil {
			h.RespondWithError(w, http.StatusBadRequest, "Could update start, end time")
			return
		}

		err = h.DB.UpdateLobbyStatus(r.Context(), database.UpdateLobbyStatusParams{
			ID:     LobbyID,
			Status: "in_progress",
		})
		if err != nil {
			h.RespondWithError(w, http.StatusInternalServerError, "Could not update lobby status")
			return
		}

		h.RespondWithJSON(w, http.StatusOK, struct{}{})
	}
}
