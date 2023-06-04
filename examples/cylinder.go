package main

import (
	"github.com/micah5/toothpaste"
)

func main() {
	w, d, h := 1.0, 1.0, 1.0
	circle := toothpaste.Circle(w, d, 8)
	node := toothpaste.NewNode(circle.To3D())
	node.Extrude(h)
	node.Center()
	node.Generate("cylinder.obj")
}
