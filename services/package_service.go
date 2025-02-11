package services

import (
	"context"
	"strconv"
	"strings"

	"github.com/Bualoi-s-Dev/backend/dto"
	"github.com/Bualoi-s-Dev/backend/models"
	repositories "github.com/Bualoi-s-Dev/backend/repositories/database"
	"github.com/Bualoi-s-Dev/backend/utils"
	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PackageService struct {
	Repo      *repositories.PackageRepository
	S3Service *S3Service
}

func NewPackageService(repo *repositories.PackageRepository, s3Service *S3Service) *PackageService {
	return &PackageService{Repo: repo, S3Service: s3Service}
}

func (s *PackageService) GetAll(ctx context.Context) ([]models.Package, error) {
	return s.Repo.GetAll(ctx)
}

func (s *PackageService) GetById(ctx context.Context, packageId string) (*models.Package, error) {
	return s.Repo.GetById(ctx, packageId)
}

func (s *PackageService) GetByList(ctx context.Context, packageIds []primitive.ObjectID) ([]models.Package, error) {
	return s.Repo.GetManyId(ctx, packageIds)
}

func (s *PackageService) CreateOne(ctx context.Context, itemInput *dto.PackageStrictRequest, ownerId primitive.ObjectID) (*models.Package, error) {
	item := s.NewPackageFromRequest(itemInput, ownerId)

	photoUrls, upErr := s.UploadPackagePhotos(itemInput.Photos, item.ID.Hex())
	if upErr != nil {
		return nil, upErr
	}
	item.PhotoUrls = photoUrls

	_, err := s.Repo.CreateOne(ctx, item)
	return item, err
}

func (s *PackageService) UpdateOne(ctx context.Context, packageId string, updates *dto.PackageRequest) (*models.Package, error) {
	pkg := &models.Package{}
	if err := copier.Copy(pkg, updates); err != nil {
		return nil, err
	}

	if updates.Photos != nil {
		// Delete old photos
		curPackage, findErr := s.GetById(ctx, packageId)
		if findErr != nil {
			return nil, findErr
		}
		delErr := s.DeletePackagePhotos(curPackage.PhotoUrls)
		if delErr != nil {
			return nil, delErr
		}

		// Upload new photos
		photoUrls, upErr := s.UploadPackagePhotos(*updates.Photos, packageId)
		if upErr != nil {
			return nil, upErr
		}
		pkg.PhotoUrls = photoUrls
	}

	item, err := utils.StructToBsonMap(pkg)
	if err != nil {
		return nil, err
	}

	_, err = s.Repo.UpdateOne(ctx, packageId, item)

	updatedItem, findErr := s.GetById(ctx, packageId)
	if findErr != nil {
		return nil, findErr
	}
	return updatedItem, err
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

func (s *PackageService) NewPackageFromRequest(req *dto.PackageStrictRequest, ownerId primitive.ObjectID) *models.Package {
	item := models.Package{
		ID:      primitive.NewObjectID(),
		OwnerID: ownerId,
		Title:   req.Title,
		Type:    req.Type,
	}
	return &item
}

func (s *PackageService) UploadPackagePhotos(photoBase64 []string, id string) ([]string, error) {
	var photoUrls []string
	for idx, photo := range photoBase64 {
		// Add / to the photo path
		photoUrl, err := s.S3Service.UploadBase64([]byte(photo), "package/"+id+"_"+strconv.Itoa(idx+1))
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

func (s *PackageService) CheckOwner(user *models.User, packageId string) bool {
	for _, id := range user.Packages {
		if id.Hex() == packageId {
			return true
		}
	}
	return false
}
