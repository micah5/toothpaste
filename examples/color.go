package main

import (
	"github.com/micah5/toothpaste"
)

func main() {
	w, d, h := 1.0, 1.0, 1.0
	node1 := toothpaste.NewNode(toothpaste.Square(w, d).To3D())
	node1.ExtrudeFlip(h)

	node2 := node1.CopyAll()
	node2.Translate(w*1.5, 0, 0)

	// Tag as seperate colors
	node1.TagAll("red")
	node2.Tag("blue")
	node1.Last().InsertAfter(node2[0])

	node1.GenerateColor("cube.obj", map[string][3]float64{
		"red":  [3]float64{1, 0, 0},
		"blue": [3]float64{0, 0, 1},
	})
}
