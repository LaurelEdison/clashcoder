package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/LaurelEdison/clashcoder/backend/handlers"
	"github.com/LaurelEdison/clashcoder/backend/routes"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func getPort(zapLogger *zap.Logger) string {
	portstring := os.Getenv("PORT")
	if portstring == "" {
		zapLogger.Warn("PORT not set, defaulting to 8080")
		return "8080"
	}
	zapLogger.Info("PORT loaded from environment", zap.String("port", portstring))
	return portstring
}

func setupCors(zapLogger *zap.Logger, router chi.Router) {

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"}, // Allow all origins for dev
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"}, // Accept all headers
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	zapLogger.Info("CORS middleware configuered")
}

func main() {

	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Failed to load .env : %v", err)
	}

	zapLogger, err := zap.NewProduction()
	fmt.Fprintf(os.Stderr, "Error starting logger: %v", err)

	defer func() {
		if err := zapLogger.Sync(); err != nil {
			if !strings.Contains(err.Error(), "already closed") {
				fmt.Fprintf(os.Stderr, "Error syncing logger: %v", err)
			}
		}
	}()

	portstring := getPort(zapLogger)
	zapLogger.Info("Starting server ", zap.String("port", portstring))

	router := chi.NewRouter()
	setupCors(zapLogger, router)

	routes.SetupRoutes(router, handlers.New(zapLogger))

	router.Mount("/clashcoder", router)

	srv := http.Server{
		Handler: router,
		Addr:    ":" + portstring,
	}

	zapLogger.Info("Server starting on port", zap.String("portstring", portstring))

	if err := srv.ListenAndServe(); err != nil {
		zapLogger.Error("Failed to start http server", zap.Error(err))
	}

	zapLogger.Info("Port", zap.String("portstring", portstring))

}
