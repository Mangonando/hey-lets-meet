package meetpoints

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type Handler struct {
	Service *Service
}

func (h Handler) Suggest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "use POST", nil)
		return
	}

	var req MeetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "BAD_JSON", "invalid json", nil)
		return
	}

	req.OriginA = strings.TrimSpace(req.OriginA)
	req.OriginB = strings.TrimSpace(req.OriginB)

	validationErrors := map[string]string{}
	if req.OriginA == "" {
		validationErrors["originA"] = "required"
	}
	if req.OriginB == "" {
		validationErrors["originB"] = "required"
	}
	if len(validationErrors) > 0 {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "invalid input", validationErrors)
		return
	}

	resp, err := h.Service.Suggest(req, time.Now(), DefaultOptions())
	if err != nil {
		writeError(w, http.StatusBadRequest, "SUGGEST_FAILED", err.Error(), nil)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

type errorBody struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, code, message string, fields map[string]string) {
	writeJSON(w, status, struct {
		Error errorBody `json:"error"`
	}{Error: errorBody{Code: code, Message: message, Fields: fields}})
}
