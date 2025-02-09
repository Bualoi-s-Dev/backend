package services

import (
	"context"
	"strconv"

	"github.com/Bualoi-s-Dev/backend/models"
	repositories "github.com/Bualoi-s-Dev/backend/repositories/database"
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

func (s *PackageService) CreateOne(ctx context.Context, itemInput *models.PackageRequest) (*models.Package, error) {
	item := MapRequestToPackage(itemInput, nil)

	photoUrls, upErr := s.UploadPackagePhotos(itemInput.Photos, item.ID.Hex())
	if upErr != nil {
		return nil, upErr
	}
	item.PhotoUrls = photoUrls

	_, err := s.Repo.CreateOne(ctx, item)
	return item, err
}

func (s *PackageService) ReplaceOne(ctx context.Context, packageId string, updates *models.PackageRequest) (*models.Package, error) {
	item := MapRequestToPackage(updates, &packageId)

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
	photoUrls, upErr := s.UploadPackagePhotos(updates.Photos, packageId)
	if upErr != nil {
		return nil, upErr
	}
	item.PhotoUrls = photoUrls

	_, err := s.Repo.ReplaceOne(ctx, packageId, item)
	return item, err
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

func MapRequestToPackage(req *models.PackageRequest, id *string) *models.Package {
	item := models.Package{
		Title: req.Title,
		Type:  req.Type,
	}
	if id != nil {
		objectId, _ := primitive.ObjectIDFromHex(*id)
		item.ID = objectId
	} else {
		item.ID = primitive.NewObjectID()
	}
	return &item
}

func (s *PackageService) UploadPackagePhotos(photoBase64 []string, id string) ([]string, error) {
	var photoUrls []string
	for idx, photo := range photoBase64 {
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
		err := s.S3Service.DeleteObject(photo)
		if err != nil {
			return err
		}
	}
	return nil
}
