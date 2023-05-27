package statistic

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	GetStatisticByUserID(ctx context.Context, userID uint) (*Statistic, error)
	GetTopXAmongUserIDsByWins(ctx context.Context, topX int, userIDs []uint) (topXUserIDs []uint, wins []int, err error)
	GetTopXAmongUserIDsByWinStreaks(ctx context.Context, topX int, userIDs []uint) (topXUserIDs []uint, winStreaks []int, err error)
	GetTopXAmongUserIDsByLossStreaks(ctx context.Context, topX int, userIDs []uint) (topXUserIDs []uint, lossStreaks []int, err error)
	GetTopXAmongUserIDsByWinLossRatios(ctx context.Context, topX int, userIDs []uint) (topXUserIDs []uint, winLossRatios []float64, err error)
	GetTopXAmongUserIDsByMatchesPlayed(ctx context.Context, topX int, userIDs []uint) (topXUserIDs []uint, matchesPlayed []int, err error)
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
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		First(&stats).Error
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

func (r *RepositoryImpl) GetTopXAmongUserIDsByWins(ctx context.Context, topX int, userIDs []uint) ([]uint, []int, error) {
	var topXUserIDs []uint
	var wins []int

	err := r.db.
		WithContext(ctx).
		Model(&Statistic{}).
		Order("wins desc").
		Limit(topX).
		Pluck("user_id", &topXUserIDs).
		Pluck("wins", &wins).
		Where("user_id IN ?", userIDs).Error
	if err != nil {
		return nil, nil, err
	}

	return topXUserIDs, wins, nil
}

func (r *RepositoryImpl) GetTopXAmongUserIDsByWinStreaks(ctx context.Context, topX int, userIDs []uint) ([]uint, []int, error) {
	var topXUserIDs []uint
	var winStreaks []int

	err := r.db.WithContext(ctx).
		Model(&Statistic{}).
		Order("win_streak desc").
		Limit(topX).
		Pluck("user_id", &topXUserIDs).
		Pluck("win_streak", &winStreaks).
		Where("user_id IN ?", userIDs).Error
	if err != nil {
		return nil, nil, err
	}

	return topXUserIDs, winStreaks, nil
}

func (r *RepositoryImpl) GetTopXAmongUserIDsByLossStreaks(ctx context.Context, topX int, userIDs []uint) ([]uint, []int, error) {
	var topXUserIDs []uint
	var lossStreaks []int

	err := r.db.WithContext(ctx).
		Model(&Statistic{}).
		Order("loss_streak desc").
		Limit(topX).
		Pluck("user_id", &topXUserIDs).
		Pluck("loss_streak", &lossStreaks).
		Where("user_id IN ?", userIDs).Error
	if err != nil {
		return nil, nil, err
	}

	return topXUserIDs, lossStreaks, nil
}

func (r *RepositoryImpl) GetTopXAmongUserIDsByWinLossRatios(ctx context.Context, topX int, userIDs []uint) ([]uint, []float64, error) {
	var topXUserIDs []uint
	var winLossRatios []float64

	err := r.db.WithContext(ctx).
		Model(&Statistic{}).
		Order("win_loss_ratio desc").
		Limit(topX).
		Pluck("user_id", &topXUserIDs).
		Pluck("win_loss_ratio", &winLossRatios).
		Where("user_id IN ?", userIDs).Error
	if err != nil {
		return nil, nil, err
	}

	return topXUserIDs, winLossRatios, nil
}

func (r *RepositoryImpl) GetTopXAmongUserIDsByMatchesPlayed(ctx context.Context, topX int, userIDs []uint) ([]uint, []int, error) {
	var topXUserIDs []uint
	var matchesPlayed []int

	err := r.db.WithContext(ctx).
		Model(&Statistic{}).
		Order("matches_played desc").
		Limit(topX).
		Pluck("user_id", &topXUserIDs).
		Pluck("matches_played", &matchesPlayed).
		Where("user_id IN ?", userIDs).Error
	if err != nil {
		return nil, nil, err
	}

	return topXUserIDs, matchesPlayed, nil
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
