package jps

import "math"

type Point struct {
	X, Y int
}

func (p Point) add(other Point) Point {
	return Point{p.X + other.X, p.Y + other.Y}
}

func (p Point) distanceTo(other Point) float64 {
	return math.Sqrt(math.Pow(float64(p.X-other.X), 2) + math.Pow(float64(p.Y-other.Y), 2))
}

func (p Point) directionTo(other Point) Point {
	return Point{sign(other.X - p.X), sign(other.Y - p.Y)}
}

func (p Point) lineTo(other Point) []Point {
	var line []Point
	dx := other.X - p.X
	dy := other.Y - p.Y
	steps := max(abs(dx), abs(dy))
	if steps == 0 {
		return append(line, p)
	}
	for i := 0; i <= steps; i++ {
		x := float64(p.X) + float64(i)*float64(dx)/float64(steps)
		y := float64(p.Y) + float64(i)*float64(dy)/float64(steps)
		line = append(line, Point{int(x), int(y)})
	}
	return line
}

type pricedPoint struct {
	p     Point
	price float64
}
