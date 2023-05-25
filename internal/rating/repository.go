package rating

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	GetRatingByUserID(ctx context.Context, userID uint) (*Rating, error)
	GetRatingsByUserIDs(ctx context.Context, userIDs []uint) ([]Rating, error)
	GetTopXUserIDsByRating(ctx context.Context, topX int) (userIDs []uint, ratings []int, err error)
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
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&rating).Error
	if err != nil {
		return nil, err
	}

	return &rating, nil
}

func (r *RepositoryImpl) GetRatingsByUserIDs(ctx context.Context, userIDs []uint) ([]Rating, error) {
	var ratings []Rating
	err := r.db.WithContext(ctx).Where("user_id IN ?", userIDs).Find(&ratings).Error
	if err != nil {
		return nil, err
	}

	return ratings, nil
}

func (r *RepositoryImpl) GetTopXUserIDsByRating(ctx context.Context, topX int) ([]uint, []int, error) {
	var userIDs []uint
	var ratings []int

	err := r.db.WithContext(ctx).Model(&Rating{}).Order("rating desc").Limit(topX).Pluck("user_id", &userIDs).Pluck("rating", &ratings).Error
	if err != nil {
		return nil, nil, err
	}

	return userIDs, ratings, nil
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
