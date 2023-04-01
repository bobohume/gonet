package base

import (
	"errors"
	"gonet/base/maps"
	"hash/crc32"
	"strconv"
	"sync"
)

const (
	REPLICASNUM = 5
)

// ErrEmptyRing is the error returned when trying to get an element when nothing has been added to hash.
var ErrEmptyRing = errors.New("empty ring")

type (
	// HashRing holds the information about the members of the consistent hash ring.
	HashRing struct {
		ringMap    map[uint32]string
		memberMap  map[string]bool
		sortedKeys *maps.Map[uint32, uint32]
		sync.RWMutex
	}

	IHashRing interface {
		Add(elt string)
		Remove(elt string)
		HasMember(elt string) bool
		Members() []string
		Get(name string) (error, string)
		Get64(val int64) (error, uint32)
	}
)

// New creates a new HashRing( object with a default setting of 20 replicas for each entry.
// To change the number of replicas, set NumberOfReplicas before adding entries.
func NewHashRing() *HashRing {
	ring := new(HashRing)
	ring.ringMap = make(map[uint32]string)
	ring.memberMap = make(map[string]bool)
	ring.sortedKeys = &maps.Map[uint32, uint32]{}
	return ring
}

func hash(key string) uint32 {
	if len(key) < 64 {
		var scratch [64]byte
		copy(scratch[:], key)
		return crc32.ChecksumIEEE(scratch[:len(key)])
	}
	return crc32.ChecksumIEEE([]byte(key))
}

// eltKey generates a string key for an element with an index.
func (h *HashRing) eltKey(elt string, idx int) string {
	// return elt + "|" + strconv.Itoa(idx)
	return strconv.Itoa(idx) + elt
}

// need c.Lock() before calling
func (h *HashRing) add(elt string) {
	for i := 0; i < REPLICASNUM; i++ {
		Id := hash(h.eltKey(elt, i))
		h.ringMap[Id] = elt
		h.sortedKeys.Put(Id, hash(elt))
	}
	h.memberMap[elt] = true
}

// need c.Lock() before calling
func (h *HashRing) remove(elt string) {
	for i := 0; i < REPLICASNUM; i++ {
		Id := hash(h.eltKey(elt, i))
		delete(h.ringMap, Id)
		h.sortedKeys.Remove(Id)
	}
	delete(h.memberMap, elt)
}

// Add inserts a string element in the consistent hash.
func (h *HashRing) Add(elt string) {
	h.Lock()
	defer h.Unlock()
	h.add(elt)
}

// Remove removes an element from the hash.
func (h *HashRing) Remove(elt string) {
	h.Lock()
	defer h.Unlock()
	h.remove(elt)
}

func (h *HashRing) HasMember(elt string) bool {
	h.RLock()
	defer h.RUnlock()
	_, bEx := h.memberMap[elt]
	return bEx
}

func (h *HashRing) Members() []string {
	h.RLock()
	defer h.RUnlock()
	var m []string
	for k := range h.memberMap {
		m = append(m, k)
	}
	return m
}

// Get returns an element close to where name hashes to in the ring.
func (h *HashRing) Get(name string) (error, string) {
	h.RLock()
	defer h.RUnlock()
	if len(h.ringMap) == 0 {
		return ErrEmptyRing, ""
	}
	key := hash(name)
	node, bOk := h.sortedKeys.Ceiling(key)
	if !bOk {
		itr := h.sortedKeys.Iterator()
		if itr.First() {
			return nil, h.ringMap[itr.Key()]
		}
		return ErrEmptyRing, ""
	}
	return nil, h.ringMap[node.Key]
}

func (h *HashRing) Get64(val int64) (error, uint32) {
	h.RLock()
	defer h.RUnlock()
	if len(h.ringMap) == 0 {
		return ErrEmptyRing, 0
	}
	key := hash(strconv.FormatInt(val, 10))
	node, bOk := h.sortedKeys.Ceiling(key)
	if !bOk {
		itr := h.sortedKeys.Iterator()
		if itr.First() {
			return nil, itr.Value()
		}
		return ErrEmptyRing, 0
	}
	return nil, node.Value
}

// use for stubring
type (
	// HashRing holds the information about the members of the consistent hash ring.
	StubHashRing struct {
		sortedKeys *maps.Map[uint32, uint32]
	}

	IStuHashRing interface {
		Init(endpoints []string)
		Get(val int64) (error, uint32)
	}
)

// eltKey generates a string key for an element with an index.
func (h *StubHashRing) eltKey(elt string, idx int) string {
	// return elt + "|" + strconv.Itoa(idx)
	return strconv.Itoa(idx) + elt
}

func (h *StubHashRing) add(elt string) {
	for i := 0; i < REPLICASNUM; i++ {
		Id := hash(h.eltKey(elt, i))
		h.sortedKeys.Put(Id, hash(elt))
	}
}

// Add inserts a string element in the consistent hash.
func (h *StubHashRing) Init(endpoints []string) {
	h.sortedKeys = &maps.Map[uint32, uint32]{}
	for _, v := range endpoints {
		h.add(v)
	}
}

func (h *StubHashRing) Get(val int64) (error, uint32) {
	key := hash(strconv.FormatInt(val, 10))
	node, bOk := h.sortedKeys.Ceiling(key)
	if !bOk {
		itr := h.sortedKeys.Iterator()
		if itr.First() {
			return nil, itr.Value()
		}
		return ErrEmptyRing, 0
	}
	return nil, node.Value
}
