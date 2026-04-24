package auth

import (
	"database/sql"
	"errors"
	"strings"
	"time"
)

var ErrEmailTaken = errors.New("email is taken already")
var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrNoSession = errors.New("no session")

type Repo struct {
	DB *sql.DB
}

func (r *Repo) CreateUser(email string, passwordHash string) (int64, error) {
	result, err := r.DB.Exec(`INSERT INTO users(email, password_hash) VALUES (?, ?)`, email, passwordHash)
	if err != nil {
		if isUniqueConstraintErr(err) {
			return 0, ErrEmailTaken
		}
		return 0, err
	}
	return result.LastInsertId()
}

func (r *Repo) GetUserByEmail(email string) (id int64, passwordHash string, err error) {
	row := r.DB.QueryRow(`SELECT id, password_hash FROM users WHERE email = ?`, email)
	if err := row.Scan(&id, &passwordHash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, "", ErrInvalidCredentials
		}
		return 0, "", err
	}
	return id, passwordHash, nil
}

func (r *Repo) GetUserByID(id int64) (*User, error) {
	row := r.DB.QueryRow(`SELECT id, email FROM users WHERE id = ?`, id)
	var user User
	if err := row.Scan(&user.ID, &user.Email); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repo) CreateSession(sessionID string, userID int64, expiresAt time.Time) error {
	_, err := r.DB.Exec(`INSERT INTO sessions(id, user_id, expires_at) VALUES (?, ?, ?)`, sessionID, userID, expiresAt.UTC().Format(time.RFC3339))
	return err
}

func (r *Repo) DeleteSession(sessionID string) error {
	_, err := r.DB.Exec(`DELETE FROM sessions WHERE id = ?`, sessionID)
	return err
}

func (r *Repo) GetSessionUser(sessionID string, now time.Time) (int64, error) {
	row := r.DB.QueryRow(`SELECT user_id FROM sessions WHERE id = ? AND expires_at > ?`, sessionID, now.UTC().Format(time.RFC3339))
	var userID int64
	if err := row.Scan(&userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrNoSession
		}
		return 0, err
	}
	return userID, nil
}

func isUniqueConstraintErr(err error) bool {
	return err != nil && (strings.Contains(err.Error(), "UNIQUE constraint failed") || strings.Contains(err.Error(), "constraint failed"))
}
