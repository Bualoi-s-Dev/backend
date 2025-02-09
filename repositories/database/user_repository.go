package repositories

import (
	"context"
	"fmt"

	"github.com/Bualoi-s-Dev/backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	Collection *mongo.Collection
}

func NewUserRepository(collection *mongo.Collection) *UserRepository {
	return &UserRepository{Collection: collection}
}

func (repo *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := repo.Collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
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

func (repo *UserRepository) UpdateUser(ctx context.Context, email string, updates *models.User) error {
	updateQuery := bson.M{"$set": updates}
	_, err := repo.Collection.UpdateOne(ctx, bson.M{"email": email}, updateQuery)
	if err != nil {
		return err
	}
	return nil
}

func (repo *UserRepository) UpdateUserField(ctx context.Context, email string, updates map[string]interface{}) error {
	updateQuery := bson.M{"$set": updates}
	_, err := repo.Collection.UpdateOne(ctx, bson.M{"email": email}, updateQuery)
	if err != nil {
		return err
	}
	return nil
}

func (repo *UserRepository) GetUserProfilePic(ctx context.Context, email string) (string, error) {
	user, err := repo.GetUserByEmail(ctx, email) 
	if err != nil {
		return "", err
	}

	if user.Profile == "" {
		return "", fmt.Errorf("no profile picture found for user")
	}

	return user.Profile, nil
}