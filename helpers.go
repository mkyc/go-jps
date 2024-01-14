package jps

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
