package main

import (
	"github.com/micah5/toothpaste"
)

func main() {
	w, d, h := 1.0, 1.0, 0.5
	outer := toothpaste.Circle(w, d, 30)
	inner := toothpaste.Circle(w/2, d/2, 20)
	inner.Translate(w/4, d/4)
	node := toothpaste.NewNode(outer.To3D(toothpaste.ZAxis), inner.To3D(toothpaste.ZAxis))
	top := node.ExtrudeFlip(h)
	node.Center()
	node.Generate("wheel.obj")
}
