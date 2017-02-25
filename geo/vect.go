package geo

import "math"

type Vect struct {
	X, Y float64
}

func (v Vect) IsOrigin() bool {
	return v.X == 0 && v.Y == 0
}

func (v Vect) Add(o Vect) Vect {
	return Vect{v.X + o.X, v.Y + o.Y}
}

func (v Vect) Subtract(o Vect) Vect {
	return Vect{v.X - o.X, v.Y - o.Y}
}

func (v Vect) Mult(n float64) Vect {
	return Vect{v.X * n, v.Y * n}
}

func (v Vect) Normalized() Vect {
	l := math.Sqrt(v.X*v.X + v.Y*v.Y)
	return Vect{v.X / l, v.Y / l}
}

func (v Vect) Perpendicular() Vect {
	// Could be either (x, y) -> (-y, x) or (y, -x)
	return Vect{-v.Y, v.X}
}

func (v Vect) RotateByAngle(theta float64) Vect {
	return v.Rotate(HeadingFromAngle(theta))
}

// The heading MUST be normalized
func (v Vect) Rotate(heading Vect) Vect {
	return Vect{v.X*heading.X - v.Y*heading.Y, v.X*heading.Y + v.Y*heading.X}
}

func (v Vect) DistanceTo(o Vect) float64 {
	return o.Subtract(v).Magnitude()
}

func (v Vect) AngleTo(o Vect) float64 {
	return math.Atan2(o.Y-v.Y, o.X-v.X)
}

func HeadingFromAngle(theta float64) Vect {
	return Vect{math.Cos(theta), math.Sin(theta)}
}

func (v Vect) AsAngle() float64 {
	return math.Atan2(v.Y, v.X)
}

func (v Vect) Dot(o Vect) float64 {
	return v.X*o.X + v.Y*o.Y
}

func (v Vect) Magnitude() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (v Vect) AngleBetween(o Vect) float64 {
	return math.Acos(v.Dot(o) / (v.Magnitude() * o.Magnitude()))
}
