package main

import (
	"github.com/micah5/toothpaste"
)

func main() {
	w, d := 1.0, 2.0
	s1 := toothpaste.NewNode(toothpaste.Square(w, d).To3D())
	s1.Translate(5, 5, 5)
	s1.Rotate(45, toothpaste.XAxis)
	s1.Rotate(90, toothpaste.YAxis)
	s1.Rotate(20, toothpaste.ZAxis)

	s2 := toothpaste.NewNode(toothpaste.Circle(w, d, 8).To3D())
	s2.Align(s1)
	s1.Last().InsertAfter(s2)

	s1.Generate("test.obj")
}
