package strheap

import (
	"container/heap"
	"encoding/binary"
	"encoding/hex"
	"math/rand"
	"sort"
	"testing"
)

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
	b.Run("strmaxheap", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			values2 := cloneStringList(values)
			b.StartTimer()

			h := MaxStr(make([]string, 0, len(values2)))
			h.Init()
			for _, v := range values2 {
				h.Push(v)
			}
			_ = h.Pop()
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
