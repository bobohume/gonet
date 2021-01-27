//
// Copyright (c) 2009-2010 Mikko Mononen memon@inside.org
//
// This software is provided 'as-is', without any express or implied
// warranty.  In no event will the authors be held liable for any damages
// arising from the use of this software.
// Permission is granted to anyone to use this software for any purpose,
// including commercial applications, and to alter it and redistribute it
// freely, subject to the following restrictions:
// 1. The origin of this software must not be misrepresented; you must not
//    claim that you wrote the original software. If you use this software
//    in a product, an acknowledgment in the product documentation would be
//    appreciated but is not required.
// 2. Altered source versions must be plainly marked as such, and must not be
//    misrepresented as being the original software.
// 3. This notice may not be removed or altered from any source distribution.
//

package detour

import "unsafe"

func DtHashRef(polyRef DtPolyRef) uint32 {
	a := uint32(polyRef)
	a += ^(a << 15)
	a ^= (a >> 10)
	a += (a << 3)
	a ^= (a >> 6)
	a += ^(a << 11)
	a ^= (a >> 16)
	return a
}

func (this *DtNodePool) constructor(maxNodes, hashSize uint32) {
	this.m_maxNodes = maxNodes
	this.m_hashSize = hashSize

	DtAssert(DtNextPow2(this.m_hashSize) == this.m_hashSize)
	// pidx is special as 0 means "none" and 1 is the first node. For that reason
	// we have 1 fewer nodes available than the number of values it can contain.
	DtAssert(this.m_maxNodes > 0 && this.m_maxNodes <= uint32(DT_NULL_IDX) && this.m_maxNodes <= (1<<DT_NODE_PARENT_BITS)-1)

	this.m_nodes = make([]DtNode, this.m_maxNodes)
	this.m_next = make([]DtNodeIndex, this.m_maxNodes)
	this.m_first = make([]DtNodeIndex, this.m_hashSize)

	DtAssert(this.m_nodes != nil)
	DtAssert(this.m_next != nil)
	DtAssert(this.m_first != nil)

	for i := 0; i < len(this.m_first); i++ {
		this.m_first[i] = DT_NULL_IDX
	}
	for i := 0; i < len(this.m_next); i++ {
		this.m_next[i] = DT_NULL_IDX
	}

	this.base = uintptr(unsafe.Pointer(&(this.m_nodes[0])))
}

func (this *DtNodePool) destructor() {
	this.m_nodes = nil
	this.m_first = nil
	this.m_next = nil
}

func (this *DtNodePool) Clear() {
	for i := 0; i < len(this.m_first); i++ {
		this.m_first[i] = DT_NULL_IDX
	}
	this.m_nodeCount = 0
}

func (this *DtNodePool) FindNodes(id DtPolyRef, nodes []*DtNode, maxNodes uint32) uint32 {
	var n uint32 = 0
	bucket := DtHashRef(id) & (this.m_hashSize - 1)
	i := this.m_first[bucket]
	for i != DT_NULL_IDX {
		if this.m_nodes[i].Id == id {
			if n >= maxNodes {
				return n
			}
			nodes[n] = &this.m_nodes[i]
			n = n + 1
		}
		i = this.m_next[i]
	}

	return n
}

func (this *DtNodePool) FindNode(id DtPolyRef, state uint8) *DtNode {
	bucket := DtHashRef(id) & (this.m_hashSize - 1)
	i := this.m_first[bucket]
	for i != DT_NULL_IDX {
		if this.m_nodes[i].Id == id && this.m_nodes[i].State == state {
			return &this.m_nodes[i]
		}
		i = this.m_next[i]
	}
	return nil
}

func (this *DtNodePool) GetNode(id DtPolyRef, state uint8) *DtNode {
	bucket := DtHashRef(id) & (this.m_hashSize - 1)
	i := this.m_first[bucket]
	var node *DtNode = nil
	for i != DT_NULL_IDX {
		if this.m_nodes[i].Id == id && this.m_nodes[i].State == state {
			return &this.m_nodes[i]
		}
		i = this.m_next[i]
	}

	if this.m_nodeCount >= this.m_maxNodes {
		return nil
	}

	i = DtNodeIndex(this.m_nodeCount)
	this.m_nodeCount++

	// Init node
	node = &this.m_nodes[i]
	node.Pidx = 0
	node.Cost = 0
	node.Total = 0
	node.Id = id
	node.State = state
	node.Flags = 0

	this.m_next[i] = this.m_first[bucket]
	this.m_first[bucket] = i

	return node
}

func (this *DtNodeQueue) constructor(n int) {
	this.m_capacity = n
	DtAssert(this.m_capacity > 0)
	this.m_heap = make([]*DtNode, this.m_capacity+1)
	DtAssert(this.m_heap != nil)
}

func (this *DtNodeQueue) destructor() {
	this.m_heap = nil
}

func (this *DtNodeQueue) bubbleUp(i int, node *DtNode) {
	parent := (i - 1) / 2
	// note: (index > 0) means there is a parent
	for (i > 0) && (this.m_heap[parent].Total > node.Total) {
		this.m_heap[i] = this.m_heap[parent]
		i = parent
		parent = (i - 1) / 2
	}
	this.m_heap[i] = node
}

func (this *DtNodeQueue) trickleDown(i int, node *DtNode) {
	child := (i * 2) + 1
	for child < this.m_size {
		if ((child + 1) < this.m_size) &&
			(this.m_heap[child].Total > this.m_heap[child+1].Total) {
			child++
		}
		this.m_heap[i] = this.m_heap[child]
		i = child
		child = (i * 2) + 1
	}
	this.bubbleUp(i, node)
}
