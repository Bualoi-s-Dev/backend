package services

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Bualoi-s-Dev/backend/apperrors"

	"github.com/Bualoi-s-Dev/backend/dto"
	"github.com/Bualoi-s-Dev/backend/models"
	repositories "github.com/Bualoi-s-Dev/backend/repositories/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AppointmentService struct {
	AppointmentRepo *repositories.AppointmentRepository
	PackageRepo     *repositories.PackageRepository
	SubpackageRepo  *repositories.SubpackageRepository
	BusyTimeRepo    *repositories.BusyTimeRepository
	UserRepo        *repositories.UserRepository
	PaymentService  *PaymentService
}

// literally just getbyID and check if the user is authorized

func NewAppointmentService(appointmentRepo *repositories.AppointmentRepository, packageRepo *repositories.PackageRepository, subpackageRepo *repositories.SubpackageRepository,
	busyTimeRepo *repositories.BusyTimeRepository, userRepo *repositories.UserRepository, paymentService *PaymentService) *AppointmentService {
	return &AppointmentService{
		AppointmentRepo: appointmentRepo,
		PackageRepo:     packageRepo,
		SubpackageRepo:  subpackageRepo,
		BusyTimeRepo:    busyTimeRepo,
		UserRepo:        userRepo,
		PaymentService:  paymentService,
	}
}

func (s *AppointmentService) GetAllAppointment(ctx context.Context, user *models.User) ([]models.Appointment, error) {
	return s.AppointmentRepo.GetAll(ctx, user.ID, user.Role)
}
func (s *AppointmentService) GetAppointmentDetailById(ctx context.Context, user *models.User, appointment *models.Appointment) (*dto.AppointmentDetail, error) {
	pkg := appointment.Package
	subpackage := appointment.Subpackage

	busyTime, err := s.BusyTimeRepo.GetById(ctx, appointment.BusyTimeID.Hex())
	if err != nil {
		fmt.Println("(GetAllAppointmentDetail) Error while getting busyTime")
		return nil, err
	}
	var customerName, photographerName string
	if user != nil && user.Role == models.Customer {
		customerName = user.Name
		photographer, err := s.UserRepo.FindUserByID(ctx, appointment.PhotographerID)
		if err != nil {
			fmt.Println("(GetAllAppointmentDetail) Error while getting photographer")
			return nil, err
		}
		photographerName = photographer.Name
	} else if user != nil && user.Role == models.Photographer {
		photographerName = user.Name
		customer, err := s.UserRepo.FindUserByID(ctx, appointment.CustomerID)
		if err != nil {
			fmt.Println("(GetAllAppointmentDetail) Error while getting customer")
			return nil, err
		}
		customerName = customer.Name
	} else {
		photographer, err := s.UserRepo.FindUserByID(ctx, appointment.PhotographerID)
		if err != nil {
			return nil, err
		}
		customer, err := s.UserRepo.FindUserByID(ctx, appointment.CustomerID)
		if err != nil {
			return nil, err
		}
		customerName = customer.Name
		photographerName = photographer.Name
	}

	detail := &dto.AppointmentDetail{
		ID:               appointment.ID,
		Package:          pkg,
		Subpackage:       subpackage,
		PhotographerID:   appointment.PhotographerID,
		CustomerID:       appointment.CustomerID,
		PackageName:      pkg.Title,
		SubpackageName:   subpackage.Title,
		CustomerName:     customerName,
		PhotographerName: photographerName,
		Price:            appointment.Price,
		StartTime:        busyTime.StartTime,
		EndTime:          busyTime.EndTime,
		Status:           appointment.Status,
		Location:         appointment.Location,
	}
	return detail, nil
}
func (s *AppointmentService) GetFilteredAppointments(ctx context.Context, user *models.User, filters map[string]string, page, limit int) ([]dto.AppointmentResponse, int, error) {
	items, err := s.AppointmentRepo.GetAll(ctx, user.ID, user.Role)
	if err != nil {
		return nil, 0, err
	}

	var totalCount int
	var appointments []dto.AppointmentResponse
	startIdx := (page - 1) * limit
	endIdx := startIdx + limit

	// Convert price filters to integers safely
	minPrice, maxPrice := 0, int(^uint(0)>>1) // Default: minPrice = 0, maxPrice = max int

	if filters["minPrice"] != "" {
		if parsedMinPrice, err := strconv.Atoi(filters["minPrice"]); err == nil {
			minPrice = parsedMinPrice
		} else {
			return nil, 0, fmt.Errorf("invalid minPrice filter: %v", err)
		}
	}

	if filters["maxPrice"] != "" {
		if parsedMaxPrice, err := strconv.Atoi(filters["maxPrice"]); err == nil {
			maxPrice = parsedMaxPrice
		} else {
			return nil, 0, fmt.Errorf("invalid maxPrice filter: %v", err)
		}
	}

	for _, item := range items {
		if filters["status"] != "" && string(item.Status) != filters["status"] {
			continue
		}
		if item.Price < minPrice || item.Price > maxPrice {
			continue
		}

		subPkg := &item.Subpackage
		if !s.passesFilters(subPkg, filters) {
			continue
		}

		if filters["name"] != "" {
			pkg := &item.Package

			ctm, err := s.UserRepo.FindUserByID(ctx, item.CustomerID)
			if err != nil {
				fmt.Println("(GetFilteredAppointments) Error while getting customer")
				return nil, 0, err
			}

			if !strings.HasPrefix(strings.ToLower(pkg.Title), strings.ToLower(filters["name"])) && !strings.HasPrefix(strings.ToLower(ctm.Name), strings.ToLower(filters["name"])) {
				continue
			}
		}

		appointmentResponse := dto.AppointmentResponse{
			ID:             item.ID,
			CustomerID:     item.CustomerID,
			PhotographerID: item.PhotographerID,
			Package:        item.Package,
			Subpackage:     item.Subpackage,
			BusyTimeID:     item.BusyTimeID,
			Status:         item.Status,
			Location:       item.Location,
			Price:          item.Price,
		}

		appointments = append(appointments, appointmentResponse)
	}

	totalCount = len(appointments)

	// Apply pagination
	if startIdx > totalCount {
		return []dto.AppointmentResponse{}, totalCount, nil
	}
	if endIdx > totalCount {
		endIdx = totalCount
	}
	return appointments[startIdx:endIdx], totalCount, nil
}

func (s *AppointmentService) passesFilters(subPkg *models.Subpackage, filters map[string]string) bool {
	return (filters["availableStartDay"] == "" || subPkg.AvailableStartDay >= filters["availableStartDay"]) &&
		(filters["availableEndDay"] == "" || subPkg.AvailableEndDay <= filters["availableEndDay"])
}

func (s *AppointmentService) GetFilteredAppointmentDetail(ctx context.Context, user *models.User, filters map[string]string) ([]dto.AppointmentDetail, error) {
	allAppointments, err := s.AppointmentRepo.GetAll(ctx, user.ID, user.Role)
	if err != nil {
		return nil, apperrors.ErrBadRequest
	}

	var filteredAppointments []dto.AppointmentDetail
	for _, appointment := range allAppointments {
		detail, err := s.GetAppointmentDetailById(ctx, user, &appointment)
		if err != nil {
			return nil, err
		}
		if matchesFilters(detail, filters) {
			filteredAppointments = append(filteredAppointments, *detail)
		}
	}

	sort.Slice(filteredAppointments, func(i, j int) bool {
		statusOrder := map[string]int{
			"Pending":   1,
			"Accepted":  2,
			"Completed": 3,
			"Rejected":  4,
			"Canceled":  5,
		}
		if statusOrder[string(filteredAppointments[i].Status)] != statusOrder[string(filteredAppointments[j].Status)] {
			return statusOrder[string(filteredAppointments[i].Status)] < statusOrder[string(filteredAppointments[j].Status)]
		}
		return filteredAppointments[i].StartTime.Before(filteredAppointments[j].StartTime)
	})

	return filteredAppointments, nil
}

func matchesFilters(detail *dto.AppointmentDetail, filters map[string]string) bool {
	stringFilters := map[string]string{
		"status": string(detail.Status),
	}
	for key, value := range stringFilters {
		if filters[key] != "" && filters[key] != value {
			return false
		}
	}

	if filters["packageName"] != "" && !strings.HasPrefix(detail.PackageName, filters["packageName"]) {
		return false
	}
	if filters["subpackageName"] != "" && !strings.HasPrefix(detail.SubpackageName, filters["subpackageName"]) {
		return false
	}
	if filters["customerName"] != "" && !strings.HasPrefix(detail.CustomerName, filters["customerName"]) {
		return false
	}
	if filters["photographerName"] != "" && !strings.HasPrefix(detail.PhotographerName, filters["photographerName"]) {
		return false
	}
	if filters["location"] != "" && !strings.HasPrefix(detail.Location, filters["location"]) {
		return false
	}

	minPrice, maxPrice := 0, int(^uint(0)>>1)
	if filters["minPrice"] != "" {
		if parsedMinPrice, err := strconv.Atoi(filters["minPrice"]); err == nil {
			minPrice = parsedMinPrice
		} else {
			return false
		}
	}
	if filters["maxPrice"] != "" {
		if parsedMaxPrice, err := strconv.Atoi(filters["maxPrice"]); err == nil {
			maxPrice = parsedMaxPrice
		} else {
			return false
		}
	}
	if detail.Price < minPrice || detail.Price > maxPrice {
		return false
	}

	for _, key := range []string{"availableStartTime", "availableEndTime"} {
		if filters[key] != "" {
			parsedTime, err := time.Parse("15:04", filters[key]) // hh:mm format
			if err != nil {
				return false
			}

			var apptTime time.Time
			if key == "availableStartTime" {
				apptTime = detail.StartTime
			} else {
				apptTime = detail.EndTime
			}

			apptMinutes := apptTime.Hour()*60 + apptTime.Minute()
			filterMinutes := parsedTime.Hour()*60 + parsedTime.Minute()
			fmt.Println(apptMinutes, filterMinutes)
			if key == "availableStartTime" && apptMinutes < filterMinutes {
				return false
			}
			if key == "availableEndTime" && apptMinutes > filterMinutes {

				return false
			}
		}
	}

	for _, key := range []string{"startTime", "endTime"} {
		if filters[key] != "" {
			if parsedTime, err := time.Parse(time.RFC3339, filters[key]); err == nil {
				if (key == "startTime" && detail.StartTime.Before(parsedTime)) || (key == "endTime" && detail.StartTime.After(parsedTime)) {
					return false
				}
			} else {
				return false
			}
		}
	}
	return true
}

func (s *AppointmentService) GetAppointmentById(ctx context.Context, user *models.User, appointmentId primitive.ObjectID) (*models.Appointment, error) {
	appointment, err := s.AppointmentRepo.GetById(ctx, appointmentId)
	if err != nil {
		return nil, apperrors.ErrBadRequest
	}

	if appointment.CustomerID != user.ID && appointment.PhotographerID != user.ID {
		return nil, apperrors.ErrUnauthorized
	}
	return appointment, nil
}

func (s *AppointmentService) CreateOneAppointment(ctx context.Context, user *models.User, subpackageId primitive.ObjectID, busyTime *models.BusyTime, req *dto.AppointmentStrictRequest) (*models.Appointment, error) {
	subpackage, err := s.SubpackageRepo.GetById(ctx, subpackageId.Hex())
	if err != nil {
		return nil, err
	}

	pkg, err := s.PackageRepo.GetById(ctx, subpackage.PackageID.Hex())
	if err != nil {
		return nil, err
	}
	appointment := req.ToModel(user, pkg, subpackage, busyTime)

	return s.AppointmentRepo.CreateAppointment(ctx, appointment)
}

func (s *AppointmentService) UpdateAppointmentStatus(ctx context.Context, user *models.User, appointment *models.Appointment, req *dto.AppointmentUpdateStatusRequest) (*models.Appointment, error) {
	appointment.Status = req.Status
	return s.AppointmentRepo.ReplaceAppointment(ctx, appointment)
}

func (s *AppointmentService) DeleteAppointment(ctx context.Context, appointmentId primitive.ObjectID, user *models.User) error {
	if _, err := s.GetAppointmentById(ctx, user, appointmentId); err != nil {
		return err
	}
	return s.AppointmentRepo.DeleteAppointment(ctx, appointmentId)
}

func (s *AppointmentService) AutoUpdateAppointmentStatus(ctx context.Context) error {

	fmt.Println("Running scheduled update...")

	// filter only start_time is grater than current time and status is "Pending"
	// TODO: Fix this curse later
	loc, _ := time.LoadLocation("Asia/Bangkok")
	t := time.Now().In(loc)
	currentTime := time.Date(
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second(),
		t.Nanosecond(), time.UTC,
	)

	go func() {
		s.AppointmentRepo.UpdateCanceledAppointment(ctx, currentTime)
	}()

	// filter only end_time is less than current time and status is "Accepted"
	// TODO: Maybe Mapped this two go routine into loop or function call
	go func() {
		updatedIds, _ := s.AppointmentRepo.UpdateCompletedAppointment(ctx, currentTime)
		for _, id := range updatedIds {
			s.PaymentService.CreatePayment(ctx, id, "", "")
		}
	}()

	return nil
}
