package repositories

import (
	"context"
	"fmt"
	"log"

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
				"as":           "photographer_packages", // Output array field
			},
		},
		{
			"$lookup": bson.M{
				"from":         "Package",           // Collection to join
				"localField":   "showcase_packages", // Field in "users"
				"foreignField": "_id",               // Field in "packages"
				"as":           "showcase_packages", // Output array field
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
		var res dto.UserResponse
		err := cursor.Decode(&res)
		if err != nil {
			return nil, err
		}
		return &res, nil
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

func (repo *UserRepository) UpdateUser(ctx context.Context, userId primitive.ObjectID, updates bson.M) (*mongo.UpdateResult, error) {
	updateQuery := bson.M{"$set": updates}
	fmt.Println("updateQuery :", updateQuery)
	fmt.Println("userId :", userId)
	// count, err := repo.Collection.CountDocuments(ctx, bson.M{"_id": userId})
	// if err != nil {
	// 	log.Println("Error checking document existence:", err)
	// } else if count == 0 {
	// 	log.Println("No document found with _id:", userId)
	// }

	res, err := repo.Collection.UpdateOne(ctx, bson.M{"_id": userId}, updateQuery)
	if err != nil {
		return nil, err
	}
	fmt.Println("res :", res)
	var updatedUser bson.M
	err = repo.Collection.FindOne(ctx, bson.M{"_id": userId}).Decode(&updatedUser)
	if err != nil {
		log.Println("Error fetching updated document:", err)
	} else {
		log.Println("Updated document:", updatedUser)
	}

	return res, nil
}
