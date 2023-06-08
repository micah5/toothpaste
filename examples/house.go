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
	bottom := toothpaste.NewNode(toothpaste.Square(w, d).To3D())
	top := bottom.Extrude(h, append([]string{"top"}, sides...)...)
	for _, tag := range sides {
		node := top.Get(tag)
		_pane := pane.Copy()
		_pane.Fit3D(node.Outer)
		node.Inner = append(node.Inner, _pane.To3D())
	}
	top.Flip()
	//bottom.Center()
	bottom.Generate("house.obj")
}
