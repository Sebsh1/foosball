package rating

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	GetRatingByUserID(ctx context.Context, userID uint) (*Rating, error)
	GetRatingsByUserIDs(ctx context.Context, userIDs []uint) ([]Rating, error)
	GetTopXAmongUserIDsByRating(ctx context.Context, topX int, userIDs []uint) (topXUserIDs []uint, ratings []int, err error)
	UpdateRatings(ctx context.Context, ratings []Rating) error
}

type RepositoryImpl struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &RepositoryImpl{db: db}
}

func (r *RepositoryImpl) GetRatingByUserID(ctx context.Context, userID uint) (*Rating, error) {
	var rating Rating

	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&rating).Error; err != nil {
		return nil, err
	}

	return &rating, nil
}

func (r *RepositoryImpl) GetRatingsByUserIDs(ctx context.Context, userIDs []uint) ([]Rating, error) {
	var ratings []Rating
	if err := r.db.WithContext(ctx).Where("user_id IN ?", userIDs).Find(&ratings).Error; err != nil {
		return nil, err
	}

	return ratings, nil
}

func (r *RepositoryImpl) GetTopXAmongUserIDsByRating(ctx context.Context, topX int, userIDs []uint) ([]uint, []int, error) {
	var topXUserIDs []uint
	var ratings []int

	err := r.db.WithContext(ctx).
		Model(&Rating{}).
		Order("rating desc").
		Limit(topX).
		Pluck("user_id", &topXUserIDs).
		Pluck("value", &ratings).
		Where("user_id IN ?", userIDs).Error
	if err != nil {
		return nil, nil, err
	}

	return topXUserIDs, ratings, nil
}

func (r *RepositoryImpl) UpdateRatings(ctx context.Context, ratings []Rating) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, rating := range ratings {
			if err := tx.WithContext(ctx).Model(&rating).Updates(rating).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
