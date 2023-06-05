package main

import (
	"github.com/micah5/toothpaste"
	"math"
)

func main() {
	w, d, h := 1.0, 1.0, 0.5
	outer := toothpaste.Circle(w, d, 10)
	inner := toothpaste.Circle(w/2, d/2, 10)
	inner.Translate(w/4, d/4)
	node := toothpaste.NewNode(outer.To3D(), inner.To3D())

	//num_iterations := 4
	//for i := 0; i < num_iterations; i++ {
	//	multiplier := math.Pow(1.2, float64(i))
	//	node = node.ExtrudeMul(h, multiplier)
	//}
	node.Flip()
	node.Center()
	node.Generate("donut.obj")
}
