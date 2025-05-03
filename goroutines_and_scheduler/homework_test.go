package main

import (
	"container/heap"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Queue []*Task

func (q Queue) Len() int {
	return len(q)
}

func (q Queue) Less(i, j int) bool {
	return q[i].Priority > q[j].Priority
}

func (q Queue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
}

func (q *Queue) Push(x interface{}) {
	task := x.(*Task)

	*q = append(*q, task)
}

func (q *Queue) Pop() interface{} {
	old := *q
	n := len(old)
	task := old[n-1]
	*q = old[0 : n-1]

	return task
}

func (q *Queue) Update(taskID int, newPriority int) {
	for i, task := range *q {
		if task.Identifier == taskID {
			task.Priority = newPriority
			heap.Fix(q, i)
			return
		}
	}
}

type Task struct {
	Identifier int
	Priority   int
}

type Scheduler struct {
	heap *Queue
}

func NewScheduler() Scheduler {
	queue := &Queue{}
	heap.Init(queue)

	return Scheduler{
		heap: queue,
	}
}

func (s *Scheduler) AddTask(task Task) {
	t := &task
	heap.Push(s.heap, t)
}

func (s *Scheduler) GetTask() Task {
	return *heap.Pop(s.heap).(*Task)
}

func (s *Scheduler) ChangeTaskPriority(taskID int, newPriority int) {
	s.heap.Update(taskID, newPriority)
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
