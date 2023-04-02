//go:generate mockgen --source=repository.go -destination=repository_mock.go -package=team
package team

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
	GetTeam(ctx context.Context, id uint) (*Team, error)
	CreateTeam(ctx context.Context, team *Team) error
	UpdateTeam(ctx context.Context, team *Team) error
	DeleteTeam(ctx context.Context, team *Team) error
}

type RepositoryImpl struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &RepositoryImpl{
		db: db,
	}
}

func (r *RepositoryImpl) GetTeam(ctx context.Context, id uint) (*Team, error) {
	var team Team

	result := r.db.WithContext(ctx).
		Where(Team{ID: id}).
		First(&team)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}

		return nil, result.Error
	}

	return &team, nil
}

func (r *RepositoryImpl) CreateTeam(ctx context.Context, team *Team) error {
	if err := r.db.WithContext(ctx).Create(&team).Error; err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == ErrCodeMySQLDuplicateEntry {
			return ErrDuplicateEntry
		}

		return err
	}

	return nil
}

func (r *RepositoryImpl) DeleteTeam(ctx context.Context, team *Team) error {
	result := r.db.WithContext(ctx).
		Where(Team{ID: team.ID}).
		Model(&Team{}).
		Delete(team)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrNotDeleted
	}

	return nil
}

func (r *RepositoryImpl) UpdateTeam(ctx context.Context, team *Team) error {
	result := r.db.WithContext(ctx).
		Where(Team{ID: team.ID}).
		Model(&Team{}).
		Updates(team)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrNotUpdated
	}

	return nil
}
