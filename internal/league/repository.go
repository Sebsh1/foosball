package league

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	CreateLeague(ctx context.Context, league *League) error
	GetLeague(ctx context.Context, id uint) (*League, error)
	UpdateLeague(ctx context.Context, league *League) error
	DeleteLeague(ctx context.Context, id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) CreateLeague(ctx context.Context, league *League) error {
	return r.db.WithContext(ctx).Create(league).Error
}

func (r *repository) GetLeague(ctx context.Context, id uint) (*League, error) {
	var league League
	result := r.db.WithContext(ctx).First(&league, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &league, nil
}

func (r *repository) UpdateLeague(ctx context.Context, league *League) error {
	return r.db.WithContext(ctx).Save(league).Error
}

func (r *repository) DeleteLeague(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&League{}, id).Error
}
