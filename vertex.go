package toothpaste

import (
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

func (v *Vertex2D) Rotate(deg int, axis Axis) {
	angle := float64(deg) * (math.Pi / 180)
	switch axis {
	case XAxis:
		v.Y = v.Y*math.Cos(angle) - v.Z*math.Sin(angle)
		v.Z = v.Y*math.Sin(angle) + v.Z*math.Cos(angle)
	case YAxis:
		v.X = v.X*math.Cos(angle) + v.Z*math.Sin(angle)
		v.Z = -v.X*math.Sin(angle) + v.Z*math.Cos(angle)
	}
}

func (v *Vertex2D) To3D(z float64) *Vertex3D {
	return &Vertex3D{v.X, z, v.Y}
}

// 3D
type Vertex3D struct {
	X, Y, Z float64
}

func (v *Vertex3D) Translate(x, y, z float64) {
	v.X += x
	v.Y += y
	v.Z += z
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

func (v *Vertex3D) Rotate(deg int, axis Axis) {
	angle := float64(deg) * (math.Pi / 180)
	switch axis {
	case XAxis:
		v.Y = v.Y*math.Cos(angle) - v.Z*math.Sin(angle)
		v.Z = v.Y*math.Sin(angle) + v.Z*math.Cos(angle)
	case YAxis:
		v.X = v.X*math.Cos(angle) + v.Z*math.Sin(angle)
		v.Z = -v.X*math.Sin(angle) + v.Z*math.Cos(angle)
	case ZAxis:
		v.X = v.X*math.Cos(angle) - v.Y*math.Sin(angle)
		v.Y = v.X*math.Sin(angle) + v.Y*math.Cos(angle)
	}
}
