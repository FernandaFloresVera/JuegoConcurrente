package models

type Duck struct {
	X         float64
	Y         float64
	VelocityX float64
	VelocityY float64
}

func NewDuck(x, y, vx, vy float64) *Duck {
	return &Duck{x, y, vx, vy}
}
