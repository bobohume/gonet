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

type RandType interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 |
		~uint32 | ~uint64
}

func RandI[T RandType](i T, n T) T {
	if i > n {
		Assert(false, "Rand::RandI: inverted range")
		return i
	}

	return i + T(RAND.Int())%(n-i+1)
}

func RandF[T ~float32 | ~float64](i T, n T) T {
	if i > n {
		Assert(false, "Rand::RandI: inverted range")
		return i
	}

	return (i + (n-i)*T(RAND.Float64()))
}

var RAND = rand.New(rand.NewSource(time.Now().UnixNano()))
