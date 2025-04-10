package repositories

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Bualoi-s-Dev/backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AppointmentRepository struct {
	AppointmentCollection *mongo.Collection
	BusyTimeCollection    *mongo.Collection
}

func NewAppointmentRepository(appointmentCollection, busyTimeCollection *mongo.Collection) *AppointmentRepository {
	return &AppointmentRepository{
		AppointmentCollection: appointmentCollection,
		BusyTimeCollection:    busyTimeCollection,
	}
}

func (repo *AppointmentRepository) UpdateCanceledAppointment(ctx context.Context, currentTime time.Time) ([]primitive.ObjectID, error) {
	pipeline := mongo.Pipeline{
		bson.D{
			{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: repo.BusyTimeCollection.Name()}, // Ensure this is a string
				{Key: "localField", Value: "busy_time_id"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "busy_time"},
			}},
		},
		bson.D{
			{Key: "$unwind", Value: "$busy_time"},
		},
		bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "status", Value: "Pending"},
				{Key: "busy_time.start_time", Value: bson.D{{Key: "$lt", Value: currentTime}}},
			}},
		},
	}

	cursor, err := repo.AppointmentCollection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Println("Aggregation error:", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var ids []primitive.ObjectID
	var count int
	for cursor.Next(ctx) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			log.Println("Error decoding document:", err)
			continue
		}
		count++
		if id, ok := result["_id"]; ok {
			if oid, ok := id.(primitive.ObjectID); ok {
				ids = append(ids, oid)
			} else {
				log.Println("Error: id is not of type primitive.ObjectID")
			}
		}
	}

	if err := cursor.Err(); err != nil {
		log.Println("Cursor error:", err)
		return nil, err
	}

	if count == 0 {
		fmt.Println("====No documents found to update from Pending to Canceled.====")
		return nil, nil
	}

	filter := bson.D{{Key: "_id", Value: bson.D{{Key: "$in", Value: ids}}}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "status", Value: "Canceled"}}}}
	result, err := repo.AppointmentCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		log.Println("Update error:", err)
		return nil, err
	}

	fmt.Printf("====AutoUpdate Pending to Canceled: %d documents updated====\n", result.ModifiedCount)
	return ids, nil
}

func (repo *AppointmentRepository) UpdateCompletedAppointment(ctx context.Context, currentTime time.Time) ([]primitive.ObjectID, error) {
	pipeline := mongo.Pipeline{
		bson.D{
			{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: repo.BusyTimeCollection.Name()}, // Ensure this is a string
				{Key: "localField", Value: "busy_time_id"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "busy_time"},
			}},
		},
		bson.D{
			{Key: "$unwind", Value: "$busy_time"},
		},
		bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "status", Value: "Accepted"},
				{Key: "busy_time.end_time", Value: bson.D{{Key: "$lt", Value: currentTime}}},
			}},
		},
	}

	cursor, err := repo.AppointmentCollection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Println("Aggregation error:", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var ids []primitive.ObjectID
	var count int
	for cursor.Next(ctx) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			log.Println("Error decoding document:", err)
			continue
		}
		count++
		if id, ok := result["_id"]; ok {
			if oid, ok := id.(primitive.ObjectID); ok {
				ids = append(ids, oid)
			} else {
				log.Println("Error: id is not of type primitive.ObjectID")
			}
		}
	}

	if err := cursor.Err(); err != nil {
		log.Println("Cursor error:", err)
		return nil, err
	}

	if count == 0 {
		fmt.Println("====No documents found to update from Accepted to Completed.====")
		return nil, nil
	}

	filter := bson.D{{Key: "_id", Value: bson.D{{Key: "$in", Value: ids}}}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "status", Value: "Completed"}}}}
	result, err := repo.AppointmentCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		log.Println("Update error:", err)
		return nil, err
	}

	fmt.Printf("====AutoUpdate Accepted to Completed: %d documents updated====\n", result.ModifiedCount)
	return ids, nil
}

func (repo *AppointmentRepository) GetAll(ctx context.Context, userID primitive.ObjectID, userRole models.UserRole) ([]models.Appointment, error) {
	var items []models.Appointment
	var fieldToFind string
	if userRole == models.Photographer {
		fieldToFind = "photographer_id"
	} else if userRole == models.Customer {
		fieldToFind = "customer_id"
	} else {
		return nil, fmt.Errorf("guest cannot have appointments") // shouldn't have this error because authorization check
	}
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

func (repo *AppointmentRepository) GetById(ctx context.Context, appointmentID primitive.ObjectID) (*models.Appointment, error) {
	var item models.Appointment

	if err := repo.AppointmentCollection.FindOne(ctx, bson.M{"_id": appointmentID}).Decode(&item); err != nil {
		return nil, err
	}

	return &item, nil
}

func (repo *AppointmentRepository) GetBySubpackageId(ctx context.Context, subpackageID primitive.ObjectID) ([]models.Appointment, error) {
	var items []models.Appointment

	cursor, err := repo.AppointmentCollection.Find(ctx, bson.M{"sub_package._id": subpackageID})
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

func (repo *AppointmentRepository) CreateAppointment(ctx context.Context, appointment *models.Appointment) (*models.Appointment, error) {
	_, err := repo.AppointmentCollection.InsertOne(ctx, appointment)
	return appointment, err
}

func (repo *AppointmentRepository) ReplaceAppointment(ctx context.Context, appointment *models.Appointment) (*models.Appointment, error) {
	_, err := repo.AppointmentCollection.ReplaceOne(ctx, bson.M{"_id": appointment.ID}, appointment)
	return appointment, err
}

func (repo *AppointmentRepository) DeleteAppointment(ctx context.Context, id primitive.ObjectID) error {
	_, err := repo.AppointmentCollection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
