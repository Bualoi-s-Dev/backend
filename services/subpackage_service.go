package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/Bualoi-s-Dev/backend/dto"
	"github.com/Bualoi-s-Dev/backend/models"
	repositories "github.com/Bualoi-s-Dev/backend/repositories/database"
	"github.com/Bualoi-s-Dev/backend/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SubpackageService struct {
	Repository            *repositories.SubpackageRepository
	PackageRepository     *repositories.PackageRepository
	BusyTimeRepository    *repositories.BusyTimeRepository
	AppointmentRepository *repositories.AppointmentRepository
}

func NewSubpackageService(repository *repositories.SubpackageRepository, packageRepository *repositories.PackageRepository, busyTimeRepository *repositories.BusyTimeRepository, appointmentRepo *repositories.AppointmentRepository) *SubpackageService {
	return &SubpackageService{Repository: repository, PackageRepository: packageRepository, BusyTimeRepository: busyTimeRepository, AppointmentRepository: appointmentRepo}
}

func (s *SubpackageService) GetAll(ctx context.Context) ([]models.Subpackage, error) {
	return s.Repository.GetAll(ctx)
}

func (s *SubpackageService) GetFilteredSubpackages(ctx context.Context, filters map[string]string, page, limit int) ([]dto.SubpackageResponse, int, error) {
	items, err := s.Repository.GetAll(ctx)
	if err != nil {
		return nil, 0, err
	}

	var totalCount int
	var responses []dto.SubpackageResponse
	startIdx := (page - 1) * limit
	endIdx := startIdx + limit

	for _, item := range items {
		pkg, err := s.PackageRepository.GetById(ctx, item.PackageID.Hex())
		if err != nil {
			return nil, 0, err
		}
		if !s.passesFilters(pkg, item, filters) {
			continue
		}

		res, err := s.MappedToSubpackageResponse(ctx, &item)
		if err != nil {
			return nil, 0, err
		}
		responses = append(responses, *res)
	}

	totalCount = len(responses)

	// Apply pagination
	if startIdx > totalCount {
		return []dto.SubpackageResponse{}, totalCount, nil
	}
	if endIdx > totalCount {
		endIdx = totalCount
	}

	return responses[startIdx:endIdx], totalCount, nil
}

func (s *SubpackageService) passesFilters(pkg *models.Package, item models.Subpackage, filters map[string]string) bool {
	if filters["packageId"] != "" && pkg.ID.Hex() != filters["packageId"] {
		return false
	}

	if filters["repeatedDay"] != "" {
		days := strings.Split(filters["repeatedDay"], ",")
		dayMap := make(map[string]bool)
		for _, day := range days {
			dayMap[strings.TrimSpace(day)] = true
		}

		match := false
		for _, itemDay := range item.RepeatedDay {
			if dayMap[string(itemDay)] {
				match = true
				break
			}
		}
		if !match {
			return false
		}
	}

	return (filters["title"] == "" || strings.Contains(strings.ToLower(item.Title), strings.ToLower(filters["title"]))) &&
		(filters["type"] == "" || string(pkg.Type) == filters["type"]) &&
		(filters["availableStartTime"] == "" || item.AvailableStartTime >= filters["availableStartTime"]) &&
		(filters["availableEndTime"] == "" || item.AvailableEndTime <= filters["availableEndTime"]) &&
		(filters["availableStartDay"] == "" || item.AvailableStartDay >= filters["availableStartDay"]) &&
		(filters["availableEndDay"] == "" || item.AvailableEndDay <= filters["availableEndDay"])
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
	if subpackage.AvailableStartTime == nil {
		return errors.New("available_start_time is required")
	}
	if subpackage.AvailableEndTime == nil {
		return errors.New("available_end_time is required")
	}

	if (subpackage.IsInf != nil && !*subpackage.IsInf) && (subpackage.AvailableStartDay == nil || subpackage.AvailableEndDay == nil) {
		return errors.New("available_start_day and available_end_day are required")
	}
	return nil
}

func (s *SubpackageService) GetIntersectBusyTime(ctx context.Context, subpackage *models.Subpackage) ([]models.BusyTime, error) {
	parentPackage, err := s.PackageRepository.GetById(ctx, subpackage.PackageID.Hex()) // ! Might be a problem
	if err != nil {
		return nil, err
	}
	ownerId := parentPackage.OwnerID
	busyTimes, err := s.BusyTimeRepository.GetByPhotographerId(ctx, ownerId)
	if err != nil {
		return nil, err
	}

	intersectBusyTime := []models.BusyTime{}
	for _, busyTime := range busyTimes {
		isIntersect, err := s.IsIntersect(ctx, subpackage, &busyTime)
		if err != nil {
			return nil, err
		}
		if isIntersect {
			intersectBusyTime = append(intersectBusyTime, busyTime)
		}
	}

	return intersectBusyTime, nil
}

func (s *SubpackageService) IsIntersect(ctx context.Context, subpackage *models.Subpackage, busyTime *models.BusyTime) (bool, error) {
	if subpackage == nil || busyTime == nil {
		return false, errors.New("invalid input: subpackage or busyTime is nil")
	}

	// Parse subpackage available start and end time
	layout := "15:04"
	_, err := time.Parse(layout, subpackage.AvailableStartTime)
	if err != nil {
		return false, errors.New("invalid available start time format")
	}
	_, err = time.Parse(layout, subpackage.AvailableEndTime)
	if err != nil {
		return false, errors.New("invalid available end time format")
	}

	// If IsInf is false, validate start and end dates
	if !subpackage.IsInf {
		subStartDate, err := time.Parse("2006-01-02", subpackage.AvailableStartDay)
		if err != nil {
			return false, errors.New("invalid available start date format")
		}
		subEndDate, err := time.Parse("2006-01-02", subpackage.AvailableEndDay)
		if err != nil {
			return false, errors.New("invalid available end date format")
		}

		// Check if BusyTime falls within the available date range
		if busyTime.EndTime.Before(subStartDate) || busyTime.StartTime.After(subEndDate) {
			return false, nil
		}
	}

	// Iterate over each day in the busy period
	for d := busyTime.StartTime; d.Before(busyTime.EndTime) || d.Equal(busyTime.EndTime); d = d.Add(24 * time.Hour) {
		for _, day := range subpackage.RepeatedDay {
			if string(day) == strings.ToUpper(d.Weekday().String()[:3]) {
				// Determine busy time range for this specific day
				var busyDayStart time.Time
				var busyDayEnd time.Time

				if d.Format("2006-01-02") == busyTime.StartTime.Format("2006-01-02") {
					// First day: busy period starts from actual start time
					busyDayStart = busyTime.StartTime
					busyDayEnd = time.Date(d.Year(), d.Month(), d.Day(), 23, 59, 59, 0, d.Location())
				} else if d.Format("2006-01-02") == busyTime.EndTime.Format("2006-01-02") {
					// Last day: busy period ends at actual end time
					busyDayStart = time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
					busyDayEnd = busyTime.EndTime
				} else {
					// Full day busy
					busyDayStart = time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
					busyDayEnd = time.Date(d.Year(), d.Month(), d.Day(), 23, 59, 59, 0, d.Location())
				}

				// Convert busy day times to strings
				busyStartTime := busyDayStart.Format("15:04")
				busyEndTime := busyDayEnd.Format("15:04")

				// Check if busy time range overlaps with subpackage available time
				if busyStartTime < subpackage.AvailableEndTime && busyEndTime > subpackage.AvailableStartTime {
					return true, nil
				}
			}
		}
	}

	return false, nil
}

func (s *SubpackageService) MappedToSubpackageResponse(ctx context.Context, subpackage *models.Subpackage) (*dto.SubpackageResponse, error) {
	busyTime, err := s.GetIntersectBusyTime(ctx, subpackage)
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
		AvailableStartTime: subpackage.AvailableStartTime,
		AvailableEndTime:   subpackage.AvailableEndTime,
		AvailableStartDay:  subpackage.AvailableStartDay,
		AvailableEndDay:    subpackage.AvailableEndDay,
		BusyTimes:          busyTime,
	}, nil
}

func (s *SubpackageService) IsSubpackageDeletable(ctx context.Context, subpackageId primitive.ObjectID) (bool, error) {
	appointments, err := s.AppointmentRepository.GetBySubpackageId(ctx, subpackageId)
	if err != nil {
		return false, err
	}

	var isDeletable bool = true
	for _, appointment := range appointments {
		if appointment.Status == models.AppointmentPending {
			isDeletable = false
			break
		}
	}
	return isDeletable, nil
}

func (s *SubpackageService) CheckDate(ctx context.Context, subpackage *dto.SubpackageRequest) error {
	// Check start time before end time
	parsedStartTime, err := time.Parse("15:04", *subpackage.AvailableStartTime)
	if err != nil {
		return err
	}
	parsedEndTime, err := time.Parse("15:04", *subpackage.AvailableEndTime)
	if err != nil {
		return err
	}
	if parsedEndTime.Before(parsedStartTime) {
		return errors.New("end time must be after start time")
	}

	// If IsInf, skip date checks
	if subpackage.IsInf != nil && *subpackage.IsInf {
		return nil
	}
	// Check start date before end date and date in the future
	parsedStartDate, err := time.Parse("2006-01-02", *subpackage.AvailableStartDay)
	if err != nil {
		return err
	}
	parsedEndDate, err := time.Parse("2006-01-02", *subpackage.AvailableEndDay)
	if err != nil {
		return err
	}
	now := time.Now()
	if parsedStartDate.Before(now) || (parsedEndDate.Before(parsedStartDate)) {
		return errors.New("start date must be after current date and end date must be after start date")
	}

	return nil
}
