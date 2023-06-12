package timer_test

import (
	"fmt"
	"gonet/base/timer"
	"sync/atomic"
	"testing"
	"time"
)

const (
	TIMERS = int64(10000000)
)

func TestTimer(t *testing.T) {
	kk := int64(0)
	for i := int64(0); i < TIMERS; i++ {
		id := new(int64)
		*id = i + 1
		timer.RegisterTimer(id, 100*time.Millisecond, func() {
			atomic.AddInt64(&kk, 1)
		})
	}
	for {
		kk1 := atomic.LoadInt64(&kk)
		if kk1 >= TIMERS-1 {
			fmt.Println(111, kk1)
			break
		} else {
			time.Sleep(time.Millisecond * 100)
		}
	}
}
