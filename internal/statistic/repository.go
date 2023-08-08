package statistic

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	GetStatisticsByUserIDs(ctx context.Context, userIDs []uint) ([]*Statistic, error)
	GetStatisticByUserID(ctx context.Context, userID uint) (*Statistic, error)
	GetTopXAmongUserIDsByWins(ctx context.Context, topX int, userIDs []uint) (topXUserIDs []uint, values []int, err error)
	GetTopXAmongUserIDsByStreak(ctx context.Context, topX int, userIDs []uint) (topXUserIDs []uint, values []int, err error)
	CreateStatistic(ctx context.Context, userID uint) error
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

func (r *RepositoryImpl) GetStatisticsByUserIDs(ctx context.Context, userIDs []uint) ([]*Statistic, error) {
	var stats []*Statistic
	result := r.db.WithContext(ctx).
		Where("user_id IN ?", userIDs).
		Find(&stats)
	if result.Error != nil {
		return nil, result.Error
	}

	return stats, nil
}

func (r *RepositoryImpl) GetStatisticByUserID(ctx context.Context, userID uint) (*Statistic, error) {
	var stats Statistic
	result := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		First(&stats)
	if result.Error != nil {
		return nil, result.Error
	}

	return &stats, nil
}

func (r *RepositoryImpl) GetTopXAmongUserIDsByWins(ctx context.Context, topX int, userIDs []uint) ([]uint, []int, error) {
	var topXUserIDs []uint
	var wins []int

	result := r.db.
		WithContext(ctx).
		Model(&Statistic{}).
		Order("wins desc").
		Limit(topX).
		Pluck("user_id", &topXUserIDs).
		Pluck("wins", &wins).
		Where("user_id IN ?", userIDs)
	if result.Error != nil {
		return nil, nil, result.Error
	}

	return topXUserIDs, wins, nil
}

func (r *RepositoryImpl) GetTopXAmongUserIDsByStreak(ctx context.Context, topX int, userIDs []uint) ([]uint, []int, error) {
	var topXUserIDs []uint
	var streaks []int

	result := r.db.WithContext(ctx).
		Model(&Statistic{}).
		Order("streak desc").
		Limit(topX).
		Where("user_id IN ?", userIDs).
		Pluck("user_id", &topXUserIDs).
		Pluck("streak", &streaks)
	if result.Error != nil {
		return nil, nil, result.Error
	}

	return topXUserIDs, streaks, nil
}

func (r *RepositoryImpl) CreateStatistic(ctx context.Context, userID uint) error {
	stat := Statistic{
		UserID: userID,
		Wins:   0,
		Draws:  0,
		Losses: 0,
		Streak: 0,
	}

	result := r.db.WithContext(ctx).
		Create(&stat)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *RepositoryImpl) UpdateStatistics(ctx context.Context, stats []Statistic) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, stat := range stats {
			result := tx.WithContext(ctx).
				Model(&stat).
				Updates(stat)
			if result.Error != nil {
				return result.Error
			}
		}
		return nil
	})
}
