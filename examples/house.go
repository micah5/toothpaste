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
	//panes := make(map[string]*toothpaste.Face2D)
	//for i := 0; i < len(sides); i++ {
	//	hole := pane.Copy()
	//	hole.Translate((w/2)*float64(i)+pad/2, (d/2)*float64(i)+pad/2)
	//	println(i, sides[i])
	//	panes[sides[i]] = hole
	//}

	// Create house
	node := toothpaste.NewNode(toothpaste.Square(w, d).To3D())
	node.ExtrudeFlip(h, append([]string{"top"}, sides...)...)

	// Add windows
	group := node.GetAll(sides...)
	group.AddHoles(pane)
	back := node.Get("back")
	back.ExtrudeInner(h / 4)

	//node.Center()
	node.Generate("house.obj")
}
