package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/LaurelEdison/clashcoder/backend/handlers"
	"github.com/LaurelEdison/clashcoder/backend/routes"
	"github.com/LaurelEdison/clashcoder/backend/utils"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func NewServer(zapLogger *zap.Logger) *http.Server {
	router := chi.NewRouter()
	subRouter := chi.NewRouter()
	utils.SetupCors(zapLogger, router)
	routes.SetupRoutes(subRouter, handlers.New(zapLogger))

	router.Mount("/clashcoder", subRouter)

	portstring := utils.GetPort(zapLogger)

	return &http.Server{
		Handler: router,
		Addr:    ":" + portstring,
	}

}

func Run() error {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Failed to load .env : %v", err)
	}

	zapLogger, err := zap.NewProduction()
	if err != nil {
		return fmt.Errorf("error starting logger: %v", err)
	}
	defer func() {
		if err := zapLogger.Sync(); err != nil && !strings.Contains(err.Error(), "already closed") {
			fmt.Fprintf(os.Stderr, "error syncing logger: %v", err)
		}
	}()

	zapLogger.Info("Starting server ", zap.String("port", portstring))

	srv := NewServer(zapLogger)
	zapLogger.Info("Server starting on port", zap.String("portstring", portstring))

	err = srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		zapLogger.Error("Failed to start http server", zap.Error(err))
	}

	zapLogger.Info("Port", zap.String("portstring", portstring))

	return err
}
