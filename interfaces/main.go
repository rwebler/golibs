package main

import (
	"fmt"
)

type triangle struct {
	h float64 // height
	b float64 // base
}

func (t triangle) getArea() float64 {
	return (t.b * t.h) / 2
}

type square struct {
	s float64 // side
}

func (s square) getArea() float64 {
	return s.s * s.s
}

type shape interface {
	getArea() float64
}

func main() {
	t := triangle{
		h: 15.0,
		b: 20.0,
	}
	s := square{
		s: 20.0,
	}
	printArea(t)
	printArea(s)
}

func printArea(s shape) {
	fmt.Println(s.getArea())
}
