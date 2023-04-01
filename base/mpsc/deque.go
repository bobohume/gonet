// Package mpsc provides an efficient implementation of a multi-producer, single-consumer lock-free queue.
//
// The Push function is safe to call from multiple goroutines. The Pop and Empty APIs must only be
// called from a single, consumer goroutine.
//
package mpsc

// This implementation is based on http://www.1024cores.net/home/lock-free-algorithms/queues/non-intrusive-mpsc-node-based-queue

import (
	"sync/atomic"
	"unsafe"
)

type node[T any] struct {
	next *node[T]
	val  T
}

type Queue[T any] struct {
	head, tail *node[T]
	Nil        T
}

func New[T any]() *Queue[T] {
	q := &Queue[T]{}
	stub := &node[T]{}
	q.head = stub
	q.tail = stub
	return q
}

// Push adds x to the back of the queue.
//
// Push can be safely called from multiple goroutines
func (q *Queue[T]) Push(x T) {
	n := new(node[T])
	n.val = x
	// current producer acquires head node
	prev := (*node[T])(atomic.SwapPointer((*unsafe.Pointer)(unsafe.Pointer(&q.head)), unsafe.Pointer(n)))

	// release node to consumer
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&prev.next)), unsafe.Pointer(n))
}

// Pop removes the item from the front of the queue or nil if the queue is empty
//
// Pop must be called from a single, consumer goroutine
func (q *Queue[T]) Pop() T {
	tail := q.tail
	next := (*node[T])(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&tail.next)))) // acquire
	if next != nil {
		q.tail = next
		v := next.val
		return v
	}
	return q.Nil
}

// Empty returns true if the queue is empty
//
// Empty must be called from a single, consumer goroutine
func (q *Queue[T]) Empty() bool {
	tail := q.tail
	next := (*node[T])(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&tail.next))))
	return next == nil
}
