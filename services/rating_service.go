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

func (s *RatingService) GetById(ctx context.Context, ratingId primitive.ObjectID) (*models.Rating, error) {
	return s.Repository.GetById(ctx, ratingId)
}

func (s *RatingService) CreateOneFromCustomer(ctx context.Context, request *dto.RatingRequest, customerId primitive.ObjectID, photographerId primitive.ObjectID) (*models.Rating, error) {
	model := request.ToModel(customerId, photographerId)
	// hasReviewed, err := s.Repository.CustomerHasReviewedPhotographer(ctx, customerId, model.PhotographerID)
	// if err != nil {
	// 	return nil, err
	// }
	// if hasReviewed {
	// 	return nil, apperrors.ErrAlreadyReviewed // Define this error to indicate duplicate review
	// }
	return nil, s.Repository.CreateOne(ctx, model)
}

func (s *RatingService) UpdateOne(ctx context.Context, userId,ratingId primitive.ObjectID, request *dto.RatingRequest) error {
	// Check if the rating exists and was written by this user
	isOwner, err := s.Repository.IsUserRatingOwner(ctx, userId, ratingId)
	if err != nil {
		return err
	}
	if !isOwner {
		return apperrors.ErrUnauthorized
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

func (s *RatingService) DeleteOne(ctx context.Context, userId,ratingId primitive.ObjectID) error {
	// Check if the rating exists and was written by this user
	isOwner, err := s.Repository.IsUserRatingOwner(ctx, userId, ratingId)
	if err != nil {
		return err
	}
	if !isOwner {
		return apperrors.ErrUnauthorized
	}

	// Proceed with deletion
	return s.Repository.DeleteOne(ctx, ratingId)
}

func (s *RatingService) MappedToRatingResponse(ctx context.Context, item *models.Rating) (*dto.RatingResponse, error) {
	ratings, err := s.RatingService.GetById(ctx, item.ID)
	if err != nil {
		return nil, err
	}
	return &dto.PackageResponse{
		ID:				item.ID,
		CustomerID:		item.CustomerID,
		PhotographerID:	item.PhotographerID,
		Rating:			item.Rating,
		Review:			item.Review,
	}, nil

}