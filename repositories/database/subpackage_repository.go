package repositories

import (
	"context"

	"github.com/Bualoi-s-Dev/backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SubpackageRepository struct {
	Collection *mongo.Collection
}

func NewSubpackageRepository(collection *mongo.Collection) *SubpackageRepository {
	return &SubpackageRepository{Collection: collection}
}

func (r *SubpackageRepository) GetAll(ctx context.Context) ([]models.Subpackage, error) {
	var items []models.Subpackage
	cursor, err := r.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &items); err != nil {
		return nil, err
	}

	if items == nil {
		items = []models.Subpackage{}
	}
	return items, nil
}

func (r *SubpackageRepository) GetById(ctx context.Context, id string) (*models.Subpackage, error) {
	var item models.Subpackage
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

func (r *SubpackageRepository) GetByPackageId(ctx context.Context, packageId string) ([]models.Subpackage, error) {
	var items []models.Subpackage
	oPackageId, err := primitive.ObjectIDFromHex(packageId)
	if err != nil {
		return nil, err
	}
	cursor, err := r.Collection.Find(ctx, bson.M{"package_id": oPackageId})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &items); err != nil {
		return nil, err
	}

	if items == nil {
		items = []models.Subpackage{}
	}
	return items, nil
}

func (r *SubpackageRepository) Create(ctx context.Context, item models.Subpackage) error {
	_, err := r.Collection.InsertOne(ctx, item)
	return err
}

func (r *SubpackageRepository) UpdateOne(ctx context.Context, id string, item bson.M) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.Collection.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": item})
	return err
}

func (r *SubpackageRepository) ReplaceOne(ctx context.Context, id string, item models.Subpackage) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.Collection.ReplaceOne(ctx, bson.M{"_id": oid}, item)
	return err
}

func (r *SubpackageRepository) DeleteOne(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.Collection.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}
