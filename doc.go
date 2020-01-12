// Package heap provides heap operations for string, int64, and uint64.
// A heap is a tree with the property that each node is the
// minimum-valued node in its subtree.
//
// The minimum element in the tree is the root, at index 0.
//
// This package is copied from the standard library container/heap
// and modified for concrete types such as string.
//
// Package heap also provides structs MaxStr, MaxInt64, and MaxUint64
// for maximum versions of heap.
//
package heap
