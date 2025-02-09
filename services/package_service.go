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
	item := models.Package{
		Title: itemInput.Title,
		Type:  itemInput.Type,
	}
	item.ID = primitive.NewObjectID()

	for idx, photo := range itemInput.Photos {
		photoUrl, err := s.S3Service.UploadBase64([]byte(photo), "package/"+item.ID.Hex()+"_"+strconv.Itoa(idx+1))
		if err != nil {
			return nil, err
		}
		item.PhotoUrls = append(item.PhotoUrls, photoUrl)
	}

	_, err := s.Repo.CreateOne(ctx, &item)
	return &item, err
}

func (s *PackageService) UpdateOne(ctx context.Context, packageId string, updates map[string]interface{}) error {
	_, err := s.Repo.UpdateOne(ctx, packageId, updates)
	return err
}

func (s *PackageService) DeleteOne(ctx context.Context, packageId string) error {
	_, err := s.Repo.DeleteOne(ctx, packageId)
	return err
}
