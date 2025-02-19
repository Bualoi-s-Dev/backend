package services

import (
	"context"
	"errors"
	"time"

	"github.com/Bualoi-s-Dev/backend/models"
	repositories "github.com/Bualoi-s-Dev/backend/repositories/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AppointmentService struct {
	Repo *repositories.AppointmentRepository
}

func NewAppointmentService(repo *repositories.AppointmentRepository) *AppointmentService {
	return &AppointmentService{Repo: repo}
}

func (s *AppointmentService) GetAllAppointment(ctx context.Context, id primitive.ObjectID) ([]models.Appointment, error) {
	return s.Repo.GetAll(ctx, id)
}

func (s *AppointmentService) GetAppointmentById(ctx context.Context, id primitive.ObjectID) (*models.Appointment, error) {
	return s.Repo.GetById(ctx, id)
}

func (s *AppointmentService) CreateAppointment(ctx context.Context, appointment *models.Appointment) error {
	if appointment.StartTime.Before(time.Now()) {
		return errors.New("start time must be in the future")
	}
	if appointment.StartTime.After(appointment.EndTime) {
		return errors.New("start time must be before end time")
	}
	return s.Repo.CreateAppointment(ctx, appointment)
}

func (s *AppointmentService) UpdateAppointment(ctx context.Context, id primitive.ObjectID, appointment *models.Appointment) error {
	if appointment.StartTime.Before(time.Now()) {
		return errors.New("start time must be in the future")
	}
	if appointment.StartTime.After(appointment.EndTime) {
		return errors.New("start time must be before end time")
	}
	return s.Repo.UpdateAppointment(ctx, id, appointment)
}

func (s *AppointmentService) DeleteAppointment(ctx context.Context, id primitive.ObjectID) error {
	return s.Repo.DeleteAppointment(ctx, id)
}
