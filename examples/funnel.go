package main

import (
	"github.com/micah5/toothpaste"
	"math"
)

func main() {
	w, d, h := 0.5, 0.5, 0.5
	funnel_perc := 0.8
	outer := toothpaste.Circle(w, d, 10)
	inner := toothpaste.Circle(w*funnel_perc, d*funnel_perc, 10)
	inner.Translate(w*(1-funnel_perc)/2, d*(1-funnel_perc)/2)
	node := toothpaste.NewNode(outer.To3D(), inner.To3D())

	num_iterations := 4
	for i := 0; i < num_iterations; i++ {
		multiplier := math.Pow(1.34, float64(i))

		// Handle outer circle
		outer2D := node.Outer.To2D()
		outer2D.Mul(multiplier)
		node.Outer = outer2D.To3D()

		// Handle inner circle
		inner2D := node.Inner[0].To2D()
		inner2D.Mul(multiplier)
		node.Inner[0] = inner2D.To3D()

		node = node.Extrude(h)
	}
	node.Flip()
	node.Center()
	node.Generate("funnel.obj")
}
