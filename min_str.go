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
// ordering for the less method, so Push adds items while Pop removes the
// highest-priority item from the queue. The Examples include such an
// implementation; the file example_pq_test.go has the complete source.
//
package heap

// StrMaxHeap is a heap for getting the minimum string value.
type MinStr []string

// Init establishes the heap invariants required by the other routines in this package.
// Init is idempotent with respect to the heap invariants
// and may be called whenever the heap invariants may have been invalidated.
// The complexity is O(n) where n = h.length().
func (h *MinStr) Init() {
	// heapify
	n := h.length()
	for i := n/2 - 1; i >= 0; i-- {
		h.down(i, n)
	}
}

// Push pushes the element x onto the heap.
// The complexity is O(log n) where n = h.length().
func (h *MinStr) Push(x string) {
	h.push(x)
	h.up(h.length() - 1)
}

// Pop removes and returns the minimum element (according to less) from the heap.
// The complexity is O(log n) where n = h.length().
// Pop is equivalent to Remove(h, 0).
func (h *MinStr) Pop() string {
	n := h.length() - 1
	h.swap(0, n)
	h.down(0, n)
	return h.pop()
}

// Remove removes and returns the element at index i from the heap.
// The complexity is O(log n) where n = h.length().
func (h *MinStr) Remove(i int) string {
	n := h.length() - 1
	if n != i {
		h.swap(i, n)
		if !h.down(i, n) {
			h.up(i)
		}
	}
	return h.pop()
}

// Fix re-establishes the heap ordering after the element at index i has changed its value.
// Changing the value of the element at index i and then calling Fix is equivalent to,
// but less expensive than, calling Remove(h, i) followed by a Push of the new value.
// The complexity is O(log n) where n = h.length().
func (h *MinStr) Fix(i int) {
	if !h.down(i, h.length()) {
		h.up(i)
	}
}

func (h *MinStr) up(j int) {
	for {
		i := (j - 1) / 2 // parent
		if i == j || (*h)[j] >= (*h)[i] {
			break
		}
		h.swap(i, j)
		j = i
	}
}

func (h *MinStr) down(i0, n int) bool {
	i := i0
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && (*h)[j2] < (*h)[j1] {
			j = j2 // = 2*i + 2  // right child
		}
		if (*h)[j] >= (*h)[i] {
			break
		}
		h.swap(i, j)
		i = j
	}
	return i > i0
}

func (h MinStr) length() int   { return len(h) }
func (h MinStr) swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *MinStr) push(x string) {
	*h = append(*h, x)
}

func (h *MinStr) pop() (x string) {
	*h, x = (*h)[:h.length()-1], (*h)[h.length()-1]
	return
}
