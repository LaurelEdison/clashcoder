package problem

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/LaurelEdison/clashcoder/backend/handlers"
	"github.com/LaurelEdison/clashcoder/backend/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func GetAllProblems(h *handlers.Handlers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		problems, err := h.DB.GetAllProblems(r.Context())
		if err != nil {
			h.RespondWithError(w, http.StatusBadRequest, "Error getting problems")
			return
		}
		h.RespondWithJSON(w, http.StatusOK, handlers.DatabaseProblemsToProblems(problems))
	}
}

func GetProblemById(h *handlers.Handlers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "id")
		problemID, err := uuid.Parse(idParam)
		if err != nil {
			h.RespondWithError(w, http.StatusBadRequest, "Could not parse problem id")
			return
		}
		problem, err := h.DB.GetProblemByID(r.Context(), problemID)

		if err != nil {
			h.RespondWithError(w, http.StatusBadRequest, "Could not get user")
			return
		}
		h.RespondWithJSON(w, http.StatusOK, handlers.DatabaseProblemToProblem(problem))
	}
}

func GetProblemByRandom(h *handlers.Handlers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		problem, err := h.DB.GetRandomProblem(r.Context())
		if err != nil {
			h.RespondWithError(w, http.StatusBadRequest, "Could not get user")
			return
		}
		h.RespondWithJSON(w, http.StatusOK, handlers.DatabaseProblemToProblem(problem))
	}
}

func CreateProblem(h *handlers.Handlers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			Title         string         `json:"title"`
			Description   string         `json:"description"`
			Difficulty    sql.NullString `json:"difficulty"`
			TimeLimit     int32          `json:"time_limit"`
			MemoryLimitMb int32          `json:"memory_limit_mb"`
		}
		decoder := json.NewDecoder(r.Body)
		params := parameters{}
		if err := decoder.Decode(&params); err != nil {
			h.RespondWithError(w, http.StatusBadRequest, "Error decoding json")
			return
		}

		problem, err := h.DB.CreateProblem(r.Context(), database.CreateProblemParams{
			ID:            uuid.New(),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			Title:         params.Title,
			Description:   params.Description,
			Difficulty:    params.Difficulty,
			TimeLimit:     params.TimeLimit,
			MemoryLimitMb: params.MemoryLimitMb,
		})
		if err != nil {
			h.RespondWithError(w, http.StatusInternalServerError, "Error creating user")
			return
		}
		h.RespondWithJSON(w, http.StatusOK, handlers.DatabaseProblemToProblem(problem))

	}
}
