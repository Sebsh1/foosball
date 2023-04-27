package invite

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	CreateInvite(ctx context.Context, userID, organizationID uint) error
	GetInvite(ctx context.Context, id uint) (Invite, error)
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

func (r *RepositoryImpl) CreateInvite(ctx context.Context, userID, organizationID uint) error {
	if err := r.db.WithContext(ctx).Create(&Invite{
		UserID:         userID,
		OrganizationID: organizationID,
	}).Error; err != nil {
		return err
	}

	return nil
}

func (r *RepositoryImpl) GetInvitesByUserID(ctx context.Context, userID uint) ([]Invite, error) {
	var invites []Invite
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&invites).Error; err != nil {
		return nil, err
	}

	return invites, nil
}

func (r *RepositoryImpl) GetInvitesByOrganizationID(ctx context.Context, organizationID uint) ([]Invite, error) {
	var invites []Invite
	if err := r.db.WithContext(ctx).Where("organization_id = ?", organizationID).Find(&invites).Error; err != nil {
		return nil, err
	}

	return invites, nil
}

func (r *RepositoryImpl) DeleteInvite(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&Invite{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (r *RepositoryImpl) GetInvite(ctx context.Context, id uint) (Invite, error) {
	var invite Invite
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&invite).Error; err != nil {
		return Invite{}, err
	}

	return invite, nil
}
