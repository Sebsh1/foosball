package rating

import "time"

const (
	maxDeviation = 2.015
	minDeviation = 0.175

	maxRating   = 3 * maxDeviation
	startRating = 0.0
	minRating   = -3 * maxDeviation

	maxVolatility   = 0.08
	startVolatility = 0.06
	minVolatility   = 0.04
	tau             = 0.5

	resultMultiplierWin  = 1.0
	resultMultiplierDraw = 0.5
	resultMultiplierLoss = 0.0
)

type Rating struct {
	ID uint `gorm:"primaryKey"`

	UserID         uint `gorm:"not null"`
	OrganizationID uint `gorm:"not null"`

	Value      float64 `gorm:"default:1000.0"`
	Deviation  float64
	Volatility float64 `gorm:"default:0.06"`

	CreatedAt time.Time
}
