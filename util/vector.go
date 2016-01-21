package util

import "math"

type Vector2 struct {
	X, Y float32
}

func (v Vector2) Distance(to Vector2) float32 {
	return float32(math.Sqrt(float64((to.X-v.X)*(to.X-v.X) + (to.Y-v.Y)*(to.Y-v.Y))))
}

func (v Vector2) Vector3() Vector3 {
	return Vector3{
		X: v.X,
		Y: v.Y,
	}
}

type Vector3 struct {
	X, Y, Z float32
}

func (v Vector3) Distance(to Vector3) float32 {
	return float32(math.Sqrt(float64((to.X-v.X)*(to.X-v.X) + (to.Y-v.Y)*(to.Y-v.Y) + (to.Z-v.Z)*(to.Z-v.Z))))
}
