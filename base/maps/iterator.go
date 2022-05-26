package maps

// Iterator holding the iterator's state
type Iterator[K OrderKey, V any] struct {
	maps     *Map[K, V]
	node     *Node[K, V]
	position position
}

type position byte

const (
	begin, between, end position = 0, 1, 2
)

// Iterator returns a stateful iterator whose elements are key/value pairs.
func (it *Map[K, V]) Iterator() Iterator[K, V] {
	return Iterator[K, V]{maps: it, node: nil, position: begin}
}

// Next moves the iterator to the next element and returns true if there was a next element in the container.
// If Next() returns true, then next element's key and value can be retrieved by Key() and Value().
// If Next() was called for the first time, then it will point the iterator to the first element if it exists.
// Modifies the state of the iterator.
func (it *Iterator[K, V]) Next() bool {
	if it.position == end {
		goto end
	}
	if it.position == begin {
		left := it.maps.Left()
		if left == nil {
			goto end
		}
		it.node = left
		goto between
	}
	if it.node.Right != nil {
		it.node = it.node.Right
		for it.node.Left != nil {
			it.node = it.node.Left
		}
		goto between
	}
	if it.node.Parent != nil {
		node := it.node
		for it.node.Parent != nil {
			it.node = it.node.Parent
			if Comparator(node.Key, it.node.Key) <= 0 {
				goto between
			}
		}
	}

end:
	it.node = nil
	it.position = end
	return false

between:
	it.position = between
	return true
}

// Prev moves the iterator to the previous element and returns true if there was a previous element in the container.
// If Prev() returns true, then previous element's key and value can be retrieved by Key() and Value().
// Modifies the state of the iterator.
func (it *Iterator[K, V]) Prev() bool {
	if it.position == begin {
		goto begin
	}
	if it.position == end {
		right := it.maps.Right()
		if right == nil {
			goto begin
		}
		it.node = right
		goto between
	}
	if it.node.Left != nil {
		it.node = it.node.Left
		for it.node.Right != nil {
			it.node = it.node.Right
		}
		goto between
	}
	if it.node.Parent != nil {
		node := it.node
		for it.node.Parent != nil {
			it.node = it.node.Parent
			if Comparator(node.Key, it.node.Key) >= 0 {
				goto between
			}
		}
	}

begin:
	it.node = nil
	it.position = begin
	return false

between:
	it.position = between
	return true
}

// Value returns the current element's value.
// Does not modify the state of the iterator.
func (it *Iterator[K, V]) Value() V {
	return it.node.Value
}

// Key returns the current element's key.
// Does not modify the state of the iterator.
func (it *Iterator[K, V]) Key() K {
	return it.node.Key
}

// Begin resets the iterator to its initial state (one-before-first)
// Call Next() to fetch the first element if any.
func (it *Iterator[K, V]) Begin() {
	it.node = nil
	it.position = begin
}

// End moves the iterator past the last element (one-past-the-end).
// Call Prev() to fetch the last element if any.
func (it *Iterator[K, V]) End() {
	it.node = nil
	it.position = end
}

// First moves the iterator to the first element and returns true if there was a first element in the container.
// If First() returns true, then first element's key and value can be retrieved by Key() and Value().
// Modifies the state of the iterator
func (it *Iterator[K, V]) First() bool {
	it.Begin()
	return it.Next()
}

// Last moves the iterator to the last element and returns true if there was a last element in the container.
// If Last() returns true, then last element's key and value can be retrieved by Key() and Value().
// Modifies the state of the iterator.
func (it *Iterator[K, V]) Last() bool {
	it.End()
	return it.Prev()
}
