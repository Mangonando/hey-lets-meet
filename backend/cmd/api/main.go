package main

import (
	"fmt"
	"log"
	"net/http"

	"hey-lets-meet/internal/db"
)

func main() {
	database, err := db.Open("hey-lets-meet.db")
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer database.SQL.Close()

	if err := db.ApplyMigrations(database.SQL, "migrations"); err != nil {
		log.Fatalf("apply migrations: %v", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprintln(w, `{"status":"ok"}`)
	})

	mux.HandleFunc("/api/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprintln(w, `{"message":"pong"}`)
	})

	addr := ":8080"
	log.Printf("API listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
