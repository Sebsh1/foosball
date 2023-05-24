package rating

import "time"

type Rating struct {
	ID uint `gorm:"primaryKey"`

	Rating     int `gorm:"default:1000"`
	Deviation  float64
	Volatility float64

	CreatedAt time.Time
}
