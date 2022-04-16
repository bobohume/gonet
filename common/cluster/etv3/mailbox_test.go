package etv3_test

import (
	"fmt"
	"gonet/common"
	"gonet/common/cluster/etv3"
	"gonet/rpc"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestMailBox(t *testing.T){
	mailBox := etv3.MailBox{}
	mailBox.Init([]string{"http://127.0.0.1:2379"}, &common.ClusterInfo{Type: rpc.SERVICE_GATE})
	locker := sync.Mutex{}
	mailBox.DeleteAll()
	mailBoxMap := map[int64] *rpc.MailBox{}
	a := sync.WaitGroup{}
	max_count  := 300
	one_count := 1000
	a.Add(max_count)
	index := int64(0)
	for i := int64(0); i < int64(max_count); i++ {
		go func() {
			for i := int64(0); i < int64(one_count); i++ {
				id := atomic.AddInt64(&index, 1)
				m := &rpc.MailBox{Id: id}
				if mailBox.Publish(m) {
					locker.Lock()
					mailBoxMap[i] = m
					locker.Unlock()
				} else {
					fmt.Println("publish failed", id)
				}
			}
			fmt.Println("publish finish")
			a.Done()
		}()
	}

	a.Wait()
	for{
		for _, v := range mailBoxMap{
			error := mailBox.Lease(v.LeaseId)
			if error != nil{
				fmt.Println(error, v.Id)
			}
		}

		fmt.Println("lease finish")

		time.Sleep(time.Second * 30)
	}
}