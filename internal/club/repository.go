package club

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
	GetClub(ctx context.Context, id uint) (*Club, error)
	GetClubs(ctx context.Context, ids []uint) ([]Club, error)
	GetUserIDsInClub(ctx context.Context, id uint) ([]uint, error)
	GetInvitesByUserID(ctx context.Context, userID uint) ([]ClubsUsers, error)
	CreateClub(ctx context.Context, Club *Club) (orgID uint, err error)
	AddUserToClub(ctx context.Context, userID uint, ClubID uint, role Role) error
	RemoveUserFromClub(ctx context.Context, userID uint, ClubID uint) error
	DeleteClub(ctx context.Context, id uint) error
	UpdateClub(ctx context.Context, id uint, name string) error
	UpdateUserRole(ctx context.Context, userID uint, ClubID uint, role Role) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetClub(ctx context.Context, id uint) (*Club, error) {
	var org Club
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

func (r *repository) GetClubs(ctx context.Context, ids []uint) ([]Club, error) {
	var orgs []Club
	result := r.db.WithContext(ctx).
		Find(&orgs, ids)
	if result.Error != nil {
		return nil, result.Error
	}

	return orgs, nil
}

func (r *repository) GetUserIDsInClub(ctx context.Context, id uint) ([]uint, error) {
	var orgUsers []ClubsUsers
	result := r.db.WithContext(ctx).
		Where("Club_id = ?", id).
		Find(&orgUsers)
	if result.Error != nil {
		return nil, result.Error
	}

	var userIDs []uint
	for _, orgUser := range orgUsers {
		userIDs = append(userIDs, orgUser.UserId)
	}

	return userIDs, nil
}

func (r *repository) GetInvitesByUserID(ctx context.Context, userID uint) ([]ClubsUsers, error) {
	var orgUsers []ClubsUsers
	result := r.db.WithContext(ctx).
		Where("user_id = ? AND accepted = ?", userID, false).
		Find(&orgUsers)
	if result.Error != nil {
		return nil, result.Error
	}

	return orgUsers, nil
}

func (r *repository) CreateClub(ctx context.Context, Club *Club) (uint, error) {
	result := r.db.WithContext(ctx).
		Create(&Club)
	if result.Error != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(result.Error, &mysqlErr) && mysqlErr.Number == ErrCodeMySQLDuplicateEntry {
			return 0, ErrDuplicateEntry
		}

		return 0, result.Error
	}

	return Club.Id, nil
}

func (r *repository) AddUserToClub(ctx context.Context, userID uint, ClubID uint, role Role) error {
	orgUser := &ClubsUsers{
		ClubId:   ClubID,
		UserId:   userID,
		Accepted: true,
		Role:     role,
	}

	result := r.db.WithContext(ctx).
		Create(&orgUser)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *repository) RemoveUserFromClub(ctx context.Context, userID uint, ClubID uint) error {
	result := r.db.WithContext(ctx).
		Where("user_id = ? AND Club_id = ?", userID, ClubID).
		Delete(&ClubsUsers{})
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *repository) DeleteClub(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).
		Delete(&Club{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *repository) UpdateClub(ctx context.Context, id uint, name string) error {
	result := r.db.WithContext(ctx).
		Model(&Club{}).
		Where("id = ?", id).
		Update("name", name)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *repository) UpdateUserRole(ctx context.Context, userID uint, ClubID uint, role Role) error {
	result := r.db.WithContext(ctx).
		Model(&ClubsUsers{}).
		Where("user_id = ? AND Club_id = ?", userID, ClubID).
		Update("role", role)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
