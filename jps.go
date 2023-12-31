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
	startCheck := checkPoints(obstacles, start)
	goalCheck := checkPoints(obstacles, goal)
	if startCheck == pointCheckOutsideMap {
		return nil, errors.New("start is outside the map")
	}
	if goalCheck == pointCheckOutsideMap {
		return nil, errors.New("goal is outside the map")
	}
	if startCheck == pointCheckInsideObst {
		return nil, errors.New("start is inside an obstacle")
	}
	if goalCheck == pointCheckInsideObst {
		return nil, errors.New("goal is inside an obstacle")
	}
	straightLine := makeLine(start, goal)
	lineCheck := checkPoints(obstacles, straightLine...)
	if lineCheck == pointCheckPassable {
		return straightLine, nil
	}
	return findPath(obstacles, start, goal)
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
