package main

import (
	"github.com/micah5/toothpaste"
	"math"
)

func main() {
	w, d, h := 1.0, 1.0, 0.5
	outer := toothpaste.Circle(w, d, 10)
	inner := toothpaste.Circle(2*w/3, 2*d/3, 10)
	inner.Translate(w/6, d/6)
	node := toothpaste.NewNode(outer.To3D(), inner.To3D())
	// ExtrudeLoop extrudes multiple times and drops any internal faces for you
	// You specify the function to apply on each iteration
	numIter := 3
	node = node.ExtrudeLoop(numIter, func(i int, node *toothpaste.Node) float64 {
		// Decrease scaling factor with each iteration to make the object bulge less and less
		node.Outer.Mul2D(math.Pow(1.15, float64(numIter-i)))
		// Decrease extrusion distance with each iteration to make the object shorter and shorter
		return (h / float64(numIter)) * math.Pow(0.5, float64(numIter-i)) * 2
	})
	node = node.ExtrudeLoop(numIter, func(i int, node *toothpaste.Node) float64 {
		node.Outer.Mul2D(math.Pow(0.98, float64(i)))
		return (h / float64(numIter)) * math.Pow(0.5, float64(i)) / 2
	})
	node.Flip()
	node.Center()
	node.Generate("donut.obj")
}
