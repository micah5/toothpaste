package main

import (
	"github.com/micah5/toothpaste"
)

func main() {
	w, d, h := 1.0, 1.0, 1.0

	// Create window shape
	sides := []string{"right", "back", "left"}
	pane := toothpaste.Square(w/2, d/2)
	pane.Translate(w/4, d/4)

	// Create door shape
	door := toothpaste.Square(w/4, 2*d/3)
	door.Translate((w-w/4)/2, 0.0)

	// Create house
	node := toothpaste.NewNode(toothpaste.Square(w, d).To3D())
	roof := node.ExtrudeFlip(h, append([]string{"top", "front"}, sides...)...)

	// Add windows
	group := node.GetAll(sides...)
	group.AddHoles(pane)
	group.ExtrudeInner(-h / 6)

	// Add door
	front := node.Get("front")
	front.AddHoles(door)
	front.ExtrudeInner(-h / 12)

	// Widen top of house for cartoony look
	roof.Mul2D(1.2)

	// Add Roof
	// Detach roof from rest of house otherwise when we resize it,
	// the rest of the house will move with it again too
	roof = roof.Detach()
	roof.Mul2D(1.2)

	//node.Center()
	node.Generate("house.obj")
}
