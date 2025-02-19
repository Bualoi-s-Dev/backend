package services

import (
	"context"
	// "errors"
	// "time"

	"github.com/Bualoi-s-Dev/backend/models"
	repositories "github.com/Bualoi-s-Dev/backend/repositories/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AppointmentService struct {
	Repo *repositories.AppointmentRepository
}

func NewAppointmentService(repo *repositories.AppointmentRepository) *AppointmentService {
	return &AppointmentService{Repo: repo}
}

func (s *AppointmentService) GetAllAppointment(ctx context.Context, userID primitive.ObjectID, userRole models.UserRole) ([]models.Appointment, error) {
	return s.Repo.GetAll(ctx, userID, userRole)
}

func (s *AppointmentService) GetAppointmentById(ctx context.Context, apopintmentID primitive.ObjectID, userID primitive.ObjectID, userRole models.UserRole) (*models.Appointment, error) {
	return s.Repo.GetById(ctx, apopintmentID, userID, userRole)
}

func (s *AppointmentService) CreateOneAppointment(ctx context.Context, appointment *models.Appointment) (*mongo.InsertOneResult, error) {

	return s.Repo.CreateAppointment(ctx, appointment)
}

func (s *AppointmentService) FindSubpackageByID(ctx context.Context, id primitive.ObjectID) (*models.Subpackage, error) {
	return s.Repo.FindSubpackageByID(ctx, id)
}

func (s *AppointmentService) FindPackageByID(ctx context.Context, id primitive.ObjectID) (*models.Package, error) {
	return s.Repo.FindPackageByID(ctx, id)
}

// func (s *AppointmentService) UpdateAppointment(ctx context.Context, appointment *models.Appointment) error {
// 	if appointment.StartTime.Before(time.Now()) {
// 		return errors.New("start time must be in the future")
// 	}
// 	if appointment.StartTime.After(appointment.EndTime) {
// 		return errors.New("start time must be before end time")
// 	}
// 	return s.Repo.UpdateAppointment(ctx, appointment)
// }

// func (s *AppointmentService) DeleteAppointment(ctx context.Context) error {
// 	return s.Repo.DeleteAppointment(ctx)
// }
