package lobby

import (
	"net/http"
	"time"

	"github.com/LaurelEdison/clashcoder/backend/handlers"
	users "github.com/LaurelEdison/clashcoder/backend/handlers/user"
	"github.com/LaurelEdison/clashcoder/backend/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func CreateLobbyUser(h *handlers.Handlers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		LobbyIDstr := chi.URLParam(r, "lobby_id")
		LobbyID, err := uuid.Parse(LobbyIDstr)
		if err != nil {
			h.RespondWithError(w, http.StatusBadRequest, "Failed decoding lobby id")
			return
		}
		UserID, ok := users.GetUserId(r.Context())
		if !ok {
			h.RespondWithError(w, http.StatusUnauthorized, "Failed getting userid")
			return
		}

		LobbyUser, err := h.DB.CreateLobbyUser(r.Context(), database.CreateLobbyUserParams{
			LobbyID:  LobbyID,
			UserID:   UserID,
			JoinedAt: time.Now(),
			Role:     "player",
		})
		if err != nil {
			h.RespondWithError(w, http.StatusInternalServerError, "Failed to create lobby user")
			return
		}
		h.RespondWithJSON(w, http.StatusOK, handlers.DatabaseLobbyUserToLobbyUser(LobbyUser))

	}
}

func CreateLobbyHost(h *handlers.Handlers, r *http.Request, LobbyID uuid.UUID, UserID uuid.UUID) (database.LobbyUser, error) {
	LobbyUser, err := h.DB.CreateLobbyUser(r.Context(), database.CreateLobbyUserParams{
		LobbyID:  LobbyID,
		UserID:   UserID,
		JoinedAt: time.Now(),
		Role:     "host",
	})

	return LobbyUser, err
}

