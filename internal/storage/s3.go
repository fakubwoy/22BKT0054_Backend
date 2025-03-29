package storage

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/fakubwoy/go-file-share/internal/config"
)

type S3Storage struct {
	uploader   *s3manager.Uploader
	downloader *s3manager.Downloader
	bucket     string
	region     string
}

func NewS3Storage(cfg *config.Config) (*S3Storage, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(cfg.S3Region),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %w", err)
	}

	return &S3Storage{
		uploader:   s3manager.NewUploader(sess),
		downloader: s3manager.NewDownloader(sess),
		bucket:     cfg.S3Bucket,
		region:     cfg.S3Region,
	}, nil
}

func (s *S3Storage) UploadFile(fileHeader *multipart.FileHeader, userID int) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	ext := filepath.Ext(fileHeader.Filename)
	key := fmt.Sprintf("%d/%s-%d%s", userID, fileHeader.Filename[:len(fileHeader.Filename)-len(ext)], time.Now().Unix(), ext)

	_, err = s.uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
		Body:   file,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	fileURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucket, s.region, key)
	return fileURL, nil
}

func (s *S3Storage) GeneratePresignedURL(key string, expires time.Duration) (string, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(s.region),
	})
	if err != nil {
		return "", fmt.Errorf("failed to create AWS session: %w", err)
	}

	svc := s3.New(sess)
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})

	urlStr, err := req.Presign(expires)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return urlStr, nil
}
