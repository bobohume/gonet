package cluster

import (
	"fmt"
	"gonet/actor"
	"gonet/base"
	"gonet/base/cluster/et"
	"gonet/rpc"
	"sync/atomic"
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
	Stub struct {
		fsm         fsm_type //状态机
		StubMailBox rpc.StubMailBox
		isRegister  int32
	}
)

const STUB_TTL_TIME = et.STUB_TTL_TIME

func (s *Stub) InitStub(stub rpc.STUB) {
	s.StubMailBox.StubType = stub
	s.StubMailBox.ClusterId = MGR.Id()
	go s.updateFsm()
}

func (s *Stub) IsRegister() bool {
	return atomic.LoadInt32(&s.isRegister) == 1
}

func (s *Stub) lease() {
	err := MGR.StubMailBox.Lease(&s.StubMailBox)
	if err != nil {
		s.fsm = fsm_idle
		atomic.StoreInt32(&s.isRegister, 0)
		actor.MGR.SendMsg(rpc.RpcHead{SendType: rpc.SEND_BOARD_CAST}, fmt.Sprintf("%s.OnStubUnRegister", s.StubMailBox.StubType.String()))
		base.LOG.Printf("stub [%s]注销成功[%d]", s.StubMailBox.StubType.String(), s.StubMailBox.Id)
	} else {
		time.Sleep(STUB_TTL_TIME / 3)
	}
}

func (s *Stub) publish() {
	s.StubMailBox.Id = (s.StubMailBox.Id + 1) % MGR.Stub.StubCount[s.StubMailBox.StubType.String()]
	if MGR.StubMailBox.Create(&s.StubMailBox) {
		s.fsm = fsm_lease
		atomic.StoreInt32(&s.isRegister, 1)
		actor.MGR.SendMsg(rpc.RpcHead{SendType: rpc.SEND_BOARD_CAST}, fmt.Sprintf("%s.OnStubRegister", s.StubMailBox.StubType.String()))
		base.LOG.Printf("stub [%s]注册成功[%d]", s.StubMailBox.StubType.String(), s.StubMailBox.Id)
		time.Sleep(STUB_TTL_TIME / 3)
	} else if MGR.IsEnoughStub(s.StubMailBox.StubType) {
		s.fsm = fsm_idle
	}
}

func (s *Stub) idle() {
	if !MGR.IsEnoughStub(s.StubMailBox.StubType) {
		s.fsm = fsm_publish
	}
}

func (s *Stub) updateFsm() {
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
