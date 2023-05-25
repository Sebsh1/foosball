package statistic

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	GetStatisticByUserID(ctx context.Context, userID uint) (*Statistic, error)
	GetTopXUserIDsByWins(ctx context.Context, topX int) (userIDs []uint, wins []int, err error)
	GetTopXUserIDsByWinStreaks(ctx context.Context, topX int) (userIDs []uint, winStreaks []int, err error)
	GetTopXUserIDsByLossStreaks(ctx context.Context, topX int) (userIDs []uint, lossStreaks []int, err error)
	GetTopXUserIDsByWinLossRatios(ctx context.Context, topX int) (userIDs []uint, winLossRatios []float64, err error)
	GetTopXUserIDsByMatchesPlayed(ctx context.Context, topX int) (userIDs []uint, matchesPlayed []int, err error)
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

func (r *RepositoryImpl) GetTopXUserIDsByWins(ctx context.Context, topX int) ([]uint, []int, error) {
	var userIDs []uint
	var wins []int

	err := r.db.WithContext(ctx).Model(&Statistic{}).Order("wins desc").Limit(topX).Pluck("user_id", &userIDs).Pluck("wins", &wins).Error
	if err != nil {
		return nil, nil, err
	}

	return userIDs, wins, nil
}

func (r *RepositoryImpl) GetTopXUserIDsByWinStreaks(ctx context.Context, topX int) ([]uint, []int, error) {
	var userIDs []uint
	var winStreaks []int

	err := r.db.WithContext(ctx).Model(&Statistic{}).Order("win_streak desc").Limit(topX).Pluck("user_id", &userIDs).Pluck("win_streak", &winStreaks).Error
	if err != nil {
		return nil, nil, err
	}

	return userIDs, winStreaks, nil
}

func (r *RepositoryImpl) GetTopXUserIDsByLossStreaks(ctx context.Context, topX int) ([]uint, []int, error) {
	var userIDs []uint
	var lossStreaks []int

	err := r.db.WithContext(ctx).Model(&Statistic{}).Order("loss_streak desc").Limit(topX).Pluck("user_id", &userIDs).Pluck("loss_streak", &lossStreaks).Error
	if err != nil {
		return nil, nil, err
	}

	return userIDs, lossStreaks, nil
}

func (r *RepositoryImpl) GetTopXUserIDsByWinLossRatios(ctx context.Context, topX int) ([]uint, []float64, error) {
	var userIDs []uint
	var winLossRatios []float64

	err := r.db.WithContext(ctx).Model(&Statistic{}).Order("win_loss_ratio desc").Limit(topX).Pluck("user_id", &userIDs).Pluck("win_loss_ratio", &winLossRatios).Error
	if err != nil {
		return nil, nil, err
	}

	return userIDs, winLossRatios, nil
}

func (r *RepositoryImpl) GetTopXUserIDsByMatchesPlayed(ctx context.Context, topX int) ([]uint, []int, error) {
	var userIDs []uint
	var matchesPlayed []int

	err := r.db.WithContext(ctx).Model(&Statistic{}).Order("matches_played desc").Limit(topX).Pluck("user_id", &userIDs).Pluck("matches_played", &matchesPlayed).Error
	if err != nil {
		return nil, nil, err
	}

	return userIDs, matchesPlayed, nil
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
