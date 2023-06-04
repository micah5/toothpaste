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
	node := toothpaste.NewNode(square.To3D())

	// Extrude the node to create a cube
	node.Extrude(h)

	// Center the cube at the origin
	node.Center()

	// Generate the .obj file
	node.Generate("cube.obj")
}
