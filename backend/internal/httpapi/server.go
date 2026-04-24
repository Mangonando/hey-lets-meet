package httpapi

import (
	"hey-lets-meet/internal/auth"
	"net/http"
)

type Server struct {
	Mux http.Handler
}

type Dependencies struct {
	AuthHandlers *auth.Handlers
	AuthService  *auth.Service
}

func New(dependencies Dependencies) *Server {
	mux := http.NewServeMux()

	// public
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	mux.HandleFunc("/auth/register", dependencies.AuthHandlers.Register)
	mux.HandleFunc("/auth/login", dependencies.AuthHandlers.Login)
	mux.HandleFunc("/auth/logout", dependencies.AuthHandlers.Logout)

	// protected
	mux.Handle("/auth/me", dependencies.AuthService.RequireAuth(http.HandlerFunc(dependencies.AuthHandlers.Me)))

	mux.Handle("/api/protected", dependencies.AuthService.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"message":"welcome"}`))
	})))
	return &Server{Mux: mux}
}
