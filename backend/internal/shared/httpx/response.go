package httpx

import (
	"encoding/json"
	"errors"
	"net/http"

	apperrors "github.com/MaiconGambini/erpGolang/backend/internal/shared/errors"
)

type ErrorBody struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}

type PaginationMeta struct {
	Limit  int   `json:"limit"`
	Offset int   `json:"offset"`
	Total  int64 `json:"total"`
}

func JSON(w http.ResponseWriter, status int, data any) {
	writeJSON(w, status, map[string]any{"data": data})
}

func Error(w http.ResponseWriter, code string, message string, status int) {
	writeJSON(w, status, map[string]any{
		"error": ErrorBody{Code: code, Message: message},
	})
}

func ErrorWithDetails(w http.ResponseWriter, code, message string, status int, details map[string]string) {
	writeJSON(w, status, map[string]any{
		"error": ErrorBody{Code: code, Message: message, Details: details},
	})
}

func Paginated(w http.ResponseWriter, data any, meta PaginationMeta) {
	writeJSON(w, http.StatusOK, map[string]any{
		"data":       data,
		"pagination": meta,
	})
}

func WriteError(w http.ResponseWriter, err error) {
	var appErr apperrors.AppError
	if errors.As(err, &appErr) {
		status := statusFromKind(appErr.Kind)
		ErrorWithDetails(w, appErr.Code, appErr.Message, status, appErr.Details)
		return
	}
	Error(w, "INTERNAL_ERROR", "internal server error", http.StatusInternalServerError)
}

func statusFromKind(kind apperrors.Kind) int {
	switch kind {
	case apperrors.Invalid:
		return http.StatusUnprocessableEntity
	case apperrors.NotFound:
		return http.StatusNotFound
	case apperrors.Conflict:
		return http.StatusConflict
	case apperrors.Forbidden:
		return http.StatusForbidden
	case apperrors.Unauthorized:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
