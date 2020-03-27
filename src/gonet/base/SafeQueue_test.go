package base

import (
	"sync/atomic"
	"testing"
)

var(
	nTimes = int32(10000000)
)

func TestSafeQueue(t *testing.T)  {
	a := SafeQueue{}
	a.Init(int64(1024))
	go func() {
		i := 0
		for i = 0; i < int(nTimes); i++ {
			a.Push(i)
		}
	}()

	bb := int32(1)
	bStop := false

	for{
		v := a.Pop()
		if v != nil{
			atomic.AddInt32(&bb, 1)
		}else{
			//time.Sleep(1)
		}
		if atomic.LoadInt32(&bb) == nTimes{
			bStop = true
		}
		if bStop{
			break
		}
	}
}

func TestChan(t *testing.T)  {
	a := make(chan int, 1024)
	bb := int32(1)
	bStop := false
	go func() {
		for i := 0; i < int(nTimes); i++ {
			a <- i
		}
	}()
	for i := 0; i <int(nTimes); i++{
		b := <-a
		b++
		atomic.AddInt32(&bb, 1)
		if atomic.LoadInt32(&bb) == nTimes{
			bStop = true
		}
		if bStop{
			break
		}
	}
}