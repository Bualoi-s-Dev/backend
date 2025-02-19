package repositories

import (
	"context"

	"github.com/Bualoi-s-Dev/backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AppointmentRepository struct {
	Collection *mongo.Collection
}

func NewAppointmentRepository(collection *mongo.Collection) *AppointmentRepository {
	return &AppointmentRepository{Collection: collection}
}

func (repo *AppointmentRepository) GetAll(ctx context.Context, id primitive.ObjectID) ([]models.Appointment, error) {
	var items []models.Appointment
	cursor, err := repo.Collection.Find(ctx, bson.M{"_id": id})
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

func (repo *AppointmentRepository) GetById(ctx context.Context, id primitive.ObjectID) (*models.Appointment, error) {
	var item models.Appointment
	err := repo.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (repo *AppointmentRepository) CreateAppointment(ctx context.Context, appointment *models.Appointment) error {
	// TODO: Available time checking
	_, err := repo.Collection.InsertOne(ctx, appointment)
	return err
}

func (repo *AppointmentRepository) UpdateAppointment(ctx context.Context, id primitive.ObjectID, appointment *models.Appointment) error {
	// TODO: Available time checking before mapping
	_, err := repo.Collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": appointment})
	return err
}

func (repo *AppointmentRepository) DeleteAppointment(ctx context.Context, id primitive.ObjectID) error {
	_, err := repo.Collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
