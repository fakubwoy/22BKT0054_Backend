package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/fakubwoy/go-file-share/internal/config"
)

type LocalStorage struct {
	baseDir string
	baseURL string
}

func NewLocalStorage(cfg *config.Config) (*LocalStorage, error) {
	if err := os.MkdirAll(cfg.LocalStorageDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	return &LocalStorage{
		baseDir: cfg.LocalStorageDir,
		baseURL: "http://localhost:" + cfg.ServerPort + "/uploads",
	}, nil
}

func (l *LocalStorage) UploadFile(fileHeader *multipart.FileHeader, userID int) (string, error) {
	userDir := filepath.Join(l.baseDir, fmt.Sprintf("%d", userID))
	if err := os.MkdirAll(userDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create user directory: %w", err)
	}

	ext := filepath.Ext(fileHeader.Filename)
	fileName := fmt.Sprintf("%s-%d%s", fileHeader.Filename[:len(fileHeader.Filename)-len(ext)], time.Now().Unix(), ext)
	filePath := filepath.Join(userDir, fileName)

	src, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("failed to copy file: %w", err)
	}

	fileURL := fmt.Sprintf("%s/%d/%s", l.baseURL, userID, fileName)
	return fileURL, nil
}

func (l *LocalStorage) GeneratePresignedURL(key string, expires time.Duration) (string, error) {
	return url.JoinPath(l.baseURL, key)
}
