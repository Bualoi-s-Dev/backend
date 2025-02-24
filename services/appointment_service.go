package services

import (
	"context"
	"time"

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
	BusyTimeRepo    *repositories.BusyTimeRepository
}

// literally just getbyID and check if the user is authorized

func NewAppointmentService(appointmentRepo *repositories.AppointmentRepository, packageRepo *repositories.PackageRepository, busyTimeRepo *repositories.BusyTimeRepository) *AppointmentService {
	return &AppointmentService{
		AppointmentRepo: appointmentRepo,
		PackageRepo:     packageRepo,
		BusyTimeRepo:    busyTimeRepo,
	}
}

func (s *AppointmentService) GetAllAppointment(ctx context.Context, user *models.User) ([]models.Appointment, error) {
	return s.AppointmentRepo.GetAll(ctx, user.ID, user.Role)
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

func (s *AppointmentService) CreateOneAppointment(ctx context.Context, user *models.User, subpackageId primitive.ObjectID, req *dto.AppointmenStrictRequest) (*models.Appointment, error) {
	subpackage, err := s.PackageRepo.GetSubpackageById(ctx, subpackageId.Hex())
	if err != nil {
		return nil, apperrors.ErrInternalServer
	}

	pkg, err := s.PackageRepo.GetById(ctx, subpackage.PackageID.Hex())
	if err != nil {
		return nil, apperrors.ErrInternalServer
	}

	appointment := req.ToModel(user, pkg, subpackage)

	// TODO: Check Schedule before create
	return s.AppointmentRepo.CreateAppointment(ctx, appointment)
}

func (s *AppointmentService) UpdateAppointment(ctx context.Context, user *models.User, appointmentId primitive.ObjectID, req *dto.AppointmentRequest) (*models.Appointment, error) {

	appointment, err := s.GetAppointmentById(ctx, user, appointmentId)
	if err != nil {
		return nil, err
	}

	// if canceled or complete it can't be edited
	if appointment.Status == "Canceled" || appointment.Status == "Completed" {
		return nil, apperrors.ErrAppointmentStatusInvalid
	}

	if req.StartTime != nil {
		loc, _ := time.LoadLocation("Asia/Bangkok")
		if req.StartTime.Before(time.Now().In(loc)) {
			return nil, apperrors.ErrBadRequest
		}
		// calculate duration Endtime - StartTime (from appointment)
		duration := busyTime.EndTime.Sub(busyTime.StartTime)
		endTime := req.StartTime.Add(duration)
		req.EndTime = &endTime
	}

	if err := copier.Copy(appointment, req); err != nil {
		return nil, apperrors.ErrBadRequest
	}

	return s.AppointmentRepo.ReplaceAppointment(ctx, appointmentId, appointment)
}

func (s *AppointmentService) UpdateAppointmentStatus(ctx context.Context, user *models.User, appointmentId primitive.ObjectID, req *dto.AppointmentUpdateStatusRequest) (*models.Appointment, error) {
	appointment, err := s.GetAppointmentById(ctx, user, appointmentId)
	if err != nil {
		return nil, err
	}

	if appointment.Status == "Canceled" || appointment.Status == "Completed" {
		return nil, apperrors.ErrAppointmentStatusInvalid
	}
	if appointment.Status == "Accepted" && *req.Status == "Canceled" && time.Now().After(appointment.StartTime) {
		return nil, apperrors.ErrAppointmentStatusTime
	}

	appointment.Status = *req.Status
	return s.AppointmentRepo.ReplaceAppointment(ctx, appointmentId, appointment)
}

func (s *AppointmentService) DeleteAppointment(ctx context.Context, appointmentId primitive.ObjectID, user *models.User) error {
	if _, err := s.GetAppointmentById(ctx, user, appointmentId); err != nil {
		return err
	}
	return s.AppointmentRepo.DeleteAppointment(ctx, appointmentId)
}
