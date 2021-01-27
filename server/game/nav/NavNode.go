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

import (
	"unsafe"
)

type DtNodeFlags uint8

const (
	DT_NODE_OPEN            DtNodeFlags = 0x01
	DT_NODE_CLOSED          DtNodeFlags = 0x02
	DT_NODE_PARENT_DETACHED DtNodeFlags = 0x04 // parent of the node is not adjacent. Found using raycast.
)

type DtNodeIndex uint16

const DT_NULL_IDX DtNodeIndex = ^DtNodeIndex(0)

const DT_NODE_PARENT_BITS uint32 = 24
const DT_NODE_STATE_BITS uint32 = 2
const DT_NODE_FLAGS_BITS uint32 = 3

type DtNode struct {
	Pos   [3]float32  ///< Position of the node.
	Cost  float32     ///< Cost from previous node to current node.
	Total float32     ///< Cost up to the node.
	Pidx  uint32      ///< Index to parent node.
	State uint8       ///< extra state information. A polyRef can have multiple nodes with different extra info. see DT_MAX_STATES_PER_NODE
	Flags DtNodeFlags ///< Node flags. A combination of dtNodeFlags.
	Id    DtPolyRef   ///< Polygon ref the node corresponds to.
}

const DT_MAX_STATES_PER_NODE int = 1 << DT_NODE_STATE_BITS // number of extra states per node. See dtNode::state

var sizeofNode = uint32(unsafe.Sizeof(DtNode{}))

type DtNodePool struct {
	m_nodes     []DtNode
	m_first     []DtNodeIndex
	m_next      []DtNodeIndex
	m_maxNodes  uint32
	m_hashSize  uint32
	m_nodeCount uint32

	base uintptr
}

func (this *DtNodePool) GetNodeIdx(node *DtNode) uint32 {
	if node == nil {
		return 0
	}
	current := uintptr(unsafe.Pointer(node))
	return (uint32)(current-this.base)/sizeofNode + 1
}

func (this *DtNodePool) GetNodeAtIdx(idx uint32) *DtNode {
	if idx == 0 {
		return nil
	}
	return &this.m_nodes[idx-1]
}

func (this *DtNodePool) GetMemUsed() uint32 {
	return uint32(unsafe.Sizeof(*this)) +
		uint32(unsafe.Sizeof(&this.m_nodes[0]))*this.m_maxNodes +
		uint32(unsafe.Sizeof(&this.m_next[0]))*this.m_maxNodes +
		uint32(unsafe.Sizeof(&this.m_first[0]))*this.m_hashSize
}

func (this *DtNodePool) GetMaxNodes() uint32             { return this.m_maxNodes }
func (this *DtNodePool) GetHashSize() uint32             { return this.m_hashSize }
func (this *DtNodePool) GetFirst(bucket int) DtNodeIndex { return this.m_first[bucket] }
func (this *DtNodePool) GetNext(i int) DtNodeIndex       { return this.m_next[i] }
func (this *DtNodePool) GetNodeCount() uint32            { return this.m_nodeCount }

func DtAllocNodePool(maxNodes, hashSize uint32) *DtNodePool {
	pool := &DtNodePool{}
	pool.constructor(maxNodes, hashSize)
	return pool
}

func DtFreeNodePool(pool *DtNodePool) {
	if pool == nil {
		return
	}
	pool.destructor()
}

type DtNodeQueue struct {
	m_heap     []*DtNode
	m_capacity int
	m_size     int
}

func (this *DtNodeQueue) Clear() { this.m_size = 0 }

func (this *DtNodeQueue) Top() *DtNode { return this.m_heap[0] }

func (this *DtNodeQueue) Pop() *DtNode {
	result := this.m_heap[0]
	this.m_size--
	this.trickleDown(0, this.m_heap[this.m_size])
	return result
}

func (this *DtNodeQueue) Push(node *DtNode) {
	this.m_size++
	this.bubbleUp(this.m_size-1, node)
}

func (this *DtNodeQueue) Modify(node *DtNode) {
	for i := 0; i < this.m_size; i++ {
		if this.m_heap[i] == node {
			this.bubbleUp(i, node)
			return
		}
	}
}

func (this *DtNodeQueue) Empty() bool { return this.m_size == 0 }

func (this *DtNodeQueue) GetMemUsed() uint32 {
	return uint32(unsafe.Sizeof(*this)) +
		uint32(unsafe.Sizeof(&this.m_heap[0]))*uint32(this.m_capacity+1)
}

func (this *DtNodeQueue) GetCapacity() int { return this.m_capacity }

func DtAllocNodeQueue(n int) *DtNodeQueue {
	queue := &DtNodeQueue{}
	queue.constructor(n)
	return queue
}

func DtFreeNodeQueue(queue *DtNodeQueue) {
	if queue == nil {
		return
	}
	queue.destructor()
}
