package storage

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"
)

var bucketName string

func InitGCS() {
	bucketName = os.Getenv("GCS_BUCKET_NAME")
	if bucketName == "" {
		fmt.Println("Warning: GCS_BUCKET_NAME is not set, uploads will be skipped or mock URLs returned if not configured properly.")
	}
}

// UploadFile uploads a multipart file to Google Cloud Storage and returns its public URL
func UploadFile(file *multipart.FileHeader) (string, error) {
	if bucketName == "" {
		return "https://mock-storage.com/mock-image.png", nil // mock if not configured
	}

	ctx := context.Background()
	client, err := storage.NewClient(ctx) // uses GOOGLE_APPLICATION_CREDENTIALS automatically
	if err != nil {
		return "", fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Generate a unique file name using only UUID and original extension
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("events/%s%s", uuid.New().String(), ext)

	bucket := client.Bucket(bucketName)
	obj := bucket.Object(filename)

	writer := obj.NewWriter(ctx)
	// Optional: set content type from the file
	writer.ContentType = file.Header.Get("Content-Type")

	if _, err := io.Copy(writer, src); err != nil {
		return "", fmt.Errorf("io.Copy: %w", err)
	}

	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("writer.Close: %w", err)
	}

	// Make the object public if the bucket doesn't have uniform bucket-level access 
	// Or just return the URL format
	// Return the public URL
	publicURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, filename)
	return publicURL, nil
}

// UploadBytes uploads raw bytes to Google Cloud Storage and returns its public URL
func UploadBytes(path string, data []byte, contentType string) (string, error) {
	if bucketName == "" {
		return "https://mock-storage.com/" + path, nil
	}

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	bucket := client.Bucket(bucketName)
	obj := bucket.Object(path)

	writer := obj.NewWriter(ctx)
	writer.ContentType = contentType

	if _, err := writer.Write(data); err != nil {
		return "", fmt.Errorf("writer.Write: %w", err)
	}

	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("writer.Close: %w", err)
	}

	publicURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, path)
	return publicURL, nil
}
