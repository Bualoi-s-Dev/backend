package services

import (
	"context"
	"fmt"

	"github.com/Bualoi-s-Dev/backend/models"
	repositories "github.com/Bualoi-s-Dev/backend/repositories/database"
)

type UserService struct {
	Repo      *repositories.UserRepository
	S3Service *S3Service
}

func NewUserService(repo *repositories.UserRepository, s3Service *S3Service) *UserService {
	return &UserService{Repo: repo, S3Service: s3Service}
}

func (s *UserService) GetUser(ctx context.Context, email string) (*models.User, error) {
	return s.Repo.GetUserByEmail(ctx, email)
}

func (s *UserService) CreateUser(ctx context.Context, user *models.User) error {
	return s.Repo.CreateUser(ctx, user)
}

func (s *UserService) UpdateUser(ctx context.Context, userId string, email string, updates *models.User) (*models.User, error) {
	if updates.Profile != "" {
		key := "profile/" + userId
		// Try to delete the existing profile picture
		fmt.Println("key :", key)
		_ = s.S3Service.DeleteObject(key)

		profileUrl, err := s.S3Service.UploadBase64([]byte(updates.Profile), key)
		if err != nil {
			return nil, err
		}

		updates.Profile = profileUrl
	}

	return s.Repo.UpdateUser(ctx, email, updates)
}
