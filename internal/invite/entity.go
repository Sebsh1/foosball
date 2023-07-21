package invite

import "time"

type Invite struct {
	ID uint `gorm:"primaryKey"`

	OrganizationID uint `gorm:"not null"`
	UserID         uint `gorm:"not null"`

	CreatedAt time.Time
}
