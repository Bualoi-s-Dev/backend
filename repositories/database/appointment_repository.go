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
	BusyTimeCollectoin    *mongo.Collection
}

func NewAppointmentRepository(appointmentCollection, busyTimeCollection *mongo.Collection) *AppointmentRepository {
	return &AppointmentRepository{
		AppointmentCollection: appointmentCollection,
		BusyTimeCollectoin:    busyTimeCollection,
	}
}

func (repo *AppointmentRepository) AutoUpdateAppointmentStatus(ctx context.Context) error {

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

	// fmt.Println("Querytime = ", currentTime)
	// FIXME: Change this to match with BusyTime model
	// TODO: Testing
	go func() {
		pipeline := mongo.Pipeline{
			bson.D{
				{Key: "$lookup", Value: bson.D{
					{Key: "from", Value: repo.BusyTimeCollectoin},
					{Key: "localField", Value: "BusyTimeID"},
					{Key: "foreignField", Value: "_id"},
					{Key: "as", Value: "busyTimeDetails"},
				}},
			},
			bson.D{
				{Key: "$unwind", Value: "$busyTimeDetails"},
			},
			bson.D{
				{Key: "$match", Value: bson.D{
					{Key: "status", Value: "Pending"},
					{Key: "busyTimeDetails.start_time", Value: bson.D{{Key: "$lt", Value: currentTime}}},
				}},
			},
			bson.D{
				{Key: "$set", Value: bson.D{
					{Key: "status", Value: "Canceled"},
				}},
			},
		}

		_, err := repo.AppointmentCollection.Aggregate(ctx, pipeline)
		if err != nil {
			log.Println(err)
		}
		fmt.Println("====AutoUpdate Pending to Canceled====")

	}()

	// filter only end_time is less than current time and status is "Accepted"
	// TODO: Change after test
	go func() {
		filter := bson.M{
			"end_time": bson.M{"$lt": currentTime},
			"status":   "Accepted",
		}
		update := bson.M{"$set": bson.M{"status": "Completed"}}
		result, err := repo.AppointmentCollection.UpdateMany(ctx, filter, update)
		if err != nil {
			log.Println("Error updating documents:", err)
		} else {
			fmt.Printf("Autoupdated Accepted to Completed: %d documents\n", result.ModifiedCount)
		}
	}()

	return nil
}

func (repo *AppointmentRepository) GetAll(ctx context.Context, userID primitive.ObjectID, userRole models.UserRole) ([]models.Appointment, error) {
	var items []models.Appointment
	var fieldToFind string
	if userRole == models.Photographer {
		fieldToFind = "photographer_id"
	} else if userRole == models.Customer {
		fieldToFind = "customer_id"
	} else {
		return nil, fmt.Errorf("Guest cannot have appointments") // shouldn't have this error because authorization check
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

func (repo *AppointmentRepository) GetById(ctx context.Context, appointmentID primitive.ObjectID, userID primitive.ObjectID, userRole models.UserRole) (*models.Appointment, error) {
	var item models.Appointment

	if err := repo.AppointmentCollection.FindOne(ctx, bson.M{"_id": appointmentID}).Decode(&item); err != nil {
		return nil, err
	}

	return &item, nil
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
