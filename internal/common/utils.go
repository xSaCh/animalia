package common

import "math"


type Vector2D struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func (v Vector2D) IsZero() bool {
	return v.X == 0 && v.Y == 0
}
func (v Vector2D) SameAs(other Vector2D) bool {
	return v.X == other.X && v.Y == other.Y
}

func (v Vector2D) Add(other Vector2D) Vector2D {
	return Vector2D{v.X + other.X, v.Y + other.Y}
}

func (v Vector2D) Subtract(other Vector2D) Vector2D {
	return Vector2D{v.X - other.X, v.Y - other.Y}
}

func (v Vector2D) Dot(other Vector2D) float64 {
	return v.X*other.X + v.Y*other.Y
}

func (v Vector2D) Magnitude() float64 {
	return (v.X*v.X + v.Y*v.Y) // square of magnitude
}
func (v Vector2D) Length() float64 {
	return math.Sqrt(v.Magnitude())
}