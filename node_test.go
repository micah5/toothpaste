package toothpaste

import (
	"reflect"
	"testing"
)

var nodes = Nodes{
	NewNode(NewFace3D(
		0, 1, 0,
		0, 1, -2,
		0, 0, -2,
		0, 0, 0,
	)),
	NewNode(NewFace3D(
		0, 1, -2,
		1, 1, -3,
		1, 0, -3,
		0, 0, -2,
	)),
	NewNode(NewFace3D(
		1, 1, -3,
		0, 1, 0,
		0, 0, 0,
		1, 0, -3,
	)),
}

func isEqualSliceUnordered(x, y []*Vertex3D) bool {
	if len(x) != len(y) {
		return false
	}

	xMap := make(map[Vertex3D]bool)
	yMap := make(map[Vertex3D]bool)

	for _, xv := range x {
		xMap[*xv] = true
	}

	for _, yv := range y {
		yMap[*yv] = true
	}

	return reflect.DeepEqual(xMap, yMap)
}

func TestUniqueVertices(t *testing.T) {
	unique := nodes.UniqueVertices()
	expected := []*Vertex3D{
		{0, 1, 0},
		{0, 1, -2},
		{0, 0, -2},
		{0, 0, 0},
		{1, 1, -3},
		{1, 0, -3},
	}
	if !isEqualSliceUnordered(unique, expected) {
		t.Errorf("Expected %v, got %v", expected, unique)
	}
}

func checkIfUnique(nodes Nodes, uniques []*Vertex3D) bool {
	// check that every vertex is now one from the uniques slice
	for _, node := range nodes {
		faces := node.Faces()
		for _, face := range faces {
			for _, vertex := range face.Vertices {
				// check that the vertex is in the uniques slice
				found := false
				for _, unique := range uniques {
					if vertex == unique {
						found = true
						break
					}
				}
				if !found {
					return false
				}
			}
		}
	}
	return true
}

func TestLinkVertices(t *testing.T) {
	uniques := nodes.UniqueVertices()
	if checkIfUnique(nodes, uniques) != false {
		t.Errorf("If this triggers then there's a problem with this test case")
	}
	nodes.LinkVertices()
	if checkIfUnique(nodes, uniques) != true {
		t.Errorf("Expected all vertices to be unique, got %v", nodes.UniqueVertices())
	}
}
