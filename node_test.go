package toothpaste

import (
	"reflect"
	"testing"
)

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
	uniques := nodes.UniqueVertices()
	if checkIfUnique(nodes, uniques) != false {
		t.Errorf("If this triggers then there's a problem with this test case")
	}
	nodes.LinkVertices()
	if checkIfUnique(nodes, uniques) != true {
		t.Errorf("Expected all vertices to be unique, got %v", nodes.UniqueVertices())
	}
}

func TestRoundVertices(t *testing.T) {
	nodes := Nodes{
		NewNode(NewFace3D(
			0, 1.001, 0,
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
	nodes.RoundVertices(1)
	uniques := nodes.UniqueVertices()
	if len(uniques) != 6 {
		t.Errorf("Expected 6 unique vertices, got %v", len(uniques))
	}
}

func TestDimensions(t *testing.T) {
	node := NewNode(NewFace3D(
		0.25, 0, 0,
		1.25, 0.75, 0,
		1.25, 0.75, 0.5,
		0.25, 0, 0.5,
	))
	width := node.Width()
	if width != 1 {
		t.Errorf("Expected width to be 1, got %v", width)
	}
	height := node.Height()
	if height != 0.75 {
		t.Errorf("Expected height to be 0.5, got %v", height)
	}
	depth := node.Depth()
	if depth != 0.5 {
		t.Errorf("Expected depth to be 0, got %v", depth)
	}
}

func TestInsertAfter(t *testing.T) {
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
	nodes.LinkVertices()
	nodes.LinkNodes()
	beforeNode := NewNode(NewFace3D(
		5, 0, 0,
		0, 0, 0,
		0, 0, 5,
		5, 0, 5,
	))
	afterNode := NewNode(NewFace3D(
		6, 0, 0,
		0, 0, 0,
		0, 0, 6,
		6, 0, 6,
	))
	beforeNode.Next = afterNode
	afterNode.Prev = beforeNode
	beforeNode.InsertAfter(nodes[0].First())
	_nodes := beforeNode.Nodes()
	if len(_nodes) != 5 {
		t.Errorf("Expected 5 nodes, got %v", len(_nodes))
	}
	if _nodes[0] != beforeNode {
		t.Errorf("Expected first node to be beforeNode, got %v", _nodes[0])
	}
	if _nodes[4] != afterNode {
		t.Errorf("Expected last node to be afterNode, got %v", _nodes[4])
	}
}

func TestDetach(t *testing.T) {
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
	nodes.LinkVertices()
	nodes.LinkNodes()
	newNode := nodes[0].Detach()
	nodes = nodes[1].Nodes()
	if len(nodes) != 3 {
		t.Errorf("Expected 3 nodes, got %v", len(nodes))
	}
	if newNode != nodes[0].Last() {
		t.Errorf("Expected newNode to be last node, got %v", newNode)
	}
}

func TestCopy(t *testing.T) {
	var nodes1 = Nodes{
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
	nodes1.LinkVertices()
	nodes1.LinkNodes()

	nodes2 := nodes1.Copy()
	if len(nodes2) != len(nodes1) {
		t.Errorf("Expected 3 nodes, got %v", len(nodes2))
	}
	// check every vertex has the same x, y, z values, but is a different pointer
	// remember, they may be in a different position in the slice
	uniques1 := nodes1.UniqueVertices()
	uniques2 := nodes2.UniqueVertices()
	if len(uniques1) != len(uniques2) {
		t.Errorf("Expected %v unique vertices, got %v", len(uniques1), len(uniques2))
	}
	for _, unique1 := range uniques1 {
		// find closest vertex in nodes2
		var closest *Vertex3D
		var closestDistance float64
		for _, _unique2 := range uniques2 {
			distance := unique1.Distance(_unique2)
			if closest == nil || distance < closestDistance {
				closest = _unique2
				closestDistance = distance
			}
		}
		if closest == nil {
			t.Errorf("Expected to find closest vertex, got nil")
		}
		if closest.X != unique1.X || closest.Y != unique1.Y || closest.Z != unique1.Z {
			t.Errorf("Expected vertices to be equal, got %v and %v", unique1, closest)
		}
		if closest == unique1 {
			t.Errorf("Expected vertices to be different pointers, got %p and %p", unique1, closest)
		}
	}
}

func TestDrop(t *testing.T) {
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
	nodes.LinkVertices()
	nodes.LinkNodes()
	nodes[0].Drop()
	nodes = nodes[1].Nodes()
	if len(nodes) != 2 {
		t.Errorf("Expected 2 nodes, got %v", len(nodes))
	}
}
