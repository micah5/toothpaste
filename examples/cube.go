package main

import (
	"github.com/micah5/toothpaste"
)

func main() {
	w, d, h := 1.0, 1.0, 1.0

	// Create a 2D shape that can be extruded
	// There are some built-in shapes, but you can also create your own
	square := toothpaste.Square(w, d)

	// A node is a face that is connected to other faces
	// Let's create a node from the square (at z=0)
	bottom := toothpaste.NewNode(square.To3D())

	// Extrude the node to create a cube
	// Extrude returns the top face of the extrusion
	top := bottom.Extrude(h)

	// Currently the top face is facing down (since it is a copy of the bottom)
	// Flip the top face so that it is facing up
	top.Flip()

	// Center the cube at the origin
	bottom.Center()

	nodes := bottom.Nodes()
	nodes2 := nodes.Copy()
	nodes2.Translate(0, 0, 4)
	bottom.Last().InsertAfter(nodes2[0].First())

	// Generate the .obj file
	bottom.Generate("cube.obj")
}
