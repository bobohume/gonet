package maps

import (
	"fmt"
)

type color bool

const (
	black, red color = true, false
)

type OrderKey interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 |
		uint32 | uint64 | uintptr | float32 | float64 | string
}

// IntComparator provides a basic comparison on int
func Comparator[K OrderKey](a, b K) int {
	switch {
	case a > b:
		return 1
	case a < b:
		return -1
	default:
		return 0
	}
}

// Map interface that all Maps implement
type IMap[K OrderKey, V any] interface {
	Empty() bool
	Size() int
	Clear()
	Values() []V
}

func assertMapImplementation() {
	var _ IMap[int, interface{}] = (*Map[int, interface{}])(nil)
}

// Map holds elements of the red-black tree
type Map[K OrderKey, V any] struct {
	Root *Node[K, V]
	size int
	Nil  V
}

// Node is a single element within the Map
type Node[K OrderKey, V any] struct {
	Key    K
	Value  V
	color  color
	Left   *Node[K, V]
	Right  *Node[K, V]
	Parent *Node[K, V]
}

// Put inserts node into the tree.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (m *Map[K, V]) Put(key K, value V) {
	var insertedNode *Node[K, V]
	if m.Root == nil {
		// Assert key is of comparator's type for initial tree
		Comparator(key, key)
		m.Root = &Node[K, V]{Key: key, Value: value, color: red}
		insertedNode = m.Root
	} else {
		node := m.Root
		loop := true
		for loop {
			compare := Comparator(key, node.Key)
			switch {
			case compare == 0:
				node.Key = key
				node.Value = value
				return
			case compare < 0:
				if node.Left == nil {
					node.Left = &Node[K, V]{Key: key, Value: value, color: red}
					insertedNode = node.Left
					loop = false
				} else {
					node = node.Left
				}
			case compare > 0:
				if node.Right == nil {
					node.Right = &Node[K, V]{Key: key, Value: value, color: red}
					insertedNode = node.Right
					loop = false
				} else {
					node = node.Right
				}
			}
		}
		insertedNode.Parent = node
	}
	m.insertCase1(insertedNode)
	m.size++
}

// Get searches the node in the tree by key and returns its value or nil if key is not found in tree.
// Second return parameter is true if key was found, otherwise false.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (m *Map[K, V]) Get(key K) (value V, found bool) {
	node := m.lookup(key)
	if node != nil {
		return node.Value, true
	}
	return m.Nil, false
}

// Remove remove the node from the tree by key.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (m *Map[K, V]) Remove(key K) {
	var child *Node[K, V]
	node := m.lookup(key)
	if node == nil {
		return
	}
	if node.Left != nil && node.Right != nil {
		pred := node.Left.maximumNode()
		node.Key = pred.Key
		node.Value = pred.Value
		node = pred
	}
	if node.Left == nil || node.Right == nil {
		if node.Right == nil {
			child = node.Left
		} else {
			child = node.Right
		}
		if node.color == black {
			node.color = nodeColor(child)
			m.deleteCase1(node)
		}
		m.replaceNode(node, child)
		if node.Parent == nil && child != nil {
			child.color = black
		}
	}
	m.size--
}

// Empty returns true if tree does not contain any nodes
func (m *Map[K, V]) Empty() bool {
	return m.size == 0
}

// Size returns number of nodes in the tree.
func (m *Map[K, V]) Size() int {
	return m.size
}

// Keys returns all keys in-order
func (m *Map[K, V]) Keys() []K {
	keys := make([]K, m.size)
	it := m.Iterator()
	for i := 0; it.Next(); i++ {
		keys[i] = it.Key()
	}
	return keys
}

// Values returns all values in-order based on the key.
func (m *Map[K, V]) Values() []V {
	values := make([]V, m.size)
	it := m.Iterator()
	for i := 0; it.Next(); i++ {
		values[i] = it.Value()
	}
	return values
}

// Left returns the left-most (min) node or nil if tree is empty.
func (m *Map[K, V]) Left() *Node[K, V] {
	var parent *Node[K, V]
	current := m.Root
	for current != nil {
		parent = current
		current = current.Left
	}
	return parent
}

// Right returns the right-most (max) node or nil if tree is empty.
func (m *Map[K, V]) Right() *Node[K, V] {
	var parent *Node[K, V]
	current := m.Root
	for current != nil {
		parent = current
		current = current.Right
	}
	return parent
}

// Floor Finds floor node of the input key, return the floor node or nil if no floor is found.
// Second return parameter is true if floor was found, otherwise false.
//
// Floor node is defined as the largest node that is smaller than or equal to the given node.
// A floor node may not be found, either because the tree is empty, or because
// all nodes in the tree are larger than the given node.
//
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (m *Map[K, V]) Floor(key K) (floor *Node[K, V], found bool) {
	found = false
	node := m.Root
	for node != nil {
		compare := Comparator(key, node.Key)
		switch {
		case compare == 0:
			return node, true
		case compare < 0:
			node = node.Left
		case compare > 0:
			floor, found = node, true
			node = node.Right
		}
	}
	if found {
		return floor, true
	}
	return nil, false
}

// Ceiling finds ceiling node of the input key, return the ceiling node or nil if no ceiling is found.
// Second return parameter is true if ceiling was found, otherwise false.
//
// Ceiling node is defined as the smallest node that is larger than or equal to the given node.
// A ceiling node may not be found, either because the tree is empty, or because
// all nodes in the tree are smaller than the given node.
//
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (m *Map[K, V]) Ceiling(key K) (ceiling *Node[K, V], found bool) {
	found = false
	node := m.Root
	for node != nil {
		compare := Comparator(key, node.Key)
		switch {
		case compare == 0:
			return node, true
		case compare < 0:
			ceiling, found = node, true
			node = node.Left
		case compare > 0:
			node = node.Right
		}
	}
	if found {
		return ceiling, true
	}
	return nil, false
}

// Clear removes all nodes from the tree.
func (m *Map[K, V]) Clear() {
	m.Root = nil
	m.size = 0
}

// String returns a string representation of container
func (m *Map[K, V]) String() string {
	str := "RedBlackTree\n"
	if !m.Empty() {
		output(m.Root, "", true, &str)
	}
	return str
}

func (node *Node[K, V]) String() string {
	return fmt.Sprintf("%v", node.Key)
}

func output[K OrderKey, V any](node *Node[K, V], prefix string, isTail bool, str *string) {
	if node.Right != nil {
		newPrefix := prefix
		if isTail {
			newPrefix += "│   "
		} else {
			newPrefix += "    "
		}
		output(node.Right, newPrefix, false, str)
	}
	*str += prefix
	if isTail {
		*str += "└── "
	} else {
		*str += "┌── "
	}
	*str += node.String() + "\n"
	if node.Left != nil {
		newPrefix := prefix
		if isTail {
			newPrefix += "    "
		} else {
			newPrefix += "│   "
		}
		output(node.Left, newPrefix, true, str)
	}
}

func (m *Map[K, V]) lookup(key K) *Node[K, V] {
	node := m.Root
	for node != nil {
		compare := Comparator(key, node.Key)
		switch {
		case compare == 0:
			return node
		case compare < 0:
			node = node.Left
		case compare > 0:
			node = node.Right
		}
	}
	return nil
}

func (node *Node[K, V]) grandparent() *Node[K, V] {
	if node != nil && node.Parent != nil {
		return node.Parent.Parent
	}
	return nil
}

func (node *Node[K, V]) uncle() *Node[K, V] {
	if node == nil || node.Parent == nil || node.Parent.Parent == nil {
		return nil
	}
	return node.Parent.sibling()
}

func (node *Node[K, V]) sibling() *Node[K, V] {
	if node == nil || node.Parent == nil {
		return nil
	}
	if node == node.Parent.Left {
		return node.Parent.Right
	}
	return node.Parent.Left
}

func (m *Map[K, V]) rotateLeft(node *Node[K, V]) {
	right := node.Right
	m.replaceNode(node, right)
	node.Right = right.Left
	if right.Left != nil {
		right.Left.Parent = node
	}
	right.Left = node
	node.Parent = right
}

func (m *Map[K, V]) rotateRight(node *Node[K, V]) {
	left := node.Left
	m.replaceNode(node, left)
	node.Left = left.Right
	if left.Right != nil {
		left.Right.Parent = node
	}
	left.Right = node
	node.Parent = left
}

func (m *Map[K, V]) replaceNode(old *Node[K, V], new *Node[K, V]) {
	if old.Parent == nil {
		m.Root = new
	} else {
		if old == old.Parent.Left {
			old.Parent.Left = new
		} else {
			old.Parent.Right = new
		}
	}
	if new != nil {
		new.Parent = old.Parent
	}
}

func (m *Map[K, V]) insertCase1(node *Node[K, V]) {
	if node.Parent == nil {
		node.color = black
	} else {
		m.insertCase2(node)
	}
}

func (m *Map[K, V]) insertCase2(node *Node[K, V]) {
	if nodeColor(node.Parent) == black {
		return
	}
	m.insertCase3(node)
}

func (m *Map[K, V]) insertCase3(node *Node[K, V]) {
	uncle := node.uncle()
	if nodeColor(uncle) == red {
		node.Parent.color = black
		uncle.color = black
		node.grandparent().color = red
		m.insertCase1(node.grandparent())
	} else {
		m.insertCase4(node)
	}
}

func (m *Map[K, V]) insertCase4(node *Node[K, V]) {
	grandparent := node.grandparent()
	if node == node.Parent.Right && node.Parent == grandparent.Left {
		m.rotateLeft(node.Parent)
		node = node.Left
	} else if node == node.Parent.Left && node.Parent == grandparent.Right {
		m.rotateRight(node.Parent)
		node = node.Right
	}
	m.insertCase5(node)
}

func (m *Map[K, V]) insertCase5(node *Node[K, V]) {
	node.Parent.color = black
	grandparent := node.grandparent()
	grandparent.color = red
	if node == node.Parent.Left && node.Parent == grandparent.Left {
		m.rotateRight(grandparent)
	} else if node == node.Parent.Right && node.Parent == grandparent.Right {
		m.rotateLeft(grandparent)
	}
}

func (node *Node[K, V]) maximumNode() *Node[K, V] {
	if node == nil {
		return nil
	}
	for node.Right != nil {
		node = node.Right
	}
	return node
}

func (m *Map[K, V]) deleteCase1(node *Node[K, V]) {
	if node.Parent == nil {
		return
	}
	m.deleteCase2(node)
}

func (m *Map[K, V]) deleteCase2(node *Node[K, V]) {
	sibling := node.sibling()
	if nodeColor(sibling) == red {
		node.Parent.color = red
		sibling.color = black
		if node == node.Parent.Left {
			m.rotateLeft(node.Parent)
		} else {
			m.rotateRight(node.Parent)
		}
	}
	m.deleteCase3(node)
}

func (m *Map[K, V]) deleteCase3(node *Node[K, V]) {
	sibling := node.sibling()
	if nodeColor(node.Parent) == black &&
		nodeColor(sibling) == black &&
		nodeColor(sibling.Left) == black &&
		nodeColor(sibling.Right) == black {
		sibling.color = red
		m.deleteCase1(node.Parent)
	} else {
		m.deleteCase4(node)
	}
}

func (m *Map[K, V]) deleteCase4(node *Node[K, V]) {
	sibling := node.sibling()
	if nodeColor(node.Parent) == red &&
		nodeColor(sibling) == black &&
		nodeColor(sibling.Left) == black &&
		nodeColor(sibling.Right) == black {
		sibling.color = red
		node.Parent.color = black
	} else {
		m.deleteCase5(node)
	}
}

func (m *Map[K, V]) deleteCase5(node *Node[K, V]) {
	sibling := node.sibling()
	if node == node.Parent.Left &&
		nodeColor(sibling) == black &&
		nodeColor(sibling.Left) == red &&
		nodeColor(sibling.Right) == black {
		sibling.color = red
		sibling.Left.color = black
		m.rotateRight(sibling)
	} else if node == node.Parent.Right &&
		nodeColor(sibling) == black &&
		nodeColor(sibling.Right) == red &&
		nodeColor(sibling.Left) == black {
		sibling.color = red
		sibling.Right.color = black
		m.rotateLeft(sibling)
	}
	m.deleteCase6(node)
}

func (m *Map[K, V]) deleteCase6(node *Node[K, V]) {
	sibling := node.sibling()
	sibling.color = nodeColor(node.Parent)
	node.Parent.color = black
	if node == node.Parent.Left && nodeColor(sibling.Right) == red {
		sibling.Right.color = black
		m.rotateLeft(node.Parent)
	} else if nodeColor(sibling.Left) == red {
		sibling.Left.color = black
		m.rotateRight(node.Parent)
	}
}

func nodeColor[K OrderKey, V any](node *Node[K, V]) color {
	if node == nil {
		return black
	}
	return node.color
}
