package routes

import (
	"github.com/LaurelEdison/clashcoder/backend/handlers"
	"github.com/LaurelEdison/clashcoder/backend/handlers/auth"
	"github.com/LaurelEdison/clashcoder/backend/handlers/lobby"
	"github.com/LaurelEdison/clashcoder/backend/handlers/problem"
	"github.com/LaurelEdison/clashcoder/backend/handlers/submission"
	users "github.com/LaurelEdison/clashcoder/backend/handlers/user"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(router chi.Router, h *handlers.Handlers) {
	router.Get("/healthz", h.HandlerHealth)
	router.Post("/users", users.SignUp(h))
	router.Post("/login", auth.Login(h))
	router.Get("/problems", problem.GetAllProblems(h))
	router.Get("/problems/{id}", problem.GetProblemById(h))
	router.Get("/problems/random", problem.GetProblemByRandom(h))

	router.Group(func(router chi.Router) {
		router.Use(auth.JWTAuthMiddleWare)
		router.Get("/me", users.FetchProfileSelf(h))
		router.Post("/submissions", submission.CreateSubmission(h))
		router.Get("/submissions/{problem_id}/latest", submission.GetSubmissionByUserID(h))
		router.Get("/submissions/{problem_id}/all", submission.GetAllSubmissionsByUserID(h))

		router.Group(func(router chi.Router) {
			router.Use(auth.RequireAdmin)
			router.Post("/problems", problem.CreateProblem(h))
			router.Post("/problems/{id}/tests", problem.CreateProblemTest(h))
			router.Get("/problems/{id}/tests/all", problem.GetProblemTestsByProblemID(h))
		})
	})
}
