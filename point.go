package jps

type Point struct {
	X, Y int
}

func add(p1, p2 Point) Point {
	return Point{p1.X + p2.X, p1.Y + p2.Y}
}

func sub(p1, p2 Point) Point {
	return Point{p1.X - p2.X, p1.Y - p2.Y}
}

type pricedPoint struct {
	p     Point
	price float64
}
