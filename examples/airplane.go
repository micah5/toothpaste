package main

import (
	"github.com/micah5/toothpaste"
)

func main() {
	w, d, h := 1.0, 1.0, 4.0

	// ***Back of plane***
	fuselage := toothpaste.NewNode(toothpaste.Circle(w, d, 10).To3D(toothpaste.ZAxis))

	// tag the backs as "inner" so we can remove it later when we are done extruding
	back := fuselage.Extrude(h, "inner")

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
	fuselage = fuselage.Extrude(-2*h/3, "inner", "", "", "", "", "", "wing", "", "", "", "wing")

	// wings
	hump := fuselage.GetAll("wing")
	hump.ExtrudeDrop(-h / 12)
	hump.Mul2D(0.5)
	//hump.Rotate(20, toothpaste.ZAxis)

	// remove the inner parts of the fuselage
	fuselage.GetAll("inner").Drop()

	fuselage.Center()
	fuselage.Generate("airplane.obj")
}
