package services

import (
	"context"
	"fmt"

	"github.com/Bualoi-s-Dev/backend/dto"
	"github.com/Bualoi-s-Dev/backend/models"
	repositories "github.com/Bualoi-s-Dev/backend/repositories/database"
	"github.com/Bualoi-s-Dev/backend/utils"
	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService struct {
	Repo      *repositories.UserRepository
	S3Service *S3Service
}

func NewUserService(repo *repositories.UserRepository, s3Service *S3Service) *UserService {
	return &UserService{Repo: repo, S3Service: s3Service}
}

func (s *UserService) FindUser(ctx context.Context, email string) (*models.User, error) {
	return s.Repo.FindUserByEmail(ctx, email)
}

func (s *UserService) GetUser(ctx context.Context, email string) (*dto.UserResponse, error) {
	return s.Repo.GetUserByEmail(ctx, email)
}

func (s *UserService) GetUserByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	return s.Repo.GetUserByID(ctx, id)
}

func (s *UserService) CreateUser(ctx context.Context, user *models.User) error {
	return s.Repo.CreateUser(ctx, user)
}

func (s *UserService) UpdateUser(ctx context.Context, userId primitive.ObjectID, email string, req *dto.UserRequest) (*dto.UserResponse, error) {
	item := &models.User{}
	if err := copier.CopyWithOption(item, req, copier.Option{IgnoreEmpty: true}); err != nil {
		return nil, err
	}

	if req.Profile != nil && *req.Profile != "" {
		key := "profile/" + userId.Hex()
		// Try to delete the existing profile picture
		fmt.Println("key :", key)
		_ = s.S3Service.DeleteObject(key)

		profileUrl, err := s.S3Service.UploadBase64([]byte(*req.Profile), key)
		if err != nil {
			return nil, err
		}

		item.Profile = profileUrl
	}

	updates, err := utils.StructToBsonMap(item)
	if err != nil {
		return nil, err
	}
	_, err = s.Repo.UpdateUser(ctx, userId, updates)
	if err != nil {
		return nil, err
	}
	res, err := s.Repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *UserService) VerifyShowcase(ctx context.Context, ownedPackages []primitive.ObjectID, checkPackages []primitive.ObjectID) bool {
	ownedMap := make(map[string]struct{}, len(ownedPackages))
	for _, pkg := range ownedPackages {
		ownedMap[pkg.Hex()] = struct{}{} // Using an empty struct{} to save memory
	}
	for _, pkgID := range checkPackages {
		if _, ok := ownedMap[pkgID.Hex()]; !ok {
			return false
		}
	}
	return true
}

func (s *UserService) UpdateOwnerPackage(ctx context.Context, userId primitive.ObjectID, req dto.UpdateUserPackageRequest) error {
	updates, err := utils.StructToBsonMap(req)
	if err != nil {
		return err
	}
	_, err = s.Repo.UpdateUser(ctx, userId, updates)
	return err
}
