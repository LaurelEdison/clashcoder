package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/LaurelEdison/clashcoder/backend/handlers"
	"github.com/LaurelEdison/clashcoder/backend/internal/database"
	"github.com/LaurelEdison/clashcoder/backend/routes"
	"github.com/LaurelEdison/clashcoder/backend/utils"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func NewServer(zapLogger *zap.Logger) *http.Server {

	portstring := utils.GetPort(zapLogger)
	dbURL := utils.GetDBUrl(zapLogger)

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		zapLogger.Error("Failed to connect to database", zap.Error(err))
	}

	queries := database.New(conn)

	router := chi.NewRouter()
	subRouter := chi.NewRouter()
	utils.SetupCors(zapLogger, router)
	routes.SetupRoutes(subRouter, handlers.New(zapLogger, queries))

	router.Mount("/clashcoder", subRouter)

	zapLogger.Info("Starting server ", zap.String("port", portstring))
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

	srv := NewServer(zapLogger)

	err = srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		zapLogger.Error("Failed to start http server", zap.Error(err))
	}

	return err
}
