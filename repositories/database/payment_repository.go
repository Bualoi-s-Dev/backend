package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/Bualoi-s-Dev/backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PaymentRepository struct {
	Collection            *mongo.Collection
	AppointmentCollection *mongo.Collection
}

func NewPaymentRepository(collection *mongo.Collection, appointmentCollection *mongo.Collection) *PaymentRepository {
	return &PaymentRepository{Collection: collection, AppointmentCollection: appointmentCollection}
}

func (repo *PaymentRepository) Create(ctx context.Context, payment *models.Payment) error {
	_, err := repo.Collection.InsertOne(ctx, payment)
	return err
}

func (repo *PaymentRepository) GetAll(ctx context.Context) ([]models.Payment, error) {
	var items []models.Payment
	cursor, err := repo.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &items); err != nil {
		return nil, err
	}
	if items == nil {
		items = []models.Payment{}
	}
	return items, nil
}

func (repo *PaymentRepository) GetById(ctx context.Context, id primitive.ObjectID) (*models.Payment, error) {
	var item models.Payment
	err := repo.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (repo *PaymentRepository) GetByAppointmentID(ctx context.Context, appointmentID primitive.ObjectID) (*models.Payment, error) {
	var item models.Payment
	err := repo.Collection.FindOne(ctx, bson.M{"appointment_id": appointmentID}).Decode(&item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (repo *PaymentRepository) GetByUserIDAndRole(ctx context.Context, role models.UserRole, userId primitive.ObjectID) ([]models.Payment, error) {
	var items []models.Payment
	var fieldToFind string
	if role == models.Photographer {
		fieldToFind = "photographer_id"
	} else if role == models.Customer {
		fieldToFind = "customer_id"
	} else {
		return nil, errors.New("guest cannot have payments")
	}
	// get all payments that have customer id match with the given customer id from the appointment collection
	pipeline := mongo.Pipeline{
		bson.D{
			{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: repo.AppointmentCollection.Name()},
				{Key: "localField", Value: "appointment_id"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "appointment"},
			}},
		},
		bson.D{
			{Key: "$unwind", Value: "$appointment"},
		},
		bson.D{
			{Key: "$match", Value: bson.D{
				{Key: fmt.Sprintf("appointment.%s", fieldToFind), Value: userId},
			}},
		},
	}

	cursor, err := repo.Collection.Aggregate(ctx, pipeline)

	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &items); err != nil {
		return nil, err
	}

	if items == nil {
		items = []models.Payment{}
	}
	return items, nil
}

func (repo *PaymentRepository) GetByCheckoutID(ctx context.Context, checkoutID string) (*models.Payment, error) {
	var item models.Payment
	err := repo.Collection.FindOne(ctx, bson.M{"customer.checkout_id": checkoutID}).Decode(&item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (repo *PaymentRepository) GetByBalanceTransactionID(ctx context.Context, balanceTransactionID string) (*models.Payment, error) {
	var item models.Payment
	err := repo.Collection.FindOne(ctx, bson.M{"photographer.balance_transaction_id": balanceTransactionID}).Decode(&item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (repo *PaymentRepository) GetByPaymentIntentID(ctx context.Context, paymentIntentID string) (*models.Payment, error) {
	var item models.Payment
	err := repo.Collection.FindOne(ctx, bson.M{"customer.payment_intent_id": paymentIntentID}).Decode(&item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (repo *PaymentRepository) Replace(ctx context.Context, id primitive.ObjectID, payment *models.Payment) error {
	_, err := repo.Collection.ReplaceOne(ctx, bson.M{"_id": id}, payment)
	return err
}

func (repo *PaymentRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := repo.Collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (repo *PaymentRepository) UpdateCustomerPayment(ctx context.Context, id string, customerPayment *models.CustomerPayment) error {
	_, err := repo.Collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"customer": customerPayment}})
	return err
}

func (repo *PaymentRepository) UpdatePhotographerPayment(ctx context.Context, id string, photographerPayment *models.PhotographerPayment) error {
	_, err := repo.Collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"photographer": photographerPayment}})
	return err
}
