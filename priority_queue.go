package jps

import "sort"

type item struct {
	point             Point
	predecessor       Point
	distanceFromStart float64
	distanceToGoal    float64
}

type priorityQueue []*item

func (pq priorityQueue) Len() int { return len(pq) }

func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].distanceFromStart+pq[i].distanceToGoal > pq[j].distanceFromStart+pq[j].distanceToGoal
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *priorityQueue) Push(x any) {
	item := x.(*item)
	*pq = append(*pq, item)
}

func (pq *priorityQueue) Pop() any {
	if !sort.IsSorted(*pq) {
		sort.Sort(*pq)
	}
	old := *pq
	n := len(old)
	result := old[n-1]
	old[n-1] = nil // avoid memory leak
	*pq = old[0 : n-1]
	return result
}
