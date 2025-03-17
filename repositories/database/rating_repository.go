package repositories

import (
	"context"

	"github.com/Bualoi-s-Dev/backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RatingRepository struct {
	Collection *mongo.Collection
}

func NewRatingRepository(collection *mongo.Collection) *RatingRepository {
	return &RatingRepository{Collection: collection}
}

func (r *RatingRepository) GetAll(ctx context.Context) ([]models.Rating, error) {
	var items []models.Rating
	cursor, err := r.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &items); err != nil {
		return nil, err
	}

	if items == nil {
		items = []models.Rating{}
	}
	return items, nil
}

func (r *RatingRepository) GetById(ctx context.Context, ratingId primitive.ObjectID) (*models.Rating, error) {
	var item models.Rating

	err := r.Collection.FindOne(ctx, bson.M{"_id": ratingId}).Decode(&item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *RatingRepository) GetByPhotographerId(ctx context.Context, photographerId primitive.ObjectID) ([]models.Rating, error) {
	var items []models.Rating
	cursor, err := r.Collection.Find(ctx, bson.M{"photographer_id": photographerId})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &items); err != nil {
		return nil, err
	}

	if items == nil {
		items = []models.Rating{}
	}
	return items, nil
}

func (r *RatingRepository) GetByCustomerId(ctx context.Context, customerId primitive.ObjectID) ([]models.Rating, error) {
	var items []models.Rating
	cursor, err := r.Collection.Find(ctx, bson.M{"customer_id": customerId})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &items); err != nil {
		return nil, err
	}

	if items == nil {
		items = []models.Rating{}
	}
	return items, nil
}

func (r *RatingRepository) GetByPhotographerIdAndRating(ctx context.Context, photographerId primitive.ObjectID, rating int) ([]models.Rating, error) {
	var items []models.Rating
	cursor, err := r.Collection.Find(ctx, bson.M{
		"photographer_id": photographerId,
		"rating": rating,
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &items); err != nil {
		return nil, err
	}

	if items == nil {
		items = []models.Rating{}
	}
	return items, nil
}

func (r *RatingRepository) CreateOne(ctx context.Context, item *models.Rating) error {
	_, err := r.Collection.InsertOne(ctx, item)
	return err
}

func (r *RatingRepository) UpdateOne(ctx context.Context, item *models.Rating) error {
	_, err := r.Collection.UpdateOne(ctx, bson.M{"_id": item.ID}, bson.M{"$set": item})
	return err
}

func (r *RatingRepository) DeleteOne(ctx context.Context, ratingId primitive.ObjectID) error {
	_, err := r.Collection.DeleteOne(ctx, bson.M{"_id": ratingId})
	return err
}

func (r *RatingRepository) CustomerHasReviewedPhotographer(ctx context.Context, customerId, photographerId primitive.ObjectID) (bool, error) {
	filter := bson.M{
		"customer_id":     customerId,
		"photographer_id": photographerId,
	}

	count, err := r.Collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *RatingRepository) IsUserRatingOwner(ctx context.Context, customerId, ratingId primitive.ObjectID) (bool, error) {
	filter := bson.M{
		"customer_id":     customerId,
		"_id": ratingId,
	}

	count, err := r.Collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

