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
	// Determine HTTP status code and log level from the error type
	var statusCode int
	var level slog.Level

	switch err {
	case realtimeforum.ErrBadRequest:
		statusCode = http.StatusBadRequest
		level = slog.LevelWarn
	case realtimeforum.ErrUnauthorized:
		statusCode = http.StatusUnauthorized
		level = slog.LevelWarn
	case realtimeforum.ErrForbidden:
		statusCode = http.StatusForbidden
		level = slog.LevelWarn
	case realtimeforum.ErrNotFound:
		statusCode = http.StatusNotFound
		level = slog.LevelWarn
	case realtimeforum.ErrMethodNotAllowed:
		statusCode = http.StatusMethodNotAllowed
		level = slog.LevelWarn
	case realtimeforum.ErrInternal:
		statusCode = http.StatusInternalServerError
		level = slog.LevelError
	default:
		// Check if the error is a known validation/business-logic error (4xx)
		switch {
		case err == realtimeforum.ErrInvalidEmail,
			err == realtimeforum.ErrEmailExists,
			err == realtimeforum.ErrNickName,
			err == realtimeforum.ErrNickNameLength,
			err == realtimeforum.ErrPasswordLength,
			err == realtimeforum.ErrPasswordsDontMatch,
			err == realtimeforum.ErrInvalidPassForm,
			err == realtimeforum.ErrInvalidAge,
			err == realtimeforum.ErrGender,
			err == realtimeforum.ErrInvalidCredentials:
			statusCode = http.StatusBadRequest
			level = slog.LevelWarn
		default:
			statusCode = http.StatusInternalServerError
			level = slog.LevelError
		}
	}

	// Log the error with request context
	re.App.Logger.LogAttrs(r.Context(), level, "request error",
		slog.String("error", err.Error()),
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
		slog.String("remote", r.RemoteAddr),
		slog.Int("status", statusCode),
	)

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorResponse{Error: err.Error()})
}
