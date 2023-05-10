package user

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
	GetUsers(ctx context.Context, ids []uint) ([]*User, error)
	GetUsersStats(ctx context.Context, ids []uint) ([]*UserStats, error)
	GetUsersInOrganization(ctx context.Context, organizationID uint) ([]User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	CreateUser(ctx context.Context, user *User) error
	DeleteUserByID(ctx context.Context, id uint) error
	UpdateUser(ctx context.Context, user *User) error
	UpdateRatings(ctx context.Context, userIDs []uint, ratings []int) error
}

type RepositoryImpl struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &RepositoryImpl{
		db: db,
	}
}

func (r *RepositoryImpl) GetUsers(ctx context.Context, ids []uint) ([]*User, error) {
	var users []*User
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (r *RepositoryImpl) GetUsersStats(ctx context.Context, ids []uint) ([]*UserStats, error) {
	var usersStats []*UserStats
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&usersStats).Error; err != nil {
		return nil, err
	}

	return usersStats, nil
}

func (r *RepositoryImpl) GetUsersInOrganization(ctx context.Context, organizationID uint) ([]User, error) {
	var users []User
	if err := r.db.WithContext(ctx).Where("organization_id = ?", organizationID).Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (r *RepositoryImpl) CreateUser(ctx context.Context, user *User) error {
	if err := r.db.WithContext(ctx).Create(&user).Omit("organization_id").Error; err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == ErrCodeMySQLDuplicateEntry {
			return ErrDuplicateEntry
		}

		return err
	}

	return nil
}

func (r *RepositoryImpl) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (r *RepositoryImpl) DeleteUserByID(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&User{}).Error; err != nil {
		return err
	}

	return nil
}

func (r *RepositoryImpl) UpdateRatings(ctx context.Context, userIDs []uint, ratings []int) error {
	tx := r.db.Begin()

	for i := range userIDs {
		if err := tx.WithContext(ctx).Model(&User{}).Where("id = ?", userIDs[i]).Update("rating", ratings[i]).Error; err != nil {
			tx.Rollback()

			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (r *RepositoryImpl) UpdateUser(ctx context.Context, user *User) error {
	if err := r.db.WithContext(ctx).Model(&User{}).Where("id = ?", user.ID).Updates(user).Error; err != nil {
		return err
	}

	return nil
}
