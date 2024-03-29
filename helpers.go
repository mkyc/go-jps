package jps

import "math"

type Point struct {
	X, Y int
}

type pricedPoint struct {
	p     Point
	price float64
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

func nearestToGoal(obstacles [][]bool, start, goal Point) Point {
	points := line(start, goal)
	points = reverse(points)
	for _, p := range points {
		if isPointPassable(obstacles, p) {
			return Point{p.X, p.Y}
		}
	}
	return start
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

func add(point Point, other Point) Point {
	return Point{point.X + other.X, point.Y + other.Y}
}

func distance(point Point, other Point) float64 {
	return math.Sqrt(math.Pow(float64(point.X-other.X), 2) + math.Pow(float64(point.Y-other.Y), 2))
}

func line(point Point, other Point) []Point {
	var line []Point
	dx := other.X - point.X
	dy := other.Y - point.Y
	steps := max(abs(dx), abs(dy))
	if steps == 0 {
		return append(line, point)
	}
	for i := 0; i <= steps; i++ {
		x := float64(point.X) + float64(i)*float64(dx)/float64(steps)
		y := float64(point.Y) + float64(i)*float64(dy)/float64(steps)
		line = append(line, Point{int(x), int(y)})
	}
	return line
}

func direction(point Point, other Point) Point {
	return Point{sign(other.X - point.X), sign(other.Y - point.Y)}
}

func isPointOutsideMap(obstacles [][]bool, point Point) bool {
	if point.X < 0 || point.Y < 0 || point.X >= len(obstacles) || point.Y >= len(obstacles[0]) {
		return true
	}
	return false
}

func isPointInsideObstacle(obstacles [][]bool, point Point) bool {
	return obstacles[point.X][point.Y]
}

func isPointPassable(obstacles [][]bool, point Point) bool {
	if isPointOutsideMap(obstacles, point) {
		return false
	}
	return !obstacles[point.X][point.Y]
}

func isLinePassable(obstacles [][]bool, line []Point) bool {
	for _, p := range line {
		if obstacles[p.X][p.Y] {
			return false
		}
	}
	return true
}

func isCornerCut(obstacles [][]bool, point, direction Point) bool {
	for _, d := range prepareSubordinatedDirections(direction) {
		if !isPointPassable(obstacles, add(point, d)) {
			return true
		}
	}
	return false
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

func prepareForcedDirections(obstacles [][]bool, current, direction Point) []Point {
	directions := make([]Point, 0)
	dx, dy := direction.X, direction.Y
	if abs(dx)+abs(dy) == 2 {
		return directions
	}
	cx, cy := current.X, current.Y
	if dy == 0 {
		if !isPointPassable(obstacles, Point{cx - dx, cy - 1}) && isPointPassable(obstacles, Point{cx, cy - 1}) {
			directions = append(directions, Point{0, -1}, Point{dx, -1})
		}
		if !isPointPassable(obstacles, Point{cx - dx, cy + 1}) && isPointPassable(obstacles, Point{cx, cy + 1}) {
			directions = append(directions, Point{0, 1}, Point{dx, 1})
		}
	}
	if dx == 0 {
		if !isPointPassable(obstacles, Point{cx - 1, cy - dy}) && isPointPassable(obstacles, Point{cx - 1, cy}) {
			directions = append(directions, Point{-1, 0}, Point{-1, dy})
		}
		if !isPointPassable(obstacles, Point{cx + 1, cy - dy}) && isPointPassable(obstacles, Point{cx + 1, cy}) {
			directions = append(directions, Point{1, 0}, Point{1, dy})
		}
	}
	return directions
}

func prepareSubordinatedDirections(direction Point) []Point {
	directions := make([]Point, 0)
	if abs(direction.X)+abs(direction.Y) == 2 {
		directions = append(directions, Point{direction.X, 0}, Point{0, direction.Y})
	}
	return directions
}

func prepareCandidate(obstacles [][]bool, current, direction, goal Point, price float64) pricedPoint {
	candidate := add(current, direction)
	if !isPointPassable(obstacles, candidate) || isCornerCut(obstacles, current, direction) {
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
