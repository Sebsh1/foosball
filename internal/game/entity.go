package game

type Game struct {
	Id   uint   `gorm:"primaryKey"`
	Name string `gorm:"not null"`
}
