package main

import "math"

type Vector2D struct {
	x float64
	y float64
}

func (v Vector2D) Add(v2 Vector2D) Vector2D {
	return Vector2D{x: v.x + v2.x, y: v.y + v2.y}
}

func (v Vector2D) Subtract(v2 Vector2D) Vector2D {
	return Vector2D{x: v.x - v2.x, y: v.y - v2.y}
}

func (v Vector2D) Multiply(v2 Vector2D) Vector2D {
	return Vector2D{x: v.x * v2.x, y: v.y * v2.y}
}

func (v Vector2D) AddV(d float64) Vector2D {
	return Vector2D{x: v.x + d, y: v.y + d}
}

func (v Vector2D) MultiplyV(d float64) Vector2D {
	return Vector2D{x: v.x * d, y: v.y * d}
}

func (v Vector2D) DivideV(d float64) Vector2D {
	return Vector2D{x: v.x / d, y: v.y / d}
}

func (v Vector2D) limit(lower, upper float64) Vector2D {
	return Vector2D{x: math.Max(v.x, lower), y: math.Min(v.y, upper)}
}

func (v Vector2D) Distance(v2 Vector2D) float64 {
	return math.Sqrt(math.Pow(v.x-v2.x, 2) + math.Pow(v.y-v2.y, 2))
}
