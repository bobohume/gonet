package mpsc

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestQueue_PushPop(t *testing.T) {
	q := New()

	q.Push(1)
	q.Push(2)
	t.Log(q.Pop())
	t.Log(q.Pop())
	t.Log(q.Empty())
}

func TestQueue_Empty(t *testing.T) {
	q := New()
	t.Log(q.Empty())
	q.Push(1)
	t.Log( q.Empty())
}

func TestQueue_PushPopOneProducer(t *testing.T) {
	expCount := 100

	var wg sync.WaitGroup
	wg.Add(1)
	q := New()
	go func() {
		i := 0
		for {
			r := q.Pop()
			if r == nil {
				runtime.Gosched()
				continue
			}
			i++
			if i == expCount {
				wg.Done()
				return
			}
		}
	}()

	var val interface{} = "foo"

	for i := 0; i < expCount; i++ {
		q.Push(val)
	}

	wg.Wait()
}

func TestMpscQueueConsistency(t *testing.T) {
	max := 1000000
	c := runtime.NumCPU() / 2
	cmax := max / c
	var wg sync.WaitGroup
	wg.Add(1)
	q := New()
	go func() {
		i := 0
		seen := make(map[string]string)
		for {
			r := q.Pop()
			if r == nil {
				runtime.Gosched()

				continue
			}
			i++
			s := r.(string)
			_, present := seen[s]
			if present {
				log.Printf("item have already been seen %v", s)
				t.FailNow()
			}
			seen[s] = s
			if i == cmax*c {
				wg.Done()
				return
			}
		}
	}()

	for j := 0; j < c; j++ {
		jj := j
		go func() {
			for i := 0; i < cmax; i++ {
				q.Push(fmt.Sprintf("%v %v", jj, i))
			}
		}()
	}

	wg.Wait()
	//time.Sleep(500 * time.Millisecond)
	// queue should be empty
	for i := 0; i < 100; i++ {
		r := q.Pop()
		if r != nil {
			log.Printf("unexpected result %+v", r)
			t.FailNow()
		}
	}
}

func benchmarkPushPop(count, c int) {
	var wg sync.WaitGroup
	wg.Add(1)
	q := New()
	go func() {
		i := 0
		for {
			r := q.Pop()
			if r == nil {
				time.Sleep(1)
				continue
			}
			i++
			if i == count {
				wg.Done()
				return
			}
		}
	}()

	var val interface{} = "foo"

	for i := 0; i < c; i++ {
		go func(n int) {
			for n > 0 {
				q.Push(val)
				n--
			}
		}(count / c)
	}

	wg.Wait()
}

func BenchmarkPushPop(b *testing.B) {
	benchmarks := []struct {
		count       int
		concurrency int
	}{
		{
			count:       1000000,
			concurrency: 1,
		},
		{
			count:       1000000,
			concurrency: 2,
		},
		{
			count:       1000000,
			concurrency: 4,
		},
		{
			count:       1000000,
			concurrency: 8,
		},
		{
			count:       1000000,
			concurrency: 16,
		},
	}
	for _, bm := range benchmarks {
		b.Run(fmt.Sprintf("%d_%d", bm.count, bm.concurrency), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				benchmarkPushPop(bm.count, bm.concurrency)
			}
		})
	}
}

func benchmarkChanPushPop(count, c int) {
	var wg sync.WaitGroup
	wg.Add(1)
	q := make(chan string, 100)
	go func() {
		i := 0
		for {
			select{
			case <-q:
				i++
				if i == count {
					wg.Done()
					return
				}
			}

		}
	}()

	var val = "foo"

	for i := 0; i < c; i++ {
		go func(n int) {
			for n > 0 {
				q <- val
				n--
			}
		}(count / c)
	}

	wg.Wait()
}

func BenchmarkChanPushPop(b *testing.B) {
	benchmarks := []struct {
		count       int
		concurrency int
	}{
		{
			count:       1000000,
			concurrency: 1,
		},
		{
			count:       1000000,
			concurrency: 2,
		},
		{
			count:       1000000,
			concurrency: 4,
		},
		{
			count:       1000000,
			concurrency: 8,
		},
		{
			count:       1000000,
			concurrency: 16,
		},
	}
	for _, bm := range benchmarks {
		b.Run(fmt.Sprintf("%d_%d", bm.count, bm.concurrency), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				benchmarkChanPushPop(bm.count, bm.concurrency)
			}
		})
	}
}

var g_MailChan =  make(chan bool, 1)
var g_bMailIn [8]int64

func benchmarkPushPopActor(count, c int) {
	var wg sync.WaitGroup
	wg.Add(1)
	q := New()
	go func() {
		i := 0
		for {
			select {
			case <- g_MailChan:
				atomic.StoreInt64(&g_bMailIn[0], 0)
				for data := q.Pop(); data != nil; data = q.Pop() {
					i++
					if i == count {
						wg.Done()
						return
					}
				}
			}
		}
	}()

	var val interface{} = "foo"

	for i := 0; i < c; i++ {
		go func(n int) {
			for n > 0 {
				q.Push(val)
				if atomic.LoadInt64(&g_bMailIn[0]) == 0 && atomic.CompareAndSwapInt64(&g_bMailIn[0], 0, 1) {
					g_MailChan <- true
				}
				n--
			}
		}(count / c)
	}

	wg.Wait()
}

func BenchmarkPushPopActor(b *testing.B) {
	benchmarks := []struct {
		count       int
		concurrency int
	}{
		{
			count:       1000000,
			concurrency: 1,
		},
		{
			count:       1000000,
			concurrency: 2,
		},
		{
			count:       1000000,
			concurrency: 4,
		},
		{
			count:       1000000,
			concurrency: 8,
		},
		{
			count:       1000000,
			concurrency: 16,
		},
	}
	for _, bm := range benchmarks {
		b.Run(fmt.Sprintf("%d_%d", bm.count, bm.concurrency), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				benchmarkPushPopActor(bm.count, bm.concurrency)
			}
		})
	}
}
