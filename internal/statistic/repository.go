package statistic

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	GetStatisticByUserID(ctx context.Context, userID uint) (*Statistic, error)
	UpdateStatistics(ctx context.Context, stats []Statistic) error
}

type RepositoryImpl struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &RepositoryImpl{
		db: db,
	}
}

func (r *RepositoryImpl) GetStatisticByUserID(ctx context.Context, userID uint) (*Statistic, error) {
	var stats Statistic
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&stats).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}

func (r *RepositoryImpl) UpdateStatistics(ctx context.Context, stats []Statistic) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, stat := range stats {
			if err := tx.WithContext(ctx).Model(&stat).Updates(stat).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
