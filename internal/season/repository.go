package season

import (
	"context"

	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const ErrCodeMySQLDuplicateEntry uint16 = 1062

var (
	ErrNotFound       = errors.New("not found")
	ErrDuplicateEntry = errors.New("already exists")
	ErrNotDeleted     = errors.New("not deleted")
	ErrNotUpdated     = errors.New("not updated")
)

type Repository interface {
	GetSeason(ctx context.Context, id uint) (*Season, error)
	CreateSeason(ctx context.Context, season *Season) error
	DeleteSeason(ctx context.Context, season *Season) error
	UpdateSeason(ctx context.Context, season *Season) error
}

type RepositoryImpl struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &RepositoryImpl{
		db: db,
	}
}

func (r *RepositoryImpl) GetSeason(ctx context.Context, id uint) (*Season, error) {
	var season Season

	result := r.db.WithContext(ctx).
		Where(Season{ID: id}).
		First(&season)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}

		return nil, result.Error
	}

	return &season, nil
}

func (r *RepositoryImpl) CreateSeason(ctx context.Context, season *Season) error {
	if err := r.db.WithContext(ctx).Create(&season).Error; err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == ErrCodeMySQLDuplicateEntry {
			return ErrDuplicateEntry
		}

		return err
	}

	return nil
}

func (r *RepositoryImpl) DeleteSeason(ctx context.Context, season *Season) error {
	result := r.db.WithContext(ctx).
		Where(Season{ID: season.ID}).
		Model(&Season{}).
		Delete(season)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrNotDeleted
	}

	return nil
}

func (r *RepositoryImpl) UpdateSeason(ctx context.Context, season *Season) error {
	result := r.db.WithContext(ctx).
		Where(Season{ID: season.ID}).
		Model(&Season{}).
		Updates(season)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrNotUpdated
	}

	return nil
}
