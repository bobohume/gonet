package mpmc

import (
	"runtime"
	"sync/atomic"
)

type(
	Cursor [8]uint64 // prevent false sharing of the sequence cursor by padding the CPU cache line with 64 *bytes* of data.

	node struct {
		sequence uint64
		val  interface{}
	}

	Queue struct{
		mWrite *Cursor	// the ring buffer has been written up to this sequence
		mRead *Cursor	// this reader has processed up to this sequence
		mBufferSize uint64
		mBufferMask uint64
		mRingBuffer []interface{}
	}
)

func (this *node) Store(value uint64) { atomic.StoreUint64(&this.sequence, value) }
func (this *node) Load() uint64       { return atomic.LoadUint64(&this.sequence) }
func (this *node) CmpAndSwap(old, new uint64) bool { return atomic.CompareAndSwapUint64(&this.sequence, old, new)}

func NewCursor() *Cursor {
	var this Cursor
	this[0] = defaultCursorValue
	return &this
}

func New(size uint64) *Queue {
	q := &Queue{}
	q.Init(size)
	return q
}

func (this *Cursor) Store(value uint64) { atomic.StoreUint64(&this[0], value) }
func (this *Cursor) Load() uint64       { return atomic.LoadUint64(&this[0]) }
func (this *Cursor) CmpAndSwap(old, new uint64) bool { return atomic.CompareAndSwapUint64(&this[0], old, new)}

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

func (this *Queue) Init(size uint64) {
	this.mBufferSize = roundUp1(size)
	this.mBufferMask = this.mBufferSize - 1
	this.mWrite = NewCursor()
	this.mRead = NewCursor()
	this.mRingBuffer = make([]interface{}, this.mBufferSize)
	for i := uint64(0); i < this.mBufferSize; i++{
		n := &node{}
		atomic.StoreUint64(&n.sequence, i)
		this.mRingBuffer[i] = n
	}
}

func (this *Queue) Push(data interface{}) {
	var n *node
	pos := this.mWrite.Load()
	for true{
		n = this.mRingBuffer[pos & this.mBufferMask].(*node)
		seq := n.Load()
		dif := int64(seq) - int64(pos)
		if dif == 0 {
			if this.mWrite.CmpAndSwap(pos, pos + 1){
				break
			}
		}else if dif < 0{
			runtime.Gosched() // LockSupport.parkNanos(1L)
		}else{
			pos = this.mWrite.Load()
		}
	}

	n.val = data
	n.Store(pos + 1)
}

func (this *Queue) Pop() interface{}{
	var n *node
	pos := this.mRead.Load()
	for true{
		n = this.mRingBuffer[pos & this.mBufferMask].(*node)
		seq := n.Load()
		dif := int64(seq) - (int64(pos + 1))
		if dif == 0{
			if this.mRead.CmpAndSwap(pos, pos + 1){
				break
			}
		}else if dif < 0{
			return nil
		}else{
			pos = this.mRead.Load()
		}
	}

	dat := n.val
	n.Store(pos + this.mBufferMask + 1)
	return dat
}