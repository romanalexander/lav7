package util

import "math"

// Vector2 is a X-Y vector, containing 2nd-dimension position.
type Vector2 struct {
	X, Y float32
}

// Distance calculates the distance between given vector.
func (v Vector2) Distance(to Vector2) float32 {
	return float32(math.Sqrt(float64((to.X-v.X)*(to.X-v.X) + (to.Y-v.Y)*(to.Y-v.Y))))
}

// Vector3 converts Vector2 to Vector3, setting X, Y as set before, but leaving Z zero.
func (v Vector2) Vector3() Vector3 {
	return Vector3{
		X: v.X,
		Y: v.Y,
	}
}

// Vector2 is a X-Y-Z vector, containing 3rd-dimension position.
type Vector3 struct {
	X, Y, Z float32
}

// Distance calculates the distance between given vector.
func (v Vector3) Distance(to Vector3) float32 {
	return float32(math.Sqrt(float64((to.X-v.X)*(to.X-v.X) + (to.Y-v.Y)*(to.Y-v.Y) + (to.Z-v.Z)*(to.Z-v.Z))))
}
