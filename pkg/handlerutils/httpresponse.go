package handlerutils

import (
	"encoding/json"
	"github.com/csotherden/prescription-parser/pkg/models"
	"go.uber.org/zap"
	"net/http"
)

// RespondWithHTML responds with an HTML payload
func RespondWithHTML(w http.ResponseWriter, logger *zap.Logger, code int, payload string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(code)
	_, err := w.Write([]byte(payload))
	if err != nil {
		logger.Error("failed to write html response", zap.Error(err))
	}
}

// RespondWithJSON responds with a JSON payload
func RespondWithJSON(w http.ResponseWriter, logger *zap.Logger, statusCode int, payload any) {
	response, err := json.Marshal(payload)
	if err != nil {
		logger.Error("failed to marshal json response", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err = w.Write(response)
	if err != nil {
		logger.Error("failed to write json response", zap.Error(err))
	}
}

// RespondWithError responds with an error message
func RespondWithError(w http.ResponseWriter, logger *zap.Logger, statusCode int, message string, err error) {
	response := models.ErrorResponse{
		Message: message,
	}

	if err != nil {
		response.Error = err.Error()
	}

	RespondWithJSON(w, logger, statusCode, response)
}

// RespondWithNoContent responds with an HTTP 204 no content
func RespondWithNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}
