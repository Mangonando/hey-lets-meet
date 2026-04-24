package auth

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type Handlers struct {
	Svc *Service
}

type creds struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handlers) Register(w http.ResponseWriter, r *http.Request) {
	var credentials creds
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}
	credentials.Email = strings.TrimSpace(strings.ToLower(credentials.Email))
	if credentials.Email == "" || len(credentials.Password) < 8 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "email required and password must be at least 8 characters"})
		return
	}

	hash, err := HashPassword(credentials.Password)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to hash password"})
		return
	}

	userID, err := h.Svc.Repo.CreateUser(credentials.Email, hash)
	if err != nil {
		if err == ErrEmailTaken {
			writeJSON(w, http.StatusConflict, map[string]string{"error": "email is already registered"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to create user"})
		return
	}
	if err := h.issueSession(w, userID); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to create session"})
		return
	}

	user := &User{ID: userID, Email: credentials.Email}
	writeJSON(w, http.StatusCreated, user)
}

func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	var credentials creds
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}
	credentials.Email = strings.TrimSpace(strings.ToLower(credentials.Email))
	if credentials.Email == "" || credentials.Password == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "email and password required"})
		return
	}

	userID, hash, err := h.Svc.Repo.GetUserByEmail(credentials.Email)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
		return
	}

	if err := CheckPassword(hash, credentials.Password); err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
		return
	}

	if err := h.issueSession(w, userID); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to create session"})
		return
	}

	user, err := h.Svc.Repo.GetUserByID(userID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to load user"})
		return
	}
	writeJSON(w, http.StatusOK, user)
}

func (h *Handlers) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(h.Svc.CookieName)
	if err == nil && cookie.Value != "" {
		_ = h.Svc.Repo.DeleteSession(cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     h.Svc.CookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   !h.Svc.CookieInsecure,
	})

	writeJSON(w, http.StatusOK, map[string]string{"message": "logged out"})
}

func (h *Handlers) Me(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromContext(r.Context())
	if userID == 0 {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	user, err := h.Svc.Repo.GetUserByID(userID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to load user"})
		return
	}
	writeJSON(w, http.StatusOK, user)
}

func (h *Handlers) issueSession(w http.ResponseWriter, userID int64) error {
	sessionID, err := NewSessionID()
	if err != nil {
		return err
	}

	expires := time.Now().Add(h.Svc.SessionTTL)
	if err := h.Svc.Repo.CreateSession(sessionID, userID, expires); err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     h.Svc.CookieName,
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   !h.Svc.CookieInsecure,
		Expires:  expires,
	})
	return nil
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
