package services

import (
	"context"
	"fmt"

	"github.com/Bualoi-s-Dev/backend/dto"
	"github.com/Bualoi-s-Dev/backend/models"
	repositories "github.com/Bualoi-s-Dev/backend/repositories/database"
	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"firebase.google.com/go/auth"
)

type UserService struct {
	Repo      *repositories.UserRepository
	S3Service *S3Service
	AuthClient *auth.Client
}

func NewUserService(repo *repositories.UserRepository, s3Service *S3Service, authClient *auth.Client) *UserService {
	return &UserService{Repo: repo, S3Service: s3Service, AuthClient: authClient}
}

func (s *UserService) FindUser(ctx context.Context, email string) (*models.User, error) {
	return s.Repo.FindUserByEmail(ctx, email)
}

func (s *UserService) GetUser(ctx context.Context, email string) (*dto.UserResponse, error) {
	return s.Repo.GetUserByEmail(ctx, email)
}

func (s *UserService) FindUserByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	return s.Repo.FindUserByID(ctx, id)
}

func (s *UserService) FindEmailByID(ctx context.Context, id primitive.ObjectID) (string, error) {
	return s.Repo.FindEmailByID(ctx, id)
}

func (s *UserService) CreateUser(ctx context.Context, user *models.User) error {
	return s.Repo.CreateUser(ctx, user)
}

func (s *UserService) UpdateUser(ctx context.Context, userId primitive.ObjectID, email string, req *dto.UserRequest) (*dto.UserResponse, error) {
	item, err := s.Repo.FindUserByID(ctx, userId)
	if err != nil {
		return nil, err
	}

	//Check if the role is changed
	roleChanged := req.Role != nil && models.UserRole(*req.Role) != item.Role


	if err := copier.Copy(item, req); err != nil {
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

	// call func change jwt role
	if roleChanged {
		newRole := models.UserRole(*req.Role)

		//find user from firebase to get firebase UID
		firebaseUser, err := s.AuthClient.GetUserByEmail(ctx, email)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch Firebase user: %v", err)
		}
		firebaseUID := firebaseUser.UID

		err = s.AuthClient.SetCustomUserClaims(ctx, firebaseUID, map[string]interface{}{
			"role": string(newRole),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to update Firebase role: %v", err)
		}

		item.Role = newRole
	}

	_, err = s.Repo.ReplaceUser(ctx, userId, item)
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
	item, err := s.Repo.FindUserByID(ctx, userId)
	if err != nil {
		return err
	}
	if err := copier.Copy(item, req); err != nil {
		return err
	}
	_, err = s.Repo.ReplaceUser(ctx, userId, item)
	return err
}
