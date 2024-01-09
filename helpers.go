package jps

import (
	"math"
)

func makeLine(start, goal Point) []Point {
	var line []Point
	dx := goal.X - start.X
	dy := goal.Y - start.Y
	steps := maxInt(absInt(dx), absInt(dy))
	if steps == 0 {
		return append(line, start)
	}
	for i := 0; i <= steps; i++ {
		x := float64(start.X) + float64(i)*float64(dx)/float64(steps)
		y := float64(start.Y) + float64(i)*float64(dy)/float64(steps)
		line = append(line, Point{int(x), int(y)})
	}
	return line
}

func estimateDistance(start, goal Point) float64 {
	return math.Sqrt(math.Pow(float64(start.X-goal.X), 2) + math.Pow(float64(start.Y-goal.Y), 2))
}

func direction(from, to Point) Point {
	d := sub(to, from)
	return Point{signInt(d.X), signInt(d.Y)}
}

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
		d := direction(current, next)
		for current != next {
			result = append(result, current)
			current = add(current, d)
		}
	}
	result = append(result, start)
	return reverse(result)
}

func prepareForcedDirections(obstacles [][]bool, current, direction Point) []Point {
	directions := make([]Point, 0)
	dx, dy := direction.X, direction.Y
	if absInt(dx)+absInt(dy) == 2 {
		return directions
	}
	cx, cy := current.X, current.Y
	if dy == 0 {
		if checkPoints(obstacles, Point{cx - dx, cy - 1}) != pointCheckPassable && checkPoints(obstacles, Point{cx, cy - 1}) == pointCheckPassable {
			directions = append(directions, Point{0, -1}, Point{dx, -1})
		}
		if checkPoints(obstacles, Point{cx - dx, cy + 1}) != pointCheckPassable && checkPoints(obstacles, Point{cx, cy + 1}) == pointCheckPassable {
			directions = append(directions, Point{0, 1}, Point{dx, 1})
		}
	}
	if dx == 0 {
		if checkPoints(obstacles, Point{cx - 1, cy - dy}) != pointCheckPassable && checkPoints(obstacles, Point{cx - 1, cy}) == pointCheckPassable {
			directions = append(directions, Point{-1, 0}, Point{-1, dy})
		}
		if checkPoints(obstacles, Point{cx + 1, cy - dy}) != pointCheckPassable && checkPoints(obstacles, Point{cx + 1, cy}) == pointCheckPassable {
			directions = append(directions, Point{1, 0}, Point{1, dy})
		}
	}
	return directions
}

func prepareSubordinatedDirections(direction Point) []Point {
	directions := make([]Point, 0, 2)
	if absInt(direction.X)+absInt(direction.Y) == 2 {
		directions = append(directions, Point{direction.X, 0}, Point{0, direction.Y})
	}
	return directions
}

func prepareDirections(obstacles [][]bool, current, predecessor Point) []Point {
	if predecessor != (Point{}) {
		d := direction(predecessor, current)
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

func prepareCandidate(obstacles [][]bool, current, direction, goal Point, price float64) pricedPoint {
	candidate := add(current, direction)
	if checkPoints(obstacles, candidate) != pointCheckPassable || isCornerCut(obstacles, current, direction) {
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
