package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

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

func main() {
	godotenv.Load(".env")
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()

	portstring := getPort(zapLogger)
	zapLogger.Info("Starting server ", zap.String("port", portstring))
}
