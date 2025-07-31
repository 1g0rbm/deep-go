package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type RWMutex struct {
	mutex sync.Mutex

	readCond *sync.Cond
	readers  int

	writeCond      *sync.Cond
	write          bool
	waitingWriters int
}

func NewRWMutex() *RWMutex {
	m := &RWMutex{}
	m.readCond = sync.NewCond(&m.mutex)
	m.writeCond = sync.NewCond(&m.mutex)

	return m
}

func (m *RWMutex) Lock() {
	m.mutex.Lock()
	defer func() {
		m.mutex.Unlock()
		m.waitingWriters--
		m.write = true
	}()

	m.waitingWriters++
	for m.readers > 0 || m.write {
		m.writeCond.Wait()
	}
}

func (m *RWMutex) Unlock() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.write = false
	if m.waitingWriters > 0 {
		m.writeCond.Signal()
	} else {
		m.readCond.Broadcast()
	}
}

func (m *RWMutex) RLock() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for m.write || m.waitingWriters > 0 {
		m.readCond.Wait()
	}

	m.readers++
}

func (m *RWMutex) RUnlock() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.readers--
	if m.readers == 0 {
		m.writeCond.Signal()
	}
}

func TestRWMutexWithWriter(t *testing.T) {
	mutex := NewRWMutex()
	mutex.Lock() // writer

	var mutualExlusionWithWriter atomic.Bool
	mutualExlusionWithWriter.Store(true)
	var mutualExlusionWithReader atomic.Bool
	mutualExlusionWithReader.Store(true)

	go func() {
		mutex.Lock() // another writer
		mutualExlusionWithWriter.Store(false)
	}()

	go func() {
		mutex.RLock() // another reader
		mutualExlusionWithReader.Store(false)
	}()

	time.Sleep(time.Second)
	assert.True(t, mutualExlusionWithWriter.Load())
	assert.True(t, mutualExlusionWithReader.Load())
}

func TestRWMutexWithReaders(t *testing.T) {
	mutex := NewRWMutex()
	mutex.RLock() // reader

	var mutualExlusionWithWriter atomic.Bool
	mutualExlusionWithWriter.Store(true)

	go func() {
		mutex.Lock() // another writer
		mutualExlusionWithWriter.Store(false)
	}()

	time.Sleep(time.Second)
	assert.True(t, mutualExlusionWithWriter.Load())
}

func TestRWMutexMultipleReaders(t *testing.T) {
	mutex := NewRWMutex()
	mutex.RLock() // reader

	var readersCount atomic.Int32
	readersCount.Add(1)

	go func() {
		mutex.RLock() // another reader
		readersCount.Add(1)
	}()

	go func() {
		mutex.RLock() // another reader
		readersCount.Add(1)
	}()

	time.Sleep(time.Second)
	assert.Equal(t, int32(3), readersCount.Load())
}

func TestRWMutexWithWriterPriority(t *testing.T) {
	mutex := NewRWMutex()
	mutex.RLock() // reader

	var mutualExlusionWithWriter atomic.Bool
	mutualExlusionWithWriter.Store(true)
	var readersCount atomic.Int32
	readersCount.Add(1)

	go func() {
		mutex.Lock() // another writer is waiting for reader
		fmt.Println("FIRST UNLOCK")
		mutualExlusionWithWriter.Store(false)
	}()

	time.Sleep(time.Second)

	go func() {
		mutex.RLock() // another reader is waiting for a higher priority writer
		fmt.Println("SECOND UNLOCK")
		readersCount.Add(1)
	}()

	go func() {
		mutex.RLock() // another reader is waiting for a higher priority writer
		fmt.Println("THIRD UNLOCK")
		readersCount.Add(1)
	}()

	time.Sleep(time.Second)

	assert.True(t, mutualExlusionWithWriter.Load())
	assert.Equal(t, int32(1), readersCount.Load())
}
