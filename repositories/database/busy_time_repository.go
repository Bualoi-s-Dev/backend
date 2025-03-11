package repositories

import (
	"context"

	"github.com/Bualoi-s-Dev/backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BusyTimeRepository struct {
	Collection *mongo.Collection
}

func NewBusyTimeRepository(collection *mongo.Collection) *BusyTimeRepository {
	return &BusyTimeRepository{Collection: collection}
}

func (r *BusyTimeRepository) GetAll(ctx context.Context) ([]models.BusyTime, error) {
	var items []models.BusyTime
	cursor, err := r.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &items); err != nil {
		return nil, err
	}

	if items == nil {
		items = []models.BusyTime{}
	}
	return items, nil
}

func (r *BusyTimeRepository) GetAllValid(ctx context.Context) ([]models.BusyTime, error) {
	var items []models.BusyTime
	cursor, err := r.Collection.Find(ctx, bson.M{
		"is_valid": true,
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &items); err != nil {
		return nil, err
	}

	if items == nil {
		items = []models.BusyTime{}
	}
	return items, nil
}

func (r *BusyTimeRepository) GetById(ctx context.Context, id string) (*models.BusyTime, error) {
	var item models.BusyTime
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	err = r.Collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *BusyTimeRepository) GetByPhotographerId(ctx context.Context, photographerId primitive.ObjectID) ([]models.BusyTime, error) {
	var items []models.BusyTime
	cursor, err := r.Collection.Find(ctx, bson.M{"photographer_id": photographerId})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &items); err != nil {
		return nil, err
	}

	if items == nil {
		items = []models.BusyTime{}
	}
	return items, nil
}

func (r *BusyTimeRepository) GetByPhotographerIdValid(ctx context.Context, photographerId primitive.ObjectID) ([]models.BusyTime, error) {
	var items []models.BusyTime
	cursor, err := r.Collection.Find(ctx, bson.M{
		"photographer_id": photographerId,
		"is_valid":        true,
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &items); err != nil {
		return nil, err
	}

	if items == nil {
		items = []models.BusyTime{}
	}
	return items, nil
}

func (r *BusyTimeRepository) Create(ctx context.Context, item *models.BusyTime) error {
	_, err := r.Collection.InsertOne(ctx, item)
	return err
}

func (r *BusyTimeRepository) DeleteOne(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.Collection.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}
