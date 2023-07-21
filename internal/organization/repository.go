package organization

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
)

type Repository interface {
	GetOrganization(ctx context.Context, id uint) (*Organization, error)
	CreateOrganization(ctx context.Context, organization *Organization) error
	DeleteOrganization(ctx context.Context, id uint) error
	UpdateOrganization(ctx context.Context, id uint, name string) error
}

type RepositoryImpl struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &RepositoryImpl{
		db: db,
	}
}

func (r *RepositoryImpl) GetOrganization(ctx context.Context, id uint) (*Organization, error) {
	var org Organization
	result := r.db.WithContext(ctx).
		First(&org, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}

		return nil, result.Error
	}

	return &org, nil
}

func (r *RepositoryImpl) CreateOrganization(ctx context.Context, organization *Organization) error {
	result := r.db.WithContext(ctx).
		Create(&organization)
	if result.Error != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(result.Error, &mysqlErr) && mysqlErr.Number == ErrCodeMySQLDuplicateEntry {
			return ErrDuplicateEntry
		}

		return result.Error
	}

	return nil
}

func (r *RepositoryImpl) DeleteOrganization(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).
		Delete(&Organization{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *RepositoryImpl) UpdateOrganization(ctx context.Context, id uint, name string) error {
	result := r.db.WithContext(ctx).
		Model(&Organization{}).
		Where("id = ?", id).
		Update("name", name)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
