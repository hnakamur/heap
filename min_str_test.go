// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package strheap

import (
	"encoding/binary"
	"encoding/hex"
	"math/rand"
	"testing"
)

func (h *MinStr) verify(t *testing.T, i int) {
	t.Helper()
	n := h.length()
	j1 := 2*i + 1
	j2 := 2*i + 2
	if j1 < n {
		if h.less(j1, i) {
			t.Errorf("heap invariant invalidated [%d] = %s > [%d] = %s", i, (*h)[i], j1, (*h)[j1])
			return
		}
		h.verify(t, j1)
	}
	if j2 < n {
		if h.less(j2, i) {
			t.Errorf("heap invariant invalidated [%d] = %s > [%d] = %s", i, (*h)[i], j1, (*h)[j2])
			return
		}
		h.verify(t, j2)
	}
}

func TestMinStrInit0(t *testing.T) {
	h := new(MinStr)
	for i := 20; i > 0; i-- {
		h.Push("0") // all elements are the same
	}
	h.Init()
	h.verify(t, 0)

	for i := 1; h.length() > 0; i++ {
		x := h.Pop()
		h.verify(t, 0)
		if x != "0" {
			t.Errorf("%d.th pop got %s; want %s", i, x, "0")
		}
	}
}

func TestMinStrInit1(t *testing.T) {
	h := new(MinStr)
	for i := 20; i > 0; i-- {
		h.Push(toHex(uint64(i))) // all elements are different
	}
	h.Init()
	h.verify(t, 0)

	for i := 1; h.length() > 0; i++ {
		x := h.Pop()
		h.verify(t, 0)
		if x != toHex(uint64(i)) {
			t.Errorf("%d.th pop got %s; want %s", i, x, toHex(uint64(i)))
		}
	}
}

func TestMinStr(t *testing.T) {
	h := new(MinStr)
	h.verify(t, 0)

	for i := 20; i > 10; i-- {
		h.push(toHex(uint64(i)))
	}
	h.Init()
	h.verify(t, 0)

	for i := 10; i > 0; i-- {
		h.Push(toHex(uint64(i)))
		h.verify(t, 0)
	}

	for i := 1; h.length() > 0; i++ {
		x := h.Pop()
		if i < 20 {
			h.Push(toHex(uint64(20 + i)))
		}
		h.verify(t, 0)
		if x != toHex(uint64(i)) {
			t.Errorf("%d.th pop got %s; want %s", i, x, toHex(uint64(i)))
		}
	}
}

func TestMinStrRemove0(t *testing.T) {
	h := new(MinStr)
	for i := 0; i < 10; i++ {
		h.push(toHex(uint64(i)))
	}
	h.verify(t, 0)

	for h.length() > 0 {
		i := h.length() - 1
		x := h.Remove(i)
		if x != toHex(uint64(i)) {
			t.Errorf("Remove(%d) got %s; want %s", i, x, toHex(uint64(i)))
		}
		h.verify(t, 0)
	}
}

func TestMinStrRemove1(t *testing.T) {
	h := new(MinStr)
	for i := 0; i < 10; i++ {
		h.push(toHex(uint64(i)))
	}
	h.verify(t, 0)

	for i := 0; h.length() > 0; i++ {
		x := h.Remove(0)
		if x != toHex(uint64(i)) {
			t.Errorf("Remove(0) got %s; want %s", x, toHex(uint64(i)))
		}
		h.verify(t, 0)
	}
}

func TestMinStrRemove2(t *testing.T) {
	N := 10

	h := new(MinStr)
	for i := 0; i < N; i++ {
		h.push(toHex(uint64(i)))
	}
	h.verify(t, 0)

	m := make(map[string]bool)
	for h.length() > 0 {
		m[h.Remove((h.length()-1)/2)] = true
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

func BenchmarkMinStrDup(b *testing.B) {
	const n = 10000
	h := make(MinStr, 0, n)
	for i := 0; i < b.N; i++ {
		for j := 0; j < n; j++ {
			h.Push("0") // all elements are the same
		}
		for h.length() > 0 {
			h.Pop()
		}
	}
}

func TestMinStrFix(t *testing.T) {
	h := new(MinStr)
	h.verify(t, 0)

	for i := 200; i > 0; i -= 10 {
		h.Push(toHex(uint64(i)))
	}
	h.verify(t, 0)

	if (*h)[0] != toHex(10) {
		t.Fatalf("Expected head to be 10, was %s", (*h)[0])
	}
	(*h)[0] = toHex(210)
	h.Fix(0)
	h.verify(t, 0)

	for i := 100; i > 0; i-- {
		elem := rand.Intn(h.length())
		if i&1 == 0 {
			(*h)[elem] = toHex(fromHex(t, (*h)[elem]) * 2)
		} else {
			(*h)[elem] = toHex(fromHex(t, (*h)[elem]) / 2)
		}
		h.Fix(elem)
		h.verify(t, 0)
	}
}

func toHex(i uint64) string {
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], i)
	return hex.EncodeToString(b[:])
}

func fromHex(t *testing.T, h string) uint64 {
	t.Helper()
	b, err := hex.DecodeString(h)
	if err != nil {
		t.Fatal(err)
	}
	return binary.BigEndian.Uint64(b)
}
