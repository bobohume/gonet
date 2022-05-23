package base

import (
	"math/rand"
	"time"
)

type (
	Rand struct {
		*rand.Rand
	}
)

func (r *Rand) RandI(i int, n int) int {
	if i > n {
		Assert(false, "Rand::RandI: inverted range")
		return i
	}

	return int(i + r.Int()%(n-i+1))
}

func (r *Rand) RandF(i float32, n float32) float32 {
	if i > n {
		Assert(false, "Rand::RandF: inverted range")
		return i
	}

	return (i + (n-i)*r.Float32())
}

var RAND = Rand{rand.New(rand.NewSource(time.Now().UnixNano()))}
