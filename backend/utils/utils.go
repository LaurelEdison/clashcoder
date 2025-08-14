package utils

import (
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"go.uber.org/zap"
)

func GetPort(zapLogger *zap.Logger) string {
	portstring := os.Getenv("PORT")
	if portstring == "" {
		zapLogger.Warn("PORT not set, defaulting to 8080")
		return "8080"
	}
	zapLogger.Info("PORT loaded from environment", zap.String("port", portstring))
	return portstring
}

func SetupCors(zapLogger *zap.Logger, router chi.Router) {

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
