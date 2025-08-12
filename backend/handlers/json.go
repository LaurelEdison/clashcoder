package handlers

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

func (h *Handlers) RespondWithJSON(w http.ResponseWriter,
	code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		h.zapLogger.Error("Error encoding json", zap.Error(err))
	}

}

func (h *Handlers) RespondWithError(w http.ResponseWriter,
	code int, msg string,
	zapLogger *zap.Logger) {
	if code > 499 {
		zapLogger.Error("Responding with server error",
			zap.Int("code", code), zap.String("message", msg))
	} else {
		zapLogger.Warn("Responding with client error",
			zap.Int("code", code), zap.String("message", msg))
	}

	type errResponse struct {
		Error string `json:"error"`
	}

	h.RespondWithJSON(w, code, errResponse{Error: msg})

}
