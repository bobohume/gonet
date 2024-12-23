package vector

import "fmt"

// minCapacity is the smallest capacity that Vector may have. Must be power of 2
// for bitwise modulus: x % n == x & (n - 1).
const minCapacity = 16

// Vector represents a single instance of the Vector data structure. A Vector
// instance contains items of the type specified by the type argument.
type Vector[T any] struct {
	buf    []T
	head   int
	tail   int
	count  int
	minCap int
}

// New creates a new Vector, optionally setting the current and minimum capacity
// when non-zero values are given for these. The Vector instance returns
// operates on items of the type specified by the type argument. For example,
// to create a Vector that contains strings,
//
//	stringVector := Vector.New[string]()
//
// To create a Vector with capacity to store 2048 ints without resizing, and
// that will not resize below space for 32 items when removing items:
//
//	d := Vector.New[int](2048, 32)
//
// To create a Vector that has not yet allocated memory, but after it does will
// never resize to have space for less than 64 items:
//
//	d := Vector.New[int](0, 64)
//
// Any size values supplied here are rounded up to the nearest power of 2.
func New[T any](size ...int) *Vector[T] {
	var capacity, minimum int
	if len(size) >= 1 {
		capacity = size[0]
		if len(size) >= 2 {
			minimum = size[1]
		}
	}

	minCap := minCapacity
	for minCap < minimum {
		minCap <<= 1
	}

	var buf []T
	if capacity != 0 {
		bufSize := minCap
		for bufSize < capacity {
			bufSize <<= 1
		}
		buf = make([]T, bufSize)
	}

	return &Vector[T]{
		buf:    buf,
		minCap: minCap,
	}
}

// Cap returns the current capacity of the Vector. If q is nil, v.Cap() is zero.
func (v *Vector[T]) Cap() int {
	if v == nil {
		return 0
	}
	return len(v.buf)
}

func (v *Vector[T]) Empty() bool {
	return v.Len() == 0
}

// Len returns the number of elements currently stored in the queue. If q is
// nil, v.Len() is zero.
func (v *Vector[T]) Len() int {
	if v == nil {
		return 0
	}
	return v.count
}

// PushBack appends an element to the back of the queue. Implements FIFO when
// elements are removed with PopFront, and LIFO when elements are removed with
// PopBack.
func (v *Vector[T]) PushBack(elem T) {
	v.growIfFull()

	v.buf[v.tail] = elem
	// Calculate new tail position.
	v.tail = v.next(v.tail)
	v.count++
}

// PushFront prepends an element to the front of the queue.
func (v *Vector[T]) PushFront(elem T) {
	v.growIfFull()

	// Calculate new head position.
	v.head = v.prev(v.head)
	v.buf[v.head] = elem
	v.count++
}

// PopFront removes and returns the element from the front of the queue.
// Implements FIFO when used with PushBack. If the queue is empty, the call
// panics.
func (v *Vector[T]) PopFront() T {
	if v.count <= 0 {
		panic("Vector: PopFront() called on empty queue")
	}
	ret := v.buf[v.head]
	var zero T
	v.buf[v.head] = zero
	// Calculate new head position.
	v.head = v.next(v.head)
	v.count--

	v.shrinkIfExcess()
	return ret
}

// PopBack removes and returns the element from the back of the queue.
// Implements LIFO when used with PushBack. If the queue is empty, the call
// panics.
func (v *Vector[T]) PopBack() T {
	if v.count <= 0 {
		panic("Vector: PopBack() called on empty queue")
	}

	// Calculate new tail position
	v.tail = v.prev(v.tail)

	// Remove value at tail.
	ret := v.buf[v.tail]
	var zero T
	v.buf[v.tail] = zero
	v.count--

	v.shrinkIfExcess()
	return ret
}

// Front returns the element at the front of the queue. This is the element
// that would be returned by PopFront. This call panics if the queue is empty.
func (v *Vector[T]) Front() T {
	if v.count <= 0 {
		panic("Vector: Front() called when empty")
	}
	return v.buf[v.head]
}

func (v *Vector[T]) withinRange(index int) bool {
	return index >= 0 && index < v.count
}

// Back returns the element at the back of the queue. This is the element that
// would be returned by PopBack. This call panics if the queue is empty.
func (v *Vector[T]) Back() T {
	if v.count <= 0 {
		panic("Vector: Back() called when empty")
	}
	return v.buf[v.prev(v.tail)]
}

// At returns the element at index i in the queue without removing the element
// from the queue. This method accepts only non-negative index values. At(0)
// refers to the first element and is the same as Front(). At(Len()-1) refers
// to the last element and is the same as Back(). If the index is invalid, the
// call panics.
//
// The purpose of At is to allow Vector to serve as a more general purpose
// circular buffer, where items are only added to and removed from the ends of
// the Vector, but may be read from any place within the Vector. Consider the
// case of a fixed-size circular log buffer: A new entry is pushed onto one end
// and when full the oldest is popped from the other end. All the log entries
// in the buffer must be readable without altering the buffer contents.
func (v *Vector[T]) Get(i int) T {
	if i < 0 || i >= v.count {
		panic(outOfRangeText(i, v.Len()))
	}
	// bitwise modulus
	return v.buf[(v.head+i)&(len(v.buf)-1)]
}

// Set assigns the item to index i in the queue. Set indexes the Vector the same
// as At but perform the opposite operation. If the index is invalid, the call
// panics.
func (v *Vector[T]) Set(i int, item T) {
	if i < 0 || i >= v.count {
		panic(outOfRangeText(i, v.Len()))
	}
	// bitwise modulus
	v.buf[(v.head+i)&(len(v.buf)-1)] = item
}

// Clear removes all elements from the queue, but retains the current capacity.
// This is useful when repeatedly reusing the queue at high frequency to avoid
// GC during reuse. The queue will not be resized smaller as long as items are
// only added. Only when items are removed is the queue subject to getting
// resized smaller.
func (v *Vector[T]) Clear() {
	var zero T
	modBits := len(v.buf) - 1
	h := v.head
	for i := 0; i < v.Len(); i++ {
		v.buf[(h+i)&modBits] = zero
	}
	v.head = 0
	v.tail = 0
	v.count = 0
}

// Index returns the index into the Deque of the first item satisfying f(item),
// or -1 if none do. If q is nil, then -1 is always returned. Search is linear
// starting with index 0.
func (v *Vector[T]) Index(f func(T) bool) int {
	if v.Len() > 0 {
		modBits := len(v.buf) - 1
		for i := 0; i < v.count; i++ {
			if f(v.buf[(v.head+i)&modBits]) {
				return i
			}
		}
	}
	return -1
}

// RIndex is the same as Index, but searches from Back to Front. The index
// returned is from Front to Back, where index 0 is the index of the item
// returned by Front().
func (v *Vector[T]) RIndex(f func(T) bool) int {
	if v.Len() > 0 {
		modBits := len(v.buf) - 1
		for i := v.count - 1; i >= 0; i-- {
			if f(v.buf[(v.head+i)&modBits]) {
				return i
			}
		}
	}
	return -1
}

// Remove removes and returns an element from the middle of the queue, at the
// specified index. Remove(0) is the same as PopFront() and Remove(Len()-1) is
// the same as PopBack(). Accepts only non-negative index values, and panics if
// index is out of range.
//
// Important: Vector is optimized for O(1) operations at the ends of the queue,
// not for operations in the the middle. Complexity of this function is
// constant plus linear in the lesser of the distances between the index and
// either of the ends of the queue.
func (v *Vector[T]) Remove(at int) T {
	if at < 0 || at >= v.Len() {
		panic(outOfRangeText(at, v.Len()))
	}

	rm := (v.head + at) & (len(v.buf) - 1)
	if at*2 < v.count {
		for i := 0; i < at; i++ {
			prev := v.prev(rm)
			v.buf[prev], v.buf[rm] = v.buf[rm], v.buf[prev]
			rm = prev
		}
		return v.PopFront()
	}
	swaps := v.count - at - 1
	for i := 0; i < swaps; i++ {
		next := v.next(rm)
		v.buf[rm], v.buf[next] = v.buf[next], v.buf[rm]
		rm = next
	}
	return v.PopBack()
}

// SetMinCapacity sets a minimum capacity of 2^minCapacityExp. If the value of
// the minimum capacity is less than or equal to the minimum allowed, then
// capacity is set to the minimum allowed. This may be called at anytime to set
// a new minimum capacity.
//
// Setting a larger minimum capacity may be used to prevent resizing when the
// number of stored items changes frequently across a wide range.
func (v *Vector[T]) SetMinCapacity(minCapacityExp uint) {
	if 1<<minCapacityExp > minCapacity {
		v.minCap = 1 << minCapacityExp
	} else {
		v.minCap = minCapacity
	}
}

// prev returns the previous buffer position wrapping around buffer.
func (v *Vector[T]) prev(i int) int {
	return (i - 1) & (len(v.buf) - 1) // bitwise modulus
}

// next returns the next buffer position wrapping around buffer.
func (v *Vector[T]) next(i int) int {
	return (i + 1) & (len(v.buf) - 1) // bitwise modulus
}

// growIfFull resizes up if the buffer is full.
func (v *Vector[T]) growIfFull() {
	if v.count != len(v.buf) {
		return
	}
	if len(v.buf) == 0 {
		if v.minCap == 0 {
			v.minCap = minCapacity
		}
		v.buf = make([]T, v.minCap)
		return
	}
	v.resize()
}

// shrinkIfExcess resize down if the buffer 1/4 full.
func (v *Vector[T]) shrinkIfExcess() {
	if len(v.buf) > v.minCap && (v.count<<2) == len(v.buf) {
		v.resize()
	}
}

// resize resizes the Vector to fit exactly twice its current contents. This is
// used to grow the queue when it is full, and also to shrink it when it is
// only a quarter full.
func (v *Vector[T]) resize() {
	newBuf := make([]T, v.count<<1)
	if v.tail > v.head {
		copy(newBuf, v.buf[v.head:v.tail])
	} else {
		n := copy(newBuf, v.buf[v.head:])
		copy(newBuf[n:], v.buf[:v.tail])
	}

	v.head = 0
	v.tail = v.count
	v.buf = newBuf
}

func (v *Vector[T]) Values() []T {
	newBuf := make([]T, v.count)
	if v.tail > v.head {
		copy(newBuf, v.buf[v.head:v.tail])
	} else {
		n := copy(newBuf, v.buf[v.head:])
		copy(newBuf[n:], v.buf[:v.tail])
	}
	return newBuf
}

func (v *Vector[T]) Swap(i, j int) {
	t := v.Get(i)
	t1 := v.Get(j)
	v.Set(i, t1)
	v.Set(j, t)
}

func (v *Vector[T]) Less(i, j int) bool {
	return true
}

func outOfRangeText(i, len int) string {
	return fmt.Sprintf("Vector: index out of range %d with length %d", i, len)
}
