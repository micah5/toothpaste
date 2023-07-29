package main

import (
	"github.com/micah5/toothpaste"
)

func main() {
	w, d, h := 1.0, 1.0, 1.0
	circle := toothpaste.Circle(w, d, 8)
	node := toothpaste.NewNode(circle.To3D())
	node.ExtrudeFlip(h) // helper function to flip the top face at the same time
	node.Center()

	// define prism to fit the cylinder into
	prism := [6]*toothpaste.Face3D{
		// the vertices are chosen to form a prism with the top face smaller than the bottom face
		&toothpaste.Face3D{Vertices: []*toothpaste.Vertex3D{{-2, -2, -0.5}, {2, -2, -0.5}, {2, 2, -0.5}, {-2, 2, -0.5}}},
		&toothpaste.Face3D{Vertices: []*toothpaste.Vertex3D{{-1, -1, 0.5}, {1, -1, 0.5}, {1, 1, 0.5}, {-1, 1, 0.5}}},
		&toothpaste.Face3D{Vertices: []*toothpaste.Vertex3D{{-2, -2, -0.5}, {2, -2, -0.5}, {1, -1, 0.5}, {-1, -1, 0.5}}},
		&toothpaste.Face3D{Vertices: []*toothpaste.Vertex3D{{2, -2, -0.5}, {2, 2, -0.5}, {1, 1, 0.5}, {1, -1, 0.5}}},
		&toothpaste.Face3D{Vertices: []*toothpaste.Vertex3D{{2, 2, -0.5}, {-2, 2, -0.5}, {-1, 1, 0.5}, {1, 1, 0.5}}},
		&toothpaste.Face3D{Vertices: []*toothpaste.Vertex3D{{-2, 2, -0.5}, {-2, -2, -0.5}, {-1, -1, 0.5}, {-1, 1, 0.5}}},
	}

	// create a node for the prism
	prismNode := toothpaste.NewSliceNode(prism[:]...)

	// "squash" the cylinder
	node.FitRectangularPrism(prism)

	node.Generate("squashed_cylinder.obj")
	prismNode.Generate("squashed_cylinder_prism.obj")
}
