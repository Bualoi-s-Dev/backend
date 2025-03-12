package services

import (
	"context"

	"github.com/Bualoi-s-Dev/backend/apperrors"

	"github.com/Bualoi-s-Dev/backend/dto"
	"github.com/Bualoi-s-Dev/backend/models"
	repositories "github.com/Bualoi-s-Dev/backend/repositories/database"
	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AppointmentService struct {
	AppointmentRepo *repositories.AppointmentRepository
	PackageRepo     *repositories.PackageRepository
	SubpackageRepo  *repositories.SubpackageRepository
	BusyTimeRepo    *repositories.BusyTimeRepository
	UserRepo        *repositories.UserRepository
}

// literally just getbyID and check if the user is authorized

func NewAppointmentService(appointmentRepo *repositories.AppointmentRepository, packageRepo *repositories.PackageRepository, subpackageRepo *repositories.SubpackageRepository, busyTimeRepo *repositories.BusyTimeRepository, userRepo *repositories.UserRepository) *AppointmentService {
	return &AppointmentService{
		AppointmentRepo: appointmentRepo,
		PackageRepo:     packageRepo,
		SubpackageRepo:  subpackageRepo,
		BusyTimeRepo:    busyTimeRepo,
		UserRepo:        userRepo,
	}
}

func (s *AppointmentService) GetAllAppointment(ctx context.Context, user *models.User) ([]models.Appointment, error) {
	return s.AppointmentRepo.GetAll(ctx, user.ID, user.Role)
}

func (s *AppointmentService) GetAllAppointmentDetail(ctx context.Context, user *models.User) ([]models.AppointmentDetail, error) {
	allAppointment, err := s.AppointmentRepo.GetAll(ctx, user.ID, user.Role)
	if err != nil {
		return nil, apperrors.ErrBadRequest
	}
	var appointmentDetails []models.AppointmentDetail
	for _, appointment := range allAppointment {
		pkg, err := s.PackageRepo.GetById(ctx, appointment.PackageID.Hex())
		if err != nil {
			return nil, apperrors.ErrInternalServer
		}
		subpackage, err := s.SubpackageRepo.GetById(ctx, appointment.SubpackageID.Hex())
		if err != nil {
			return nil, apperrors.ErrInternalServer
		}
		busyTime, err := s.BusyTimeRepo.GetById(ctx, appointment.BusyTimeID.Hex())
		if err != nil {
			return nil, apperrors.ErrInternalServer
		}
		var customerName, photographerName string
		if user.Role == models.Customer {
			customerName = user.Name
			photographer, err := s.UserRepo.FindUserByID(ctx, appointment.PhotographerID)
			if err != nil {
				return nil, apperrors.ErrInternalServer
			}
			photographerName = photographer.Name
		} else if user.Role == models.Photographer {
			photographerName = user.Name
			customer, err := s.UserRepo.FindUserByID(ctx, appointment.CustomerID)
			if err != nil {
				return nil, apperrors.ErrInternalServer
			}
			customerName = customer.Name
		}

		detail := models.AppointmentDetail{
			ID:               appointment.ID,
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
		appointmentDetails = append(appointmentDetails, detail)
	}
	return appointmentDetails, nil
}

func (s *AppointmentService) GetAppointmentById(ctx context.Context, user *models.User, appointmentId primitive.ObjectID) (*models.Appointment, error) {
	appointment, err := s.AppointmentRepo.GetById(ctx, appointmentId, user.ID, user.Role)
	if err != nil {
		return nil, apperrors.ErrBadRequest
	}

	if appointment.CustomerID != user.ID && appointment.PhotographerID != user.ID {
		return nil, apperrors.ErrUnauthorized
	}
	return appointment, nil
}

func (s *AppointmentService) CreateOneAppointment(ctx context.Context, user *models.User, subpackageId primitive.ObjectID, busyTime *models.BusyTime, req *dto.AppointmenStrictRequest) (*models.Appointment, error) {
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

func (s *AppointmentService) UpdateAppointment(ctx context.Context, user *models.User, appointment *models.Appointment, req *dto.AppointmentRequest) (*models.Appointment, error) {
	if err := copier.Copy(appointment, req); err != nil {
		return nil, err
	}

	return s.AppointmentRepo.ReplaceAppointment(ctx, appointment)
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
