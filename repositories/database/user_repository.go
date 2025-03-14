package repositories

import (
	"context"

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

func (repo *UserRepository) FindUserByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	var user models.User
	err := repo.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *UserRepository) FindEmailByID(ctx context.Context, id primitive.ObjectID) (string, error) {
	var user models.User
	err := repo.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return "", err
	}
	return user.Email, nil
}

func (repo *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	if user.ShowcasePackages == nil {
		user.ShowcasePackages = []primitive.ObjectID{}
	}
	_, err := repo.Collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (repo *UserRepository) UpdateUser(ctx context.Context, userId primitive.ObjectID, updates bson.M) (*mongo.UpdateResult, error) {
	updateQuery := bson.M{"$set": updates}

	res, err := repo.Collection.UpdateOne(ctx, bson.M{"_id": userId}, updateQuery)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (repo *UserRepository) ReplaceUser(ctx context.Context, userId primitive.ObjectID, newUser *models.User) (*mongo.UpdateResult, error) {
	res, err := repo.Collection.ReplaceOne(ctx, bson.M{"_id": userId}, newUser)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (repo *UserRepository) FindPhotographers(ctx context.Context) ([]models.User, error) {
	var users []models.User
	cursor, err := repo.Collection.Find(ctx, bson.M{"role": "Photographer"})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}
