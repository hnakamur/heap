// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package strheap

import (
	"math/rand"
	"testing"
)

func (h *MaxStrHeap) verify(t *testing.T, i int) {
	t.Helper()
	n := h.Len()
	j1 := 2*i + 1
	j2 := 2*i + 2
	if j1 < n {
		if h.Less(j1, i) {
			t.Errorf("heap invariant invalidated [%d] = %s > [%d] = %s", i, (*h)[i], j1, (*h)[j1])
			return
		}
		h.verify(t, j1)
	}
	if j2 < n {
		if h.Less(j2, i) {
			t.Errorf("heap invariant invalidated [%d] = %s > [%d] = %s", i, (*h)[i], j1, (*h)[j2])
			return
		}
		h.verify(t, j2)
	}
}

func TestMaxStrHeapInit0(t *testing.T) {
	h := new(MaxStrHeap)
	for i := 20; i > 0; i-- {
		h.Push("0") // all elements are the same
	}
	h.Init()
	h.verify(t, 0)

	for i := 1; h.Len() > 0; i++ {
		x := h.Pop()
		h.verify(t, 0)
		if x != "0" {
			t.Errorf("%d.th pop got %s; want %s", i, x, "0")
		}
	}
}

func TestMaxStrHeapInit1(t *testing.T) {
	h := new(MaxStrHeap)
	for i := 20; i > 0; i-- {
		h.Push(toHex(uint64(i))) // all elements are different
	}
	h.Init()
	h.verify(t, 0)

	for i := 20; h.Len() > 0; i-- {
		x := h.Pop()
		h.verify(t, 0)
		if x != toHex(uint64(i)) {
			t.Errorf("%d.th pop got %s; want %s", i, x, toHex(uint64(i)))
		}
	}
}

func TestMaxStrHeap(t *testing.T) {
	h := new(MaxStrHeap)
	h.verify(t, 0)

	for i := 30; i > 20; i-- {
		h.push(toHex(uint64(i)))
	}
	h.Init()
	h.verify(t, 0)

	for i := 20; i > 10; i-- {
		h.Push(toHex(uint64(i)))
		h.verify(t, 0)
	}

	for i := 30; h.Len() > 0; i-- {
		x := h.Pop()
		if i < 10 {
			h.Push(toHex(uint64(i)))
		}
		h.verify(t, 0)
		if x != toHex(uint64(i)) {
			t.Errorf("%d.th pop got %s; want %s", i, x, toHex(uint64(i)))
		}
	}
}

func TestMaxStrHeapRemove0(t *testing.T) {
	h := new(MaxStrHeap)
	for i := 9; i >= 0; i-- {
		h.push(toHex(uint64(i)))
	}
	h.verify(t, 0)

	for h.Len() > 0 {
		i := h.Len() - 1
		x := h.Remove(i)
		if x != toHex(uint64(9-i)) {
			t.Errorf("Remove(%d) got %s; want %s", i, x, toHex(uint64(9-i)))
		}
		h.verify(t, 0)
	}
}

func TestMaxStrHeapRemove1(t *testing.T) {
	h := new(MaxStrHeap)
	for i := 9; i >= 0; i-- {
		h.push(toHex(uint64(i)))
	}
	h.verify(t, 0)

	for i := 0; h.Len() > 0; i++ {
		x := h.Remove(0)
		if x != toHex(uint64(9-i)) {
			t.Errorf("Remove(0) got %s; want %s", x, toHex(uint64(9-i)))
		}
		h.verify(t, 0)
	}
}

func TestMaxStrHeapRemove2(t *testing.T) {
	N := 10

	h := new(MaxStrHeap)
	for i := N - 1; i >= 0; i-- {
		h.push(toHex(uint64(i)))
	}
	h.verify(t, 0)

	m := make(map[string]bool)
	for h.Len() > 0 {
		m[h.Remove((h.Len()-1)/2)] = true
		h.verify(t, 0)
	}

	if len(m) != N {
		t.Errorf("len(m) = %d; want %d", len(m), N)
	}
	for i := 0; i < len(m); i++ {
		k := toHex(uint64(i))
		if !m[k] {
			t.Errorf("m[%s] doesn't exist", k)
		}
	}
}

func BenchmarkMaxStrHeapDup(b *testing.B) {
	const n = 10000
	h := make(MaxStrHeap, 0, n)
	for i := 0; i < b.N; i++ {
		for j := 0; j < n; j++ {
			h.Push("0") // all elements are the same
		}
		for h.Len() > 0 {
			h.Pop()
		}
	}
}

func TestMaxStrHeapFix(t *testing.T) {
	h := new(MaxStrHeap)
	h.verify(t, 0)

	for i := 200; i > 0; i -= 10 {
		h.Push(toHex(uint64(i)))
	}
	h.verify(t, 0)

	if (*h)[0] != toHex(200) {
		t.Fatalf("Expected head to be 200, was %s", (*h)[0])
	}
	(*h)[0] = toHex(210)
	h.Fix(0)
	h.verify(t, 0)

	for i := 100; i > 0; i-- {
		elem := rand.Intn(h.Len())
		if i&1 == 0 {
			(*h)[elem] = toHex(fromHex(t, (*h)[elem]) * 2)
		} else {
			(*h)[elem] = toHex(fromHex(t, (*h)[elem]) / 2)
		}
		h.Fix(elem)
		h.verify(t, 0)
	}
}
