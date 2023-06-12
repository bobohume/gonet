package et_test

import (
	"fmt"
	"gonet/base"
	"gonet/base/cluster/et"
	"gonet/rpc"
	"sync"
	"testing"
	"time"
)

func TestMailBox(t *testing.T) {
	mailBox := et.MailBox{}
	mailBox.Init([]string{"http://127.0.0.1:2379"}, &rpc.ClusterInfo{Type: rpc.SERVICE_GATE})
	//mailBox.DeleteAll()
	a := sync.WaitGroup{}
	max_count := 100
	one_count := 1000
	a.Add(max_count)
	time1 := time.Now().UnixMilli()
	i := base.RandI(int64(0), 2000)
	fmt.Println(i)
	base.UUID.Init(i)
	for i := int64(0); i < int64(max_count); i++ {
		go func() {
			for i := int64(0); i < int64(one_count); i++ {
				id := base.UUID.UUID()
				m := &rpc.MailBox{Id: id}
				if mailBox.Create(m) {

				} else {
					fmt.Println("publish failed", id)
				}
			}
			fmt.Println("publish finish")
			a.Done()
		}()
	}

	a.Wait()
	fmt.Println("past time", time.Now().UnixMilli()-time1)
	for {
		fmt.Println("lease finish", mailBox.Len())

		time.Sleep(time.Second * 10)
	}
}
