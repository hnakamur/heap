// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package heap

// MinInt64 is a heap for getting the minimum int64 value.
type MinInt64 []int64

// Init establishes the heap invariants required by the other routines in this package.
// Init is idempotent with respect to the heap invariants
// and may be called whenever the heap invariants may have been invalidated.
// The complexity is O(n) where n = len(*h).
func (h *MinInt64) Init() {
	// heapify
	n := h.length()
	for i := n/2 - 1; i >= 0; i-- {
		h.down(i, n)
	}
}

// Push pushes the element x onto the heap.
// The complexity is O(log n) where n = len(*h).
func (h *MinInt64) Push(x int64) {
	h.push(x)
	h.up(h.length() - 1)
}

// Pop removes and returns the minimum element from the heap.
// The complexity is O(log n) where n = len(*h).
// Pop is equivalent to Remove(h, 0).
func (h *MinInt64) Pop() int64 {
	n := h.length() - 1
	h.swap(0, n)
	h.down(0, n)
	return h.pop()
}

// Remove removes and returns the element at index i from the heap.
// The complexity is O(log n) where n = len(*h).
func (h *MinInt64) Remove(i int) int64 {
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
// The complexity is O(log n) where n = len(*h).
func (h *MinInt64) Fix(i int) {
	if !h.down(i, h.length()) {
		h.up(i)
	}
}

func (h *MinInt64) up(j int) {
	for {
		i := (j - 1) / 2 // parent
		if i == j || !h.less(j, i) {
			break
		}
		h.swap(i, j)
		j = i
	}
}

func (h *MinInt64) down(i0, n int) bool {
	i := i0
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && h.less(j2, j1) {
			j = j2 // = 2*i + 2  // right child
		}
		if !h.less(j, i) {
			break
		}
		h.swap(i, j)
		i = j
	}
	return i > i0
}

func (h MinInt64) length() int        { return len(h) }
func (h MinInt64) less(i, j int) bool { return h[i] < h[j] }
func (h MinInt64) swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MinInt64) push(x int64) {
	*h = append(*h, x)
}

func (h *MinInt64) pop() (x int64) {
	*h, x = (*h)[:h.length()-1], (*h)[h.length()-1]
	return
}
