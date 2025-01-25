package repositories

import (
	"context"

	"github.com/Bualoi-s-Dev/backend/models"
	"go.mongodb.org/mongo-driver/bson"
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

func (repo *PackageRepository) CreateOne(ctx context.Context, item *models.Package) (*mongo.InsertOneResult, error) {
	return repo.Collection.InsertOne(ctx, item)
}
