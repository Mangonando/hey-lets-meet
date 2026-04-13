package db

import (
	"database/sql"
	"time"

	_ "modernc.org/sqlite"
)

type DB struct {
	SQL *sql.DB
}

func Open(path string) (*DB, error) {
	dsn := "file:" + path + "?_pragma=foreign_keys(1)"

	sqlDB, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	if err := sqlDB.Ping(); err != nil {
		_ = sqlDB.Close()
		return nil, err
	}
	return &DB{SQL: sqlDB}, nil
}
