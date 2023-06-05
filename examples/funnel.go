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

	// For a simplier solution, we could use ExtrudeLoop here
	// but we'll do it manually to demonstrate the process
	num_iterations := 4
	for i := 0; i < num_iterations; i++ {
		multiplier := math.Pow(1.34, float64(i))

		// Handle changing size of inner and outer circles
		node.Mul2D(multiplier)

		if i > 0 {
			// ExtrudeDrop is a helper function to extrude and drop the bottom face
			// It's useful for the second iteration and beyond because
			// the bottom faces are hidden within the mesh from this point on
			node = node.ExtrudeDrop(h)
		} else {
			// We can't use it for the first iteration because we need the bottom
			node = node.Extrude(h)
		}
	}
	node.Flip()
	node.Center()
	node.Generate("funnel.obj")
}
