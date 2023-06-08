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

func (n *Node) ExtrudeLoop(numIter int, fn func(int, *Node) float64, tags ...string) *Node {
	var node *Node = n
	for i := 0; i < numIter; i++ {
		height := fn(i, node)
		face := node.Outer.Copy()
		node2 := extrude(node, face, height, tags...)
		if i > 0 {
			node.Drop()
		}
		node = node2
	}
	fn(numIter, node)
	return node
}

func (n *Node) ExtrudeDrop(height float64, tags ...string) *Node {
	top := n.Outer.Copy()
	res := extrude(n, top, height, tags...)
	n.Drop()
	return res
}

func (n *Node) ExtrudeFlip(height float64, tags ...string) *Node {
	top := n.Outer.Copy()
	res := extrude(n, top, height, tags...)
	res.Flip()
	return res
}

func (n *Node) Extrude(height float64, tags ...string) *Node {
	top := n.Outer.Copy()
	return extrude(n, top, height, tags...)
}

func (n *Node) Faces() []*Face3D {
	faces := []*Face3D{n.Outer}
	faces = append(faces, n.Inner...)
	return faces
}

func (n *Node) Nodes() Nodes {
	nodes := Nodes{n}
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

func (n *Node) Drop() {
	n.Prev.Next = n.Next
	n.Next.Prev = n.Prev
}

func (n *Node) Remove() {
	n.Drop()
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

func (n *Node) Mul(m float64) {
	faces := n.Faces()
	for _, f := range faces {
		f.Mul(m)
	}
}

func (n *Node) Translate2D(x, y float64) {
	faces := n.Faces()
	for _, f := range faces {
		f.Translate2D(x, y)
	}
}

func (n *Node) Rotate2D(deg int) {
	faces := n.Faces()
	for _, f := range faces {
		f.Rotate2D(deg)
	}
}

func (n *Node) Scale2D(x, y float64) {
	faces := n.Faces()
	for _, f := range faces {
		f.Scale2D(x, y)
	}
}

func (n *Node) Mul2D(m float64) {
	faces := n.Faces()
	for _, f := range faces {
		f.Mul2D(m)
	}
}

func (n *Node) Flip() {
	faces := n.Faces()
	for _, f := range faces {
		f.Flip()
	}
}

func (n *Node) AddHoles(holes2D ...*Face2D) {
	holes3D := make([]*Face3D, len(holes2D))
	for i, h := range holes2D {
		holes3D[i] = h.To3D(true)
	}
	n.Inner = append(n.Inner, holes3D...)
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

func (n *Node) GetAll(tag string) Nodes {
	nodes := n.Nodes()
	var matches []*Node
	for _, node := range nodes {
		if node.Tag == tag {
			matches = append(matches, node)
		}
	}
	return matches
}

func (n *Node) Centroid() *Vertex3D {
	nodes := n.Nodes()
	var sumX, sumY, sumZ float64
	for _, node := range nodes {
		sumX += node.Outer.Centroid().X
		sumY += node.Outer.Centroid().Y
		sumZ += node.Outer.Centroid().Z
	}
	return &Vertex3D{sumX / float64(len(nodes)), sumY / float64(len(nodes)), sumZ / float64(len(nodes))}
}

func (n *Node) Center() {
	ns := n.Nodes()
	centroid := ns.Centroid()
	ns.Translate(-centroid.X, -centroid.Y, -centroid.Z)
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

type Nodes []*Node

func (ns Nodes) Centroid() *Vertex3D {
	var sumX, sumY, sumZ float64
	for _, node := range ns {
		sumX += node.Outer.Centroid().X
		sumY += node.Outer.Centroid().Y
		sumZ += node.Outer.Centroid().Z
	}
	return &Vertex3D{sumX / float64(len(ns)), sumY / float64(len(ns)), sumZ / float64(len(ns))}
}

func (ns Nodes) UniqueVertices() []*Vertex3D {
	unique := make(map[string]*Vertex3D)
	for _, node := range ns {
		for _, v := range node.Outer.Vertices {
			unique[v.String()] = v
		}
	}
	var vertices []*Vertex3D
	for _, v := range unique {
		vertices = append(vertices, v)
	}
	return vertices
}

func (ns Nodes) Translate(x, y, z float64) {
	uniques := ns.UniqueVertices()
	for _, v := range uniques {
		v.Translate(x, y, z)
	}
}

// utils
func getTag(idx int, tags []string) string {
	if idx >= len(tags) {
		return ""
	}
	return tags[idx]
}

func extrude(n *Node, top *Face3D, height float64, tags ...string) *Node {
	// Negate the normal vector components to flip the direction
	normal := n.Outer.Normal()
	normal.Mul(-1)

	// Create the top face
	top.Translate(normal.X*height, normal.Y*height, normal.Z*height)
	holes := make([]*Face3D, 0)
	for _, f := range n.Inner {
		hole := f.Copy()
		hole.Translate(normal.X*height, normal.Y*height, normal.Z*height)
		holes = append(holes, hole)
	}
	topN := NewTaggedNode(getTag(0, tags), top, holes...)

	// Create the sides
	var cur *Node = n
	faces := n.Faces()
	for k, f := range faces {
		for i := range f.Vertices {
			i1, i2 := i, (i+1)%len(f.Vertices)
			v1 := f.Vertices[i1]
			v2 := f.Vertices[i2]
			var topV1, topV2 *Vertex3D
			if k == 0 {
				topV1 = top.Vertices[i1]
				topV2 = top.Vertices[i2]
			} else {
				topV1 = holes[k-1].Vertices[i1]
				topV2 = holes[k-1].Vertices[i2]
			}
			sideFace := &Face3D{
				Vertices: []*Vertex3D{
					topV1,
					topV2,
					v2,
					v1,
				},
			}
			if k != 0 {
				sideFace.Flip()
			}
			newNode := NewTaggedNode(getTag(i+1, tags), sideFace)
			cur.Next = newNode
			newNode.Prev = cur
			cur = newNode
		}
	}
	cur.Next = topN
	topN.Prev = cur
	cur = topN

	return topN
}
