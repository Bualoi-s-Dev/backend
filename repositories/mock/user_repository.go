package repositories_mock

import (
	"context"

	"github.com/Bualoi-s-Dev/backend/models"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockUserRepository is a mock implementation of UserRepositoryInterface
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindUserByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.User), args.Error(1)
}

// func (m *MockUserRepository) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
// 	args := m.Called(ctx, email)
// 	return nil, args.Error(1)
// }

// func (m *MockUserRepository) FindEmailByID(ctx context.Context, id primitive.ObjectID) (string, error) {
// 	args := m.Called(ctx, id)
// 	return "", args.Error(1)
// }

// func (m *MockUserRepository) CreateUser(ctx context.Context, user *models.User) error {
// 	args := m.Called(ctx, user)
// 	return args.Error(1)
// }

// func (m *MockUserRepository) UpdateUser(ctx context.Context, userId primitive.ObjectID, updates bson.M) (*mongo.UpdateResult, error) {
// 	args := m.Called(ctx, userId, updates)
// 	return args.Get(0).(*mongo.UpdateResult), args.Error(1)
// }

// func (m *MockUserRepository) ReplaceUser(ctx context.Context, userId primitive.ObjectID, newUser *models.User) (*mongo.UpdateResult, error) {
// 	args := m.Called(ctx, userId, newUser)
// 	return args.Get(0).(*mongo.UpdateResult), args.Error(1)
// }

// func (m *MockUserRepository) FindPhotographers(ctx context.Context) ([]models.User, error) {
// 	args := m.Called(ctx)
// 	return args.Get(0).([]models.User), args.Error(1)
// }
