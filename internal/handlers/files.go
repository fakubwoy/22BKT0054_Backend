package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/fakubwoy/go-file-share/internal/auth"
	"github.com/fakubwoy/go-file-share/internal/config"
	"github.com/fakubwoy/go-file-share/internal/models"
	"github.com/fakubwoy/go-file-share/internal/storage"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

type FileResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Size      int64     `json:"size"`
	Type      string    `json:"type"`
	URL       string    `json:"url"`
	IsPublic  bool      `json:"is_public"`
	CreatedAt time.Time `json:"created_at"`
}

type UploadResponse struct {
	Message string       `json:"message"`
	File    FileResponse `json:"file"`
}

type ShareResponse struct {
	ShareURL string `json:"share_url"`
}

func UploadHandler(db *sql.DB, cfg *config.Config, storage storage.Storage, rdb *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID").(int)

		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
			return
		}

		file, fileHeader, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Failed to get file from form", http.StatusBadRequest)
			return
		}
		defer file.Close()

		resultChan := make(chan *models.File)
		errChan := make(chan error)

		go func() {
			fileURL, err := storage.UploadFile(fileHeader, userID)
			if err != nil {
				errChan <- err
				return
			}

			newFile := &models.File{
				UserID:     userID,
				Name:       fileHeader.Filename,
				Size:       fileHeader.Size,
				Type:       fileHeader.Header.Get("Content-Type"),
				S3URL:      fileURL,
				IsPublic:   false,
				ShareToken: "",
			}

			if err := newFile.Create(db); err != nil {
				errChan <- err
				return
			}

			resultChan <- newFile
		}()

		select {
		case newFile := <-resultChan:
			ctx := context.Background()
			cacheKey := fmt.Sprintf("user_files:%d", userID)
			rdb.Del(ctx, cacheKey)

			response := UploadResponse{
				Message: "File uploaded successfully",
				File: FileResponse{
					ID:        newFile.ID,
					Name:      newFile.Name,
					Size:      newFile.Size,
					Type:      newFile.Type,
					URL:       newFile.S3URL,
					IsPublic:  newFile.IsPublic,
					CreatedAt: newFile.CreatedAt,
				},
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)

		case err := <-errChan:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func ListFilesHandler(db *sql.DB, rdb *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID").(int)
		ctx := context.Background()
		cacheKey := fmt.Sprintf("user_files:%d", userID)

		cachedFiles, err := rdb.Get(ctx, cacheKey).Result()
		if err == nil {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(cachedFiles))
			return
		}

		files, err := models.GetFilesByUser(db, userID)
		if err != nil {
			http.Error(w, "Failed to get files", http.StatusInternalServerError)
			return
		}

		var response []FileResponse
		for _, f := range files {
			var url string
			if f.S3URL != "" {
				url = f.S3URL
			} else {
				url = f.LocalPath
			}

			response = append(response, FileResponse{
				ID:        f.ID,
				Name:      f.Name,
				Size:      f.Size,
				Type:      f.Type,
				URL:       url,
				IsPublic:  f.IsPublic,
				CreatedAt: f.CreatedAt,
			})
		}

		jsonResponse, err := json.Marshal(response)
		if err == nil {
			rdb.Set(ctx, cacheKey, jsonResponse, 5*time.Minute)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func SearchFilesHandler(db *sql.DB, rdb *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID").(int)
		query := r.URL.Query().Get("q")

		if query == "" {
			http.Error(w, "Search query is required", http.StatusBadRequest)
			return
		}

		files, err := models.SearchFiles(db, userID, query)
		if err != nil {
			http.Error(w, "Failed to search files", http.StatusInternalServerError)
			return
		}

		var response []FileResponse
		for _, f := range files {
			var url string
			if f.S3URL != "" {
				url = f.S3URL
			} else {
				url = f.LocalPath
			}

			response = append(response, FileResponse{
				ID:        f.ID,
				Name:      f.Name,
				Size:      f.Size,
				Type:      f.Type,
				URL:       url,
				IsPublic:  f.IsPublic,
				CreatedAt: f.CreatedAt,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func ShareFileHandler(db *sql.DB, cfg *config.Config, storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID").(int)
		vars := mux.Vars(r)
		fileID, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid file ID", http.StatusBadRequest)
			return
		}

		_, err = models.GetFileByID(db, fileID, userID)
		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}

		token := auth.GenerateRandomString(32)
		expiresAt := time.Now().Add(24 * time.Hour)

		if err := models.MakeFilePublic(db, fileID, userID, token, expiresAt); err != nil {
			http.Error(w, "Failed to share file", http.StatusInternalServerError)
			return
		}

		shareURL := fmt.Sprintf("%s/share/%s", cfg.ServerBaseURL, token)
		response := ShareResponse{ShareURL: shareURL}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func DeleteFileHandler(db *sql.DB, rdb *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID").(int)
		vars := mux.Vars(r)
		fileID, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid file ID", http.StatusBadRequest)
			return
		}

		if err := models.DeleteFile(db, fileID, userID); err != nil {
			http.Error(w, "Failed to delete file", http.StatusInternalServerError)
			return
		}

		ctx := context.Background()
		cacheKey := fmt.Sprintf("user_files:%d", userID)
		rdb.Del(ctx, cacheKey)

		w.WriteHeader(http.StatusNoContent)
	}
}
