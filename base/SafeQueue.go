package base

import (
	"runtime"
	"sync/atomic"
)

type Cursor [8]int64 // prevent false sharing of the sequence cursor by padding the CPU cache line with 64 *bytes* of data.

func NewCursor() *Cursor {
	var this Cursor
	this[0] = defaultCursorValue
	return &this
}

func (this *Cursor) Store(value int64) { atomic.StoreInt64(&this[0], value) }
func (this *Cursor) Load() int64       { return atomic.LoadInt64(&this[0]) }
func (this *Cursor) CmpAndSwap(old, new int64) bool { return atomic.CompareAndSwapInt64(&this[0], old, new)}

const defaultCursorValue = 0

type (
	SafeQueue struct{
		mWrite *Cursor	// the ring buffer has been written up to this sequence
		mRead *Cursor	// this reader has processed up to this sequence
		mBufferSize int64
		mBufferMask int64
		mRingBuffer []interface{}
	}

	ISafeQueue interface {
		Init(size int64)
		Push(data interface{})
		Pop() interface{}
	}
)

func roundUp1(v int64) int64 {
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

func (this *SafeQueue) Init(size int64) {
	this.mBufferSize = roundUp1(size)
	this.mBufferMask = this.mBufferSize - 1
	this.mRingBuffer = make([]interface{}, this.mBufferSize)
	this.mWrite = NewCursor()
	this.mRead = NewCursor()
}

func (this *SafeQueue) Push(data interface{}) {
	lower := this.mWrite.Load()
	upper := lower+1
	for !this.mWrite.CmpAndSwap(lower, upper){
		runtime.Gosched() // LockSupport.parkNanos(1L)
		lower = this.mWrite.Load()
		upper = lower+1
	}
	this.mRingBuffer[lower & this.mBufferMask] = data
}

func (this *SafeQueue) Pop() interface{}{
	lower := this.mRead.Load()
	upper := this.mWrite.Load()
	bSucess := false
	for lower+1 < upper{
		bSucess = this.mRead.CmpAndSwap(lower, lower+1)
		if !bSucess{
			runtime.Gosched() // LockSupport.parkNanos(1L)
			lower = this.mRead.Load()
			upper = this.mWrite.Load()
		}else{
			return this.mRingBuffer[lower & this.mBufferMask]
		}
	}
	return nil
}

func NewSafeQueue(size int64) *SafeQueue{
	queue := &SafeQueue{}
	queue.Init(size)
	return queue
}
