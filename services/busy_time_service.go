package services

import (
	"context"
	"time"

	"github.com/Bualoi-s-Dev/backend/dto"
	"github.com/Bualoi-s-Dev/backend/models"
	repositories "github.com/Bualoi-s-Dev/backend/repositories/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BusyTimeService struct {
	Repository         *repositories.BusyTimeRepository
	SubpackageRepo     *repositories.SubpackageRepository
	PackageRepo        *repositories.PackageRepository
}

func NewBusyTimeService(repository *repositories.BusyTimeRepository, subpackageRepo *repositories.SubpackageRepository, packageRepo *repositories.PackageRepository) *BusyTimeService {
	return &BusyTimeService{
		Repository:     repository,
		SubpackageRepo: subpackageRepo,
		PackageRepo:    packageRepo,
	}
}

func (s *BusyTimeService) GetAll(ctx context.Context) ([]models.BusyTime, error) {
	return s.Repository.GetAll(ctx)
}

func (s *BusyTimeService) GetById(ctx context.Context, id string) (*models.BusyTime, error) {
	return s.Repository.GetById(ctx, id)
}

func (s *BusyTimeService) GetByPhotographerId(ctx context.Context, photographerId primitive.ObjectID) ([]models.BusyTime, error) {
	return s.Repository.GetByPhotographerId(ctx, photographerId)
}

func (s *BusyTimeService) Create(ctx context.Context, request *dto.BusyTimeRequest, photographerID primitive.ObjectID) error {
	model := request.ToModel(photographerID)
	return s.Repository.Create(ctx, *model)
}

func (s *BusyTimeService) Delete(ctx context.Context, id string) error {
	return s.Repository.DeleteOne(ctx, id)
}

func (s *BusyTimeService) IsSubpackageAvailable(ctx context.Context, subpackageID primitive.ObjectID, startTime, endTime time.Time) (bool, error) {
	subpackage, err := s.SubpackageRepo.GetById(ctx, subpackageID.Hex())
	if err != nil {
		return false, err
	}

	pkg, err := s.PackageRepo.GetById(ctx, subpackage.PackageID.Hex())
	if err != nil {
		return false, err
	}

	busyTimes, err := s.Repository.GetByPhotographerId(ctx, pkg.OwnerID)
	if err != nil {
		return false, err
	}

	for _, busy := range busyTimes {
		//check overlap
		if (startTime.Before(busy.EndTime) && endTime.After(busy.StartTime)) ||
			(startTime.Equal(busy.StartTime) || endTime.Equal(busy.EndTime)) {
			return false, nil
		}
	}

	return true, nil
}
