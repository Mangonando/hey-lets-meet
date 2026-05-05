package httpapi

import (
	"net/http"

	"hey-lets-meet/internal/auth"
)

type Server struct {
	Mux http.Handler
}

type Dependencies struct {
	AuthHandlers      *auth.Handlers
	AuthService       *auth.Service
	MeetpointsHandler http.Handler
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func New(deps Dependencies) *Server {
	mux := http.NewServeMux()

	// Public
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	mux.HandleFunc("/auth/register", deps.AuthHandlers.Register)
	mux.HandleFunc("/auth/login", deps.AuthHandlers.Login)
	mux.HandleFunc("/auth/logout", deps.AuthHandlers.Logout)

	// Protected
	mux.Handle("/auth/me", deps.AuthService.RequireAuth(http.HandlerFunc(deps.AuthHandlers.Me)))

	mux.Handle("/api/protected", deps.AuthService.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"message":"welcome"}`))
	})))

	mux.Handle("/api/meetpoints/suggest",
		deps.AuthService.RequireAuth(deps.MeetpointsHandler),
	)

	return &Server{Mux: corsMiddleware(mux)}
}
