package problem

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/LaurelEdison/clashcoder/backend/handlers"
	"github.com/LaurelEdison/clashcoder/backend/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func GetProblemTestsByProblemID(h *handlers.Handlers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		problemIDstr := chi.URLParam(r, "id")
		problemID, err := uuid.Parse(problemIDstr)
		if err != nil {
			h.RespondWithError(w, http.StatusBadRequest, "Error parsing problemID")
			return
		}

		problemTests, err := h.DB.GetProblemTestsByProblemID(r.Context(), problemID)
		if err != nil {
			h.RespondWithError(w, http.StatusBadRequest, "Error getting problem tests")
			return
		}

		h.RespondWithJSON(w, http.StatusOK, handlers.DatabaseProblemTestsToProblemTests(problemTests))
	}
}

func CreateProblemTest(h *handlers.Handlers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		problemIDstr := chi.URLParam(r, "id")
		problemID, err := uuid.Parse(problemIDstr)
		if err != nil {
			h.RespondWithError(w, http.StatusBadRequest, "Error parsing problemID")
			return
		}
		type parameters struct {
			TestCode string `json:"test_code"`
		}

		decoder := json.NewDecoder(r.Body)
		params := parameters{}
		if err := decoder.Decode(&params); err != nil {
			h.RespondWithError(w, http.StatusBadRequest, "Error decoding json")
		}

		problemTest, err := h.DB.CreateProblemTests(r.Context(), database.CreateProblemTestsParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			ProblemID: problemID,
			TestCode:  params.TestCode,
		})

		if err != nil {
			h.RespondWithError(w, http.StatusInternalServerError, "Error adding problem test")
		}

		h.RespondWithJSON(w, http.StatusOK, handlers.DatabaseProblemTestToProblemTest(problemTest))

	}
}
