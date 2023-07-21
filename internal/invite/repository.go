package invite

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

var (
	ErrNotFound = errors.New("not found")
)

type Repository interface {
	CreateInvites(ctx context.Context, userIDs []uint, organizationID uint) error
	GetInvite(ctx context.Context, id uint) (*Invite, error)
	GetInvitesByUserID(ctx context.Context, userID uint) ([]Invite, error)
	GetInvitesByOrganizationID(ctx context.Context, organizationID uint) ([]Invite, error)
	DeleteInvite(ctx context.Context, id uint) error
}

type RepositoryImpl struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &RepositoryImpl{
		db: db,
	}
}

func (r *RepositoryImpl) CreateInvites(ctx context.Context, userIDs []uint, organizationID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, userID := range userIDs {
			result := tx.WithContext(ctx).
				Create(&Invite{
					UserID:         userID,
					OrganizationID: organizationID,
				})
			if result.Error != nil {
				tx.Rollback()
				return result.Error
			}
		}

		return nil
	})
}

func (r *RepositoryImpl) GetInvitesByUserID(ctx context.Context, userID uint) ([]Invite, error) {
	var invites []Invite
	result := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&invites)
	if result.Error != nil {
		return nil, result.Error
	}

	return invites, nil
}

func (r *RepositoryImpl) GetInvitesByOrganizationID(ctx context.Context, organizationID uint) ([]Invite, error) {
	var invites []Invite
	result := r.db.WithContext(ctx).
		Where("organization_id = ?", organizationID).
		Find(&invites)
	if result.Error != nil {
		return nil, result.Error
	}

	return invites, nil
}

func (r *RepositoryImpl) DeleteInvite(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).
		Delete(&Invite{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *RepositoryImpl) GetInvite(ctx context.Context, id uint) (*Invite, error) {
	var invite *Invite
	result := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&invite)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}

		return nil, result.Error
	}

	return invite, nil
}
