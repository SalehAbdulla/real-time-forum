package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	realtimeforum "real-time-forum"
)

type errorResponse struct {
	Error string `json:"error"`
}

func (re *Repository) parseForm(w http.ResponseWriter, r *http.Request) bool {
	if err := r.ParseForm(); err != nil {
		re.HandleError(w, r, realtimeforum.ErrBadRequest)
		return false
	}
	return true
}

// HandleError logs the error with request context and sends an appropriate HTTP response.
func (re *Repository) HandleError(w http.ResponseWriter, r *http.Request, err error) {
	// Determine HTTP status code from the error type
	var statusCode int
	switch err {
	case realtimeforum.ErrBadRequest:
		statusCode = http.StatusBadRequest
	case realtimeforum.ErrUnauthorized:
		statusCode = http.StatusUnauthorized
	case realtimeforum.ErrForbidden:
		statusCode = http.StatusForbidden
	case realtimeforum.ErrNotFound:
		statusCode = http.StatusNotFound
	case realtimeforum.ErrMethodNotAllowed:
		statusCode = http.StatusMethodNotAllowed
	default:
		statusCode = http.StatusInternalServerError
	}

	// Log the error with request context
	re.App.Logger.LogAttrs(r.Context(), slog.LevelError, "request error",
		slog.String("error", err.Error()),
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
		slog.String("remote", r.RemoteAddr),
		slog.Int("status", statusCode),
	)

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorResponse{Error: err.Error()})
}
