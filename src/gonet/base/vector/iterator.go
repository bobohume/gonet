package vector

import "gonet/base/containers"

func assertIteratorImplementation() {
	var _ containers.ReverseIteratorWithIndex = (*Iterator)(nil)
}

// Iterator holding the iterator's state
type Iterator struct {
	vec  *Vector
	index int
}

// Iterator returns a stateful iterator whose values can be fetched by an index.
func (this *Vector) Iterator() Iterator {
	return Iterator{vec: this, index: -1}
}

// Next moves the iterator to the next element and returns true if there was a next element in the container.
// If Next() returns true, then next element's index and value can be retrieved by Index() and Value().
// If Next() was called for the first time, then it will point the iterator to the first element if it exists.
// Modifies the state of the iterator.
func (this *Iterator) Next() bool {
	if this.index < this.vec.mElementCount{
		this.index++
	}
	return this.vec.withinRange(this.index)
}

// Prev moves the iterator to the previous element and returns true if there was a previous element in the container.
// If Prev() returns true, then previous element's index and value can be retrieved by Index() and Value().
// Modifies the state of the iterator.
func (this *Iterator) Prev() bool {
	if this.index >= 0 {
		this.index--
	}
	return this.vec.withinRange(this.index)
}

// Value returns the current element's value.
// Does not modify the state of the iterator.
func (this *Iterator) Value() interface{} {
	return this.vec.Get(this.index)
}

// Index returns the current element's index.
// Does not modify the state of the iterator.
func (this *Iterator) Index() int {
	return this.index
}

// Begin resets the iterator to its initial state (one-before-first)
// Call Next() to fetch the first element if any.
func (this *Iterator) Begin() {
	this.index = -1
}

// End moves the iterator past the last element (one-past-the-end).
// Call Prev() to fetch the last element if any.
func (this *Iterator) End() {
	this.index = this.vec.mElementCount
}

// First moves the iterator to the first element and returns true if there was a first element in the container.
// If First() returns true, then first element's index and value can be retrieved by Index() and Value().
// Modifies the state of the iterator.
func (this *Iterator) First() bool {
	this.Begin()
	return this.Next()
}

// Last moves the iterator to the last element and returns true if there was a last element in the container.
// If Last() returns true, then last element's index and value can be retrieved by Index() and Value().
// Modifies the state of the iterator.
func (this *Iterator) Last() bool {
	this.End()
	return this.Prev()
}
