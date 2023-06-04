package toothpaste

import (
	"math"
)

func Square(w, d float64) *Face2D {
	return NewFace2D(
		0, 0,
		w, 0,
		w, d,
		0, d,
	)
}

func Circle(w, h float64, resolution int) *Face2D {
	var vertices []float64
	for i := 0; i < resolution; i++ {
		angle := float64(i) * (2 * math.Pi / float64(resolution))
		vertices = append(vertices, (w/2)*math.Cos(angle)+w/2, (h/2)*math.Sin(angle)+h/2)
	}
	return NewFace2D(vertices...)
}
