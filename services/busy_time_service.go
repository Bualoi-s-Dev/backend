package services

import (
	"context"
	"time"

	"github.com/Bualoi-s-Dev/backend/apperrors"
	"github.com/Bualoi-s-Dev/backend/dto"
	"github.com/Bualoi-s-Dev/backend/models"
	repositories "github.com/Bualoi-s-Dev/backend/repositories/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BusyTimeService struct {
	Repository     *repositories.BusyTimeRepository
	SubpackageRepo *repositories.SubpackageRepository
	PackageRepo    *repositories.PackageRepository
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

func (s *BusyTimeService) Create(ctx context.Context, request *dto.BusyTimeRequest, photographerID primitive.ObjectID) (primitive.ObjectID, error) {
	model := request.ToModel(photographerID)
	isAvailable, err := s.IsPhotographerAvailable(ctx, photographerID, model.StartTime, model.EndTime)
	if err != nil {
		return primitive.NilObjectID, err
	}
	if !isAvailable {
		return primitive.NilObjectID, apperrors.ErrTimeOverlapped
	}
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

	return s.IsPhotographerAvailable(ctx, pkg.OwnerID, startTime, endTime)
}

func (s *BusyTimeService) IsPhotographerAvailable(ctx context.Context, photographerID primitive.ObjectID, startTime, endTime time.Time) (bool, error) {
	busyTimes, err := s.Repository.GetByPhotographerId(ctx, photographerID)
	if err != nil {
		return false, err
	}

	for _, busy := range busyTimes {
		// Check overlap
		if (startTime.Before(busy.EndTime) && endTime.After(busy.StartTime)) ||
			(startTime.Equal(busy.StartTime) || endTime.Equal(busy.EndTime)) {
			return false, nil
		}
	}

	return true, nil
}
