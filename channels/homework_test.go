package main

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

const chanSizeMultiplier = 2

type WorkerPool struct {
	tasksCh      chan func()
	wg           *sync.WaitGroup
	shutdownOnce *sync.Once
}

func NewWorkerPool(workersNumber int) *WorkerPool {
	wp := &WorkerPool{
		tasksCh:      make(chan func(), chanSizeMultiplier*workersNumber),
		wg:           &sync.WaitGroup{},
		shutdownOnce: &sync.Once{},
	}

	wp.wg.Add(workersNumber)
	for i := 0; i < workersNumber; i++ {
		go func() {
			defer wp.wg.Done()
			for task := range wp.tasksCh {
				task()
			}
		}()
	}

	return wp
}

// Return an error if the pool is full
func (wp *WorkerPool) AddTask(task func()) error {
	select {
	case wp.tasksCh <- task:
		return nil
	default:
		return errors.New("pool is full")
	}
}

// Shutdown all workers and wait for all
// tasks in the pool to complete
func (wp *WorkerPool) Shutdown() {
	wp.shutdownOnce.Do(func() {
		close(wp.tasksCh)
		wp.wg.Wait()
	})
}

func TestWorkerPool(t *testing.T) {
	var counter atomic.Int32
	task := func() {
		time.Sleep(time.Millisecond * 500)
		counter.Add(1)
	}

	pool := NewWorkerPool(2)
	_ = pool.AddTask(task)
	_ = pool.AddTask(task)
	_ = pool.AddTask(task)

	time.Sleep(time.Millisecond * 600)
	assert.Equal(t, int32(2), counter.Load())

	time.Sleep(time.Millisecond * 600)
	assert.Equal(t, int32(3), counter.Load())

	_ = pool.AddTask(task)
	_ = pool.AddTask(task)
	_ = pool.AddTask(task)
	pool.Shutdown() // wait tasks

	assert.Equal(t, int32(6), counter.Load())
}
