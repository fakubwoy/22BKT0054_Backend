package models

import (
	"database/sql"
	"time"
)

type File struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	Name       string    `json:"name"`
	Size       int64     `json:"size"`
	Type       string    `json:"type"`
	S3URL      string    `json:"s3_url,omitempty"`
	LocalPath  string    `json:"local_path,omitempty"`
	IsPublic   bool      `json:"is_public"`
	ShareToken string    `json:"share_token,omitempty"`
	ExpiresAt  time.Time `json:"expires_at,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (f *File) Create(db *sql.DB) error {
	query := `INSERT INTO files (user_id, name, size, type, s3_url, local_path, is_public, share_token, expires_at) 
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) 
              RETURNING id, created_at, updated_at`
	return db.QueryRow(query, f.UserID, f.Name, f.Size, f.Type, f.S3URL, f.LocalPath,
		f.IsPublic, f.ShareToken, f.ExpiresAt).Scan(&f.ID, &f.CreatedAt, &f.UpdatedAt)
}

func GetFileByID(db *sql.DB, fileID, userID int) (*File, error) {
	f := &File{}
	query := `SELECT id, user_id, name, size, type, s3_url, local_path, is_public, 
              share_token, expires_at, created_at, updated_at 
              FROM files WHERE id = $1 AND user_id = $2`
	err := db.QueryRow(query, fileID, userID).Scan(
		&f.ID, &f.UserID, &f.Name, &f.Size, &f.Type, &f.S3URL, &f.LocalPath,
		&f.IsPublic, &f.ShareToken, &f.ExpiresAt, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func GetFilesByUser(db *sql.DB, userID int) ([]*File, error) {
	query := `SELECT id, user_id, name, size, type, s3_url, local_path, is_public, 
              share_token, expires_at, created_at, updated_at 
              FROM files WHERE user_id = $1`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []*File
	for rows.Next() {
		f := &File{}
		err := rows.Scan(
			&f.ID, &f.UserID, &f.Name, &f.Size, &f.Type, &f.S3URL, &f.LocalPath,
			&f.IsPublic, &f.ShareToken, &f.ExpiresAt, &f.CreatedAt, &f.UpdatedAt)
		if err != nil {
			return nil, err
		}
		files = append(files, f)
	}
	return files, nil
}

func SearchFiles(db *sql.DB, userID int, query string) ([]*File, error) {
	sqlQuery := `SELECT id, user_id, name, size, type, s3_url, local_path, is_public, 
                share_token, expires_at, created_at, updated_at 
                FROM files WHERE user_id = $1 AND name LIKE '%' || $2 || '%'`
	rows, err := db.Query(sqlQuery, userID, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []*File
	for rows.Next() {
		f := &File{}
		err := rows.Scan(
			&f.ID, &f.UserID, &f.Name, &f.Size, &f.Type, &f.S3URL, &f.LocalPath,
			&f.IsPublic, &f.ShareToken, &f.ExpiresAt, &f.CreatedAt, &f.UpdatedAt)
		if err != nil {
			return nil, err
		}
		files = append(files, f)
	}
	return files, nil
}

func DeleteFile(db *sql.DB, fileID, userID int) error {
	query := `DELETE FROM files WHERE id = $1 AND user_id = $2`
	_, err := db.Exec(query, fileID, userID)
	return err
}

func MakeFilePublic(db *sql.DB, fileID, userID int, token string, expiresAt time.Time) error {
	query := `UPDATE files SET is_public = true, share_token = $1, expires_at = $2 
              WHERE id = $3 AND user_id = $4`
	_, err := db.Exec(query, token, expiresAt, fileID, userID)
	return err
}

func GetFileByShareToken(db *sql.DB, token string) (*File, error) {
	f := &File{}
	query := `SELECT id, user_id, name, size, type, s3_url, local_path, is_public, 
              share_token, expires_at, created_at, updated_at 
              FROM files WHERE share_token = $1 AND is_public = true AND (expires_at IS NULL OR expires_at > NOW())`
	err := db.QueryRow(query, token).Scan(
		&f.ID, &f.UserID, &f.Name, &f.Size, &f.Type, &f.S3URL, &f.LocalPath,
		&f.IsPublic, &f.ShareToken, &f.ExpiresAt, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return f, nil
}
