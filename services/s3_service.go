package services

import (
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"strings"

	repositories "github.com/Bualoi-s-Dev/backend/repositories/s3"
)

type S3Service struct {
	Repo *repositories.S3Repository
}

func NewS3Service(repo *repositories.S3Repository) *S3Service {
	return &S3Service{Repo: repo}
}

func (s *S3Service) UploadFile(file *multipart.FileHeader, key string) (string, error) {
	return s.Repo.UploadFile(file, key)
}

func (s *S3Service) UploadBase64(fileBytes []byte, key string) (string, error) {
	_, ext, err := DetectMimeType(string(fileBytes))
	if err != nil {
		return "", err
	}
	fmt.Println("ext :", ext)

	imageData, err := base64.StdEncoding.DecodeString(strings.Split(string(fileBytes), ",")[1])
	if err != nil {
		return "", err
	}

	return s.Repo.UploadBase64(imageData, "/"+key, ext)
}

func (s *S3Service) DeleteObject(key string) error {
	return s.Repo.DeleteObject(key)
}

func DetectMimeType(base64Str string) (string, string, error) {
	parts := strings.SplitN(base64Str, ",", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid base64 format")
	}

	header := parts[0] // Example: "data:image/png;base64"
	// data := parts[1]   // The actual base64-encoded string

	// Extract MIME type from header
	mimeType := strings.TrimPrefix(strings.Split(header, ";")[0], "data:")
	ext := ""

	switch mimeType {
	case "image/jpeg":
		ext = "jpg"
	case "image/png":
		ext = "png"
	case "image/gif":
		ext = "gif"
	case "image/webp":
		ext = "webp"
	default:
		return "", "", fmt.Errorf("unsupported image type")
	}

	return mimeType, ext, nil
}
