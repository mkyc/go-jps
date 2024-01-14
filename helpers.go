package jps

import (
	"math"
)

func reverse(points []Point) []Point {
	for i := len(points)/2 - 1; i >= 0; i-- {
		opp := len(points) - 1 - i
		points[i], points[opp] = points[opp], points[i]
	}
	return points
}

func reconstructPath(predecessors map[Point]Point, start, goal Point) []Point {
	result := make([]Point, 0)
	current := goal
	for current != start {
		next := predecessors[current]
		d := current.directionTo(next)
		for current != next {
			result = append(result, current)
			current = current.add(d)
		}
	}
	result = append(result, start)
	return reverse(result)
}

func prepareForcedDirections(obstacles obstacles, current, direction Point) []Point {
	directions := make([]Point, 0)
	dx, dy := direction.X, direction.Y
	if abs(dx)+abs(dy) == 2 {
		return directions
	}
	cx, cy := current.X, current.Y
	if dy == 0 {
		if !obstacles.isPointPassable(Point{cx - dx, cy - 1}) && obstacles.isPointPassable(Point{cx, cy - 1}) {
			directions = append(directions, Point{0, -1}, Point{dx, -1})
		}
		if !obstacles.isPointPassable(Point{cx - dx, cy + 1}) && obstacles.isPointPassable(Point{cx, cy + 1}) {
			directions = append(directions, Point{0, 1}, Point{dx, 1})
		}
	}
	if dx == 0 {
		if !obstacles.isPointPassable(Point{cx - 1, cy - dy}) && obstacles.isPointPassable(Point{cx - 1, cy}) {
			directions = append(directions, Point{-1, 0}, Point{-1, dy})
		}
		if !obstacles.isPointPassable(Point{cx + 1, cy - dy}) && obstacles.isPointPassable(Point{cx + 1, cy}) {
			directions = append(directions, Point{1, 0}, Point{1, dy})
		}
	}
	return directions
}

func prepareSubordinatedDirections(direction Point) []Point {
	directions := make([]Point, 0, 2)
	if abs(direction.X)+abs(direction.Y) == 2 {
		directions = append(directions, Point{direction.X, 0}, Point{0, direction.Y})
	}
	return directions
}

func prepareDirections(obstacles [][]bool, current, predecessor Point) []Point {
	if predecessor != (Point{}) {
		d := predecessor.directionTo(current)
		directions := []Point{d}
		directions = append(directions, prepareForcedDirections(obstacles, current, d)...)
		directions = append(directions, prepareSubordinatedDirections(d)...)
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

func prepareCandidate(obstacles obstacles, current, direction, goal Point, price float64) pricedPoint {
	candidate := current.add(direction)
	if !obstacles.isPointPassable(candidate) || obstacles.isCornerCut(current, direction) {
		return pricedPoint{}
	}
	currentPrice := price + math.Sqrt(math.Abs(float64(direction.X))+math.Abs(float64(direction.Y)))
	if candidate == goal {
		return pricedPoint{candidate, currentPrice}
	}
	if len(prepareForcedDirections(obstacles, candidate, direction)) > 0 {
		return pricedPoint{candidate, currentPrice}
	}
	for _, d := range prepareSubordinatedDirections(direction) {
		if c := prepareCandidate(obstacles, candidate, d, goal, currentPrice); c != (pricedPoint{}) {
			return pricedPoint{candidate, currentPrice}
		}
	}
	return prepareCandidate(obstacles, candidate, direction, goal, currentPrice)
}

func prepareCandidates(obstacles [][]bool, current, predecessor, goal Point) []pricedPoint {
	directions := prepareDirections(obstacles, current, predecessor)
	candidates := make([]pricedPoint, 0)
	for _, d := range directions {
		candidate := prepareCandidate(obstacles, current, d, goal, 0)
		if candidate != (pricedPoint{}) {
			candidates = append(candidates, candidate)
		}
	}
	return candidates
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func sign(i int) int {
	switch {
	case i > 0:
		return 1
	case i < 0:
		return -1
	}
	return 0
}
