package services

import (
	"mime/multipart"

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

func (s *S3Service) DeleteObject(key string) error {
	return s.Repo.DeleteObject(key)
}
