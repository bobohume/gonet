package main

import (
	"base"
	"fmt"
	"sync"
	"time"
)

func testLock(){
	mm := make(chan int, 1024*1024)
	wait := sync.WaitGroup{}
	test := func() {
		wait.Add(8)
		go func() {
			for i := 0; i < 1000000; i++{
				mm <- i
			}
			wait.Done()
		}()

		go func() {
			for i := 0; i < 1000000; i++{
				<- mm
			}
			wait.Done()
		}()
		go func() {
			for i := 0; i < 1000000; i++{
				mm <- i
			}
			wait.Done()
		}()

		go func() {
			for i := 0; i < 1000000; i++{
				<- mm
			}
			wait.Done()
		}()
		go func() {
			for i := 0; i < 1000000; i++{
				mm <- i
			}
			wait.Done()
		}()

		go func() {
			for i := 0; i < 1000000; i++{
				<- mm
			}
			wait.Done()
		}()
		go func() {
			for i := 0; i < 1000000; i++{
				mm <- i
			}
			wait.Done()
		}()

		go func() {
			for i := 0; i < 1000000; i++{
				<- mm
			}
			wait.Done()
		}()
	}

	time1 := time.Now().UnixNano()
	test()
	test()
	time2 := time.Now().UnixNano()
	wait.Wait()
	fmt.Println(time1, time2, time2 - time1)
}

func testNoblack(){
	mm := base.NewRingBuffer(1024*1024)
	wait := &sync.WaitGroup{}
	test := func() {
		wait.Add(8)
		go func() {
			for i := 0; i < 1000000; i++{
				mm.Put(i)
			}
			wait.Done()
		}()

		go func() {
			for i := 0; i < 1000000; i++{
				mm.Get()
			}
			wait.Done()
		}()
		go func() {
			for i := 0; i < 1000000; i++{
				mm.Put(i)
			}
			wait.Done()
		}()

		go func() {
			for i := 0; i < 1000000; i++{
				mm.Get()
			}
			wait.Done()
		}()
		go func() {
			for i := 0; i < 1000000; i++{
				mm.Put(i)
			}
			wait.Done()
		}()

		go func() {
			for i := 0; i < 1000000; i++{
				mm.Get()
			}
			wait.Done()
		}()
		go func() {
			for i := 0; i < 1000000; i++{
				mm.Put(i)
			}
			wait.Done()
		}()

		go func() {
			for i := 0; i < 1000000; i++{
				mm.Get()
			}
			wait.Done()
		}()
	}

	time1 := time.Now().UnixNano()
	test()
	test()
	time2 := time.Now().UnixNano()
	wait.Wait()
	fmt.Println("no lock ", time1, time2, time2 - time1)
}

func main() {
	for i := 0; i < 30; i++{
		testNoblack()
		testLock()
	}

	for{
		i := 0
		i++
	}
}
