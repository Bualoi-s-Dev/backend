package services

import (
	"context"
	"errors"
	"math/rand"
	"strings"

	"github.com/Bualoi-s-Dev/backend/dto"
	"github.com/Bualoi-s-Dev/backend/models"
	repositories "github.com/Bualoi-s-Dev/backend/repositories/database"
	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PackageService struct {
	Repo              *repositories.PackageRepository
	S3Service         *S3Service
	SubpackageService *SubpackageService
	UserRepo          *repositories.UserRepository
}

func NewPackageService(repo *repositories.PackageRepository, s3Service *S3Service, subpackageService *SubpackageService, userRepo *repositories.UserRepository) *PackageService {
	return &PackageService{Repo: repo, S3Service: s3Service, SubpackageService: subpackageService, UserRepo: userRepo}
}

func (s *PackageService) GetAll(ctx context.Context) ([]models.Package, error) {
	return s.Repo.GetAll(ctx)
}

func (s *PackageService) GetAllRecommended(ctx context.Context, size int) ([]models.Package, error) {
	pkgs, err := s.Repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	if len(pkgs) < size {
		size = len(pkgs)
	}
	// TODO: Integrate with rating system later
	rand.Shuffle(len(pkgs), func(i, j int) { pkgs[i], pkgs[j] = pkgs[j], pkgs[i] })

	return pkgs[:size], nil

}

func (s *PackageService) GetById(ctx context.Context, packageId string) (*models.Package, error) {
	return s.Repo.GetById(ctx, packageId)
}

func (s *PackageService) GetByList(ctx context.Context, packageIds []primitive.ObjectID) ([]models.Package, error) {
	if len(packageIds) == 0 {
		return []models.Package{}, nil
	}
	return s.Repo.GetManyId(ctx, packageIds)
}

func (s *PackageService) GetByOwnerId(ctx context.Context, ownerId primitive.ObjectID) ([]models.Package, error) {
	return s.Repo.GetByOwnerId(ctx, ownerId)
}

func (s *PackageService) CreateOne(ctx context.Context, itemInput *dto.PackageRequest, ownerId primitive.ObjectID) (*models.Package, error) {
	item := itemInput.ToModel(ownerId)
	item.ID = primitive.NewObjectID()

	photoUrls, upErr := s.UploadPackagePhotos(*itemInput.Photos, item.ID.Hex())
	if upErr != nil {
		return nil, upErr
	}
	item.PhotoUrls = photoUrls

	_, err := s.Repo.CreateOne(ctx, item)
	return item, err
}

func (s *PackageService) UpdateOne(ctx context.Context, packageId string, updates *dto.PackageRequest) (*models.Package, error) {
	// Fetch current package
	pkg, err := s.Repo.GetById(ctx, packageId)
	if err != nil {
		return nil, err
	}
	// Replace with new values
	if err := copier.Copy(pkg, updates); err != nil {
		return nil, err
	}

	// Upload new photos if any
	if updates.Photos != nil {
		// Delete old photos
		delErr := s.DeletePackagePhotos(pkg.PhotoUrls)
		if delErr != nil {
			return nil, delErr
		}
		pkg.PhotoUrls = []string{}

		// Upload new photos
		photoUrls, upErr := s.UploadPackagePhotos(*updates.Photos, packageId)
		if upErr != nil {
			return nil, upErr
		}
		pkg.PhotoUrls = photoUrls
	}

	_, err = s.Repo.ReplaceOne(ctx, packageId, pkg)
	return pkg, err
}

func (s *PackageService) DeleteOne(ctx context.Context, packageId string) error {
	// Delete photos
	curPackage, findErr := s.GetById(ctx, packageId)
	if findErr != nil {
		return findErr
	}
	delErr := s.DeletePackagePhotos(curPackage.PhotoUrls)
	if delErr != nil {
		return delErr
	}

	_, err := s.Repo.DeleteOne(ctx, packageId)
	return err
}

// Helper function

func (s *PackageService) UploadPackagePhotos(photoBase64 []string, id string) ([]string, error) {
	photoUrls := []string{}
	for _, photo := range photoBase64 {
		// Add / to the photo path
		photoUrl, err := s.S3Service.UploadBase64([]byte(photo), "package/"+id)
		if err != nil {
			return nil, err
		}
		photoUrls = append(photoUrls, photoUrl)
	}
	return photoUrls, nil
}

func (s *PackageService) DeletePackagePhotos(photoUrls []string) error {
	for _, photo := range photoUrls {
		// Remove / from the photo path
		cleanedPath := strings.TrimPrefix(photo, "/")
		err := s.S3Service.DeleteObject(cleanedPath)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *PackageService) CheckOwner(ctx context.Context, user *models.User, packageId string) (bool, error) {
	pkg, err := s.Repo.GetById(ctx, packageId)
	if err != nil {
		return false, err
	}
	return pkg.OwnerID == user.ID, nil
}

func (s *PackageService) CheckPackageExist(ctx context.Context, packageId string) error {
	_, err := s.Repo.GetById(ctx, packageId)
	return err
}

func (s *PackageService) VerifyStrictRequest(ctx context.Context, req *dto.PackageRequest) error {
	if req.Title == nil {
		return errors.New("title is required")
	}
	if req.Type == nil {
		return errors.New("type is required")
	}
	if req.Photos == nil {
		return errors.New("photos is required")
	}
	return nil
}

func (s *PackageService) MappedToPackageResponse(ctx context.Context, item *models.Package) (*dto.PackageResponse, error) {
	subpackages, err := s.SubpackageService.GetByPackageId(ctx, item.ID)
	if err != nil {
		return nil, err
	}

	return &dto.PackageResponse{
		ID:          item.ID,
		OwnerID:     item.OwnerID,
		Title:       item.Title,
		Type:        item.Type,
		PhotoUrls:   item.PhotoUrls,
		SubPackages: subpackages,
	}, nil

}

func (s *PackageService) FilterPackage(ctx context.Context, item *models.Package, searchTitle, searchOwnerName string, searchType models.PackageType) (bool, error) {
	hasSearchType, hasSearchTitle, hasSearchOwnerName := searchType != "", searchTitle != "", searchOwnerName != ""
	if hasSearchTitle && !strings.HasPrefix(strings.ToLower(item.Title), strings.ToLower(searchTitle)) {
		return false, nil
	}
	if hasSearchOwnerName {
		ownerUser, err := s.UserRepo.FindUserByID(ctx, item.OwnerID)
		if err != nil {
			return false, err
		}
		if !strings.HasPrefix(strings.ToLower(ownerUser.Name), strings.ToLower(searchOwnerName)) {
			return false, nil
		}
	}
	if hasSearchType && item.Type != searchType {
		return false, nil
	}
	return true, nil
}
