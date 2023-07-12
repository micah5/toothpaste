package main

import (
	"fmt"
	"github.com/micah5/toothpaste"
)

func main() {
	w, d, h := 1.0, 1.0, 1.0
	h = h
	pad := 0.1
	windowWidth, windowHeight := 0.8, 0.35

	square := toothpaste.Square(w, d)
	frame := toothpaste.NewNode(toothpaste.NewFace3D(
		0.010000, 0.407143, 0.037500, 0.075000, 0.407143, 0.037500, 0.075000, 0.155714, 0.037500, 0.010000, 0.155714, 0.037500,
	))
	frame := toothpaste.NewNode(toothpaste.NewFace3D(
		0.265000, 0.407143, 0.037500, 0.330000, 0.407143, 0.037500, 0.330000, 0.155714, 0.037500, 0.265000, 0.155714, 0.037500,
	))
	frame := toothpaste.NewNode(square.To3D(toothpaste.ZAxis))
	frame.Translate(1, 1, 1)
	frame.Rotate(-30, toothpaste.YAxis)
	windowBottom := toothpaste.Square(windowWidth, windowHeight)
	windowBottom.Translate(pad, pad)
	windowTop := toothpaste.Square(windowWidth, windowHeight)
	windowTop.Translate(pad, pad*2+windowHeight)
	frame.Inner = append(frame.Inner, windowBottom.To3D(toothpaste.ZAxis), windowTop.To3D(toothpaste.ZAxis))
	end := frame.ExtrudeInner(h)
	end.Flip()
	frame.Center()
	frame.Generate("window.obj")
}
