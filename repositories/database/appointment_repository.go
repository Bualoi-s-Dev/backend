package repositories

import (
	"context"
	"errors"
	"strings"

	"github.com/Bualoi-s-Dev/backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AppointmentRepository struct {
	AppointmentCollection *mongo.Collection
	PackageCollection     *mongo.Collection
}

func NewAppointmentRepository(appointmentCollection, packageCollection *mongo.Collection) *AppointmentRepository {
	return &AppointmentRepository{
		AppointmentCollection: appointmentCollection,
		PackageCollection:     packageCollection,
	}
}

func (repo *AppointmentRepository) GetAll(ctx context.Context, userID primitive.ObjectID, userRole models.UserRole) ([]models.Appointment, error) {
	var items []models.Appointment
	fieldToFind := strings.ToLower(string(userRole)) + "_id" // photographer_id or customer_id
	cursor, err := repo.AppointmentCollection.Find(ctx, bson.M{
		fieldToFind: userID,
	})

	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &items); err != nil {
		return nil, err
	}

	if items == nil {
		items = []models.Appointment{}
	}
	return items, nil
}

func (repo *AppointmentRepository) GetById(ctx context.Context, appointmentID primitive.ObjectID, userID primitive.ObjectID, userRole models.UserRole) (*models.Appointment, error) {
	var item models.Appointment
	// FIXME: This is not correct rn

	if err := repo.AppointmentCollection.FindOne(ctx, bson.M{"_id": appointmentID}).Decode(&item); err != nil {
		return nil, err
	}

	if item.CustomerID != userID && item.PhotographerID != userID {
		return nil, errors.New("You are not authorized to view this appointment")
	}

	return &item, nil
}

func (repo *AppointmentRepository) FindSubpackageByID(ctx context.Context, id primitive.ObjectID) (*models.Subpackage, error) {
	var item models.Subpackage
	if err := repo.PackageCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&item); err != nil {
		return nil, err
	}
	return &item, nil
}

func (repo *AppointmentRepository) FindPackageByID(ctx context.Context, id primitive.ObjectID) (*models.Package, error) {
	var item models.Package
	if err := repo.PackageCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&item); err != nil {
		return nil, err
	}
	return &item, nil
}

func (repo *AppointmentRepository) CreateAppointment(ctx context.Context, appointment *models.Appointment) (*mongo.InsertOneResult, error) {
	// TODO: Available time checking
	item, err := repo.AppointmentCollection.InsertOne(ctx, appointment)
	return item, err
}

func (repo *AppointmentRepository) UpdateAppointment(ctx context.Context, id primitive.ObjectID, appointment *models.Appointment) error {
	// TODO: Available time checking before mapping
	_, err := repo.AppointmentCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": appointment})
	return err
}

func (repo *AppointmentRepository) DeleteAppointment(ctx context.Context, id primitive.ObjectID) error {
	_, err := repo.AppointmentCollection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
