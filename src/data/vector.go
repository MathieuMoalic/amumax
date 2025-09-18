package data

import "math"

// Vector represents a 3D vector.
type Vector [3]float64

func (v Vector) X() float64 { return v[0] }
func (v Vector) Y() float64 { return v[1] }
func (v Vector) Z() float64 { return v[2] }

// Mul returns a*v.
func (v Vector) Mul(a float64) Vector {
	return Vector{a * v[0], a * v[1], a * v[2]}
}

// Div returns (1/a)*v.
func (v Vector) Div(a float64) Vector {
	return v.Mul(1 / a)
}

// Add returns a+b.
func (v Vector) Add(b Vector) Vector {
	return Vector{v[0] + b[0], v[1] + b[1], v[2] + b[2]}
}

// MAdd Returns a+s*b.
func (v Vector) MAdd(s float64, b Vector) Vector {
	return Vector{v[0] + s*b[0], v[1] + s*b[1], v[2] + s*b[2]}
}

// Sub Returns a-b.
func (v Vector) Sub(b Vector) Vector {
	return Vector{v[0] - b[0], v[1] - b[1], v[2] - b[2]}
}

// Len Returns the norm of v.
func (v Vector) Len() float64 {
	len2 := v.Dot(v)
	return math.Sqrt(len2)
}

// Dot Returns the dot (inner) product a.b.
func (v Vector) Dot(b Vector) float64 {
	return v[0]*b[0] + v[1]*b[1] + v[2]*b[2]
}

// Cross Returns the cross (vector) product a x b
// in a right-handed coordinate system.
func (v Vector) Cross(b Vector) Vector {
	x := v[1]*b[2] - v[2]*b[1]
	y := v[2]*b[0] - v[0]*b[2]
	z := v[0]*b[1] - v[1]*b[0]
	return Vector{x, y, z}
}

const (
	X = 0
	Y = 1
	Z = 2
)
