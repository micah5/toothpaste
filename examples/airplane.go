package main

import (
	"github.com/micah5/toothpaste"
)

func main() {
	w, d, h := 1.0, 1.0, 4.0

	// ***Back of plane***
	fuselage := toothpaste.NewNode(toothpaste.Circle(w, d, 10).To3D(toothpaste.ZAxis))

	// tag the backs as "inner" so we can remove it later when we are done extruding
	back := fuselage.Extrude(h/2, "inner")

	// **tailfin**
	// body of airplane
	back = back.Extrude(h/3, "inner", "", "", "tailfin")
	back.Mul2D(0.5)
	back.Translate2D(0, 0.1)
	tail := back.Get("tailfin")

	// tail
	tailfin := tail.ExtrudeDrop(-h / 3)
	tailfin.Mul2D(0.3)
	tailfin.Translate(0, 0, -1)

	// **flaps**
	// body of airplane
	back = back.Extrude(h/6, "back", "flap", "", "", "", "flap")
	back.Mul2D(0.3)
	back.Translate2D(0, 0.1)
	back.Flip()

	// flaps
	flaps := back.GetAll("flap")
	flapEnd := flaps.ExtrudeDrop(-h / 8)
	flapEnd.Mul2D(0.5)

	// ***middle of plane***
	fuselage = fuselage.ExtrudeDrop(-h, "inner", "", "", "", "", "", "wing", "", "", "", "wing")

	// wings
	hump := fuselage.GetAll("wing")
	hump = hump.ExtrudeDrop(-h / 16)
	hump.Mul2D(0.5)
	hump[0].Rotate(16, toothpaste.ZAxis)
	hump[1].Rotate(-16, toothpaste.ZAxis)
	wing := hump.ExtrudeDrop(-h/3, "", "", "b2", "", "b1")
	wing.Mul2D(0.8)
	wing.Translate(0, 0.2, -0.5)
	wing2 := wing.ExtrudeDrop(-2 * h / 3)
	wing2.Scale2D(0.3, 1)
	wing2.Translate(0, 0.5, -1.0)
	tip := wing2.ExtrudeDrop(-0.2, "", "", "top")
	tip.Scale2D(0.5, 1)
	tip.Translate(0, 0.2, -0.1)

	// ***engines***
	// mounting
	// the wings are mirrored so the bottom is on different sides
	// hence the use of "b1" and "b2" tags
	var mount toothpaste.Nodes
	for i, tag := range []string{"b1", "b2"} {
		bottom := wing.GetAll(tag)[i]
		_mount := bottom.ExtrudeDrop(-h / 18)
		_mount.Mul2D(0.3)
		_mount.Translate(0, 0.0, 0.75)
		mount = append(mount, _mount)
	}

	// engine
	outer := toothpaste.Circle(w/3, w/3, 8)
	inner := toothpaste.Circle(w/4, w/4, 8)
	inner.Translate(w/24, w/24)
	engine := toothpaste.NewNode(outer.To3D(toothpaste.ZAxis), inner.To3D(toothpaste.ZAxis))
	intake := engine.Extrude(w / 8)
	intake.Mul2D(1.15)
	engine = intake.ExtrudeOuter(w / 8)
	intake.Outer = intake.Inner[0]
	intake.Inner = nil
	engine.Mul2D(1.15)
	engine = engine.ExtrudeDrop(w / 4)
	engine.Mul2D(0.7)
	engine.Flip()
	// the engine starts off as a seperate node
	// so we need to attach it to the other nodes
	// so that it will be included in the final model
	println(len(fuselage.Nodes()), len(engine.Nodes()))
	mount.Attach(engine)
	println(len(fuselage.Nodes()))

	// remove the inner parts of the fuselage
	fuselage.GetAll("inner").Drop()

	fuselage.Center()
	fuselage.Generate("airplane.obj")
}
