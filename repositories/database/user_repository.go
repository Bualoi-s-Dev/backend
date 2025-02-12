package repositories

import (
	"context"

	"github.com/Bualoi-s-Dev/backend/dto"
	"github.com/Bualoi-s-Dev/backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	Collection *mongo.Collection
}

func NewUserRepository(collection *mongo.Collection) *UserRepository {
	return &UserRepository{Collection: collection}
}

func (repo *UserRepository) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := repo.Collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *UserRepository) GetUserByEmail(ctx context.Context, email string) (*dto.UserResponse, error) {
	// var user models.User
	pipeline := []bson.M{
		{
			"$match": bson.M{"email": email}, // Filter user
		},
		{
			"$lookup": bson.M{
				"from":         "Package",               // Collection to join
				"localField":   "photographer_packages", // Field in "users"
				"foreignField": "_id",                   // Field in "packages"
				"as":           "package_details",       // Output array field
			},
		},
	}

	// Run aggregation
	cursor, err := repo.Collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx) // Ensure cursor is closed after function exits

	// Decode result
	if cursor.Next(ctx) {
		var res *dto.UserResponse
		err := cursor.Decode(&res)
		if err != nil {
			return nil, err
		}
		return res, nil
	}

	return nil, mongo.ErrNoDocuments
}

func (repo *UserRepository) GetUserByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	var user models.User
	err := repo.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	_, err := repo.Collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (repo *UserRepository) UpdateUser(ctx context.Context, email string, updates *models.User) (*models.User, error) {
	updateQuery := bson.M{"$set": updates}
	_, err := repo.Collection.UpdateOne(ctx, bson.M{"email": email}, updateQuery)
	if err != nil {
		return nil, err
	}
	return updates, nil
}
