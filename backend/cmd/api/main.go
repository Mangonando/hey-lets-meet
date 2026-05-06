package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"hey-lets-meet/internal/auth"
	"hey-lets-meet/internal/db"
	"hey-lets-meet/internal/httpapi"
	"hey-lets-meet/internal/meetpoints"
)

func main() {
	_ = os.MkdirAll("./data", 0o755)

	dbPath := filepath.Join("data", "app.db")
	database, err := db.Open(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer database.SQL.Close()

	if err := db.ApplyMigrations(database.SQL, "./migrations"); err != nil {
		log.Fatal(err)
	}

	authRepo := &auth.Repo{DB: database.SQL}
	authSvc := &auth.Service{
		Repo:           authRepo,
		SessionTTL:     7 * 24 * time.Hour,
		CookieName:     "hlm_session",
		CookieInsecure: true,
	}
	authHandlers := &auth.Handlers{Svc: authSvc}

	meetpointsCache := meetpoints.CacheRepo{DB: database.SQL}

	orsKey := os.Getenv("ORS_API_KEY")
	orsBaseURL := os.Getenv("ORS_BASE_URL")
	if orsBaseURL == "" {
		orsBaseURL = "https://api.openrouteservice.org"
	}

	httpClient := &http.Client{Timeout: 10 * time.Second}

	var meetpointsService *meetpoints.Service
	if orsKey != "" {
		orsConfig := meetpoints.ORSConfig{
			BaseURL: orsBaseURL,
			APIKey:  orsKey,
			Timeout: 10 * time.Second,
		}
		meetpointsService = &meetpoints.Service{
			Geocoder: meetpoints.ORSGeocoder{HTTP: httpClient, Config: orsConfig, Cache: meetpointsCache},
			Router:   meetpoints.ORSRouter{HTTP: httpClient, Config: orsConfig, Cache: meetpointsCache},
		}
		log.Println("Meetpoints: using OpenRouteService providers")
	} else {
		meetpointsService = &meetpoints.Service{
			Geocoder: meetpoints.MockGeocoder{},
			Router:   meetpoints.MockRouter{WalkingSpeedMps: 1.4},
		}
		log.Println("Meetpoints: ORS_API_KEY not set, using mock providers")
	}

	meetpointsHandler := http.HandlerFunc(meetpoints.Handler{Service: meetpointsService}.Suggest)

	server := httpapi.New(httpapi.Dependencies{
		AuthHandlers:      authHandlers,
		AuthService:       authSvc,
		MeetpointsHandler: meetpointsHandler,
	})

	address := ":8080"
	log.Printf("API listening on %s", address)
	log.Fatal(http.ListenAndServe(address, server.Mux))
}
