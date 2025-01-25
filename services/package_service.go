package services

import (
	"context"

	"github.com/Bualoi-s-Dev/backend/models"
	repositories "github.com/Bualoi-s-Dev/backend/repositories/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PackageService struct {
	Repo *repositories.PackageRepository
}

func NewPackageService(repo *repositories.PackageRepository) *PackageService {
	return &PackageService{Repo: repo}
}

func (s *PackageService) GetAll(ctx context.Context) ([]models.Package, error) {
	return s.Repo.GetAll(ctx)
}

func (s *PackageService) CreateOne(ctx context.Context, item *models.Package) error {
	item.ID = primitive.NewObjectID()
	_, err := s.Repo.CreateOne(ctx, item)
	return err
}
