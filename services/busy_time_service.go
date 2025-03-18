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
func (s *BusyTimeService) CreateFromUser(ctx context.Context, request *dto.BusyTimeStrictRequest, photographerId primitive.ObjectID) (*models.BusyTime, error) {
	model := request.ToModel(photographerId)
	isAvailable, err := s.IsPhotographerAvailable(ctx, photographerId, model.StartTime, model.EndTime, nil)
	if err != nil {
		return nil, err
	}
	if !isAvailable {
		return nil, apperrors.ErrTimeOverlapped
	}
	return nil, s.Repository.Create(ctx, model)
}

// func (s *BusyTimeService) CreateForUpdate(ctx context.Context, request *dto.BusyTimeStrictRequest, oldID, photographerId primitive.ObjectID) error {
// 	model := request.ToModelUpdate(oldID, photographerId)
// 	return s.CreateFromModel(ctx, photographerId, model)
// }

func (s *BusyTimeService) CreateFromSubpackage(ctx context.Context, request *dto.BusyTimeStrictRequest, subpackageId primitive.ObjectID) (*models.BusyTime, error) {
	subpackage, err := s.SubpackageRepo.GetById(ctx, subpackageId.Hex())
	if err != nil {
		return nil, err
	}

	// subpackage.Duration // minute
	// set end time = start time + duration(in minute)
	EndTime := request.StartTime.Add(time.Duration(subpackage.Duration) * time.Minute)
	request.EndTime = EndTime

	pkg, err := s.PackageRepo.GetById(ctx, subpackage.PackageID.Hex())
	if err != nil {
		return nil, err
	}

	photographerId := pkg.OwnerID
	model := request.ToModel(photographerId)                                                                // when customer create first time
	isAvailable, err := s.IsPhotographerAvailable(ctx, photographerId, model.StartTime, model.EndTime, nil) // always validate when create new busyTime via appointment
	if err != nil {
		return nil, err
	}
	if !isAvailable {
		return nil, apperrors.ErrTimeOverlapped
	}
	return model, s.Repository.Create(ctx, model)
}

func (s *BusyTimeService) UpdateValidStatus(ctx context.Context, busyTime *models.BusyTime) error {
	if busyTime.IsValid { // need to validate to reserve a photographer time
		isAvailable, err := s.IsPhotographerAvailable(ctx, busyTime.PhotographerID, busyTime.StartTime, busyTime.EndTime, &busyTime.ID) // ignore current busyTime when accept offer
		if err != nil {
			return err
		}
		if !isAvailable {
			return apperrors.ErrTimeOverlapped
		}
	}
	return s.Repository.UpdateOne(ctx, busyTime)
}

func (s *BusyTimeService) Delete(ctx context.Context, id string) error {
	return s.Repository.DeleteOne(ctx, id)
}

func (s *BusyTimeService) IsPhotographerAvailable(ctx context.Context, photographerId primitive.ObjectID, startTime, endTime time.Time, ignoreBusyTime *primitive.ObjectID) (bool, error) {
	busyTimes, err := s.Repository.GetByPhotographerIdValid(ctx, photographerId)
	if err != nil {
		return false, err
	}

	for _, busy := range busyTimes {
		if ignoreBusyTime != nil && busy.ID == *ignoreBusyTime {
			continue
		}
		// Check overlap
		if (startTime.Before(busy.EndTime) && endTime.After(busy.StartTime)) ||
			(startTime.Equal(busy.StartTime) || endTime.Equal(busy.EndTime)) {
			return false, nil
		}
	}

	return true, nil
}

// TODO: AutoUpdate on overlapped appointment case
// e.g. 1,2 has overlapped appointment
// photogrpaher accepted 1 (so photographer can't accept 2)
// then 1 canceled (after photographer accepted)
// then photogrpaher can accept 2
