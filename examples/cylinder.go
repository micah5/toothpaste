package main

import (
	"github.com/micah5/toothpaste"
)

func main() {
	w, d, h := 1.0, 1.0, 1.0
	circle := toothpaste.Circle(w, d, 8)
	node := toothpaste.NewNode(circle.To3D())
	node.ExtrudeFlip(h) // helper function to flip the top face at the same time
	node.Center()
	node.Generate("cylinder.obj")
}
