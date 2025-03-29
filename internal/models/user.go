package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (u *User) Create(db *sql.DB) error {
	query := `INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id, created_at, updated_at`
	return db.QueryRow(query, u.Email, u.PasswordHash).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
}

func GetUserByEmail(db *sql.DB, email string) (*User, error) {
	u := &User{}
	query := `SELECT id, email, password_hash, created_at, updated_at FROM users WHERE email = $1`
	err := db.QueryRow(query, email).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}
