package services

import (
	"context"

	"github.com/Bualoi-s-Dev/backend/apperrors"
	"github.com/Bualoi-s-Dev/backend/dto"
	"github.com/Bualoi-s-Dev/backend/models"
	repositories "github.com/Bualoi-s-Dev/backend/repositories/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RatingService struct {
	Repository	*repositories.RatingRepository
}

func NewRatingService(repository *repositories.RatingRepository) *RatingService {
	return &RatingService{
		Repository:	repository,
	}
}

func (s *RatingService) GetAll(ctx context.Context) ([]models.Rating, error) {
	return s.Repository.GetAll(ctx)
}

func (s *RatingService) GetByPhotographerId(ctx context.Context, photographerId primitive.ObjectID) ([]models.Rating, error) {
	return s.Repository.GetByPhotographerId(ctx, photographerId)
}

func (s *RatingService) GetById(ctx context.Context, photographerId, ratingId primitive.ObjectID) (*models.Rating, error) {
	var item *models.Rating
	isPhotographerOwner, err := s.Repository.IsPhotographerRatingOwner(ctx, photographerId, ratingId)
	if err != nil {
		return item, err
	}
	if !isPhotographerOwner {
		return item, apperrors.ErrPhotographerRatingMismatched
	}

	return s.Repository.GetById(ctx, ratingId)
}

func (s *RatingService) CreateOneFromCustomer(ctx context.Context, request *dto.RatingRequest, customerId primitive.ObjectID, photographerId primitive.ObjectID) error {
	model := request.ToModel(customerId, photographerId)
	// hasReviewed, err := s.Repository.CustomerHasReviewedPhotographer(ctx, customerId, model.PhotographerID)
	// if err != nil {
	// 	return nil, err
	// }
	// if hasReviewed {
	// 	return nil, apperrors.ErrAlreadyReviewed // Define this error to indicate duplicate review
	// }
	return s.Repository.CreateOne(ctx, model)
}

func (s *RatingService) UpdateOne(ctx context.Context, customerId, photographerId, ratingId primitive.ObjectID, request *dto.RatingRequest) error {
	if err := s.IsOwner(ctx, customerId, photographerId, ratingId); err != nil {
		return err
	}

	// Fetch the existing rating
	existingRating, err := s.Repository.GetById(ctx, ratingId)
	if err != nil {
		return err
	}

	// Update fields
	existingRating.Rating = request.Rating
	if request.Review != nil {
		existingRating.Review = *request.Review
	}

	return s.Repository.UpdateOne(ctx, existingRating)
}

func (s *RatingService) DeleteOne(ctx context.Context, customerId, photographerId, ratingId primitive.ObjectID) error {
	if err := s.IsOwner(ctx, customerId, photographerId, ratingId); err != nil {
		return err
	}

	// Proceed with deletion
	return s.Repository.DeleteOne(ctx, ratingId)
}

func (s *RatingService) MappedToRatingResponse(ctx context.Context, item *models.Rating) (*dto.RatingResponse, error) {
	return &dto.RatingResponse{
		ID:				item.ID,
		CustomerID:		item.CustomerID,
		PhotographerID:	item.PhotographerID,
		Rating:			item.Rating,
		Review:			item.Review,
	}, nil

}

func (s *RatingService) IsOwner(ctx context.Context, customerId, photographerId, ratingId primitive.ObjectID,) error {
	// Check if the rating exists and was written by this user
	isCustomerOwner, err := s.Repository.IsCustomerRatingOwner(ctx, customerId, ratingId)
	if err != nil {
		return err
	}
	if !isCustomerOwner {
		return apperrors.ErrCustomerRatingMismatched
	}
	
	// Check if the rating exists and was written by this user
	isPhotographerOwner, err := s.Repository.IsPhotographerRatingOwner(ctx, photographerId, ratingId)
	if err != nil {
		return err
	}
	if !isPhotographerOwner {
		return apperrors.ErrPhotographerRatingMismatched
	}

	return nil

}