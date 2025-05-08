package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Node struct {
	k int
	v int
	l *Node
	r *Node
}

type OrderedMap struct {
	rt *Node
	sz int
}

func NewOrderedMap() OrderedMap {
	return OrderedMap{}
}

func (m *OrderedMap) ins(n *Node, k, v int) {
	if k < n.k {
		if n.l == nil {
			n.l = &Node{
				k: k,
				v: v,
			}
			m.sz++
			return
		}
		m.ins(n.l, k, v)
	}
	if k > n.k {
		if n.r == nil {
			n.r = &Node{
				k: k,
				v: v,
			}
			m.sz++
			return
		}
		m.ins(n.r, k, v)
	}
}

func (m *OrderedMap) Insert(key, value int) {
	if m.rt == nil {
		m.rt = &Node{
			k: key,
			v: value,
		}
		m.sz++
		return
	}

	m.ins(m.rt, key, value)
}

func getsuc(n *Node) *Node {
	n = n.r
	for n != nil && n.l != nil {
		n = n.l
	}
	return n
}

func (m *OrderedMap) del(n *Node, k int) *Node {
	if n == nil {
		return nil
	}

	if n.k < k {
		n.r = m.del(n.r, k)
	} else if n.k > k {
		n.l = m.del(n.l, k)
	} else {
		m.sz--
		if n.l == nil {
			n = n.r
		} else if n.r == nil {
			n = n.l
		} else {
			suc := getsuc(n)
			n.k, n.v = suc.k, suc.v
			m.del(n.r, suc.k)
		}
	}

	return n
}

func (m *OrderedMap) Erase(key int) {
	if m.rt == nil {
		return
	}
	m.rt = m.del(m.rt, key)
}

func find(n *Node, k int) bool {
	if n == nil {
		return false
	} else if n.k == k {
		return true
	}
	return find(n.l, k) || find(n.r, k)
}

func (m *OrderedMap) Contains(key int) bool {
	return find(m.rt, key)
}

func (m *OrderedMap) Size() int {
	return m.sz
}

func inord(n *Node, action func(int, int)) {
	if n == nil {
		return
	}
	inord(n.l, action)
	action(n.k, n.v)
	inord(n.r, action)
}

func (m *OrderedMap) ForEach(action func(int, int)) {
	if m.rt == nil {
		return
	}
	inord(m.rt, action)
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
