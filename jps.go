package jps

import (
	"container/heap"
	"errors"
)

// Find returns a path from start to goal, or an error in following cases:
// - start or goal is outside the map
// - start or goal is inside an obstacle
// - there is no path from start to goal
func Find(obstacles [][]bool, start, goal Point) ([]Point, error) {
	//check if start and goal are inside the map
	if pointsOutsideMap(obstacles, start, goal) {
		return nil, errors.New("start or goal is outside the map")
	}
	//check if start and goal are inside an obstacle
	if pointsInsideObstacle(obstacles, start, goal) {
		return nil, errors.New("start or goal is inside an obstacle")
	}
	//check if path is straight line
	straightLine := makeLine(start, goal)
	if isLinePassable(obstacles, straightLine) {
		return straightLine, nil
	}
	//check if there is a path from start to goal
	path, err := findPath(obstacles, start, goal)
	if err != nil {
		return nil, err
	}
	return path, nil
}

// TryFind returns a path from start to goal, or an error in following cases:
// - start or goal is outside the map
// If start or goal is inside an obstacle, it will try to find nearest point on the edge of the obstacle.
// If there is no path from start to goal, it will try to find nearest point on the edge of the obstacle.
func TryFind(obstacles [][]bool, start, goal Point) ([]Point, error) {
	panic("TODO")
}

// MustFind returns a path from start to goal, or panics in following cases:
// - start or goal is outside the map
// If start or goal is inside an obstacle, it will try to find nearest point on the edge of the obstacle.
// If there is no path from start to goal, it will try to find nearest point on the edge of the obstacle.
func MustFind(obstacles [][]bool, start, goal Point) []Point {
	panic("TODO")
}

func findPath(obstacles [][]bool, start, goal Point) ([]Point, error) {
	m := make(map[Point]Point)
	q := make(priorityQueue, 0)
	heap.Init(&q)
	heap.Push(&q, &item{
		point:             start,
		predecessor:       Point{},
		distanceFromStart: 0,
		distanceToGoal:    estimateDistance(start, goal),
	})
	for q.Len() > 0 {
		current := heap.Pop(&q).(*item)
		if _, ok := m[current.point]; ok {
			continue
		}
		m[current.point] = current.predecessor
		if current.point == goal {
			return reconstructPath(m, start, goal), nil
		}
		candidates := prepareCandidates(obstacles, current.point, current.predecessor, goal)
		for _, candidate := range candidates {
			if _, ok := m[candidate.p]; !ok {
				heap.Push(&q, &item{
					point:             candidate.p,
					predecessor:       current.point,
					distanceFromStart: current.distanceFromStart + candidate.price,
					distanceToGoal:    estimateDistance(candidate.p, goal),
				})
			}
		}
	}
	return nil, errors.New("no path found")
}
