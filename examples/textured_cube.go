package main

import (
	"github.com/micah5/toothpaste"
)

func main() {
	w, d, h := 1.0, 1.0, 1.0
	bottom := toothpaste.NewNode(toothpaste.Square(w, d).To3D())
	top := bottom.ExtrudeFlip(h)

	// You can either specify the texture coordinates manually...
	top.AddTexture("examples/media/apple.png",
		toothpaste.NewVertex2D(0, 0),
		toothpaste.NewVertex2D(1, 0),
		toothpaste.NewVertex2D(1, 1),
		toothpaste.NewVertex2D(0, 1),
	)

	// ...or you can use the built-in texture coordinate generator
	bottom.AddTexture("examples/media/banana.png")

	bottom.Center()
	bottom.GenerateColor("textured_cube.obj")
}
