package cluster_test

import (
	"gonet/base"
	"gonet/base/cluster"
	"gonet/rpc"
	"sync/atomic"
	"testing"
	"time"
)

type fsm_type uint32

const (
	fsm_idle    fsm_type = iota //空闲
	fsm_publish fsm_type = iota //注册
	fsm_lease   fsm_type = iota //ttl
) //fsm_type

type (
	//stub容器
	stub struct {
		fsm         fsm_type //状态机
		StubMailBox rpc.StubMailBox
		isRegister  int32
	}
)

var (
	g_stubCount [8]int64
	max_count   = int64(10)
)

const stub_ttl_time = 30 * time.Second

func addStubCount() bool {
	count := atomic.LoadInt64(&g_stubCount[0])
	if count != max_count {
		return atomic.CompareAndSwapInt64(&g_stubCount[0], count, count+1)
	}
	return false
}

func subStubCount() {
	atomic.AddInt64(&g_stubCount[0], -1)
}

func isMaxStubCount() bool {
	return atomic.LoadInt64(&g_stubCount[0]) == max_count
}

func isEnoughStub(rpc.STUB) bool {
	return false
}

func (s *stub) InitStub(stub rpc.STUB) {
	s.StubMailBox.StubType = stub
	go s.updateFsm()
}

func (s *stub) IsRegister() bool {
	return atomic.LoadInt32(&s.isRegister) == 1
}

func (s *stub) lease() {
	err := cluster.MGR.StubMailBox.Lease(&s.StubMailBox)
	if err != nil {
		s.fsm = fsm_idle
		atomic.StoreInt32(&s.isRegister, 0)
		subStubCount()
		base.LOG.Printf("stub [%s]注销成功[%d]", s.StubMailBox.StubType.String(), s.StubMailBox.Id)
	} else {
		time.Sleep(stub_ttl_time / 3)
	}
}

func (s *stub) publish() {
	if addStubCount() {
		//s.StubMailBox.Id = (s.StubMailBox.Id + 1) % max_count
		if cluster.MGR.StubMailBox.Create(&s.StubMailBox) {
			s.fsm = fsm_lease
			atomic.StoreInt32(&s.isRegister, 1)
			base.LOG.Printf("stub [%s]注册成功[%d]", s.StubMailBox.StubType.String(), s.StubMailBox.Id)
			time.Sleep(stub_ttl_time / 3)
		} else if isMaxStubCount() || isEnoughStub(s.StubMailBox.StubType) {
			s.fsm = fsm_idle
			subStubCount()
		} else {
			subStubCount()
		}
	} else if isMaxStubCount() || isEnoughStub(s.StubMailBox.StubType) {
		s.fsm = fsm_idle
	}
}

func (s *stub) idle() {
	if !isMaxStubCount() && !isEnoughStub(s.StubMailBox.StubType) {
		s.fsm = fsm_publish
	}
}

func (s *stub) updateFsm() {
	for {
		switch s.fsm {
		case fsm_idle:
			s.idle()
		case fsm_publish:
			s.publish()
		case fsm_lease:
			s.lease()
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func TestMailBox(t *testing.T) {
	cluster.MGR.StubMailBox.Init([]string{"http://127.0.0.1:2379"}, &rpc.ClusterInfo{})
	for i := 0; i < 10000; i++ {
		stub := stub{}
		stub.StubMailBox.Id = int64(i)
		stub.InitStub(rpc.STUB_AccountMgr)
	}
	for {
		time.Sleep(time.Millisecond * 100)
	}
}
