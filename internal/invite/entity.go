package invite

import "time"

type Invite struct {
	ID uint `gorm:"primaryKey"`

	OrganizationID uint
	UserID         uint

	CreatedAt time.Time
}
