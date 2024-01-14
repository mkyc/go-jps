package jps

import "math"

type obstacles [][]bool

func (o obstacles) isPointOutsideMap(point Point) bool {
	if point.X < 0 || point.Y < 0 || point.X >= len(o) || point.Y >= len(o[0]) {
		return true
	}
	return false
}

func (o obstacles) isPointInsideObstacle(point Point) bool {
	return o[point.X][point.Y]
}

func (o obstacles) isPointPassable(p Point) bool {
	if p.X < 0 || p.Y < 0 || p.X >= len(o) || p.Y >= len(o[0]) {
		return false
	}
	return !o[p.X][p.Y]
}

func (o obstacles) isCornerCut(point, direction Point) bool {
	for _, d := range o.prepareSubordinatedDirections(direction) {
		if !o.isPointPassable(point.add(d)) {
			return true
		}
	}
	return false
}

func (o obstacles) isLinePassable(line []Point) bool {
	for _, p := range line {
		if o[p.X][p.Y] {
			return false
		}
	}
	return true
}

func (o obstacles) prepareDirections(current, predecessor Point) []Point {
	if predecessor != (Point{}) {
		d := predecessor.directionTo(current)
		directions := []Point{d}
		directions = append(directions, o.prepareForcedDirections(current, d)...)
		directions = append(directions, o.prepareSubordinatedDirections(d)...)
		return directions
	} else {
		return []Point{
			{0, -1},
			{1, -1},
			{1, 0},
			{1, 1},
			{0, 1},
			{-1, 1},
			{-1, 0},
			{-1, -1},
		}
	}
}

func (o obstacles) prepareForcedDirections(current, direction Point) []Point {
	directions := make([]Point, 0)
	dx, dy := direction.X, direction.Y
	if abs(dx)+abs(dy) == 2 {
		return directions
	}
	cx, cy := current.X, current.Y
	if dy == 0 {
		if !o.isPointPassable(Point{cx - dx, cy - 1}) && o.isPointPassable(Point{cx, cy - 1}) {
			directions = append(directions, Point{0, -1}, Point{dx, -1})
		}
		if !o.isPointPassable(Point{cx - dx, cy + 1}) && o.isPointPassable(Point{cx, cy + 1}) {
			directions = append(directions, Point{0, 1}, Point{dx, 1})
		}
	}
	if dx == 0 {
		if !o.isPointPassable(Point{cx - 1, cy - dy}) && o.isPointPassable(Point{cx - 1, cy}) {
			directions = append(directions, Point{-1, 0}, Point{-1, dy})
		}
		if !o.isPointPassable(Point{cx + 1, cy - dy}) && o.isPointPassable(Point{cx + 1, cy}) {
			directions = append(directions, Point{1, 0}, Point{1, dy})
		}
	}
	return directions
}

func (o obstacles) prepareSubordinatedDirections(direction Point) []Point {
	directions := make([]Point, 0)
	if abs(direction.X)+abs(direction.Y) == 2 {
		directions = append(directions, Point{direction.X, 0}, Point{0, direction.Y})
	}
	return directions
}

func (o obstacles) prepareCandidate(current, direction, goal Point, price float64) pricedPoint {
	candidate := current.add(direction)
	if !o.isPointPassable(candidate) || o.isCornerCut(current, direction) {
		return pricedPoint{}
	}
	currentPrice := price + math.Sqrt(math.Abs(float64(direction.X))+math.Abs(float64(direction.Y)))
	if candidate == goal {
		return pricedPoint{candidate, currentPrice}
	}
	if len(o.prepareForcedDirections(candidate, direction)) > 0 {
		return pricedPoint{candidate, currentPrice}
	}
	for _, d := range o.prepareSubordinatedDirections(direction) {
		if c := o.prepareCandidate(candidate, d, goal, currentPrice); c != (pricedPoint{}) {
			return pricedPoint{candidate, currentPrice}
		}
	}
	return o.prepareCandidate(candidate, direction, goal, currentPrice)
}

func (o obstacles) prepareCandidates(current, predecessor, goal Point) []pricedPoint {
	directions := o.prepareDirections(current, predecessor)
	candidates := make([]pricedPoint, 0)
	for _, d := range directions {
		candidate := o.prepareCandidate(current, d, goal, 0)
		if candidate != (pricedPoint{}) {
			candidates = append(candidates, candidate)
		}
	}
	return candidates
}
