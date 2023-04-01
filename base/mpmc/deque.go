package mpmc

import (
	"runtime"
	"sync/atomic"
)

type (
	Cursor [8]uint64 // prevent false sharing of the sequence cursor by padding the CPU cache line with 64 *bytes* of data.

	node[T any] struct {
		sequence uint64
		val      T
	}

	Queue[T any] struct {
		write      *Cursor // the ring buffer has been written up to q sequence
		read       *Cursor // q reader has processed up to q sequence
		bufferSize uint64
		bufferMask uint64
		ringBuffer []*node[T]
		Nil        T
	}
)

func (n *node[T]) Store(value uint64) { atomic.StoreUint64(&n.sequence, value) }
func (n *node[T]) Load() uint64       { return atomic.LoadUint64(&n.sequence) }
func (n *node[T]) CmpAndSwap(old, new uint64) bool {
	return atomic.CompareAndSwapUint64(&n.sequence, old, new)
}

func NewCursor() *Cursor {
	var c Cursor
	c[0] = defaultCursorValue
	return &c
}

func New[T any](size uint64) *Queue[T] {
	q := &Queue[T]{}
	q.Init(size)
	return q
}

func (c *Cursor) Store(value uint64) { atomic.StoreUint64(&c[0], value) }
func (c *Cursor) Load() uint64       { return atomic.LoadUint64(&c[0]) }
func (c *Cursor) CmpAndSwap(old, new uint64) bool {
	return atomic.CompareAndSwapUint64(&c[0], old, new)
}

const defaultCursorValue = 0

func roundUp1(v uint64) uint64 {
	v--
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	v |= v >> 32
	v++
	return v
}

func (q *Queue[T]) Init(size uint64) {
	q.bufferSize = roundUp1(size)
	q.bufferMask = q.bufferSize - 1
	q.write = NewCursor()
	q.read = NewCursor()
	q.ringBuffer = make([]*node[T], q.bufferSize)
	for i := uint64(0); i < q.bufferSize; i++ {
		n := &node[T]{}
		atomic.StoreUint64(&n.sequence, i)
		q.ringBuffer[i] = n
	}
}

func (q *Queue[T]) Push(data T) {
	var n *node[T]
	pos := q.write.Load()
	for true {
		n = q.ringBuffer[pos&q.bufferMask]
		seq := n.Load()
		dif := int64(seq) - int64(pos)
		if dif == 0 {
			if q.write.CmpAndSwap(pos, pos+1) {
				break
			}
		} else if dif < 0 {
			runtime.Gosched() // LockSupport.parkNanos(1L)
		} else {
			pos = q.write.Load()
		}
	}

	n.val = data
	n.Store(pos + 1)
}

func (q *Queue[T]) Pop() T {
	var n *node[T]
	pos := q.read.Load()
	for true {
		n = q.ringBuffer[pos&q.bufferMask]
		seq := n.Load()
		dif := int64(seq) - (int64(pos + 1))
		if dif == 0 {
			if q.read.CmpAndSwap(pos, pos+1) {
				break
			}
		} else if dif < 0 {
			return q.Nil
		} else {
			pos = q.read.Load()
		}
	}

	dat := n.val
	n.Store(pos + q.bufferMask + 1)
	return dat
}
