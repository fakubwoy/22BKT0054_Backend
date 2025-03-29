package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/fakubwoy/go-file-share/internal/config"
	_ "github.com/lib/pq"
)

func NewPostgresDB(cfg *config.Config) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	var db *sql.DB
	var err error

	maxAttempts := 5
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		db, err = sql.Open("postgres", connStr)
		if err == nil {
			err = db.Ping()
			if err == nil {
				break
			}
		}

		if attempt < maxAttempts {
			time.Sleep(time.Duration(attempt) * time.Second)
			log.Printf("Retrying database connection (attempt %d/%d)", attempt, maxAttempts)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", maxAttempts, err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	log.Println("Successfully connected to PostgreSQL database")
	return db, nil
}
