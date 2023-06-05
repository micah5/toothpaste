package main

import (
	"github.com/micah5/toothpaste"
)

func main() {
	w, d, h := 1.0, 1.0, 1.0
	node := toothpaste.NewTaggedNode("bottom", toothpaste.Square(w, d).To3D())
	node.Extrude(h, "top", "front", "right", "back", "left")
	node.Get("left").Rotate(30, toothpaste.ZAxis)
	node.Get("right").Rotate(30, toothpaste.ZAxis)
	node.Get("top").Flip()
	node.Center()
	node.Generate("parallelogram.obj")
}
