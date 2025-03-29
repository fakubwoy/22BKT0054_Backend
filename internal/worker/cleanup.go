package worker

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/fakubwoy/go-file-share/internal/storage"
)

type CleanupWorker struct {
	db       *sql.DB
	storage  storage.Storage
	interval time.Duration
}

func NewCleanupWorker(db *sql.DB, storage storage.Storage, interval time.Duration) *CleanupWorker {
	return &CleanupWorker{
		db:       db,
		storage:  storage,
		interval: interval,
	}
}

func (w *CleanupWorker) Start() {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for range ticker.C {
		w.cleanupExpiredFiles()
	}
}

func (w *CleanupWorker) cleanupExpiredFiles() {
	ctx := context.Background()

	// Get expired files
	rows, err := w.db.QueryContext(ctx,
		"SELECT id, s3_url, local_path FROM files WHERE expires_at IS NOT NULL AND expires_at < NOW()")
	if err != nil {
		log.Printf("Failed to query expired files: %v", err)
		return
	}
	defer rows.Close()

	var files []struct {
		ID        int
		S3URL     string
		LocalPath string
	}

	for rows.Next() {
		var f struct {
			ID        int
			S3URL     string
			LocalPath string
		}
		if err := rows.Scan(&f.ID, &f.S3URL, &f.LocalPath); err != nil {
			log.Printf("Failed to scan file: %v", err)
			continue
		}
		files = append(files, f)
	}

	// Delete files from storage and database
	for _, f := range files {
		// Delete from storage (implementation depends on storage type)
		if f.S3URL != "" {
			// Delete from S3 (implementation needed based on your storage)
		} else {
			// Delete local file
		}

		// Delete from database
		if _, err := w.db.ExecContext(ctx, "DELETE FROM files WHERE id = $1", f.ID); err != nil {
			log.Printf("Failed to delete file %d: %v", f.ID, err)
			continue
		}

		log.Printf("Deleted expired file %d", f.ID)
	}
}
