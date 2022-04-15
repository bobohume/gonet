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
		m_RingMap    map[uint32]string
		m_MemberMap  map[string]bool
		m_SortedKeys *maps.Map
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
	pRing := new(HashRing)
	pRing.m_RingMap = make(map[uint32]string)
	pRing.m_MemberMap = make(map[string]bool)
	pRing.m_SortedKeys = maps.NewWithUInt32Comparator()
	return pRing
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
func (this *HashRing) eltKey(elt string, idx int) string {
	// return elt + "|" + strconv.Itoa(idx)
	return strconv.Itoa(idx) + elt
}

// need c.Lock() before calling
func (this *HashRing) add(elt string) {
	for i := 0; i < REPLICASNUM; i++ {
		Id := hash(this.eltKey(elt, i))
		this.m_RingMap[Id] = elt
		this.m_SortedKeys.Put(Id, hash(elt))
	}
	this.m_MemberMap[elt] = true
}

// need c.Lock() before calling
func (this *HashRing) remove(elt string) {
	for i := 0; i < REPLICASNUM; i++ {
		Id := hash(this.eltKey(elt, i))
		delete(this.m_RingMap, Id)
		this.m_SortedKeys.Remove(Id)
	}
	delete(this.m_MemberMap, elt)
}

// Add inserts a string element in the consistent hash.
func (this *HashRing) Add(elt string) {
	this.Lock()
	defer this.Unlock()
	this.add(elt)
}

// Remove removes an element from the hash.
func (this *HashRing) Remove(elt string) {
	this.Lock()
	defer this.Unlock()
	this.remove(elt)
}

func (this *HashRing) HasMember(elt string) bool {
	this.RLock()
	defer this.RUnlock()
	_, bEx := this.m_MemberMap[elt]
	return bEx
}

func (this *HashRing) Members() []string {
	this.RLock()
	defer this.RUnlock()
	var m []string
	for k := range this.m_MemberMap {
		m = append(m, k)
	}
	return m
}

// Get returns an element close to where name hashes to in the ring.
func (this *HashRing) Get(name string) (error, string) {
	this.RLock()
	defer this.RUnlock()
	if len(this.m_RingMap) == 0 {
		return ErrEmptyRing, ""
	}
	key := hash(name)
	node, bOk := this.m_SortedKeys.Ceiling(key)
	if !bOk {
		itr := this.m_SortedKeys.Iterator()
		if itr.First() {
			return nil, this.m_RingMap[itr.Key().(uint32)]
		}
		return ErrEmptyRing, ""
	}
	return nil, this.m_RingMap[node.Key.(uint32)]
}

func (this *HashRing) Get64(val int64) (error, uint32) {
	this.RLock()
	defer this.RUnlock()
	if len(this.m_RingMap) == 0 {
		return ErrEmptyRing, 0
	}
	key := hash(strconv.FormatInt(val, 10))
	node, bOk := this.m_SortedKeys.Ceiling(key)
	if !bOk {
		itr := this.m_SortedKeys.Iterator()
		if itr.First() {
			return nil, itr.Value().(uint32)
		}
		return ErrEmptyRing, 0
	}
	return nil, node.Value.(uint32)
}

// use for stubring
type (
	// HashRing holds the information about the members of the consistent hash ring.
	StubHashRing struct {
		m_SortedKeys *maps.Map
	}

	IStuHashRing interface {
		Init(endpoints []string)
		Get(val int64) (error, uint32)
	}
)

// eltKey generates a string key for an element with an index.
func (this *StubHashRing) eltKey(elt string, idx int) string {
	// return elt + "|" + strconv.Itoa(idx)
	return strconv.Itoa(idx) + elt
}

func (this *StubHashRing) add(elt string) {
	for i := 0; i < REPLICASNUM; i++ {
		Id := hash(this.eltKey(elt, i))
		this.m_SortedKeys.Put(Id, hash(elt))
	}
}

// Add inserts a string element in the consistent hash.
func (this *StubHashRing) Init(endpoints []string) {
	this.m_SortedKeys = maps.NewWithUInt32Comparator()
	for _, v := range endpoints {
		this.add(v)
	}
}

func (this *StubHashRing) Get(val int64) (error, uint32) {
	key := hash(strconv.FormatInt(val, 10))
	node, bOk := this.m_SortedKeys.Ceiling(key)
	if !bOk {
		itr := this.m_SortedKeys.Iterator()
		if itr.First() {
			return nil, itr.Value().(uint32)
		}
		return ErrEmptyRing, 0
	}
	return nil, node.Value.(uint32)
}
