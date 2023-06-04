package toothpaste

import (
	"github.com/micah5/earcut-3d"
	"github.com/micah5/exhaustive-fitter"
)

type ProjectionDetails struct {
	Basis    []float64
	RefPoint [3]float64
}

// 2D
type Face2D struct {
	Vertices []*Vertex2D
	PD       *ProjectionDetails
}

func NewFace2D(vertices ...float64) *Face2D {
	f := make(Face2D, len(vertices)/2)
	for i := 0; i < len(vertices); i += 2 {
		v[i/2] = &Vertex2D{vertices[i], vertices[i+1]}
	}
	return &v
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
	for _, vertex := range f.Vertices {
		vertex.Scale(x, y)
	}
	f.Translate(cen.X, cen.Y)
}

func (f *Face2D) Rotate(deg int, axis Axis) {
	cen := f.Centroid()
	f.Translate(-cen.X, -cen.Y)
	for _, vertex := range f.Vertices {
		vertex.Rotate(deg, axis)
	}
	f.Translate(cen.X, cen.Y)
}

func (f *Face2D) Flatten() []float64 {
	flattened := make([]float64, len(v)*2)
	for i, vertex := range f.Vertices {
		flattened[i*2] = vertex.X
		flattened[i*2+1] = vertex.Y
	}
	return flattened
}

func (inner *Face2D) Fit2D(outer *Face2D) {
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
}

func (f *Face2D) Copy() *Face2D {
	copy := NewFace2D(f.Flatten()...)
	if f.PD != nil {
		copy.PD = f.PD
	}
	return copy
}

func (f *Face2D) To3D() *Face3D {
	var face3D *Face3D
	if f.PD != nil {
		points3D := ProjectShapeTo3D(f.Flatten(), f.PD.Basis, f.PD.RefPoint)
		face3D = NewFace3D(points3D...)
	} else {
		face3D = NewFace3D()
		for _, vertex := range f.Vertices {
			face3D.Vertices = append(face3D.Vertices, vertex.To3D(0))
		}
	}
	return face3D
}

// 3D
type Face3D struct {
	Vertices []*Vertex3D
}

func NewFace3D(vertices ...float64) *Face3D {
	f := make(Face3D, len(vertices)/3)
	for i := 0; i < len(vertices); i += 3 {
		v[i/3] = &Vertex3D{vertices[i], vertices[i+1], vertices[i+2]}
	}
	return &v
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

	return normal
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

func (f *Face3D) Scale(x, y, z float64) {
	cen := f.Centroid()
	f.Translate(-cen.X, -cen.Y, -cen.Z)
	for _, vertex := range f.Vertices {
		vertex.Scale(x, y, z)
	}
	f.Translate(cen.X, cen.Y, cen.Z)
}

func (f *Face3D) Rotate(deg int, axis Axis) {
	cen := f.Centroid()
	f.Translate(-cen.X, -cen.Y, -cen.Z)
	for _, vertex := range f.Vertices {
		vertex.Rotate(deg, axis)
	}
	f.Translate(cen.X, cen.Y, cen.Z)
}

func (f *Face3D) Flatten() []float64 {
	flattened := make([]float64, len(v)*3)
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

func (f *Face3D) Flip() {
	for i, j := 0, len(f.Vertices)-1; i < j; i, j = i+1, j-1 {
		f.Vertices[i], f.Vertices[j] = f.Vertices[j], f.Vertices[i]
	}
}

func (f *Face3D) Reverse() {
	f.Flip()
}

func (f *Face3D) To2D() *Face2D {
	inputFace := f.Flatten()
	basis := earcut3d.FindBasis(inputFace)
	points2D := earcut3d.ProjectShapeTo2D(inputFace, basis)
	face2D := NewFace2D(points2D...)
	face2D.PD = &ProjectionDetails{basis, inputFace[:3]}
	return face2D
}
