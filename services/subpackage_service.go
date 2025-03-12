package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Bualoi-s-Dev/backend/dto"
	"github.com/Bualoi-s-Dev/backend/models"
	repositories "github.com/Bualoi-s-Dev/backend/repositories/database"
	"github.com/Bualoi-s-Dev/backend/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SubpackageService struct {
	Repository         *repositories.SubpackageRepository
	PackageRepository  *repositories.PackageRepository
	BusyTimeRepository *repositories.BusyTimeRepository
}

func NewSubpackageService(repository *repositories.SubpackageRepository, packageRepository *repositories.PackageRepository, busyTimeRepository *repositories.BusyTimeRepository) *SubpackageService {
	return &SubpackageService{Repository: repository, PackageRepository: packageRepository, BusyTimeRepository: busyTimeRepository}
}

func (s *SubpackageService) GetAll(ctx context.Context) ([]models.Subpackage, error) {
	return s.Repository.GetAll(ctx)
}

func (s *SubpackageService) GetById(ctx context.Context, id string) (*models.Subpackage, error) {
	return s.Repository.GetById(ctx, id)
}

func (s *SubpackageService) GetByPackageId(ctx context.Context, packageId primitive.ObjectID) ([]models.Subpackage, error) {
	return s.Repository.GetByPackageId(ctx, packageId)
}

func (s *SubpackageService) Create(ctx context.Context, subpackage *models.Subpackage) error {
	subpackage.ID = primitive.NewObjectID()
	return s.Repository.Create(ctx, *subpackage)
}

func (s *SubpackageService) Update(ctx context.Context, id string, subpackage *dto.SubpackageRequest) error {
	bsonSubpackage, err := utils.StructToBsonMap(subpackage)
	if err != nil {
		return err
	}
	return s.Repository.UpdateOne(ctx, id, bsonSubpackage)
}

func (s *SubpackageService) Replace(ctx context.Context, id string, subpackage *models.Subpackage) error {
	return s.Repository.ReplaceOne(ctx, id, *subpackage)
}

func (s *SubpackageService) Delete(ctx context.Context, id string) error {
	return s.Repository.DeleteOne(ctx, id)
}

func (s *SubpackageService) VerifyStrictRequest(ctx context.Context, subpackage *dto.SubpackageRequest) error {
	if subpackage.Title == nil {
		return errors.New("title is required")
	}
	if subpackage.Description == nil {
		return errors.New("description is required")
	}
	if subpackage.Price == nil {
		return errors.New("price is required")
	}
	if subpackage.IsInf == nil {
		return errors.New("is_inf is required")
	}
	if subpackage.RepeatedDay == nil {
		return errors.New("repeated_day is required")
	}
	if subpackage.AvaliableStartTime == nil {
		return errors.New("avaliable_start_time is required")
	}
	if subpackage.AvaliableEndTime == nil {
		return errors.New("avaliable_end_time is required")
	}

	if (subpackage.IsInf != nil && !*subpackage.IsInf) && (subpackage.AvaliableStartDay == nil || subpackage.AvaliableEndDay == nil) {
		return errors.New("avaliable_start_day and avaliable_end_day are required")
	}
	return nil
}

func (s *SubpackageService) FindIntersectBusyTime(ctx context.Context, subpackage *models.Subpackage) ([]models.BusyTime, error) {
	parentPackage, err := s.PackageRepository.GetById(ctx, subpackage.PackageID.Hex())
	if err != nil {
		return nil, err
	}
	ownerId := parentPackage.OwnerID
	busyTimes, err := s.BusyTimeRepository.GetByPhotographerId(ctx, ownerId)
	if err != nil {
		return nil, err
	}

	// TODO: Add intersect busy time
	// intersectBusyTime := []models.BusyTime{}
	// for _, busyTime := range busyTimes {
	// 	isIntersect, err := s.IsIntersect(subpackage, &busyTime)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	if isIntersect {
	// 		intersectBusyTime = append(intersectBusyTime, busyTime)
	// 	}
	// }

	return busyTimes, nil
}

func (s *SubpackageService) IsIntersect(subpackage *models.Subpackage, busyTime *models.BusyTime) (bool, error) {
	// TODO: Implement this method
	// // Check date
	// avaliableStartDay, _ := time.Parse("2006-01-02", subpackage.AvaliableStartDay)
	// avaliableEndDay, _ := time.Parse("2006-01-02", subpackage.AvaliableEndDay)
	// if !subpackage.IsInf && (busyTime.StartTime.Before(avaliableStartDay) || busyTime.EndTime.After(avaliableEndDay)) {
	// 	continue
	// }

	// // Check weekday
	// dayBusyTime := strings.ToUpper(busyTime.StartTime.Weekday().String()[0:3])
	// if !slices.Contains(subpackage.RepeatedDay, models.DayName(dayBusyTime)) {
	// 	continue
	// }

	// // Check time
	// avaliableStartMinute := utils.TimeToMinutes(subpackage.AvaliableStartTime)
	// avaliableEndMinute := utils.TimeToMinutes(subpackage.AvaliableEndTime)

	// busyStartMinute := utils.TimeToMinutes(busyTime.StartTime.Format("15:04"))
	// busyEndMinute := utils.TimeToMinutes(busyTime.EndTime.Format("15:04"))
	// if (busyStartMinute < avaliableStartMinute) || (busyEndMinute > avaliableEndMinute) {
	// 	continue
	// }
	return false, nil
}

func (s *SubpackageService) MappedToSubpackageResponse(ctx context.Context, subpackage *models.Subpackage) (*dto.SubpackageResponse, error) {
	busyTime, err := s.FindIntersectBusyTime(ctx, subpackage)
	if err != nil {
		return nil, err
	}
	busyTimeMap, err := s.GetBusyTimeDateMap(ctx, *subpackage, busyTime)
	if err != nil {
		return nil, err
	}
	return &dto.SubpackageResponse{
		ID:                 subpackage.ID,
		PackageID:          subpackage.PackageID,
		Title:              subpackage.Title,
		Description:        subpackage.Description,
		Price:              subpackage.Price,
		Duration:           subpackage.Duration,
		IsInf:              subpackage.IsInf,
		RepeatedDay:        subpackage.RepeatedDay,
		AvaliableStartTime: subpackage.AvaliableStartTime,
		AvaliableEndTime:   subpackage.AvaliableEndTime,
		AvaliableStartDay:  subpackage.AvaliableStartDay,
		AvaliableEndDay:    subpackage.AvaliableEndDay,
		// TODO: Change this to busyTimes
		BusyTimes:   []models.BusyTime{},
		BusyTimeMap: busyTimeMap,
	}, nil
}

// TODO: Remove this temp function
func (s *SubpackageService) GetBusyTimeDateMap(ctx context.Context, subpackage models.Subpackage, busyTimes []models.BusyTime) (map[string][]models.BusyTime, error) {
	var mapBusyTime = make(map[string][]models.BusyTime)

	var dayRange int
	var startDate time.Time
	if subpackage.IsInf {
		dayRange = 30
		startDate = time.Now()
	} else {
		startDate, _ = time.Parse("2006-01-02", subpackage.AvaliableStartDay)
		endDate, _ := time.Parse("2006-01-02", subpackage.AvaliableEndDay)
		dayRange = int(endDate.Sub(startDate).Hours()/24) + 1
	}
	fmt.Println("dayRange", dayRange)
	fmt.Println("startDate", startDate)
	fmt.Println("busyTimes", busyTimes)
	for _, busyTime := range busyTimes {
		for i := 0; i < dayRange; i++ {
			date := startDate.AddDate(0, 0, i).Format("2006-01-02")

			if !(busyTime.StartTime.Before(startDate) && busyTime.EndTime.After(startDate)) {
				fmt.Println("skip this time", busyTime.StartTime, busyTime.EndTime, startDate)
				continue
			}

			fmt.Println("add this time", busyTime.StartTime, busyTime.EndTime, startDate)
			if _, ok := mapBusyTime[date]; !ok {
				mapBusyTime[date] = []models.BusyTime{}
			}
			mapBusyTime[date] = append(mapBusyTime[date], busyTime)
		}
	}
	fmt.Println("mapBusyTime", mapBusyTime)
	return mapBusyTime, nil
}
