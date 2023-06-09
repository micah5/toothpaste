package main

import (
	"github.com/micah5/toothpaste"
)

func main() {
	w, d, h := 1.0, 1.0, 1.0

	// Create window shape
	sides := []string{"front", "right", "back", "left"}
	pane := toothpaste.Square(w/2, d/2)
	pane.Translate(w/4, d/4)

	// Create house
	node := toothpaste.NewNode(toothpaste.Square(w, d).To3D())
	roof := node.ExtrudeFlip(h, append([]string{"top"}, sides...)...)

	// Add windows
	group := node.GetAll(sides...)
	group.AddHoles(pane)
	group.ExtrudeInner(h / 6)

	// Add Roof
	// Detach roof from rest of house otherwise when we resize it, the rest of the house will resize too
	//roof.Detach()
	roof.Mul2D(1.2)
	roof = roof

	node.Center()
	node.Generate("house.obj")
}
