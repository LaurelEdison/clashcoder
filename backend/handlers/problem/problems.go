package problem

import (
	"net/http"

	"github.com/LaurelEdison/clashcoder/backend/handlers"
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
