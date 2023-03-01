package player

import (
	"context"
	"foosball/internal/models"

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
	GetPlayer(ctx context.Context, id uint) (*models.Player, error)
	GetPlayers(ctx context.Context, ids []uint) (*[]models.Player, error)
	CreatePlayer(ctx context.Context, player *models.Player) error
	DeletePlayer(ctx context.Context, player *models.Player) error
	UpdatePlayers(ctx context.Context, playesr []*models.Player) error
	GetTopPlayersByRating(ctx context.Context, topX int) ([]*models.Player, error)
}

type RepositoryImpl struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &RepositoryImpl{
		db: db,
	}
}

func (r *RepositoryImpl) GetPlayer(ctx context.Context, id uint) (*models.Player, error) {
	var player models.Player

	result := r.db.WithContext(ctx).
		Where(models.Player{ID: id}).
		First(&player)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}

		return nil, result.Error
	}

	return &player, nil
}

func (r *RepositoryImpl) GetPlayers(ctx context.Context, ids []uint) (*[]models.Player, error) {
	var players []models.Player

	result := r.db.WithContext(ctx).
		Find(&players, ids)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}

		return nil, result.Error
	}

	return &players, nil
}

func (r *RepositoryImpl) CreatePlayer(ctx context.Context, player *models.Player) error {
	if err := r.db.WithContext(ctx).Create(&player).Error; err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == ErrCodeMySQLDuplicateEntry {
			return ErrDuplicateEntry
		}

		return err
	}

	return nil
}

func (r *RepositoryImpl) DeletePlayer(ctx context.Context, player *models.Player) error {
	result := r.db.WithContext(ctx).
		Where(models.Player{ID: player.ID}).
		Model(&models.Player{}).
		Delete(player)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrNotDeleted
	}

	return nil
}

func (r *RepositoryImpl) UpdatePlayers(ctx context.Context, players []*models.Player) error {
	tx := r.db.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return err
	}

	for _, p := range players {
		result := tx.WithContext(ctx).
			Where(models.Player{ID: p.ID}).
			Model(&models.Player{}).
			Select("rating").
			Updates(p)
		if result.Error != nil {
			tx.Rollback()
			return result.Error
		}
	}

	if tx.RowsAffected == 0 {
		tx.Rollback()
		return ErrNotUpdated
	}

	tx.Commit()

	return nil
}

func (r *RepositoryImpl) GetTopPlayersByRating(ctx context.Context, topX int) ([]*models.Player, error) {
	var players []*models.Player

	result := r.db.WithContext(ctx).
		Select(&players).
		Order("rating DESC").
		Limit(topX)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}

		return nil, result.Error
	}

	return players, nil
}
