package testing_runner

import (
	"context"
	"errors"
	"testing"

	"github.com/Bualoi-s-Dev/backend/models"
	repositories_mock "github.com/Bualoi-s-Dev/backend/repositories/mock"
	"github.com/Bualoi-s-Dev/backend/services"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestUnitTest(t *testing.T) {
	ctx := context.Background()
	userRepo := &repositories_mock.MockUserRepository{}
	// Pass nil for other repos, assuming they're not used in FilterPackage
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
				userRepo.On("FindUserByID", ctx, mockOwnerId).Return(&models.User{Name: "Dave Yok"}, nil)
			},
			expectedResult: false,
			expectedError:  nil,
		},
		{
			name:         "search string with user repo error",
			item:         &models.Package{Title: "Test Package", Type: models.PackageType("WEDDING_BLISS"), OwnerID: mockOwnerId},
			searchString: "xyz",
			searchType:   models.PackageType("WEDDING_BLISS"),
			setupMock: func() {
				userRepo.On("FindUserByID", ctx, mockOwnerId).Return((*models.User)(nil), errors.New("user not found"))
			},
			expectedResult: false,
			expectedError:  errors.New("user not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo.Calls = nil
			userRepo.ExpectedCalls = nil
			tt.setupMock()
			result, err := service.FilterPackage(ctx, tt.item, tt.searchString, tt.searchType)
			assert.Equal(t, tt.expectedResult, result)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
			userRepo.AssertExpectations(t) // Verify all expected mock calls were made
		})
	}
}
