package repositories

import (
	"context"
	"errors"

	"github.com/Bualoi-s-Dev/backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PackageRepository struct {
	Collection *mongo.Collection
}

func NewPackageRepository(collection *mongo.Collection) *PackageRepository {
	return &PackageRepository{Collection: collection}
}

func (repo *PackageRepository) GetAll(ctx context.Context) ([]models.Package, error) {
	var items []models.Package
	cursor, err := repo.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &items); err != nil {
		return nil, err
	}

	if items == nil {
		items = []models.Package{}
	}
	return items, nil
}

func (repo *PackageRepository) GetById(ctx context.Context, id string) (*models.Package, error) {
	var item models.Package
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	err = repo.Collection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (repo *PackageRepository) CreateOne(ctx context.Context, item *models.Package) (*mongo.InsertOneResult, error) {
	return repo.Collection.InsertOne(ctx, item)
}
func (repo *PackageRepository) UpdateOne(ctx context.Context, id string, updates bson.M) (*mongo.UpdateResult, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	// Check no empty update
	if len(updates) == 0 {
		return nil, errors.New("no fields to update")
	}

	update := bson.M{
		"$set": updates,
	}

	return repo.Collection.UpdateOne(ctx, bson.M{"_id": objectId}, update)
}
