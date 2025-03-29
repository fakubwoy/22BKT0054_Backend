package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/fakubwoy/go-file-share/internal/models"
	"github.com/fakubwoy/go-file-share/internal/storage"
	"github.com/gorilla/mux"
)

func GetSharedFileHandler(db *sql.DB, storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		token := vars["token"]

		file, err := models.GetFileByShareToken(db, token)
		if err != nil {
			http.Error(w, "File not found or expired", http.StatusNotFound)
			return
		}

		var fileURL string
		if file.S3URL != "" {
			presignedURL, err := storage.GeneratePresignedURL(file.S3URL, 15*time.Minute)
			if err != nil {
				http.Error(w, "Failed to generate file URL", http.StatusInternalServerError)
				return
			}
			fileURL = presignedURL
		} else {
			fileURL = file.LocalPath
		}

		http.Redirect(w, r, fileURL, http.StatusFound)
	}
}
