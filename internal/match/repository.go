package match

import (
	"context"

	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const ErrCodeMySQLDuplicateEntry uint16 = 1062

var (
	ErrDuplicateEntry = errors.New("already exists")
	ErrNotFound       = errors.New("not found")
	ErrNotUpdated     = errors.New("not updated")
	ErrNotDeleted     = errors.New("not deleted")
)

type Repository interface {
	GetMatch(ctx context.Context, id uint) (*Match, error)
	GetMatchesWithPlayerID(ctx context.Context, id uint) ([]*Match, error)
	CreateMatch(ctx context.Context, match *Match) error
	DeleteMatch(ctx context.Context, match *Match) error
}

type RepositoryImpl struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &RepositoryImpl{
		db: db,
	}
}

func (r *RepositoryImpl) CreateMatch(ctx context.Context, match *Match) error {
	if err := r.db.WithContext(ctx).Create(&match).Error; err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == ErrCodeMySQLDuplicateEntry {
			return ErrDuplicateEntry
		}

		return err
	}

	return nil
}

func (r *RepositoryImpl) GetMatch(ctx context.Context, id uint) (*Match, error) {
	var match Match
	if err := r.db.WithContext(ctx).First(&match, id).Error; err != nil {
		return nil, err
	}

	return &match, nil
}

func (r *RepositoryImpl) DeleteMatch(ctx context.Context, match *Match) error {
	result := r.db.WithContext(ctx).
		Where(Match{ID: match.ID}).
		Model(&Match{}).
		Delete(match)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrNotDeleted
	}

	return nil
}

func (r *RepositoryImpl) GetMatchesWithPlayerID(ctx context.Context, id uint) ([]*Match, error) {
	var matchesAsTeamA []*Match
	var matchesAsTeamB []*Match
	tx := r.db.Begin()

	err := tx.WithContext(ctx).
		Model(&Match{}).
		Where("teamA = ?", id).
		Association("Players").
		Find(&matchesAsTeamA)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.WithContext(ctx).
		Model(&Match{}).
		Where("teamB = ?", id).
		Association("Players").
		Append(&matchesAsTeamB)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return append(matchesAsTeamA, matchesAsTeamB...), nil
}
