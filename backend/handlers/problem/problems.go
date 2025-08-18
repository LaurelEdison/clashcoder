package problem

import (
	"encoding/json"
	"net/http"

	"github.com/LaurelEdison/clashcoder/backend/handlers"
	"github.com/google/uuid"
)

func GetAll(h *handlers.Handlers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		problems, err := h.DB.GetAllProblems(r.Context())
		if err != nil {
			h.RespondWithError(w, http.StatusBadRequest, "Error getting problems")
			return
		}
		h.RespondWithJSON(w, http.StatusOK, handlers.DatabaseProblemsToProblems(problems))
	}
}

func GetById(h *handlers.Handlers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			ID uuid.UUID `json:"id"`
		}
		decoder := json.NewDecoder(r.Body)
		params := parameters{}

		if err := decoder.Decode(&params); err != nil {
			h.RespondWithError(w, http.StatusBadRequest, "Error parsing json")
			return
		}
		problem, err := h.DB.GetProblemByID(r.Context(), params.ID)

		if err != nil {
			h.RespondWithError(w, http.StatusBadRequest, "Could not get user")
			return
		}
		h.RespondWithJSON(w, http.StatusOK, handlers.DatabaseProblemToProblem(problem))
	}
}
