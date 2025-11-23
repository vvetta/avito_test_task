package http

import (
	"net/http"
	"encoding/json"

	"avito_test_task/internal/adapters/http/openapi"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	resp := openapi.ErrorResponse{}

	resp.Error.Code = openapi.ErrorResponseErrorCode(code)
	resp.Error.Message = message

	writeJSON(w, status, resp)
}
