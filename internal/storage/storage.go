package storage

import (
	"mime/multipart"
	"time"
)

type Storage interface {
	UploadFile(fileHeader *multipart.FileHeader, userID int) (string, error)
	GeneratePresignedURL(key string, expires time.Duration) (string, error)
}
