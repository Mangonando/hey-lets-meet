package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func ApplyMigrations(sqlDB *sql.DB, migrationsDir string) error {
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("read migrations dir: %w", err)
	}

	var files []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(name, ".sql") {
			files = append(files, name)
		}
	}
	sort.Strings(files)

	if err := ensureSchemaMigrations(sqlDB); err != nil {
		return err
	}

	applied, err := appliedMigrations(sqlDB)
	if err != nil {
		return err
	}

	for _, filename := range files {
		if applied[filename] {
			continue
		}

		full := filepath.Join(migrationsDir, filename)
		sqlContent, err := os.ReadFile(full)
		if err != nil {
			return fmt.Errorf("read %s: %w", filename, err)
		}

		tx, err := sqlDB.Begin()
		if err != nil {
			return err
		}

		if _, err := tx.Exec(string(sqlContent)); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("apply %s: %w", filename, err)
		}

		if _, err := tx.Exec(`INSERT INTO schema_migrations(filename) VALUES (?)`, filename); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("record migration %s: %w", filename, err)
		}

		if err := tx.Commit(); err != nil {
			return err
		}
	}
	return nil
}

func ensureSchemaMigrations(sqlDB *sql.DB) error {
	_, err := sqlDB.Exec(`
CREATE TABLE IF NOT EXISTS schema_migrations (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  filename TEXT NOT NULL UNIQUE,
  applied_at TEXT NOT NULL DEFAULT (datetime('now'))
);
`)
	return err
}

func appliedMigrations(sqlDB *sql.DB) (map[string]bool, error) {
	rows, err := sqlDB.Query(`SELECT filename FROM schema_migrations`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var filename string
		if err := rows.Scan(&filename); err != nil {
			return nil, err
		}
		applied[filename] = true
	}
	return applied, rows.Err()
}
