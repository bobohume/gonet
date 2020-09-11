package base_test

import (
	"fmt"
	"gonet/base"
	"testing"
)

const(
	nHashRingTimes = int32(100000)
	nHashRingSize = 1000
)

func TestHashRing(t *testing.T)  {
	c := base.NewHashRing()
	for i := 0; i < nHashRingSize; i++ {
		c.Add(fmt.Sprintf("%d", i))
	}

	for i := 0; i < int(nHashRingTimes); i++ {
		c.Get("1")
	}
}

func TestRandom(t *testing.T)  {
	vec := base.NewVector()
	for i := 0; i < nHashRingSize * base.REPLICASNUM; i++ {
		vec.Push_back(i)
	}
	for i := 0; i < int(nHashRingTimes); i++ {
		nRand := base.RAND.RandI(0, vec.Len()-1)
		vec.Get(nRand)
	}
}

func TestMod(t *testing.T)  {
	vec := base.NewVector()
	for i := 0; i < nHashRingSize * base.REPLICASNUM; i++ {
		vec.Push_back(i)
	}
	for i := 0; i < int(nHashRingTimes); i++ {
		nRand := i % vec.Size()
		nRand++
	}
}