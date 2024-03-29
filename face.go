package toothpaste

import (
	"github.com/micah5/earcut-3d"
	"github.com/micah5/exhaustive-fitter"
	"math"
)

type ProjectionDetails struct {
	Basis    []float64
	RefPoint [3]float64
	Face3D   *Face3D
}

// 2D
type Face2D struct {
	Vertices  []*Vertex2D
	PD        *ProjectionDetails
	PercShape *Face2D
}

func NewFace2D(vertices ...float64) *Face2D {
	var v []*Vertex2D
	for i := 0; i < len(vertices); i += 2 {
		v = append(v, NewVertex2D(vertices[i], vertices[i+1]))
	}
	return &Face2D{
		Vertices: v,
	}
}

func (f *Face2D) Centroid() *Vertex2D {
	var x, y float64
	for _, vertex := range f.Vertices {
		x += vertex.X
		y += vertex.Y
	}
	return NewVertex2D(x/float64(len(f.Vertices)), y/float64(len(f.Vertices)))
}

func (f *Face2D) Translate(x, y float64) {
	for _, vertex := range f.Vertices {
		vertex.Translate(x, y)
	}
}

func (f *Face2D) Scale(x, y float64) {
	cen := f.Centroid()
	f.Translate(-cen.X, -cen.Y)
	f.ScaleFixed(x, y)
	f.Translate(cen.X, cen.Y)
}

func (f *Face2D) ScaleFixed(x, y float64) {
	for _, vertex := range f.Vertices {
		vertex.Scale(x, y)
	}
}

func (f *Face2D) Rotate(deg float64) {
	cen := f.Centroid()
	f.Translate(-cen.X, -cen.Y)
	f.RotateFixed(deg)
	f.Translate(cen.X, cen.Y)
}

func (f *Face2D) RotateFixed(deg float64) {
	for _, vertex := range f.Vertices {
		vertex.Rotate(deg)
	}
}

func (f *Face2D) ContainsExact(v *Vertex2D) bool {
	for _, vertex := range f.Vertices {
		if vertex.X == v.X && vertex.Y == v.Y {
			return true
		}
	}
	return false
}

func (f *Face2D) Flatten() []float64 {
	flattened := make([]float64, len(f.Vertices)*2)
	for i, vertex := range f.Vertices {
		flattened[i*2] = vertex.X
		flattened[i*2+1] = vertex.Y
	}
	return flattened
}

func (f *Face2D) MulFixed(magnitude float64) {
	for _, vertex := range f.Vertices {
		vertex.Mul(magnitude)
	}
}

func (f *Face2D) Mul(magnitude float64) {
	cen := f.Centroid()
	f.Translate(-cen.X, -cen.Y)
	f.MulFixed(magnitude)
	f.Translate(cen.X, cen.Y)
}

func (inner *Face2D) Fit2D(outer *Face2D) {
	inner.PercShape = inner.Copy()
	// Flip the inner shape upside down
	for _, vertex := range inner.Vertices {
		vertex.Y = 1.0 - vertex.Y
	}
	result, err := fitter.Transform(inner.Flatten(), outer.Flatten())
	if err != nil {
		println(err)
	}
	for i, vertex := range inner.Vertices {
		vertex.X = result[i*2]
		vertex.Y = result[i*2+1]
	}
}

func (inner *Face2D) Fit3D(face3D *Face3D) {
	outer := face3D.To2D()
	inner.Fit2D(outer)
	inner.PD = outer.PD
}

func (f *Face2D) Copy() *Face2D {
	copy := NewFace2D(f.Flatten()...)
	if f.PD != nil {
		copy.PD = f.PD
	}
	return copy
}

func (f *Face2D) ToProjection3D(createNewFace bool) *Face3D {
	var face3D *Face3D
	if f.PD == nil {
		println("No projection details")
	} else {
		points3D := earcut3d.ProjectShapeTo3D(f.Flatten(), f.PD.Basis, f.PD.RefPoint[:])
		if createNewFace {
			face3D = NewFace3D(points3D...)
		} else {
			face3D = f.PD.Face3D
			for i := 0; i < len(points3D); i += 3 {
				face3D.Vertices[i/3].X = points3D[i]
				face3D.Vertices[i/3].Y = points3D[i+1]
				face3D.Vertices[i/3].Z = points3D[i+2]
			}
		}
	}
	return face3D
}

func (f *Face2D) MinMax() (min, max *Vertex2D) {
	min = NewVertex2D(math.MaxFloat64, math.MaxFloat64)
	max = NewVertex2D(-math.MaxFloat64, -math.MaxFloat64)
	for _, vertex := range f.Vertices {
		if vertex.X < min.X {
			min.X = vertex.X
		}
		if vertex.Y < min.Y {
			min.Y = vertex.Y
		}
		if vertex.X > max.X {
			max.X = vertex.X
		}
		if vertex.Y > max.Y {
			max.Y = vertex.Y
		}
	}
	return
}

func (f *Face2D) Height() float64 {
	min, max := f.MinMax()
	return max.Y - min.Y
}

func (f *Face2D) Width() float64 {
	min, max := f.MinMax()
	return max.X - min.X
}

func (f *Face2D) TopLeftVertex() *Vertex2D {
	// Return a vertex with the smallest x and y values
	// It does not need to be an existing vertex
	minX, minY := math.MaxFloat64, math.MaxFloat64
	var topLeft *Vertex2D
	for _, vertex := range f.Vertices {
		if vertex.X < minX {
			minX = vertex.X
		}
		if vertex.Y < minY {
			minY = vertex.Y
		}
	}
	topLeft = NewVertex2D(minX, minY)
	return topLeft
}

func (f *Face2D) ToFixed3D(axis ...Axis) *Face3D {
	face3D := NewFace3D()
	for _, vertex := range f.Vertices {
		face3D.Vertices = append(face3D.Vertices, vertex.To3D(axis...))
	}
	return face3D
}

func (f *Face2D) To3D(params ...interface{}) *Face3D {
	if len(params) == 0 {
		return f.ToFixed3D()
	}
	param := params[0]
	switch param.(type) {
	case bool:
		face3D := f.ToProjection3D(param.(bool))
		face3D.PercShape = f.PercShape
		return face3D
	case Axis:
		return f.ToFixed3D(param.(Axis))
	default:
		return nil
	}
}

func (f *Face2D) Bounds() (minX, minY, maxX, maxY float64) {
	minX, minY = f.Vertices[0].X, f.Vertices[0].Y
	maxX, maxY = minX, minY

	for _, v := range f.Vertices {
		if v.X < minX {
			minX = v.X
		}
		if v.Y < minY {
			minY = v.Y
		}
		if v.X > maxX {
			maxX = v.X
		}
		if v.Y > maxY {
			maxY = v.Y
		}
	}

	return
}

func (f *Face2D) MapUVs(minX, minY, maxX, maxY float64) {
	pMinX, pMinY, pMaxX, pMaxY := f.Bounds()

	// Map the polygon to the specified bounds
	for i, v := range f.Vertices {
		normalizedX := (v.X - pMinX) / (pMaxX - pMinX)
		normalizedY := (v.Y - pMinY) / (pMaxY - pMinY)

		f.Vertices[i].U = lerp(minX, maxX, normalizedX)
		f.Vertices[i].V = lerp(minY, maxY, normalizedY)
	}
}

func (f *Face2D) Find(label string) *Vertex2D {
	for _, vertex := range f.Vertices {
		if vertex.Label == label {
			return vertex
		}
	}
	return nil
}

// 3D
type Face3D struct {
	Vertices  []*Vertex3D
	PercShape *Face2D
}

func NewFace3D(vertices ...float64) *Face3D {
	var v []*Vertex3D
	for i := 0; i < len(vertices); i += 3 {
		v = append(v, NewVertex3D(vertices[i], vertices[i+1], vertices[i+2]))
	}
	return &Face3D{
		Vertices: v,
	}
}

// Normal calculates the normal vector for the Face3D using all vertices
func (f *Face3D) Normal() *Vertex3D {
	vertexCount := len(f.Vertices)
	if vertexCount < 3 {
		// Not enough vertices to define a plane
		return NewVertex3D(0, 0, 0)
	}

	// Calculate the centroid of the face
	centroid := NewVertex3D(0, 0, 0)
	for _, vertex := range f.Vertices {
		centroid = centroid.Add(vertex)
	}
	centroid.X /= float64(vertexCount)
	centroid.Y /= float64(vertexCount)
	centroid.Z /= float64(vertexCount)

	// Initialize the normal vector
	normal := NewVertex3D(0, 0, 0)

	// Iterate over each edge and compute the partial normal
	for i := 0; i < vertexCount; i++ {
		currentVertex := f.Vertices[i]
		nextVertex := f.Vertices[(i+1)%vertexCount] // Wrap-around

		// Edge vector and vector from centroid to current vertex
		edgeVector := nextVertex.Subtract(currentVertex)
		centroidVector := currentVertex.Subtract(centroid)

		// Update the normal with the cross product of centroidVector and edgeVector
		normal = normal.Add(centroidVector.Cross(edgeVector))
	}

	// Normalize the normal vector
	return normal.Normalize()
}

func (f *Face3D) Centroid() *Vertex3D {
	var x, y, z float64
	for _, vertex := range f.Vertices {
		x += vertex.X
		y += vertex.Y
		z += vertex.Z
	}
	return NewVertex3D(x/float64(len(f.Vertices)), y/float64(len(f.Vertices)), z/float64(len(f.Vertices)))
}

func (f *Face3D) Translate(x, y, z float64) {
	for _, vertex := range f.Vertices {
		vertex.Translate(x, y, z)
	}
}

func (f *Face3D) Align(f2 *Face3D) {
	cen := f.Centroid()
	cen2 := f2.Centroid()
	f.Translate(cen.X-cen2.X, cen.Y-cen2.Y, cen.Z-cen2.Z)
}

func (f *Face3D) Snap(point *Vertex3D) {
	// find closest vertex
	var closest *Vertex3D
	var closestDist float64
	for _, vertex := range f.Vertices {
		dist := vertex.Distance(point)
		if closest == nil || dist < closestDist {
			closest = vertex
			closestDist = dist
		}
	}
	// snap face to point
	f.Translate(point.X-closest.X, point.Y-closest.Y, point.Z-closest.Z)
}

func (f *Face3D) MoveTo(x, y, z float64) {
	cen := f.Centroid()
	f.Translate(x-cen.X, y-cen.Y, z-cen.Z)
}

func (f *Face3D) ScaleFixed(x, y, z float64) {
	for _, vertex := range f.Vertices {
		vertex.Scale(x, y, z)
	}
}

func (f *Face3D) Scale(x, y, z float64) {
	cen := f.Centroid()
	f.Translate(-cen.X, -cen.Y, -cen.Z)
	f.ScaleFixed(x, y, z)
	f.Translate(cen.X, cen.Y, cen.Z)
}

func (f *Face3D) MulFixed(magnitude float64) {
	for _, vertex := range f.Vertices {
		vertex.Mul(magnitude)
	}
}

func (f *Face3D) Mul(magnitude float64) {
	cen := f.Centroid()
	f.Translate(-cen.X, -cen.Y, -cen.Z)
	f.MulFixed(magnitude)
	f.Translate(cen.X, cen.Y, cen.Z)
}

func (f *Face3D) RotateFixed(deg float64, axis Axis) {
	for _, vertex := range f.Vertices {
		vertex.Rotate(deg, axis)
	}
}

func (f *Face3D) Rotate(deg float64, axis Axis) {
	cen := f.Centroid()
	f.Translate(-cen.X, -cen.Y, -cen.Z)
	f.RotateFixed(deg, axis)
	f.Translate(cen.X, cen.Y, cen.Z)
}

func (f *Face3D) Translate2D(x, y float64) {
	face2D := f.To2D()
	face2D.Translate(x, y)
	f = face2D.To3D(false)
}

func (f *Face3D) Scale2D(x, y float64) {
	face2D := f.To2D()
	face2D.Scale(x, y)
	f = face2D.To3D(false)
}

func (f *Face3D) Mul2D(magnitude float64) {
	face2D := f.To2D()
	face2D.Mul(magnitude)
	f = face2D.To3D(false)
}

func (f *Face3D) Rotate2D(deg float64) {
	face2D := f.To2D()
	face2D.Rotate(deg)
	f = face2D.To3D(false)
}

func (f *Face3D) RoundVertices(precision int) {
	for _, vertex := range f.Vertices {
		vertex.Round(precision)
	}
}

func (f *Face3D) MinMax(axis Axis) (float64, float64) {
	min := math.MaxFloat64
	max := -math.MaxFloat64
	for _, vertex := range f.Vertices {
		switch axis {
		case XAxis:
			if vertex.X < min {
				min = vertex.X
			}
			if vertex.X > max {
				max = vertex.X
			}
		case YAxis:
			if vertex.Y < min {
				min = vertex.Y
			}
			if vertex.Y > max {
				max = vertex.Y
			}
		case ZAxis:
			if vertex.Z < min {
				min = vertex.Z
			}
			if vertex.Z > max {
				max = vertex.Z
			}
		}
	}
	return min, max
}

func (f *Face3D) Width() float64 {
	min, max := f.MinMax(XAxis)
	return max - min
}

func (f *Face3D) Height() float64 {
	min, max := f.MinMax(YAxis)
	return max - min
}

func (f *Face3D) Depth() float64 {
	min, max := f.MinMax(ZAxis)
	return max - min
}

func (f *Face3D) Mirror(axis Axis) {
	for _, v := range f.Vertices {
		v.Mirror(axis)
	}
}

func (f *Face3D) Flatten() []float64 {
	flattened := make([]float64, len(f.Vertices)*3)
	for i, vertex := range f.Vertices {
		flattened[i*3] = vertex.X
		flattened[i*3+1] = vertex.Y
		flattened[i*3+2] = vertex.Z
	}
	return flattened
}

func (f *Face3D) Copy() *Face3D {
	vertices := make([]*Vertex3D, len(f.Vertices))
	for i, vertex := range f.Vertices {
		vertices[i] = vertex.Copy()
	}
	copy := &Face3D{Vertices: vertices}
	//copy := NewFace3D(f.Flatten()...)
	return copy
}

func (f *Face3D) ContainsExact(vertex *Vertex3D) bool {
	for _, v := range f.Vertices {
		if v == vertex {
			return true
		}
	}
	return false
}

func (f *Face3D) Flip() {
	// Reverse the order of the entire vertices slice to flip the face
	for i, j := 0, len(f.Vertices)-1; i < j; i, j = i+1, j-1 {
		f.Vertices[i], f.Vertices[j] = f.Vertices[j], f.Vertices[i]
	}
}

func (f *Face3D) ShareVertices(other *Face3D) bool {
	for _, vertex := range f.Vertices {
		for _, otherVertex := range other.Vertices {
			if vertex == otherVertex {
				return true
			}
		}
	}
	return false
}

func (f *Face3D) Reverse() {
	f.Flip()
}

func (f *Face3D) To2D() *Face2D {
	inputFace := f.Flatten()
	basis := earcut3d.FindBasis(inputFace)
	points2D := earcut3d.ProjectShapeTo2D(inputFace, basis)
	face2D := NewFace2D(points2D...)
	refPoint := f.Vertices[0]

	face2D.PD = &ProjectionDetails{basis, [3]float64{refPoint.X, refPoint.Y, refPoint.Z}, f}
	return face2D
}

func (f *Face3D) Find(label string) *Vertex3D {
	for _, vertex := range f.Vertices {
		if vertex.Label == label {
			return vertex
		}
	}
	return nil
}

func (f *Face3D) AddTexture(uvCoords ...*Vertex2D) {
	if len(uvCoords) > 0 {
		for i, vertex := range f.Vertices {
			vertex.UV(uvCoords[i].X, uvCoords[i].Y)
		}
	} else if len(f.Vertices) == 3 {
		f.Vertices[0].UV(0, 0)
		f.Vertices[1].UV(1, 0)
		f.Vertices[2].UV(0, 1)
	} else if len(f.Vertices) == 4 {
		f.Vertices[0].UV(0, 0)
		f.Vertices[1].UV(1, 0)
		f.Vertices[2].UV(1, 1)
		f.Vertices[3].UV(0, 1)
	} else {
		f2D := f.To2D()

		// scale f2D so that it fits in a positive 0-1 square
		min, max := f2D.MinMax()

		diffX := max.X - min.X
		diffY := max.Y - min.Y

		for _, vertex := range f2D.Vertices {
			vertex.X = (vertex.X - min.X) / diffX
			vertex.Y = (vertex.Y - min.Y) / diffY
		}

		percFace := NewFace2D(
			0, 0,
			1, 0,
			1, 1,
			0, 1,
		)
		f2D.Fit2D(percFace)
		for i, vertex := range f.Vertices {
			vertex.UV(f2D.Vertices[i].X, f2D.Vertices[i].Y)
		}
	}
}

func lerp(a, b, t float64) float64 {
	return a + (b-a)*t
}
