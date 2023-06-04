package main

import (
	"github.com/micah5/toothpaste"
)

func main() {
	w, d, h := 1.0, 1.0, 1.0
	node := toothpaste.NewNode(toothpaste.Square(w, d))
	node.Extrude(h)
	node.Generate("cube.obj")
}
