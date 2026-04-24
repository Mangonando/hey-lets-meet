package auth

import (
	"context"
	"net/http"
	"time"
)

type ctxKey string

const userIDKey ctxKey = "userID"

func UserIDFromContext(ctx context.Context) int64 {
	value := ctx.Value(userIDKey)
	if value == nil {
		return 0
	}
	userID, _ := value.(int64)
	return userID
}

func (s *Service) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(s.CookieName)
		if err != nil || cookie.Value == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		userID, err := s.Repo.GetSessionUser(cookie.Value, time.Now())
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
