// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package heap provides heap operations for any type that implements
// heap.Interface. A heap is a tree with the property that each node is the
// minimum-valued node in its subtree.
//
// The minimum element in the tree is the root, at index 0.
//
// A heap is a common way to implement a priority queue. To build a priority
// queue, implement the Heap interface with the (negative) priority as the
// ordering for the Less method, so Push adds items while Pop removes the
// highest-priority item from the queue. The Examples include such an
// implementation; the file example_pq_test.go has the complete source.
//
package strheap

// MaxStr is a heap for getting the maximum string value.
type MaxStr []string

// Len implements sort.Interface.
func (h MaxStr) Len() int { return len(h) }

// Less implements sort.Interface.
func (h MaxStr) Less(i, j int) bool { return h[i] > h[j] }

// Swap implements sort.Interface.
func (h MaxStr) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

// Init establishes the heap invariants required by the other routines in this package.
// Init is idempotent with respect to the heap invariants
// and may be called whenever the heap invariants may have been invalidated.
// The complexity is O(n) where n = h.Len().
func (h *MaxStr) Init() {
	// heapify
	n := h.Len()
	for i := n/2 - 1; i >= 0; i-- {
		down(h, i, n)
	}
}

// Push pushes the element x onto the heap.
// The complexity is O(log n) where n = h.Len().
func (h *MaxStr) Push(x string) {
	h.push(x)
	up(h, h.Len()-1)
}

// Pop removes and returns the minimum element (according to Less) from the heap.
// The complexity is O(log n) where n = h.Len().
// Pop is equivalent to Remove(h, 0).
func (h *MaxStr) Pop() string {
	n := h.Len() - 1
	h.Swap(0, n)
	down(h, 0, n)
	return h.pop()
}

// Remove removes and returns the element at index i from the heap.
// The complexity is O(log n) where n = h.Len().
func (h *MaxStr) Remove(i int) string {
	n := h.Len() - 1
	if n != i {
		h.Swap(i, n)
		if !down(h, i, n) {
			up(h, i)
		}
	}
	return h.pop()
}

// Fix re-establishes the heap ordering after the element at index i has changed its value.
// Changing the value of the element at index i and then calling Fix is equivalent to,
// but less expensive than, calling Remove(h, i) followed by a Push of the new value.
// The complexity is O(log n) where n = h.Len().
func (h *MaxStr) Fix(i int) {
	if !down(h, i, h.Len()) {
		up(h, i)
	}
}

func (h *MaxStr) push(x string) {
	*h = append(*h, x)
}

func (h *MaxStr) pop() (x string) {
	*h, x = (*h)[:h.Len()-1], (*h)[h.Len()-1]
	return
}