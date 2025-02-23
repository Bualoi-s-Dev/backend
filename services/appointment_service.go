package services

import (
	"context"
	"time"

	"github.com/jinzhu/copier"

	"errors"
	// "time"

	"github.com/Bualoi-s-Dev/backend/dto"
	"github.com/Bualoi-s-Dev/backend/models"
	repositories "github.com/Bualoi-s-Dev/backend/repositories/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrBadRequest     = errors.New("Invalid request data")
	ErrInternalServer = errors.New("Internal server error")
	ErrUnauthorized   = errors.New("Unauthorized")

	ErrStatusInvalid = errors.New("Invalid status")
	ErrStatusTime    = errors.New("Invalid status time")
)

type AppointmentService struct {
	Repo *repositories.AppointmentRepository
}

// literally just getbyID and check if the user is authorized

func NewAppointmentService(repo *repositories.AppointmentRepository) *AppointmentService {
	return &AppointmentService{Repo: repo}
}

func (s *AppointmentService) GetAllAppointment(ctx context.Context, user *models.User) ([]models.Appointment, error) {
	return s.Repo.GetAll(ctx, user.ID, user.Role)
}

func (s *AppointmentService) GetAppointmentById(ctx context.Context, user *models.User, appointmentId primitive.ObjectID) (*models.Appointment, error) {
	appointment, err := s.Repo.GetById(ctx, appointmentId, user.ID, user.Role)
	if err != nil {
		return nil, ErrBadRequest
	}

	if appointment.CustomerID != user.ID && appointment.PhotographerID != user.ID {
		return nil, ErrUnauthorized
	}
	return appointment, nil
}

func (s *AppointmentService) CreateOneAppointment(ctx context.Context, user *models.User, subpackageId primitive.ObjectID, req *dto.AppointmenStrictRequest) (*models.Appointment, error) {
	subpackage, err := s.Repo.FindSubpackageByID(ctx, subPackageId)
	if err != nil {
		return nil, ErrInternalServer
	}

	pkg, err := s.Repo.FindPackageByID(ctx, subpackage.PackageID)
	if err != nil {
		return nil, ErrInternalServer
	}

	appointment := models.Appointment{
		ID:             primitive.NewObjectID(),
		CustomerID:     user.ID,
		PhotographerID: pkg.OwnerID,
		PackageID:      pkg.ID,
		SubPackageID:   subpackageId,
		StartTime:      req.StartTime,
		EndTime:        req.StartTime.Add(time.Duration(subpackage.Duration) * time.Minute),
		Status:         "Pending",
		Location:       req.Location,
	}

	// TODO: Check Schedule before create
	return s.Repo.CreateAppointment(ctx, &appointment)
}

func (s *AppointmentService) UpdateAppointment(ctx context.Context, user *models.User, appointmentId primitive.ObjectID, req *dto.AppointmentRequest) (*models.Appointment, error) {

	appointment, err := s.GetAppointmentById(ctx, user, appointmentId)
	if err != nil {
		return nil, err
	}

	// if canceled or complete it can't be edited
	if appointment.Status == "Canceled" || appointment.Status == "Completed" {
		return nil, ErrStatusInvalid
	}

	if req.StartTime != nil {
		loc, _ := time.LoadLocation("Asia/Bangkok")
		if req.StartTime.Before(time.Now().In(loc)) {
			return nil, ErrBadRequest
		}
		// calculate duration Endtime - StartTime (from appointment)
		duration := appointment.EndTime.Sub(appointment.StartTime)
		endTime := req.StartTime.Add(duration)
		req.EndTime = &endTime
	}
	// use copier
	copier.Copy(appointment, req)
	// fmt.Print(appointment)

	return s.Repo.ReplaceAppointment(ctx, appointmentId, appointment)
}

func (s *AppointmentService) UpdateAppointmentStatus(ctx context.Context, user *models.User, appointmentId primitive.ObjectID, req *dto.AppointmentRequest) (*models.Appointment, error) {
	appointment, err := s.GetAppointmentById(ctx, user, appointmentId)
	if err != nil {
		return nil, err
	}

	if appointment.Status == "Canceled" || appointment.Status == "Completed" {
		return nil, ErrStatusInvalid
	}
	if appointment.Status == "Accepted" && *req.Status == "Canceled" && time.Now().After(appointment.StartTime) {
		return nil, ErrStatusTime
	}

	appointment.Status = *req.Status
	return s.Repo.ReplaceAppointment(ctx, appointmentId, appointment)
}

func (s *AppointmentService) DeleteAppointment(ctx context.Context, appointmentId primitive.ObjectID, user *models.User) error {
	if _, err := s.GetAppointmentById(ctx, user, appointmentId); err != nil {
		return err
	}
	return s.Repo.DeleteAppointment(ctx, appointmentId)
}
