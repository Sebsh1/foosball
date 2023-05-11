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
	UpdateOrganization(ctx context.Context, id uint, name, ratingMethod string) error
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
	if err := r.db.WithContext(ctx).First(&org, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return &org, nil
}

func (r *RepositoryImpl) CreateOrganization(ctx context.Context, organization *Organization) error {
	if err := r.db.WithContext(ctx).Create(&organization).Error; err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == ErrCodeMySQLDuplicateEntry {
			return ErrDuplicateEntry
		}

		return err
	}

	return nil
}

func (r *RepositoryImpl) DeleteOrganization(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&Organization{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (r *RepositoryImpl) UpdateOrganization(ctx context.Context, id uint, name, ratingMethod string) error {
	if err := r.db.WithContext(ctx).Model(&Organization{}).Where("id = ?", id).Updates(map[string]interface{}{
		"name":          name,
		"rating_method": ratingMethod,
	}).Error; err != nil {
		return err
	}

	return nil
}
