package main

import (
	"bytes"
	"encoding/json"
	"hey-lets-meet/internal/auth"
	"hey-lets-meet/internal/db"
	"hey-lets-meet/internal/httpapi"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestRegisterLoginAndProtectRoute(t *testing.T) {
	tmpDir := t.TempDir()
	dbFile := filepath.Join(tmpDir, "test.db")

	database, err := db.Open(dbFile)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer func() { _ = database.SQL.Close() }()

	// migrations need a real folder path. ergo copy migrations into temp
	if err := copyMigrations(tmpDir); err != nil {
		t.Fatalf("copy migrations: %v", err)
	}
	migrationsDir := filepath.Join(tmpDir, "migrations")

	if err := db.ApplyMigrations(database.SQL, migrationsDir); err != nil {
		t.Fatalf("apply migrations: %v", err)
	}

	authRepo := &auth.Repo{DB: database.SQL}
	authService := &auth.Service{
		Repo:           authRepo,
		SessionTTL:     24 * time.Hour,
		CookieName:     "html_session",
		CookieInsecure: true,
	}
	authHandlers := &auth.Handlers{Svc: authService}

	server := httpapi.New(httpapi.Dependencies{
		AuthHandlers: authHandlers,
		AuthService:  authService,
	})

	testServer := httptest.NewServer(server.Mux)
	defer testServer.Close()

	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	// register auto log in
	{
		body := map[string]string{"email": "test@example.com", "password": "pwd123456"}
		response := postJSON(t, client, testServer.URL+"/auth/register", body)
		if response.StatusCode != http.StatusCreated {
			t.Fatalf("register status = %d, want %d", response.StatusCode, http.StatusCreated)
		}
		_ = response.Body.Close()
	}

	// work protected
	{
		response, err := client.Get(testServer.URL + "/api/protected")
		if err != nil {
			t.Fatalf("get protected: %v", err)
		}
		if response.StatusCode != http.StatusOK {
			t.Fatalf("protected status = %d, want %d", response.StatusCode, http.StatusOK)
		}
		_ = response.Body.Close()
	}

	// logout
	{
		response := postJSON(t, client, testServer.URL+"/auth/logout", map[string]string{})
		if response.StatusCode != http.StatusOK {
			t.Fatalf("logout status = %d, want %d", response.StatusCode, http.StatusOK)
		}
		_ = response.Body.Close()
	}

	// now fail protected
	{
		response, err := client.Get(testServer.URL + "/api/protected")
		if err != nil {
			t.Fatalf("get protected after logout: %v", err)
		}
		if response.StatusCode != http.StatusUnauthorized {
			t.Fatalf("protected after logout status = %d, want %d", response.StatusCode, http.StatusUnauthorized)
		}
		_ = response.Body.Close()
	}

	// log in
	{
		response := postJSON(t, client, testServer.URL+"/auth/login", map[string]string{"email": "test@example.com", "password": "pwd123456"})
		if response.StatusCode != http.StatusOK {
			t.Fatalf("login status = %d. want %d", response.StatusCode, http.StatusOK)
		}
		_ = response.Body.Close()
	}

	// work protected again
	{
		response, err := client.Get(testServer.URL + "/api/protected")
		if err != nil {
			t.Fatalf("get protected after login: %v", err)
		}
		if response.StatusCode != http.StatusOK {
			t.Fatalf("protected after login status = %d, want %d", response.StatusCode, http.StatusOK)
		}
		_ = response.Body.Close()
	}
}

func copyMigrations(tmpDir string) error {
	source := filepath.Join("..", "..", "migrations")
	destination := filepath.Join(tmpDir, "migrations")

	if err := os.MkdirAll(destination, 0o755); err != nil {
		return err
	}

	entries, err := os.ReadDir(source)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		content, err := os.ReadFile(filepath.Join(source, entry.Name()))
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(destination, entry.Name()), content, 0o644); err != nil {
			return err
		}
	}
	return nil
}

func postJSON(t *testing.T, client *http.Client, url string, body any) *http.Response {
	t.Helper()

	content, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(content))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	response, err := client.Do(req)
	if err != nil {
		t.Fatalf("do request: %v", err)
	}
	return response
}
