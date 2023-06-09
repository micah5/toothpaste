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
		node2 := extrude(node, []*Face3D{node.Outer}, height, true, tags...)
		if i > 0 {
			node.Drop()
		}
		node = node2
	}
	fn(numIter, node)
	return node
}

func (n *Node) ExtrudeDrop(height float64, tags ...string) *Node {
	res := extrude(n, n.Faces(), height, true, tags...)
	n.Drop()
	return res
}

func (n *Node) ExtrudeFlip(height float64, tags ...string) *Node {
	res := extrude(n, n.Faces(), height, true, tags...)
	res.Flip()
	return res
}

func (n *Node) Extrude(height float64, tags ...string) *Node {
	return extrude(n, n.Faces(), height, true, tags...)
}

func (n *Node) ExtrudeInner(height float64, tags ...string) Nodes {
	tops := make(Nodes, len(n.Inner))
	for _, f := range n.Inner {
		res := extrude(n, []*Face3D{f}, height, false, tags...)
		res.Flip()
		tops = append(tops, res)
	}
	return tops
}

func (n *Node) ExtrudePoint(height float64, tags ...string) *Node {
	cen := n.Outer.Centroid()
	normal := n.Outer.Normal()
	normal.Mul(-1)
	cen.Translate(normal.X*height, normal.Y*height, normal.Z*height)

	// loop through every pair of vertices and create a triangle face
	// with the centroid
	var cur *Node = n
	for i := 0; i < len(n.Outer.Vertices); i++ {
		v1 := n.Outer.Vertices[i]
		v2 := n.Outer.Vertices[(i+1)%len(n.Outer.Vertices)]
		f := &Face3D{
			Vertices: []*Vertex3D{cen, v2, v1},
		}
		newNode := NewTaggedNode(getTag(i+1, tags), f)
		cur.Next = newNode
		newNode.Prev = cur
		cur = newNode
	}

	return n
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

func (n *Node) Connected() Nodes {
	nodes := n.Nodes()
	connected := map[*Node]bool{}
	for _, node := range nodes {
		for _, v := range n.Outer.Vertices {
			if node.Outer.ContainsExact(v) && node != n {
				connected[node] = true
			}
		}
	}
	res := Nodes{}
	for node, _ := range connected {
		res = append(res, node)
	}
	return res
}

func (n *Node) RealignConnectedInner() {
	connected := n.Connected()
	for _, node := range connected {
		for _, hole := range node.Inner {
			percShape := hole.PercShape.Copy()
			percShape.Fit3D(node.Outer)
			res := percShape.To3D(true)
			for i, v := range hole.Vertices {
				v.X = res.Vertices[i].X
				v.Y = res.Vertices[i].Y
				v.Z = res.Vertices[i].Z
			}
		}
	}
}

func (n *Node) Detach() *Node {
	n.Drop()
	tmp := n.Copy()
	n.Prev.InsertAfter(tmp)
	return tmp
}

func (n *Node) Copy() *Node {
	holes := make([]*Face3D, len(n.Inner))
	for i, f := range n.Inner {
		holes[i] = f.Copy()
	}
	return NewTaggedNode(n.Tag, n.Outer.Copy(), holes...)
}

func (n *Node) Drop() {
	if n.Prev != nil {
		n.Prev.Next = n.Next
	}
	if n.Next != nil {
		n.Next.Prev = n.Prev
	}
}

func (n *Node) InsertAfter(node *Node) {
	if n.Next != nil {
		n.Next.Prev = node
	}
	node.Next = n.Next
	node.Prev = n
	n.Next = node
}

func (n *Node) InsertBefore(node *Node) {
	if n.Prev != nil {
		n.Prev.Next = node
	}
	node.Prev = n.Prev
	node.Next = n
	n.Prev = node
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
	n.RealignConnectedInner()
}

func (n *Node) Rotate2D(deg int) {
	faces := n.Faces()
	for _, f := range faces {
		f.Rotate2D(deg)
	}
	n.RealignConnectedInner()
}

func (n *Node) Scale2D(x, y float64) {
	faces := n.Faces()
	for _, f := range faces {
		f.Scale2D(x, y)
	}
	n.RealignConnectedInner()
}

func (n *Node) Mul2D(m float64) {
	faces := n.Faces()
	for _, f := range faces {
		f.Mul2D(m)
	}
	n.RealignConnectedInner()
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
		tmp := h.Copy()
		tmp.Fit3D(n.Outer)
		holes3D[i] = tmp.To3D(true)
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

func (n *Node) GetPrev(num int) Nodes {
	nodes := Nodes{}
	cur := n.Prev
	for i := 0; i < num; i++ {
		nodes = append(nodes, cur)
		cur = cur.Prev
	}
	return nodes
}

func (n *Node) GetNext(num int) Nodes {
	nodes := Nodes{}
	cur := n.Next
	for i := 0; i < num; i++ {
		nodes = append(nodes, cur)
		cur = cur.Next
	}
	return nodes
}

func (n *Node) GetAll(tags ...string) Nodes {
	nodes := n.Nodes()
	var matches []*Node
	for _, node := range nodes {
		for _, tag := range tags {
			if node.Tag == tag {
				matches = append(matches, node)
			}
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
		faces := node.Faces()
		for _, f := range faces {
			for _, v := range f.Vertices {
				unique[v.String()] = v
			}
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

func (ns Nodes) Mul(magnitude float64) {
	uniques := ns.UniqueVertices()
	for _, v := range uniques {
		v.Mul(magnitude)
	}
}

func (ns Nodes) Scale(x, y, z float64) {
	uniques := ns.UniqueVertices()
	for _, v := range uniques {
		v.Scale(x, y, z)
	}
}

func (ns Nodes) Rotate(deg int, axis Axis) {
	uniques := ns.UniqueVertices()
	for _, v := range uniques {
		v.Rotate(deg, axis)
	}
}

func (ns Nodes) Flip() {
	for _, node := range ns {
		node.Flip()
	}
}

func (ns Nodes) Extrude(height float64, tags ...string) Nodes {
	var nodes Nodes
	for _, node := range ns {
		nodes = append(nodes, node.Extrude(height, tags...))
	}
	return nodes
}

func (ns Nodes) ExtrudeFlip(height float64, tags ...string) Nodes {
	var nodes Nodes
	for _, node := range ns {
		nodes = append(nodes, node.ExtrudeFlip(height, tags...))
	}
	return nodes
}

func (ns Nodes) ExtrudeDrop(height float64, tags ...string) Nodes {
	var nodes Nodes
	for _, node := range ns {
		nodes = append(nodes, node.ExtrudeDrop(height, tags...))
	}
	return nodes
}

func (ns Nodes) Drop() {
	for _, node := range ns {
		node.Drop()
	}
}

func (ns Nodes) ExtrudeInner(height float64, tags ...string) Nodes {
	var nodes Nodes
	for _, node := range ns {
		nodes = append(nodes, node.ExtrudeInner(height, tags...)...)
	}
	return nodes
}

func (ns Nodes) AddHoles(holes2D ...*Face2D) {
	for _, node := range ns {
		node.AddHoles(holes2D...)
	}
}

func (ns Nodes) Mul2D(magnitude float64) {
	for _, node := range ns {
		node.Mul2D(magnitude)
	}
}

func (ns Nodes) Scale2D(x, y float64) {
	for _, node := range ns {
		node.Scale2D(x, y)
	}
}

func (ns Nodes) Rotate2D(deg int) {
	for _, node := range ns {
		node.Rotate2D(deg)
	}
}

func (ns Nodes) Translate2D(x, y float64) {
	for _, node := range ns {
		node.Translate2D(x, y)
	}
}

func (ns Nodes) Get(tag string) Nodes {
	return ns[0].GetAll(tag)
}

func (ns Nodes) GetAll(tag string) Nodes {
	return ns[0].GetAll(tag)
}

// utils
func getTag(idx int, tags []string) string {
	if idx >= len(tags) {
		return ""
	}
	return tags[idx]
}

func extrude(n *Node, faces []*Face3D, height float64, addHoles bool, tags ...string) *Node {
	var next *Node
	if n.Next != nil {
		next = n.Next
	}

	// Negate the normal vector components to flip the direction
	normal := faces[0].Normal()
	normal.Mul(-1)

	// Create the top face
	top := faces[0].Copy()
	top.Translate(normal.X*height, normal.Y*height, normal.Z*height)
	holes := make([]*Face3D, 0)
	if addHoles {
		for _, f := range n.Inner {
			hole := f.Copy()
			hole.Translate(normal.X*height, normal.Y*height, normal.Z*height)
			holes = append(holes, hole)
		}
	}
	topN := NewTaggedNode(getTag(0, tags), top, holes...)

	// Create the sides
	var cur *Node = n
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
			vertices := []*Vertex3D{
				topV1,
				topV2,
				v2,
				v1,
			}
			if height < 0 {
				vertices = []*Vertex3D{
					v1,
					v2,
					topV2,
					topV1,
				}
			}
			sideFace := &Face3D{
				Vertices: vertices,
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

	if next != nil {
		cur.Next = next
		next.Prev = cur
	}

	return topN
}
