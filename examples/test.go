package main

import (
	"github.com/micah5/toothpaste"
)

func main() {
	w := 1.0
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
	engine.Center()
	engine.Generate("test.obj")
}
