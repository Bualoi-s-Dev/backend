package services

import (
	"context"
	"fmt"

	"firebase.google.com/go/auth"
	"github.com/Bualoi-s-Dev/backend/dto"
	"github.com/Bualoi-s-Dev/backend/models"
	repositories "github.com/Bualoi-s-Dev/backend/repositories/database"
	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService struct {
	Repo           *repositories.UserRepository
	S3Service      *S3Service
	PackageService *PackageService
	AuthClient     *auth.Client
	RatingService  *RatingService
}

func NewUserService(repo *repositories.UserRepository, s3Service *S3Service, packageService *PackageService, authClient *auth.Client, ratingService *RatingService) *UserService {
	return &UserService{Repo: repo, S3Service: s3Service, PackageService: packageService, AuthClient: authClient, RatingService: ratingService}
}

func (s *UserService) FindUser(ctx context.Context, email string) (*models.User, error) {
	return s.Repo.FindUserByEmail(ctx, email)
}

func (s *UserService) FindEmailByID(ctx context.Context, id primitive.ObjectID) (string, error) {
	return s.Repo.FindEmailByID(ctx, id)
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*dto.UserResponse, error) {
	user, err := s.Repo.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return s.mappedToUserResponse(ctx, user)
}

func (s *UserService) GetUserByID(ctx context.Context, userId primitive.ObjectID) (*dto.UserResponse, error) {
	user, err := s.Repo.FindUserByID(ctx, userId)
	if err != nil {
		return nil, err
	}

	return s.mappedToUserResponse(ctx, user)
}

func (s *UserService) CreateUser(ctx context.Context, user *models.User) error {
	return s.Repo.CreateUser(ctx, user)
}

func (s *UserService) UpdateUser(ctx context.Context, userId primitive.ObjectID, email string, req *dto.UserRequest) (*dto.UserResponse, error) {
	item, err := s.Repo.FindUserByID(ctx, userId)
	if err != nil {
		return nil, err
	}

	// Check if the role is changed
	roleChanged := req.Role != nil && models.UserRole(*req.Role) != item.Role

	if err := copier.Copy(item, req); err != nil {
		return nil, err
	}

	if req.Profile != nil && *req.Profile != "" {
		key := "profile/" + userId.Hex()
		// Try to delete the existing profile picture
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

		// find user from firebase to get firebase UID
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

	// map to dto.response
	res, err := s.mappedToUserResponse(ctx, item)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *UserService) VerifyShowcase(ctx context.Context, ownerId primitive.ObjectID, checkPackages []primitive.ObjectID) (bool, error) {
	ownedPackages, err := s.PackageService.GetByOwnerId(ctx, ownerId)
	if err != nil {
		return false, err
	}

	ownedMap := make(map[string]struct{}, len(ownedPackages))
	for _, pkg := range ownedPackages {
		ownedMap[pkg.ID.Hex()] = struct{}{} // Using an empty struct{} to save memory
	}
	for _, pkgID := range checkPackages {
		if _, ok := ownedMap[pkgID.Hex()]; !ok {
			return false, nil
		}
	}
	return true, nil
}

func (s *UserService) mappedToUserResponse(ctx context.Context, user *models.User) (*dto.UserResponse, error) {
	packages, err := s.PackageService.GetByOwnerId(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	packageResponse := []dto.PackageResponse{}
	for _, pkg := range packages {
		packageRes, err := s.PackageService.MappedToPackageResponse(ctx, &pkg)
		if err != nil {
			return nil, err
		}
		packageResponse = append(packageResponse, *packageRes)
	}

	showcasePackages, err := s.PackageService.GetByList(ctx, user.ShowcasePackages)
	if err != nil {
		return nil, err
	}
	showcasePackageResponse := []dto.PackageResponse{}
	for _, pkg := range showcasePackages {
		packageRes, err := s.PackageService.MappedToPackageResponse(ctx, &pkg)
		if err != nil {
			return nil, err
		}
		showcasePackageResponse = append(showcasePackageResponse, *packageRes)
	}

	ratingResponse := []dto.RatingResponse{}
	if user.Role == models.Photographer {
		ratings, err := s.RatingService.GetByPhotographerId(ctx, user.ID)
		if err != nil {
			return nil, err
		}

		for _, rating := range ratings {
			ratingRes, err := s.RatingService.MappedToRatingResponse(ctx, &rating)
			if err != nil {
				return nil, err
			}
			ratingResponse = append(ratingResponse, *ratingRes)
		}
	}

	return &dto.UserResponse{
		ID:               user.ID,
		Email:            user.Email,
		Name:             user.Name,
		Gender:           user.Gender,
		Profile:          user.Profile,
		Phone:            user.Phone,
		Location:         user.Location,
		Role:             user.Role,
		Description:      user.Description,
		BankName:         user.BankName,
		BankAccount:      user.BankAccount,
		LineID:           user.LineID,
		Facebook:         user.Facebook,
		Instagram:        user.Instagram,
		ShowcasePackages: showcasePackageResponse,
		Packages:         packageResponse,
		Ratings:		  ratingResponse,
	}, nil
}

func (s *UserService) GetUserRoleByID(ctx context.Context, userId primitive.ObjectID) (*models.UserRole, error) {
	user, err := s.GetUserByID(ctx, userId)
	if err != nil {
		return nil, err
	}

	return &user.Role, nil
}

func (s *UserService) IsPhotographerByUserId(ctx context.Context, userId primitive.ObjectID) (bool, error) {
	userRole, err := s.GetUserRoleByID(ctx, userId)
	if err != nil {
		return false, err
	}

	return *userRole == models.Photographer, nil
}
