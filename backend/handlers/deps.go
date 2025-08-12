package handlers

import "go.uber.org/zap"

type Handlers struct {
	zapLogger *zap.Logger
}

func New(zapLogger *zap.Logger) *Handlers {
	return &Handlers{
		zapLogger: zapLogger,
	}
}
