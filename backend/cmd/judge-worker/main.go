package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/LaurelEdison/clashcoder/backend/handlers/auth"
	"github.com/LaurelEdison/clashcoder/backend/internal/database"
	"github.com/LaurelEdison/clashcoder/backend/utils"
	"go.uber.org/zap"
)

func main() {
	zapLogger, err := zap.NewProduction()
	if err != nil {
		panic(fmt.Sprintf("Error creating zap logger: %v", err))
	}
	dbURL := utils.GetDBUrl(zapLogger)

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("JWT_SECRET is not set")
	}
	auth.InitJWT(secret)

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		zapLogger.Error("Failed to connect to database", zap.Error(err))
	}

	if err := conn.Ping(); err != nil {
		zapLogger.Error("Failed to connect to database", zap.Error(err))
	}

	queries := database.New(conn)

	zapLogger.Info("Judge worker started")

}
