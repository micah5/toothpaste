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
	fuselage = fuselage.Extrude(-h, "inner", "", "", "", "", "", "wing", "", "", "", "wing")

	// wings
	hump := fuselage.GetAll("wing")
	hump = hump.ExtrudeDrop(-h / 16)
	hump.Mul2D(0.5)
	hump[0].Rotate(16, toothpaste.ZAxis)
	hump[1].Rotate(-16, toothpaste.ZAxis)
	wing := hump.ExtrudeDrop(-h / 3)
	wing.Mul2D(0.8)
	wing.Translate(0, 0.2, -0.5)
	wing = wing.ExtrudeDrop(-2 * h / 3)
	wing.Scale2D(0.3, 1)
	wing.Translate(0, 0.5, -1.0)
	tip := wing.ExtrudeDrop(-0.1, "", "top")
	tip.Get("top").ExtrudeDrop(-1.0)

	// remove the inner parts of the fuselage
	fuselage.GetAll("inner").Drop()

	fuselage.Center()
	fuselage.Generate("airplane.obj")
}
