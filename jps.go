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
	if isPointOutsideMap(obstacles, start) {
		return nil, errors.New("start is outside the map")
	}
	if isPointOutsideMap(obstacles, goal) {
		return nil, errors.New("goal is outside the map")
	}
	if isPointInsideObstacle(obstacles, start) {
		return nil, errors.New("start is inside an obstacle")
	}
	if isPointInsideObstacle(obstacles, goal) {
		return nil, errors.New("goal is inside an obstacle")
	}
	if straightLine := line(start, goal); isLinePassable(obstacles, straightLine) {
		return straightLine, nil
	}
	return findPath(obstacles, start, goal)
}

// TryFind returns a path from start to goal, or an error in following cases:
// - start or goal is outside the map
// - there is no path from start to goal
// If start or goal is inside an obstacle, it will try to find nearest point on the edge of the obstacle.
func TryFind(obstacles [][]bool, start, goal Point) ([]Point, error) {
	if isPointOutsideMap(obstacles, start) {
		return nil, errors.New("start is outside the map")
	}
	if isPointOutsideMap(obstacles, goal) {
		return nil, errors.New("goal is outside the map")
	}
	if isPointInsideObstacle(obstacles, start) {
		start = nearestToGoal(obstacles, goal, start)
	}
	if isPointInsideObstacle(obstacles, goal) {
		goal = nearestToGoal(obstacles, start, goal)
	}
	if straightLine := line(start, goal); isLinePassable(obstacles, straightLine) {
		return straightLine, nil
	}
	return findPath(obstacles, start, goal)
}

// MustFind returns a path from start to goa.
// If start or goal is outside the map, it will try to find nearest point on the edge of the map.
// If start or goal is inside an obstacle, it will try to find nearest point on the edge of the obstacle.
// If there is no path from start to goal (or the new ones because of previous points) it will return path
// with just start point.
func MustFind(obstacles [][]bool, start, goal Point) []Point {
	if isPointOutsideMap(obstacles, start) {
		start = nearestToGoal(obstacles, goal, start)
	}
	if isPointOutsideMap(obstacles, goal) {
		goal = nearestToGoal(obstacles, start, goal)
	}
	if isPointInsideObstacle(obstacles, start) {
		start = nearestToGoal(obstacles, goal, start)
	}
	if isPointInsideObstacle(obstacles, goal) {
		goal = nearestToGoal(obstacles, start, goal)
	}
	if straightLine := line(start, goal); isLinePassable(obstacles, straightLine) {
		return straightLine
	}
	path, err := findPath(obstacles, start, goal)
	if err != nil {
		return []Point{{start.X, start.Y}}
	}
	return path
}

func findPath(obstacles [][]bool, start, goal Point) ([]Point, error) {
	m := make(map[Point]Point)
	q := make(priorityQueue, 0)
	heap.Init(&q)
	heap.Push(&q, &item{
		point:             start,
		predecessor:       Point{},
		distanceFromStart: 0,
		distanceToGoal:    distance(start, goal),
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
					distanceToGoal:    distance(candidate.p, goal),
				})
			}
		}
	}
	return nil, errors.New("no path found")
}
