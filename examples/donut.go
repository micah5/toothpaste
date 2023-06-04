package main

import (
	"github.com/micah5/toothpaste"
)

func main() {
	w, d, h := 1.0, 1.0, 1.0
	outer := toothpaste.Circle(w, d, 10)
	inner := toothpaste.Circle(w/2, d/2, 10)
	inner.Translate(w/4, d/4)
	node := toothpaste.NewNode(outer.To3D(), inner.To3D())
	node.Extrude(h)
	//node.Center()
	node.Generate("donut.obj")
}
