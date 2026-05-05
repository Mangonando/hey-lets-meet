package main

import (
	"log"
	"net/http"
	"time"

	"hey-lets-meet/internal/auth"
	"hey-lets-meet/internal/db"
	"hey-lets-meet/internal/httpapi"
	"hey-lets-meet/internal/meetpoints"
)

func main() {
	database, err := db.Open("hey-lets-meet.db")
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer func() { _ = database.SQL.Close() }()

	if err := db.ApplyMigrations(database.SQL, "migrations"); err != nil {
		log.Fatalf("apply migrations: %v", err)
	}

	authRepo := &auth.Repo{DB: database.SQL}
	authService := &auth.Service{
		Repo:           authRepo,
		SessionTTL:     7 * 24 * time.Hour,
		CookieName:     "session",
		CookieInsecure: true,
	}
	authHandlers := &auth.Handlers{Svc: authService}

	meetpointsHandler := http.HandlerFunc(meetpoints.Handler{
		Service: &meetpoints.Service{
			Geocoder: meetpoints.MockGeocoder{},
			Router:   meetpoints.MockRouter{},
		},
	}.Suggest)

	server := httpapi.New(httpapi.Dependencies{
		AuthHandlers:      authHandlers,
		AuthService:       authService,
		MeetpointsHandler: meetpointsHandler,
	})

	address := ":8080"
	log.Printf("API listening on %s", address)
	log.Fatal(http.ListenAndServe(address, server.Mux))
}
