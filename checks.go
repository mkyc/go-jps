package jps

func pointsOutsideMap(obstacles [][]bool, points ...Point) bool {
	for _, p := range points {
		if p.X < 0 || p.Y < 0 || p.X >= len(obstacles) || p.Y >= len(obstacles[0]) {
			return true
		}
	}
	return false
}

func pointsInsideObstacle(obstacles [][]bool, points ...Point) bool {
	for _, p := range points {
		if obstacles[p.X][p.Y] {
			return true
		}
	}
	return false
}

func isLinePassable(obstacles [][]bool, line []Point) bool {
	for _, p := range line {
		if obstacles[p.X][p.Y] {
			return false
		}
	}
	return true
}

func isPointPassable(obstacles [][]bool, p Point) bool {
	if p.X < 0 || p.Y < 0 || p.X >= len(obstacles) || p.Y >= len(obstacles[0]) {
		return false
	}
	return !obstacles[p.X][p.Y]
}

func isCornerCut(obstacles [][]bool, point, direction Point) bool {
	for _, v := range prepareSubordinatedDirections(direction) {
		if !isPointPassable(obstacles, add(point, v)) {
			return true
		}
	}
	return false
}
