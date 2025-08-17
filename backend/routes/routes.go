package routes

import (
	"github.com/LaurelEdison/clashcoder/backend/handlers"
	"github.com/LaurelEdison/clashcoder/backend/handlers/auth"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(router chi.Router, h *handlers.Handlers) {
	router.Get("/healthz", h.HandlerHealth)
	router.Post("/users", h.SignUp)
	router.With(auth.JWTAuthMiddleWare).Get("/me", h.FetchProfileSelf)
}
