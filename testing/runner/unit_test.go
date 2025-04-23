package testing_runner

import (
	"context"
	"errors"
	"testing"

	"github.com/Bualoi-s-Dev/backend/models"
	"github.com/Bualoi-s-Dev/backend/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestUnitTest(t *testing.T) {
	ctx := context.Background()
	userRepo := &MockUserRepository{}
	service := services.NewPackageService(nil, nil, nil, userRepo)

	mockOwnerId, _ := primitive.ObjectIDFromHex("123")
	tests := []struct {
		name           string
		item           *models.Package
		searchString   string
		searchType     models.PackageType
		setupMock      func()
		expectedResult bool
		expectedError  error
	}{
		{
			name:           "no search string, no search type",
			item:           &models.Package{Title: "Test Package", Type: models.PackageType("WEDDING_BLISS"), OwnerID: mockOwnerId},
			searchString:   "",
			searchType:     "",
			setupMock:      func() {},
			expectedResult: true,
			expectedError:  nil,
		},
		{
			name:           "search type matches",
			item:           &models.Package{Title: "Test Package", Type: models.PackageType("WEDDING_BLISS"), OwnerID: mockOwnerId},
			searchString:   "",
			searchType:     models.PackageType("WEDDING_BLISS"),
			setupMock:      func() {},
			expectedResult: true,
			expectedError:  nil,
		},
		{
			name:           "search type does not match",
			item:           &models.Package{Title: "Test Package", Type: models.PackageType("WEDDING_BLISS"), OwnerID: mockOwnerId},
			searchString:   "",
			searchType:     models.PackageType("OTHER"),
			setupMock:      func() {},
			expectedResult: false,
			expectedError:  nil,
		},
		{
			name:           "search string matches title prefix",
			item:           &models.Package{Title: "Test Package", Type: models.PackageType("WEDDING_BLISS"), OwnerID: mockOwnerId},
			searchString:   "test",
			searchType:     "",
			setupMock:      func() {},
			expectedResult: true,
			expectedError:  nil,
		},
		{
			name:         "search string matches owner name prefix",
			item:         &models.Package{Title: "Test Package", Type: models.PackageType("WEDDING_BLISS"), OwnerID: mockOwnerId},
			searchString: "john",
			searchType:   "",
			setupMock: func() {
				userRepo.On("FindUserByID", ctx, mockOwnerId).Return(&models.User{Name: "John Doe"}, nil)
			},
			expectedResult: true,
			expectedError:  nil,
		},
		{
			name:         "search string matches neither title nor owner name",
			item:         &models.Package{Title: "Test Package", Type: models.PackageType("WEDDING_BLISS"), OwnerID: mockOwnerId},
			searchString: "xyz",
			searchType:   "",
			setupMock: func() {
				userRepo.On("FindUserByID", ctx, mockOwnerId).Return(&models.User{Name: "John Doe"}, nil)
			},
			expectedResult: false,
			expectedError:  nil,
		},
		{
			name:         "search string with user repo error",
			item:         &models.Package{Title: "Test Package", Type: models.PackageType("WEDDING_BLISS"), OwnerID: mockOwnerId},
			searchString: "xyz",
			searchType:   "",
			setupMock: func() {
				userRepo.On("FindUserByID", ctx, mockOwnerId).Return((*models.User)(nil), errors.New("user not found"))
			},
			expectedResult: false,
			expectedError:  errors.New("user not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo.ExpectedCalls = nil // Reset mock expectations
			tt.setupMock()
			result, err := service.FilterPackage(ctx, tt.item, tt.searchString, tt.searchType)
			assert.Equal(t, tt.expectedResult, result)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

// MockUserRepository is a mock implementation of the UserRepository interface
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindUserByID(ctx context.Context, id int) (*models.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	return nil, args.Error(1)
}

func (m *MockUserRepository) FindEmailByID(ctx context.Context, id primitive.ObjectID) (string, error) {
	args := m.Called(ctx, id)
	return "", args.Error(1)
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(1)
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, userId primitive.ObjectID, updates bson.M) (*mongo.UpdateResult, error) {
	args := m.Called(ctx, userId, updates)

	return args.Get(0).(*mongo.UpdateResult), args.Error(1)
}

func (m *MockUserRepository) ReplaceUser(ctx context.Context, userId primitive.ObjectID, newUser *models.User) (*mongo.UpdateResult, error) {
	args := m.Called(ctx, userId, newUser)
	return args.Get(0).(*mongo.UpdateResult), args.Error(1)
}

func (m *MockUserRepository) FindPhotographers(ctx context.Context) ([]models.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.User), args.Error(1)
}

func makeUserFromName(name string) *models.User {
	return &models.User{

		Name: name,
	}
}
