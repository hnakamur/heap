package strheap

import (
	"container/heap"
	"encoding/binary"
	"encoding/hex"
	"math/rand"
	"sort"
	"testing"
)

func TestHeapAndSort(t *testing.T) {
	const n = 100_000
	values := make([]string, n)
	for i := 0; i < n; i++ {
		var b [8]byte
		binary.BigEndian.PutUint64(b[:], uint64(i))
		values[i] = hex.EncodeToString(b[:])
	}
	want := values[len(values)-1]
	rand.Shuffle(n, func(i, j int) {
		values[i], values[j] = values[j], values[i]
	})

	t.Run("sort", func(t *testing.T) {
		values2 := cloneStringList(values)
		sort.Strings(values2)
		got := values2[len(values2)-1]
		if got != want {
			t.Errorf("max value unmatch, got=%s, want=%s", got, want)
		}
	})
	t.Run("heap", func(t *testing.T) {
		h := strMaxHeap(make([]string, 0, len(values)))
		heap.Init(&h)
		for _, v := range values {
			heap.Push(&h, v)
		}
		got := heap.Pop(&h).(string)
		if got != want {
			t.Errorf("max value unmatch, got=%s, want=%s", got, want)
		}
	})
}

func BenchmarkHeapAndSort(b *testing.B) {
	const n = 100_000
	values := make([]string, n)
	for i := 0; i < n; i++ {
		var b [8]byte
		binary.BigEndian.PutUint64(b[:], uint64(i))
		values[i] = hex.EncodeToString(b[:])
	}
	rand.Shuffle(n, func(i, j int) {
		values[i], values[j] = values[j], values[i]
	})

	b.Run("sort", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			values2 := cloneStringList(values)
			b.StartTimer()

			sort.Strings(values2)
		}
	})
	b.Run("heap", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			values2 := cloneStringList(values)
			b.StartTimer()

			h := strMaxHeap(make([]string, 0, len(values2)))
			heap.Init(&h)
			for _, v := range values2 {
				heap.Push(&h, v)
			}
			_ = heap.Pop(&h).(string)
		}
	})
}

func cloneStringList(values []string) []string {
	ret := make([]string, len(values))
	copy(ret, values)
	return ret
}

type strMaxHeap []string

func (h strMaxHeap) Len() int           { return len(h) }
func (h strMaxHeap) Less(i, j int) bool { return h[i] > h[j] }
func (h strMaxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *strMaxHeap) Push(x interface{}) {
	*h = append(*h, x.(string))
}

func (h *strMaxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
