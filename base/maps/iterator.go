package maps

import "gonet/base/containers"

func assertIteratorImplementation() {
	var _ containers.ReverseIteratorWithKey = (*Iterator)(nil)
}

// Iterator holding the iterator's state
type Iterator struct {
	maps     *Map
	node     *Node
	position position
}

type position byte

const (
	begin, between, end position = 0, 1, 2
)

// Iterator returns a stateful iterator whose elements are key/value pairs.
func (this *Map) Iterator() Iterator {
	return Iterator{maps: this, node: nil, position: begin}
}

// Next moves the iterator to the next element and returns true if there was a next element in the container.
// If Next() returns true, then next element's key and value can be retrieved by Key() and Value().
// If Next() was called for the first time, then it will point the iterator to the first element if it exists.
// Modifies the state of the iterator.
func (this *Iterator) Next() bool {
	if this.position == end {
		goto end
	}
	if this.position == begin {
		left := this.maps.Left()
		if left == nil {
			goto end
		}
		this.node = left
		goto between
	}
	if this.node.Right != nil {
		this.node = this.node.Right
		for this.node.Left != nil {
			this.node = this.node.Left
		}
		goto between
	}
	if this.node.Parent != nil {
		node := this.node
		for this.node.Parent != nil {
			this.node = this.node.Parent
			if this.maps.Comparator(node.Key, this.node.Key) <= 0 {
				goto between
			}
		}
	}

end:
	this.node = nil
	this.position = end
	return false

between:
	this.position = between
	return true
}

// Prev moves the iterator to the previous element and returns true if there was a previous element in the container.
// If Prev() returns true, then previous element's key and value can be retrieved by Key() and Value().
// Modifies the state of the iterator.
func (this *Iterator) Prev() bool {
	if this.position == begin {
		goto begin
	}
	if this.position == end {
		right := this.maps.Right()
		if right == nil {
			goto begin
		}
		this.node = right
		goto between
	}
	if this.node.Left != nil {
		this.node = this.node.Left
		for this.node.Right != nil {
			this.node = this.node.Right
		}
		goto between
	}
	if this.node.Parent != nil {
		node := this.node
		for this.node.Parent != nil {
			this.node = this.node.Parent
			if this.maps.Comparator(node.Key, this.node.Key) >= 0 {
				goto between
			}
		}
	}

begin:
	this.node = nil
	this.position = begin
	return false

between:
	this.position = between
	return true
}

// Value returns the current element's value.
// Does not modify the state of the iterator.
func (this *Iterator) Value() interface{} {
	return this.node.Value
}

// Key returns the current element's key.
// Does not modify the state of the iterator.
func (this *Iterator) Key() interface{} {
	return this.node.Key
}

// Begin resets the iterator to its initial state (one-before-first)
// Call Next() to fetch the first element if any.
func (this *Iterator) Begin() {
	this.node = nil
	this.position = begin
}

// End moves the iterator past the last element (one-past-the-end).
// Call Prev() to fetch the last element if any.
func (this *Iterator) End() {
	this.node = nil
	this.position = end
}

// First moves the iterator to the first element and returns true if there was a first element in the container.
// If First() returns true, then first element's key and value can be retrieved by Key() and Value().
// Modifies the state of the iterator
func (this *Iterator) First() bool {
	this.Begin()
	return this.Next()
}

// Last moves the iterator to the last element and returns true if there was a last element in the container.
// If Last() returns true, then last element's key and value can be retrieved by Key() and Value().
// Modifies the state of the iterator.
func (this *Iterator) Last() bool {
	this.End()
	return this.Prev()
}
