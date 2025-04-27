package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Heap struct {
	queue         []Task
	idxToPosition map[int]int
}

func (h *Heap) Push(task Task) {
	h.queue = append(h.queue, task)
	h.idxToPosition[task.Identifier] = len(h.queue) - 1
	h.up(len(h.queue) - 1)
}

func (h *Heap) Pop() Task {
	if len(h.queue) == 0 {
		return Task{}
	}

	task := h.queue[0]
	h.queue[0] = h.queue[len(h.queue)-1]
	h.queue = h.queue[:len(h.queue)-1]
	delete(h.idxToPosition, task.Identifier)

	h.down(0)

	return task
}

func (h *Heap) Update(taskID int, newPriority int) {
	i, ok := h.idxToPosition[taskID]
	if !ok {
		return
	}

	oldPriority := h.queue[i].Priority
	h.queue[i].Priority = newPriority
	if newPriority > oldPriority {
		h.up(i)
	} else if newPriority < oldPriority {
		h.down(i)
	}
}

func (h *Heap) down(currIdx int) {
	for currIdx < len(h.queue) {
		leftIdx := 2*currIdx + 1
		rightIdx := 2*currIdx + 2

		if leftIdx > len(h.queue)-1 {
			break
		}

		var toSwapIdx int
		if rightIdx >= len(h.queue) || h.queue[leftIdx].Priority > h.queue[rightIdx].Priority {
			toSwapIdx = leftIdx
		} else {
			toSwapIdx = rightIdx
		}

		if h.queue[currIdx].Priority > h.queue[toSwapIdx].Priority {
			break
		}

		h.idxToPosition[h.queue[currIdx].Identifier] = toSwapIdx
		h.idxToPosition[h.queue[toSwapIdx].Identifier] = currIdx
		h.queue[toSwapIdx], h.queue[currIdx] = h.queue[currIdx], h.queue[toSwapIdx]

		currIdx = toSwapIdx
	}
}

func (h *Heap) up(currIdx int) {
	for currIdx > 0 {
		parentIdx := (currIdx - 1) / 2
		current := h.queue[currIdx]
		parent := h.queue[parentIdx]
		if parent.Priority > current.Priority {
			break
		}

		h.idxToPosition[current.Identifier] = parentIdx
		h.idxToPosition[parent.Identifier] = currIdx
		h.queue[parentIdx], h.queue[currIdx] = h.queue[currIdx], h.queue[parentIdx]

		currIdx = parentIdx
	}
}

type Task struct {
	Identifier int
	Priority   int
}

type Scheduler struct {
	heap Heap
}

func NewScheduler() Scheduler {
	return Scheduler{
		heap: Heap{
			queue:         []Task{},
			idxToPosition: make(map[int]int),
		},
	}
}

func (s *Scheduler) AddTask(task Task) {
	s.heap.Push(task)
}

func (s *Scheduler) ChangeTaskPriority(taskID int, newPriority int) {
	s.heap.Update(taskID, newPriority)
}

func (s *Scheduler) GetTask() Task {
	return s.heap.Pop()
}

func TestTrace(t *testing.T) {
	task1 := Task{Identifier: 1, Priority: 10}
	task2 := Task{Identifier: 2, Priority: 20}
	task3 := Task{Identifier: 3, Priority: 30}
	task4 := Task{Identifier: 4, Priority: 40}
	task5 := Task{Identifier: 5, Priority: 50}

	scheduler := NewScheduler()
	scheduler.AddTask(task1)
	scheduler.AddTask(task2)
	scheduler.AddTask(task3)
	scheduler.AddTask(task4)
	scheduler.AddTask(task5)

	task := scheduler.GetTask()
	assert.Equal(t, task5, task)

	task = scheduler.GetTask()
	assert.Equal(t, task4, task)

	scheduler.ChangeTaskPriority(1, 100)

	task = scheduler.GetTask()
	assert.Equal(t, Task{Identifier: 1, Priority: 100}, task)

	task = scheduler.GetTask()
	assert.Equal(t, task3, task)
}
