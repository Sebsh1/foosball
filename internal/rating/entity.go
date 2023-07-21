package rating

import "time"

type Rating struct {
	ID uint `gorm:"primaryKey"`

	UserID uint `gorm:"not null"`

	Value      float64 `gorm:"default:1000.0"`
	Deviation  float64
	Volatility float64

	CreatedAt time.Time
}
