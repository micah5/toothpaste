package toothpaste

import (
	"fmt"
	"math"
)

type Axis int

const (
	XAxis Axis = iota
	YAxis
	ZAxis
)

// 2D
type Vertex2D struct {
	X, Y float64
}

func (v *Vertex2D) Copy() *Vertex2D {
	return &Vertex2D{v.X, v.Y}
}

func (v *Vertex2D) Translate(x, y float64) {
	v.X += x
	v.Y += y
}

func (v *Vertex2D) Scale(x, y float64) {
	v.X *= x
	v.Y *= y
}

func (v *Vertex2D) Mul(m float64) {
	v.X *= m
	v.Y *= m
}

func (v *Vertex2D) Rotate(deg float64) {
	angle := float64(deg) * (math.Pi / 180)
	v.X = v.X*math.Cos(angle) - v.Y*math.Sin(angle)
	v.Y = v.X*math.Sin(angle) + v.Y*math.Cos(angle)
}

func (v *Vertex2D) MoveTo(x, y float64) {
	v.X = x
	v.Y = y
}

func (v *Vertex2D) String() string {
	return fmt.Sprintf("{%f, %f}", v.X, v.Y)
}

func (v *Vertex2D) Distance(v2 *Vertex2D) float64 {
	return math.Sqrt(math.Pow(v.X-v2.X, 2) + math.Pow(v.Y-v2.Y, 2))
}

func (v *Vertex2D) To3D(_axis ...Axis) *Vertex3D {
	var axis Axis
	if len(_axis) == 0 {
		axis = XAxis
	} else {
		axis = _axis[0]
	}
	switch axis {
	case XAxis:
		return &Vertex3D{v.X, 0, v.Y}
	case YAxis:
		return &Vertex3D{0, v.Y, v.X}
	case ZAxis:
		return &Vertex3D{v.X, v.Y, 0}
	default:
		return nil
	}
}

// 3D
type Vertex3D struct {
	X, Y, Z float64
}

func (v *Vertex3D) Copy() *Vertex3D {
	return &Vertex3D{v.X, v.Y, v.Z}
}

func (v *Vertex3D) Round(precision int) {
	multiplier := math.Pow(10, float64(precision))
	v.X = math.Round(v.X*multiplier) / multiplier
	v.Y = math.Round(v.Y*multiplier) / multiplier
	v.Z = math.Round(v.Z*multiplier) / multiplier
}

func (v *Vertex3D) Distance(v2 *Vertex3D) float64 {
	return math.Sqrt(math.Pow(v.X-v2.X, 2) + math.Pow(v.Y-v2.Y, 2) + math.Pow(v.Z-v2.Z, 2))
}

func (v *Vertex3D) Translate(x, y, z float64) {
	v.X += x
	v.Y += y
	v.Z += z
}

func (v *Vertex3D) Subtract(v2 *Vertex3D) {
	v.X -= v2.X
	v.Y -= v2.Y
	v.Z -= v2.Z
}

func (v *Vertex3D) Scale(x, y, z float64) {
	v.X *= x
	v.Y *= y
	v.Z *= z
}

func (v *Vertex3D) Mul(m float64) {
	v.X *= m
	v.Y *= m
	v.Z *= m
}

func (v *Vertex3D) Mirror(axis Axis) {
	switch axis {
	case XAxis:
		v.X = -v.X
	case YAxis:
		v.Y = -v.Y
	case ZAxis:
		v.Z = -v.Z
	}
}

func (v *Vertex3D) Normalize() Vertex3D {
	mag := math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
	return Vertex3D{
		X: v.X / mag,
		Y: v.Y / mag,
		Z: v.Z / mag,
	}
}

func (v *Vertex3D) Rotate(deg float64, axis Axis) {
	angle := float64(deg) * (math.Pi / 180)
	var newX, newY, newZ float64
	switch axis {
	case XAxis:
		newY = v.Y*math.Cos(angle) - v.Z*math.Sin(angle)
		newZ = v.Y*math.Sin(angle) + v.Z*math.Cos(angle)
		v.Y = newY
		v.Z = newZ
	case YAxis:
		newX = v.X*math.Cos(angle) + v.Z*math.Sin(angle)
		newZ = -v.X*math.Sin(angle) + v.Z*math.Cos(angle)
		v.X = newX
		v.Z = newZ
	case ZAxis:
		newX = v.X*math.Cos(angle) - v.Y*math.Sin(angle)
		newY = v.X*math.Sin(angle) + v.Y*math.Cos(angle)
		v.X = newX
		v.Y = newY
	}
}

func (v *Vertex3D) MoveTo(x, y, z float64) {
	v.X = x
	v.Y = y
	v.Z = z
}

func (v *Vertex3D) Angle(v2 *Vertex3D) float64 {
	return math.Atan2(v2.Y-v.Y, v2.X-v.X) * (180 / math.Pi)
}

func (v *Vertex3D) Cross(v2 *Vertex3D) float64 {
	return v.X*v2.Y - v.Y*v2.X
}

func (v *Vertex3D) String() string {
	return fmt.Sprintf("{%f, %f, %f}", v.X, v.Y, v.Z)
}

func (v *Vertex3D) Equals(v2 *Vertex3D) bool {
	return v.X == v2.X && v.Y == v2.Y && v.Z == v2.Z
}
