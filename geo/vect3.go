package geo

import "math"

type Vect3 struct {
	X, Y, Z float64
}

func (v Vect3) Dot(o Vect3) float64 {
	return v.X*o.X + v.Y*o.Y + v.Z*o.Z
}

func (v Vect3) Cross(o Vect3) Vect3 {
	return Vect3{
		v.Y*o.Z - v.Z*o.Y,
		v.Z*o.X - v.X*o.Z,
		v.X*o.Y - v.Y*o.X}
}

func (v Vect3) Magnitude() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

func (v Vect3) AngleBetween(o Vect3) float64 {
	return math.Acos(v.Dot(o) / (v.Magnitude() * o.Magnitude()))
}
