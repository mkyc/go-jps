package jps

type obstacles [][]bool

func (o obstacles) pointOutsideMap(point Point) bool {
	if point.X < 0 || point.Y < 0 || point.X >= len(o) || point.Y >= len(o[0]) {
		return true
	}
	return false
}

func (o obstacles) pointInsideObstacle(point Point) bool {
	return o[point.X][point.Y]
}

func (o obstacles) isPointPassable(p Point) bool {
	if p.X < 0 || p.Y < 0 || p.X >= len(o) || p.Y >= len(o[0]) {
		return false
	}
	return !o[p.X][p.Y]
}

func (o obstacles) isCornerCut(point, direction Point) bool {
	for _, d := range prepareSubordinatedDirections(direction) {
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
