package base

import (
	"errors"
	"hash/crc32"
	"gonet/base/maps"
	"strconv"
	"sync"
)

const(
	REPLICASNUM = 5
)

// ErrEmptyRing is the error returned when trying to get an element when nothing has been added to hash.
var ErrEmptyRing = errors.New("empty ring")

type (
	// HashRing holds the information about the members of the consistent hash ring.
	HashRing struct {
		m_RingMap           map[uint32] string
		m_MemberMap         map[string] bool
		m_SortedKeys    	*maps.Map
		m_Scratch          	[64]byte// prevent false sharing of the sequence cursor by padding the CPU cache line with 64 *bytes* of data.
		sync.RWMutex
	}

	IHashRing interface {
		Add(elt string)
		Remove(elt string)
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

// eltKey generates a string key for an element with an index.
func (this *HashRing) eltKey(elt string, idx int) string {
	// return elt + "|" + strconv.Itoa(idx)
	return strconv.Itoa(idx) + elt
}

// need c.Lock() before calling
func (this *HashRing) add(elt string) {
	for i := 0; i < REPLICASNUM; i++ {
		Id := this.hashKey(this.eltKey(elt, i))
		this.m_RingMap[Id] = elt
		this.m_SortedKeys.Put(Id, true)
	}
	this.m_MemberMap[elt] = true
}

// need c.Lock() before calling
func (this *HashRing) remove(elt string) {
	for i := 0; i < REPLICASNUM; i++ {
		Id := this.hashKey(this.eltKey(elt, i))
		delete(this.m_RingMap, Id)
		this.m_SortedKeys.Remove(Id)
	}
	delete(this.m_MemberMap, elt)
}

func (this *HashRing) hashKey(key string) uint32 {
	if len(key) < 64 {
		var scratch [64]byte
		copy(scratch[:], key)
		return crc32.ChecksumIEEE(scratch[:len(key)])
	}
	return crc32.ChecksumIEEE([]byte(key))
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
	key := this.hashKey(name)
	if !bOk{
		return ErrEmptyRing, ""
	}
}

func (this *HashRing) Get64(val int64)  (error, uint32) {
	this.RLock()
	defer this.RUnlock()
	if len(this.m_RingMap) == 0 {
		return ErrEmptyRing, 0
	}
	key := this.hashKey(strconv.FormatInt(val, 10))
	if !bOk{
		return ErrEmptyRing, 0
	}
}