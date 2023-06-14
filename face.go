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
		v = append(v, &Vertex2D{vertices[i], vertices[i+1]})
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
	return &Vertex2D{x / float64(len(f.Vertices)), y / float64(len(f.Vertices))}
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

// 3D
type Face3D struct {
	Vertices  []*Vertex3D
	PercShape *Face2D
}

func NewFace3D(vertices ...float64) *Face3D {
	var v []*Vertex3D
	for i := 0; i < len(vertices); i += 3 {
		v = append(v, &Vertex3D{vertices[i], vertices[i+1], vertices[i+2]})
	}
	return &Face3D{
		Vertices: v,
	}
}

func (f *Face3D) Normal() *Vertex3D {
	var normal Vertex3D

	v := f.Vertices
	for i := range v {
		v1 := *v[i]
		v2 := *v[(i+1)%len(v)]

		normal.X += (v1.Y - v2.Y) * (v1.Z + v2.Z)
		normal.Y += (v1.Z - v2.Z) * (v1.X + v2.X)
		normal.Z += (v1.X - v2.X) * (v1.Y + v2.Y)
	}

	// Normalize the normal vector
	length := math.Sqrt(normal.X*normal.X + normal.Y*normal.Y + normal.Z*normal.Z)
	normal.X /= length
	normal.Y /= length
	normal.Z /= length

	return &normal
}

func (f *Face3D) Centroid() *Vertex3D {
	var x, y, z float64
	for _, vertex := range f.Vertices {
		x += vertex.X
		y += vertex.Y
		z += vertex.Z
	}
	return &Vertex3D{x / float64(len(f.Vertices)), y / float64(len(f.Vertices)), z / float64(len(f.Vertices))}
}

func (f *Face3D) Translate(x, y, z float64) {
	for _, vertex := range f.Vertices {
		vertex.Translate(x, y, z)
	}
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
	copy := NewFace3D(f.Flatten()...)
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
