package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Node struct {
	Key   int
	Value int

	Left  *Node
	Right *Node
}

func inorderTraverse(node *Node, action func(int, int)) {
	if node == nil {
		return
	}

	inorderTraverse(node.Left, action)
	action(node.Key, node.Value)
	inorderTraverse(node.Right, action)
}

type OrderedMap struct {
	bst *Node

	size int
}

func NewOrderedMap() OrderedMap {
	return OrderedMap{
		size: 0,
	}
}

func (m *OrderedMap) Insert(key, value int) {
	if m.bst == nil {
		m.bst = &Node{
			Key:   key,
			Value: value,
		}
	} else {
		curr := m.bst
		for curr != nil {
			if key > curr.Key {
				if curr.Right == nil {
					curr.Right = &Node{Key: key, Value: value}
					break
				}
				curr = curr.Right
			} else if key < curr.Key {
				if curr.Left == nil {
					curr.Left = &Node{Key: key, Value: value}
					break
				}
				curr = curr.Left
			} else {
				curr.Value = value
				return
			}
		}
	}

	m.size += 1
}

func (m *OrderedMap) Erase(key int) {
	var (
		curr, parent *Node = m.bst, nil
	)
	for curr != nil && curr.Key != key {
		parent = curr
		if key > curr.Key {
			curr = curr.Right
		} else if key < curr.Key {
			curr = curr.Left
		}
	}

	if curr == nil {
		return
	}

	m.size -= 1

	if curr.Right == nil && curr.Left == nil {
		if parent == nil {
			m.bst = nil
		} else if parent.Left == curr {
			parent.Left = nil
		} else {
			parent.Right = nil
		}
		return
	}

	if curr.Right == nil || curr.Left != nil {
		child := &Node{}
		if curr.Left != nil {
			child = curr.Left
		} else {
			child = curr.Right
		}

		if parent == nil {
			m.bst = child
		} else if parent.Left == curr {
			parent.Left = child
		} else {
			parent.Right = child
		}

		return
	}

	var (
		min, minParent *Node = curr.Right, curr
	)
	for min.Left != nil {
		minParent = min
		min = min.Left
	}

	curr.Key = min.Key
	curr.Value = min.Value
	if minParent.Left == min {
		minParent.Left = min.Right
	} else {
		minParent.Right = min.Right
	}
}

func (m *OrderedMap) Contains(key int) bool {
	next := m.bst
	for next != nil {
		if key > next.Key {
			next = next.Right
		} else if key < next.Key {
			next = next.Left
		} else {
			return true
		}
	}

	return false
}

func (m *OrderedMap) Size() int {
	return m.size
}

func (m *OrderedMap) ForEach(action func(int, int)) {
	inorderTraverse(m.bst, action)
}

func TestCircularQueue(t *testing.T) {
	data := NewOrderedMap()
	assert.Zero(t, data.Size())

	data.Insert(10, 10)
	data.Insert(5, 5)
	data.Insert(15, 15)
	data.Insert(2, 2)
	data.Insert(4, 4)
	data.Insert(12, 12)
	data.Insert(14, 14)

	assert.Equal(t, 7, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(3))
	assert.False(t, data.Contains(13))

	var keys []int
	expectedKeys := []int{2, 4, 5, 10, 12, 14, 15}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))

	data.Erase(15)
	data.Erase(14)
	data.Erase(2)

	assert.Equal(t, 4, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(2))
	assert.False(t, data.Contains(14))

	keys = nil
	expectedKeys = []int{4, 5, 10, 12}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))
}
