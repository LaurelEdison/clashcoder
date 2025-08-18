package routes

import (
	"github.com/LaurelEdison/clashcoder/backend/handlers"
	"github.com/LaurelEdison/clashcoder/backend/handlers/auth"
	"github.com/LaurelEdison/clashcoder/backend/handlers/problem"
	users "github.com/LaurelEdison/clashcoder/backend/handlers/user"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(router chi.Router, h *handlers.Handlers) {
	router.Get("/healthz", h.HandlerHealth)
	router.Post("/users", users.SignUp(h))
	router.Post("/login", auth.Login(h))
	router.With(auth.JWTAuthMiddleWare).Get("/me", users.FetchProfileSelf(h))
	router.Get("/problems", problem.GetAllProblems(h))
	router.Post("/problems", problem.CreateProblem(h))
	router.Get("/problems/{id}", problem.GetProblemById(h))
	router.Get("/problems", problem.GetProblemByRandom(h))
}
