package game

type Game struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"not null"`
}
