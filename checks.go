package jps

type pointCheck string

const (
	pointCheckPassable   pointCheck = "passable"
	pointCheckOutsideMap pointCheck = "outside"
	pointCheckInsideObst pointCheck = "obstacle"
)

func checkPoints(obstacles [][]bool, points ...Point) pointCheck {
	for _, p := range points {
		if p.X < 0 || p.Y < 0 || p.X >= len(obstacles) || p.Y >= len(obstacles[0]) {
			return pointCheckOutsideMap
		}
		if obstacles[p.X][p.Y] {
			return pointCheckInsideObst
		}
	}
	return pointCheckPassable
}

func isCornerCut(obstacles [][]bool, point, direction Point) bool {
	for _, v := range prepareSubordinatedDirections(direction) {
		if checkPoints(obstacles, add(point, v)) != pointCheckPassable {
			return true
		}
	}
	return false
}
