package services

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/Bualoi-s-Dev/backend/models"
	repositories "github.com/Bualoi-s-Dev/backend/repositories/database"
)

type UserService struct {
	Repo *repositories.UserRepository
	S3Service S3Service
}

func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{Repo: repo}
}

func (s *UserService) GetUserProfile(ctx context.Context, email string) (*models.User, error) {
	return s.Repo.GetUserByEmail(ctx, email)
}

func (s *UserService) CreateUser(ctx context.Context, user *models.User) error {
	return s.Repo.CreateUser(ctx, user)
}

func (s *UserService) UpdateUser(ctx context.Context, email string, updates *models.User) error {
	return s.Repo.UpdateUser(ctx, email, updates)
}

func (s *UserService) GetUserProfilePic(ctx context.Context, email string) (string, error) {
	imageKey, err := s.Repo.GetUserProfilePic(ctx, email)

	if err != nil {
		return "", err
	}
	
	profilePicURL := fmt.Sprintf("https://s3.amazonaws.com/%s/%s", s.S3Service.Repo.BucketName, imageKey)
	return profilePicURL, nil
}

func (s *UserService) UpdateUserProfilePic(ctx context.Context, email string, file *multipart.FileHeader) error {
	imageKey := fmt.Sprintf("profile_pictures/%s", email)

	s3URL, err := s.S3Service.UploadFile(file, imageKey)
	if err != nil {
		return err 
	}

	updates := map[string]interface{}{"profileImage": s3URL}
	return s.Repo.UpdateUserField(ctx, email, updates)
}