package toothpaste

import (
	"github.com/micah5/earcut-3d"
)

type Node struct {
	Tag   string
	Outer *Face3D
	Inner []*Face3D
	Prev  *Node
	Next  *Node
}

func NewNode(outer *Face3D, inner ...*Face3D) *Node {
	return &Node{"", outer, inner, nil, nil}
}

func NewTaggedNode(tag string, outer *Face3D, inner ...*Face3D) *Node {
	return &Node{tag, outer, inner, nil, nil}
}

func (n *Node) Extrude(height float64, tags ...string) {
	faces := n.Faces()
	for _, f := range faces {
		// Negate the normal vector components to flip the direction
		normal := f.Normal()
		normal.Mul(-1)

		// Create the top face
		top := f.Copy()
		top.Translate(normal.X*height, normal.Y*height, normal.Z*height)
		topN := NewTaggedNode(getTag(0, tags), top)

		// Create the sides
		var cur *Node = n
		for i := range f.Vertices {
			i1, i2 := i, (i+1)%len(f.Vertices)
			v1 := f.Vertices[i1]
			v2 := f.Vertices[i2]
			sideFace := NewFace3D(
				top.Vertices[i1].X, top.Vertices[i1].Y, top.Vertices[i1].Z,
				top.Vertices[i2].X, top.Vertices[i2].Y, top.Vertices[i2].Z,
				v2.X, v2.Y, v2.Z,
				v1.X, v1.Y, v1.Z,
			)
			newNode := NewTaggedNode(getTag(i+1, tags), sideFace)
			cur.Next = newNode
			newNode.Prev = cur
			cur = newNode
		}
		top.Flip()
		cur.Next = topN
		topN.Prev = cur
	}
}

func (n *Node) Faces() []*Face3D {
	faces := []*Face3D{n.Outer}
	faces = append(faces, n.Inner...)
	return faces
}

func (n *Node) Nodes() []*Node {
	nodes := []*Node{n}
	cur := n.Next
	for cur != nil {
		nodes = append(nodes, cur)
		cur = cur.Next
	}
	cur = n.Prev
	for cur != nil {
		nodes = append(nodes, cur)
		cur = cur.Prev
	}
	return nodes
}

func (n *Node) Remove() {
	n.Prev.Next = n.Next
	n.Next.Prev = n.Prev
}

func (n *Node) Translate(x, y, z float64) {
	faces := n.Faces()
	for _, f := range faces {
		f.Translate(x, y, z)
	}
}

func (n *Node) Rotate(deg int, axis Axis) {
	faces := n.Faces()
	for _, f := range faces {
		f.Rotate(deg, axis)
	}
}

func (n *Node) Scale(x, y, z float64) {
	faces := n.Faces()
	for _, f := range faces {
		f.Scale(x, y, z)
	}
}

func (n *Node) Flip() {
	faces := n.Faces()
	for _, f := range faces {
		f.Flip()
	}
}

func (n *Node) Get(tag string) *Node {
	nodes := n.Nodes()
	for _, node := range nodes {
		if node.Tag == tag {
			return node
		}
	}
	return nil
}

func (n *Node) GetAll(tag string) []*Node {
	nodes := n.Nodes()
	var matches []*Node
	for _, node := range nodes {
		if node.Tag == tag {
			matches = append(matches, node)
		}
	}
	return matches
}

func (n *Node) Reverse() {
	n.Flip()
}

func (n *Node) Generate(filename string) {
	nodes := n.Nodes()
	var faces [][]float64
	var holes [][][]float64
	for _, node := range nodes {
		faces = append(faces, node.Outer.Flatten())
		_holes := make([][]float64, 0)
		for _, inner := range node.Inner {
			_holes = append(_holes, inner.Flatten())
		}
		holes = append(holes, _holes)
	}
	triangles := earcut3d.Earcut(faces, holes...)
	earcut3d.CreateObjFile(filename, triangles)
}

// utils
func getTag(idx int, tags []string) string {
	if idx >= len(tags) {
		return ""
	}
	return tags[idx]
}
