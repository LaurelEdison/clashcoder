package submission

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/LaurelEdison/clashcoder/backend/handlers"
	users "github.com/LaurelEdison/clashcoder/backend/handlers/user"
	"github.com/LaurelEdison/clashcoder/backend/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func CreateSubmission(h *handlers.Handlers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			ProblemID uuid.UUID `json:"problem_id"`
			Code      string    `json:"code"`
			Language  string    `json:"language"`
		}

		UserID, ok := users.GetUserId(r.Context())
		if !ok {
			h.RespondWithError(w, http.StatusUnauthorized, "Failed getting userid")
		}

		decoder := json.NewDecoder(r.Body)
		params := parameters{}
		if err := decoder.Decode(&params); err != nil {
			h.RespondWithError(w, http.StatusBadRequest, "Error decoding json")
			return
		}

		submission, err := h.DB.CreateSubmission(r.Context(), database.CreateSubmissionParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UserID:    UserID,
			ProblemID: params.ProblemID,
			Code:      params.Code,
			Language:  params.Language,
		})
		if err != nil {
			h.RespondWithError(w, http.StatusInternalServerError, "Error creating submission")
			return
		}
		h.RespondWithJSON(w, http.StatusOK, handlers.DatabaseSubmissiontoSubmission(submission))

	}
}

func GetSubmissionByID(h *handlers.Handlers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		IDparam := chi.URLParam(r, "id")
		submissionID, err := uuid.Parse(IDparam)
		if err != nil {
			h.RespondWithError(w, http.StatusBadRequest, "Failed parsing submission id")
		}
		submission, err := h.DB.GetSubmissionByID(r.Context(), submissionID)
		if err != nil {
			h.RespondWithError(w, http.StatusInternalServerError, "Error getting submission")
			return
		}
		h.RespondWithJSON(w, http.StatusOK, handlers.DatabaseSubmissiontoSubmission(submission))
	}
}

func GetSubmissionByUserID(h *handlers.Handlers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		problemIDstr := chi.URLParam(r, "problem_id")
		problemID, err := uuid.Parse(problemIDstr)
		if err != nil {
			h.RespondWithError(w, http.StatusUnauthorized, "Failed getting problem id")
			return
		}
		UserID, ok := users.GetUserId(r.Context())
		if !ok {
			h.RespondWithError(w, http.StatusUnauthorized, "Failed getting userid")
			return
		}

		submission, err := h.DB.GetLatestSubmissionByUserID(r.Context(), database.GetLatestSubmissionByUserIDParams{
			UserID:    UserID,
			ProblemID: problemID,
		})

		if err != nil {
			h.RespondWithError(w, http.StatusBadRequest, "Error fetching submissions")
			return
		}

		h.RespondWithJSON(w, http.StatusOK, handlers.DatabaseSubmissiontoSubmission(submission))
	}
}

func GetAllSubmissionsByUserID(h *handlers.Handlers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		problemIDstr := chi.URLParam(r, "problem_id")
		problemID, err := uuid.Parse(problemIDstr)
		if err != nil {
			h.RespondWithError(w, http.StatusUnauthorized, "Failed getting problem id")
			return
		}
		UserID, ok := users.GetUserId(r.Context())
		if !ok {
			h.RespondWithError(w, http.StatusUnauthorized, "Failed getting userid")
		}

		submissions, err := h.DB.GetAllSubmissionByUserID(r.Context(), database.GetAllSubmissionByUserIDParams{
			UserID:    UserID,
			ProblemID: problemID,
		})

		if err != nil {
			h.RespondWithError(w, http.StatusBadRequest, "Could not find user")
			return
		}
		h.RespondWithJSON(w, http.StatusOK, handlers.DatabaseSubmissionsToSubmissions(submissions))

	}
}
