package toothpaste

import (
	"reflect"
	"testing"
)

func TestUniqueVertices(t *testing.T) {
	nodes := Nodes{
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
	if !reflect.DeepEqual(unique, expected) {
		t.Errorf("Expected %v, got %v", expected, unique)
	}
}
