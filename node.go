package toothpaste

import (
	"fmt"
	"github.com/micah5/earcut-3d"
	"math"
	"os"
	"reflect"
	//"strings"
)

type Node struct {
	Tag          string
	Outer        *Face3D
	Inner        []*Face3D
	Prev         *Node
	Next         *Node
	ImageTexture bool
	Meta         map[string]interface{}
}

func NewNode(outer *Face3D, inner ...*Face3D) *Node {
	return &Node{"", outer, inner, nil, nil, false, nil}
}

func NewTaggedNode(tag string, outer *Face3D, inner ...*Face3D) *Node {
	return &Node{tag, outer, inner, nil, nil, false, nil}
}

func NewSliceNode(outers ...*Face3D) *Node {
	// create multiple nodes, each with a single outer face
	// and then link them together
	var prev *Node
	for _, outer := range outers {
		node := NewNode(outer)
		if prev != nil {
			prev.Next = node
			node.Prev = prev
		}
		prev = node
	}
	return prev
}

func NewLinkedNodes(nodes ...*Node) *Node {
	var prev *Node
	for _, node := range nodes {
		if prev != nil {
			prev.Next = node
			node.Prev = prev
		}
		prev = node
	}
	return prev
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
	res.GetPrev(n.CountInnerVertices()).Flip()
	n.Drop()
	return res
}

func (n *Node) ExtrudeFlip(height float64, tags ...string) *Node {
	res := extrude(n, n.Faces(), height, true, tags...)
	res.GetPrev(n.CountInnerVertices()).Flip()
	res.Flip()
	return res
}

func (n *Node) Extrude(height float64, tags ...string) *Node {
	res := extrude(n, n.Faces(), height, true, tags...)
	res.GetPrev(n.CountInnerVertices()).Flip()
	return res
}

func (n *Node) ExtrudeInner(height float64, tags ...string) Nodes {
	tops := make(Nodes, 0)
	for _, f := range n.Inner {
		res := extrude(n, []*Face3D{f}, height, false, tags...)
		res.GetPrev(len(f.Vertices)).Flip()
		res.Flip()
		tops = append(tops, res)
	}
	return tops
}

func (n *Node) ExtrudeOuter(height float64, tags ...string) *Node {
	return extrude(n, []*Face3D{n.Outer}, height, false, tags...)
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

func (n *Node) CountInnerVertices() int {
	count := 0
	for _, f := range n.Inner {
		count += len(f.Vertices)
	}
	return count
}

func (n *Node) Width() float64 {
	return n.Outer.Width()
}

func (n *Node) Height() float64 {
	return n.Outer.Height()
}

func (n *Node) Depth() float64 {
	return n.Outer.Depth()
}

func (n *Node) Faces() []*Face3D {
	faces := []*Face3D{n.Outer}
	faces = append(faces, n.Inner...)
	return faces
}

func (n *Node) Merge(nodes ...*Node) {
	for _, node := range nodes {
		n.Last().InsertAfter(node.First())
	}
}

func (n *Node) Nodes() Nodes {
	nodes := Nodes{}
	cur := n.First()
	for cur != nil {
		nodes = append(nodes, cur)
		cur = cur.Next
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

func (n *Node) Align(n2 *Node) {
	faces := n.Faces()
	faces2 := n2.Faces()
	for i, f := range faces {
		f.Align(faces2[i])
	}
}

func (n *Node) Detach() *Node {
	last := n.Last()
	if last == n {
		last = n.Prev
	}
	n.Drop()
	tmp := n.Copy()
	last.Next = tmp
	tmp.Prev = last
	return tmp
}

func (n *Node) DetachVertices() {
	// Keep the node atatched to the other nodes,
	// but detach the vertices from the node
	// so that they can be moved independently
	for _, face := range n.Faces() {
		newVertices := make([]*Vertex3D, 0)
		for _, vertex := range face.Vertices {
			newVertex := vertex.Copy()
			newVertices = append(newVertices, newVertex)
		}
		face.Vertices = newVertices
	}
}

func (n *Node) Copy() *Node {
	holes := make([]*Face3D, len(n.Inner))
	for i, f := range n.Inner {
		holes[i] = f.Copy()
	}
	copyNode := NewTaggedNode(n.Tag, n.Outer.Copy(), holes...)
	copyNode.ImageTexture = n.ImageTexture
	copyNode.Meta = n.Meta
	return copyNode
}

func (n *Node) CopyAll() Nodes {
	nodes := n.Nodes()
	res := nodes.CopyAll()
	return res
}

func (n *Node) Drop() {
	if n.Prev != nil {
		n.Prev.Next = n.Next
	}
	if n.Next != nil {
		n.Next.Prev = n.Prev
	}
	n.Next = nil
	n.Prev = nil
}

func (n *Node) SetMeta(key string, value interface{}) {
	if n.Meta == nil {
		n.Meta = map[string]interface{}{}
	}
	n.Meta[key] = value
}

func (n *Node) SetMetaAll(key string, value interface{}) {
	nodes := n.Nodes()
	nodes.SetMeta(key, value)
}

func (n *Node) GetMeta(key string) interface{} {
	if n.Meta == nil {
		return nil
	}
	return n.Meta[key]
}

func (n *Node) GetNodesMeta(key string, value interface{}) Nodes {
	nodes := make(Nodes, 0)
	for _, node := range n.Nodes() {
		if node.GetMeta(key) == value {
			nodes = append(nodes, node)
		}
	}
	return nodes
}

func (n *Node) InsertAfter(node *Node) {
	// Save reference to next Node
	originalNext := n.Next

	// Make the connection from the current node to the inserted node
	n.Next = node
	node.Prev = n

	// Find the last node of the inserted nodes
	currentNode := node
	for currentNode.Next != nil {
		currentNode = currentNode.Next
	}

	// Connect the last inserted node to the original next node
	currentNode.Next = originalNext
	if originalNext != nil {
		originalNext.Prev = currentNode
	}
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

func (n *Node) Rotate(deg float64, axis Axis) {
	faces := n.Faces()
	for _, f := range faces {
		f.Rotate(deg, axis)
	}
}

func (n *Node) RotateFixed(deg float64, axis Axis) {
	faces := n.Faces()
	for _, f := range faces {
		f.RotateFixed(deg, axis)
	}
}

func (n *Node) RoundVertices(precision int) {
	faces := n.Faces()
	for _, f := range faces {
		f.RoundVertices(precision)
	}
}

func (n *Node) Scale(x, y, z float64) {
	faces := n.Faces()
	for _, f := range faces {
		f.Scale(x, y, z)
	}
}

func (n *Node) ScaleFixed(x, y, z float64) {
	faces := n.Faces()
	for _, f := range faces {
		f.ScaleFixed(x, y, z)
	}
}

func (n *Node) MoveTo(x, y, z float64) {
	faces := n.Faces()
	for _, f := range faces {
		f.MoveTo(x, y, z)
	}
}

func (n *Node) Mul(m float64) {
	faces := n.Faces()
	for _, f := range faces {
		f.Mul(m)
	}
}

func (n *Node) Snap(point *Vertex3D) {
	faces := n.Faces()
	for _, f := range faces {
		f.Snap(point)
	}
}

func (n *Node) Mirror(axis Axis) {
	faces := n.Faces()
	for _, f := range faces {
		f.Mirror(axis)
	}
}

func (n *Node) MulFixed(m float64) {
	faces := n.Faces()
	for _, f := range faces {
		f.MulFixed(m)
	}
}

func (n *Node) Translate2D(x, y float64) {
	faces := n.Faces()
	for _, f := range faces {
		f.Translate2D(x, y)
	}
	n.RealignConnectedInner()
}

func (n *Node) Rotate2D(deg float64) {
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

func (n *Node) AddHolesMiddleFixed(holes2D ...*Face2D) {
	holes3D := make([]*Face3D, len(holes2D))
	for i, f2D := range holes2D {
		holeW, holeH := f2D.Width(), f2D.Height()
		nodeW, nodeH := n.Outer.Width(), n.Outer.Height()
		x := (nodeW - holeW) / 2.0
		y := (nodeH - holeH) / 2.0
		wPerc, hPerc := holeW/nodeW, holeH/nodeH
		xPerc, yPerc := x/nodeW, y/nodeH
		tmp := Square(wPerc, hPerc)
		tmp.Translate(xPerc, yPerc)
		tmp.Fit3D(n.Outer)
		holes3D[i] = tmp.To3D(true)
	}
	n.Inner = append(n.Inner, holes3D...)
}

func (n *Node) DetachHoles() Nodes {
	tmpFaces := make([]*Face3D, len(n.Inner))
	for i, f := range n.Inner {
		tmpFaces[i] = f.Copy()
	}
	retNodes := make(Nodes, 0)
	for _, f := range tmpFaces {
		_n := NewNode(f)
		n.Attach(_n)
		retNodes = append(retNodes, _n)
	}
	return retNodes
}

func (n *Node) DetachHole(index int) *Node {
	tmpFaces := make([]*Face3D, len(n.Inner))
	for i, f := range n.Inner {
		tmpFaces[i] = f.Copy()
	}
	_n := NewNode(tmpFaces[index])
	n.Attach(_n)
	return _n
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

func (n *Node) GetPrevByIndex(index int) *Node {
	cur := n
	for i := 0; i < index; i++ {
		cur = cur.Prev
	}
	return cur
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

func (n *Node) GetNextByIndex(index int) *Node {
	cur := n
	for i := 0; i < index; i++ {
		cur = cur.Next
	}
	return cur
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
	return NewVertex3D(sumX/float64(len(nodes)), sumY/float64(len(nodes)), sumZ/float64(len(nodes)))
}

func (n *Node) AddTexture(texturePath string, uvCoords ...*Vertex2D) {
	n.Outer.AddTexture(uvCoords...)
	n.Tag = texturePath
	n.ImageTexture = true
}

func (n *Node) Center() {
	ns := n.Nodes()
	centroid := ns.Centroid()
	ns.Translate(-centroid.X, -centroid.Y, -centroid.Z)
}

func (n *Node) Attach(nodes ...*Node) {
	n.First().InsertAfter(nodes[0].First())
	//node.Align(n)
	//node.First().InsertBefore(n.Last())
}

func (n *Node) First() *Node {
	cur := n
	for cur.Prev != nil {
		cur = cur.Prev
	}
	return cur
}

func (n *Node) Last() *Node {
	cur := n
	for cur.Next != nil {
		cur = cur.Next
	}
	return cur
}

func (n *Node) TagAll(tag string) {
	nodes := n.Nodes()
	for _, node := range nodes {
		node.Tag = tag
	}
}

func (n *Node) TagAllUntagged(tag string) {
	nodes := n.Nodes()
	for _, node := range nodes {
		if node.Tag == "" {
			node.Tag = tag
		}
	}
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

func (node *Node) GenerateColor(name string, _colors ...map[string][3]float64) {
	colors := map[string][3]float64{}
	if len(_colors) > 0 {
		colors = _colors[0]
	}
	nodes := node.Nodes()
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
	faces2 := earcut3d.EarcutFaces(faces, holes...)

	// find all the unique uv coordinates
	uniqueUVs := make(map[[2]float64]int)
	for _, n := range nodes {
		for _, vertex := range n.Outer.Vertices {
			uniqueUVs[[2]float64{vertex.U, vertex.V}] = 1
		}
	}
	uvIndices := make(map[[2]float64]int)
	for uvCoord := range uniqueUVs {
		uvIndices[uvCoord] = len(uvIndices) + 1
	}

	// Create obj file
	f, err := os.Create(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Create a map to store unique vertices and their indices
	vertexIndices := make(map[[3]float64]int)
	currentIndex := 1

	// Write header
	mtlFilename := name[:len(name)-4]
	f.WriteString("mtllib " + mtlFilename + ".mtl\n")
	colors["default"] = [3]float64{1.0, 1.0, 1.0}

	// Write triangles
	for _, face := range faces2 {
		for _, triangleArray := range face {
			for i := 0; i < len(triangleArray); i += 3 {
				// If the vertex hasn't been seen before, write it and store its index
				key := [3]float64{triangleArray[i], triangleArray[i+1], triangleArray[i+2]}
				if _, seen := vertexIndices[key]; !seen {
					f.WriteString(fmt.Sprintf("v %f %f %f\n", triangleArray[i], triangleArray[i+1], triangleArray[i+2]))
					vertexIndices[key] = currentIndex
					currentIndex++
				}
			}
		}
	}

	// Organise triangles by tag
	type TypedTag struct {
		ImageTexture bool
		Path         string
		Index        int
	}
	trianglesByTag := make(map[string][][]float64)
	nodeIndicesByTag := make(map[string][]int)
	metaTag := make(map[string]TypedTag)
	normalIndices := make(map[int]int)
	for i, face := range faces2 {
		tag := nodes[i].Tag
		if tag == "" {
			tag = "default"
		}
		imageTexture := nodes[i].ImageTexture
		if imageTexture {
			path := nodes[i].Tag
			// split on / and take the last element, also remove the .png or .jpg or whatever
			//tag = strings.Split(tag, "/")[len(strings.Split(tag, "/"))-1]
			//tag = strings.Split(tag, ".")[0]
			tag := path
			metaTag[tag] = TypedTag{
				ImageTexture: imageTexture,
				Path:         path,
				Index:        i,
			}
		}
		for _, triangleArray := range face {
			trianglesByTag[tag] = append(trianglesByTag[tag], triangleArray)
			nodeIndicesByTag[tag] = append(nodeIndicesByTag[tag], i)
		}

		// Write normal
		normal := nodes[i].Outer.Normal()
		f.WriteString(fmt.Sprintf("vn %f %f %f\n", normal.X, normal.Y, normal.Z))
		normalIndices[i] = i + 1
	}

	// Write texture coordinates
	sortedUvIndices := make([][2]float64, len(uvIndices))
	for uvCoord, index := range uvIndices {
		sortedUvIndices[index-1] = uvCoord
	}
	for _, uvCoord := range sortedUvIndices {
		f.WriteString(fmt.Sprintf("vt %f %f\n", uvCoord[0], uvCoord[1]))
	}

	// Write faces
	for tag, triangles := range trianglesByTag {
		f.WriteString(fmt.Sprintf("usemtl %s\n", tag))
		for t_i, triangleArray := range triangles {
			f.WriteString("f")
			for i := 0; i < len(triangleArray); i += 3 {
				key := [3]float64{triangleArray[i], triangleArray[i+1], triangleArray[i+2]}

				// check if the vertex has a uv coordinate
				uvIndex := -1
				meta := metaTag[tag]
				if meta.ImageTexture {
					for _, idx := range nodeIndicesByTag[tag] {
						for _, n := range nodes[idx].Outer.Vertices {
							if n.X == key[0] && n.Y == key[1] && n.Z == key[2] {
								uvIndex = uvIndices[[2]float64{n.U, n.V}]
								break
							}
						}
					}
				}

				// Write vertex, texture, and normal indices
				normalIndex := normalIndices[nodeIndicesByTag[tag][t_i]]
				if uvIndex != -1 {
					f.WriteString(fmt.Sprintf(" %d/%d/%d", vertexIndices[key], uvIndex, normalIndex))
				} else {
					f.WriteString(fmt.Sprintf(" %d//%d", vertexIndices[key], normalIndex))
				}
				normalIndex += 1
			}
			f.WriteString("\n")
		}
	}

	// Create mtl file
	// remove .obj from name
	f, err = os.Create(mtlFilename + ".mtl")
	if err != nil {
		panic(err)
	}

	// Write materials
	for tag, _ := range trianglesByTag {
		f.WriteString(fmt.Sprintf("newmtl %s\n", tag))
		meta := metaTag[tag]
		if meta.ImageTexture {
			f.WriteString(fmt.Sprintf("map_Kd %s\n", meta.Path))
		} else {
			color := colors[tag]
			f.WriteString(fmt.Sprintf("Kd %f %f %f\n", color[0], color[1], color[2]))
		}
	}
}

func (n *Node) getBounds() (minX, minY, minZ, maxX, maxY, maxZ float64) {
	nodes := n.Nodes()
	minX, minY, minZ = math.Inf(1), math.Inf(1), math.Inf(1)
	maxX, maxY, maxZ = math.Inf(-1), math.Inf(-1), math.Inf(-1)

	for _, node := range nodes {
		for _, vertex := range node.Outer.Vertices {
			minX = math.Min(minX, vertex.X)
			minY = math.Min(minY, vertex.Y)
			minZ = math.Min(minZ, vertex.Z)
			maxX = math.Max(maxX, vertex.X)
			maxY = math.Max(maxY, vertex.Y)
			maxZ = math.Max(maxZ, vertex.Z)
		}
	}
	return minX, minY, minZ, maxX, maxY, maxZ
}

func (n *Node) FitRectangularPrism(prism [6]*Face3D) {
	// Get bounds of the rectangular prism
	var minX, minY, minZ, maxX, maxY, maxZ float64
	for _, face := range prism {
		for _, v := range face.Vertices {
			minX = math.Min(minX, v.X)
			minY = math.Min(minY, v.Y)
			minZ = math.Min(minZ, v.Z)
			maxX = math.Max(maxX, v.X)
			maxY = math.Max(maxY, v.Y)
			maxZ = math.Max(maxZ, v.Z)
		}
	}

	// Get bounds of the node
	nodeMinX, nodeMinY, nodeMinZ, nodeMaxX, nodeMaxY, nodeMaxZ := n.getBounds()

	// Calculate node and prism dimensions
	nodeDimX := nodeMaxX - nodeMinX
	nodeDimY := nodeMaxY - nodeMinY
	nodeDimZ := nodeMaxZ - nodeMinZ
	prismDimX := maxX - minX
	prismDimY := maxY - minY
	prismDimZ := maxZ - minZ

	// Calculate scale factors
	scaleX := prismDimX / nodeDimX
	scaleY := prismDimY / nodeDimY
	scaleZ := prismDimZ / nodeDimZ

	// Translate and scale all connected nodes
	nodes := n.Nodes()
	nodes.Translate(-nodeMinX, -nodeMinY, -nodeMinZ)
	nodes.Scale(scaleX, scaleY, scaleZ)
	nodes.Translate((minX+maxX)/2, (minY+maxY)/2, (minZ+maxZ)/2)
}

type Nodes []*Node

func (ns Nodes) Centroid() *Vertex3D {
	var sumX, sumY, sumZ float64
	for _, node := range ns {
		sumX += node.Outer.Centroid().X
		sumY += node.Outer.Centroid().Y
		sumZ += node.Outer.Centroid().Z
	}
	return NewVertex3D(sumX/float64(len(ns)), sumY/float64(len(ns)), sumZ/float64(len(ns)))
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

func (ns Nodes) Mirror(axis Axis) {
	uniques := ns.UniqueVertices()
	for _, v := range uniques {
		v.Mirror(axis)
	}
}

func (ns Nodes) Translate(x, y, z float64) {
	uniques := ns.UniqueVertices()
	for _, v := range uniques {
		v.Translate(x, y, z)
	}
}

func (ns Nodes) Mul(magnitude float64) {
	uniques := ns.UniqueVertices()
	cen := ns.Centroid()
	for _, v := range uniques {
		v.Translate(-cen.X, -cen.Y, -cen.Z)
		v.Mul(magnitude)
		v.Translate(cen.X, cen.Y, cen.Z)
	}
}

func (ns Nodes) Scale(x, y, z float64) {
	uniques := ns.UniqueVertices()
	cen := ns.Centroid()
	for _, v := range uniques {
		v.Translate(-cen.X, -cen.Y, -cen.Z)
		v.Scale(x, y, z)
		v.Translate(cen.X, cen.Y, cen.Z)
	}
}

func (ns Nodes) ScaleFixed(x, y, z float64) {
	uniques := ns.UniqueVertices()
	for _, v := range uniques {
		v.Scale(x, y, z)
	}
}

func (ns Nodes) Rotate(deg float64, axis Axis) {
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

func (ns Nodes) SetMeta(key string, value interface{}) {
	for _, node := range ns {
		node.SetMeta(key, value)
	}
}

func (ns Nodes) Width() float64 {
	var w float64
	for _, n := range ns {
		if n == nil {
			continue
		}
		_w := 0.0
		for _, v := range n.Outer.Vertices {
			if v.X > _w {
				_w = v.X
			}
		}
		if _w > w {
			w = _w
		}

	}
	return w
}

func (ns Nodes) Height() float64 {
	var h float64
	for _, n := range ns {
		if n == nil {
			continue
		}
		_h := 0.0
		for _, v := range n.Outer.Vertices {
			if v.Y > _h {
				_h = v.Y
			}
		}
		if _h > h {
			h = _h
		}
	}
	return h
}

func (ns Nodes) Contains(node *Node) bool {
	for _, n := range ns {
		if n == node {
			return true
		}
	}
	return false
}

func (ns Nodes) Filter(tags ...string) Nodes {
	var nodes Nodes
	for _, node := range ns {
		for _, tag := range tags {
			if node.Tag == tag {
				nodes = append(nodes, node)
			}
		}
	}
	return nodes
}

func (ns Nodes) Tag(tag string) {
	for _, node := range ns {
		node.Tag = tag
	}
}

func (ns Nodes) RenameTag(prev, next string) {
	for _, node := range ns {
		if node.Tag == prev {
			node.Tag = next
		}
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

func (nodes Nodes) CopyAll() Nodes {
	uniques := nodes.UniqueVertices()

	// copy vertices
	uniques2 := make([]*Vertex3D, 0)
	for _, v := range uniques {
		uniques2 = append(uniques2, v.Copy())
	}
	res := Nodes{}
	var prev *Node
	for _, node := range nodes {
		// lookup new vertices
		verts := make([]*Vertex3D, len(node.Outer.Vertices))
		for i, v := range node.Outer.Vertices {
			for j, v2 := range uniques {
				if v == v2 {
					verts[i] = uniques2[j]
				}
			}
		}
		holes := make([]*Face3D, len(node.Inner))
		for i, f := range node.Inner {
			verts := make([]*Vertex3D, len(f.Vertices))
			for i, v := range f.Vertices {
				for j, v2 := range uniques {
					if v == v2 {
						verts[i] = uniques2[j]
					}
				}
			}
			holes[i] = &Face3D{Vertices: verts}
		}
		_node := NewTaggedNode(node.Tag, &Face3D{Vertices: verts}, holes...)
		_node.ImageTexture = node.ImageTexture
		_node.Meta = deepCopyMap(node.Meta)
		if prev != nil {
			prev.Next = _node
			_node.Prev = prev
		}
		prev = _node
		res = append(res, _node)
	}

	return res
}

func deepCopyMap(originalMap map[string]interface{}) map[string]interface{} {
	newMap := make(map[string]interface{})
	for key, value := range originalMap {
		newValue := deepCopyValue(value)
		newMap[key] = newValue
	}
	return newMap
}

func deepCopyValue(value interface{}) interface{} {
	if value == nil {
		return nil
	}

	// Use reflection to handle arbitrary types
	originalValue := reflect.ValueOf(value)

	switch originalValue.Kind() {
	case reflect.Map:
		// Create a new map with the same type as the original
		newMap := reflect.MakeMap(originalValue.Type())
		for _, key := range originalValue.MapKeys() {
			originalMapValue := originalValue.MapIndex(key)
			// Recursively deep copy the map value
			newMap.SetMapIndex(key, reflect.ValueOf(deepCopyValue(originalMapValue.Interface())))
		}
		return newMap.Interface()
	case reflect.Slice:
		return deepCopySlice(originalValue)
	default:
		return value
	}
}

func deepCopySlice(original reflect.Value) interface{} {
	copy := reflect.MakeSlice(original.Type(), original.Len(), original.Cap())
	for i := 0; i < original.Len(); i++ {
		copy.Index(i).Set(reflect.ValueOf(deepCopyValue(original.Index(i).Interface())))
	}
	return copy.Interface()
}

func (ns Nodes) Copy() Nodes {
	ns2 := make(Nodes, 0)
	for _, node := range ns {
		ns2 = append(ns2, node.Copy())
	}
	ns2.LinkVertices()
	ns2.LinkNodes()
	return ns2
}

func (ns Nodes) Detach() Nodes {
	var nodes Nodes
	for _, node := range ns {
		_node := node.Detach()
		nodes = append(nodes, _node)
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

func (ns Nodes) Rotate2D(deg float64) {
	for _, node := range ns {
		node.Rotate2D(deg)
	}
}

func (ns Nodes) Translate2D(x, y float64) {
	for _, node := range ns {
		node.Translate2D(x, y)
	}
}

func (ns Nodes) Snap(point *Vertex3D) {
	// find closest vertex
	var closest *Vertex3D
	var closestDist float64
	uniques := ns.UniqueVertices()
	for _, vertex := range uniques {
		dist := vertex.Distance(point)
		if closest == nil || dist < closestDist {
			closest = vertex
			closestDist = dist
		}
	}

	// snap all to closest vertex
	ns.Translate(point.X-closest.X, point.Y-closest.Y, point.Z-closest.Z)
}

func (ns Nodes) Get(tag string) Nodes {
	return ns[0].GetAll(tag)
}

func (ns Nodes) GetAll(tag string) Nodes {
	return ns[0].GetAll(tag)
}

func (ns Nodes) JoinIfWithin(distance float64) {
	// find the closest nodes and join them if they are within distance
	// and if they are not already joined
	// this is a naive implementation and will not work for all cases
	// but it is good enough for simple cases

	// find all the unique vertices
	uniques := ns.UniqueVertices()

	// find the closest vertices
	// and join them if they are within distance
	for i, v := range uniques {
		for j := i + 1; j < len(uniques); j++ {
			v2 := uniques[j]
			if v.Distance(v2) < distance {
				for _, node := range ns {
					faces := node.Faces()
					for _, face := range faces {
						for k, v3 := range face.Vertices {
							if v3.Equals(v2) {
								face.Vertices[k] = v
							}
						}
					}
				}
			}
		}
	}

	ns.LinkVertices()
}

func (ns Nodes) LinkVertices(_setEqual ...bool) {
	setEqual := true
	if len(_setEqual) > 0 {
		setEqual = _setEqual[0]
	}
	uniques := ns.UniqueVertices()
	for _, node := range ns {
		// check if vertices are in uniques
		// and if so replace them with the unique vertex
		faces := node.Faces()
		for _, face := range faces {
			for i, v := range face.Vertices {
				for _, v2 := range uniques {
					if v.Equals(v2) {
						if setEqual {
							face.Vertices[i] = v2
						} else {
							face.Vertices[i].X = v2.X
							face.Vertices[i].Y = v2.Y
							face.Vertices[i].Z = v2.Z
						}
					}
				}
			}
		}
	}
}

func (ns Nodes) LinkNodes() {
	// Check if there is at least one node
	if len(ns) < 1 {
		return
	}

	// Iterate over all nodes, setting the next node for each one
	for i := 0; i < len(ns)-1; i++ {
		ns[i].Next = ns[i+1]
		ns[i+1].Prev = ns[i]
	}
}

func (ns Nodes) RoundVertices(precision int) {
	for _, node := range ns {
		node.RoundVertices(precision)
	}
}

func (ns Nodes) Attach(node *Node) {
	for _, n := range ns {
		_node := node.CopyAll()[0]
		n.Attach(_node)
	}
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
