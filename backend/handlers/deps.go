package handlers

import (
	"github.com/LaurelEdison/clashcoder/backend/internal/database"
	"go.uber.org/zap"
)

type Handlers struct {
	zapLogger *zap.Logger
	DB        *database.Queries
}

func New(zapLogger *zap.Logger, DB *database.Queries) *Handlers {
	return &Handlers{
		zapLogger: zapLogger,
		DB:        DB,
	}
}
