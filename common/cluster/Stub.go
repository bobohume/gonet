package cluster

import (
	"fmt"
	"gonet/actor"
	"gonet/base"
	"gonet/common"
	"gonet/common/cluster/et"
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
		m_fsm        fsm_type //状态机
		StubMailBox  common.StubMailBox
		m_ActorName  string
		m_isRegister int32
	}
)

const STUB_TTL_TIME = et.STUB_TTL_TIME

func (this *Stub) InitStub(actorName string, stub rpc.STUB) {
	this.StubMailBox.StubType = stub
	this.StubMailBox.ClusterId = MGR.Id()
	this.m_ActorName = actorName
	go this.updateFsm()
}

func (this *Stub) IsRegister() bool {
	return atomic.LoadInt32(&this.m_isRegister) == 1
}

func (this *Stub) lease() {
	err := MGR.StubMailBox.Lease(&this.StubMailBox)
	if err != nil {
		this.m_fsm = fsm_idle
		atomic.StoreInt32(&this.m_isRegister, 0)
		if this.m_ActorName != "" {
			actor.MGR.SendMsg(rpc.RpcHead{}, fmt.Sprintf("%s.Stub_On_UnRegister", this.m_ActorName), this.StubMailBox.Id)
		}
		base.LOG.Printf("stub [%s]注销成功[%d]", this.StubMailBox.StubType.String(), this.StubMailBox.Id)
	} else {
		time.Sleep(STUB_TTL_TIME / 3)
	}
}

func (this *Stub) publish() {
	this.StubMailBox.Id = (this.StubMailBox.Id + 1) % int64(MGR.Stub.StubCount[this.StubMailBox.StubType])
	if MGR.StubMailBox.Publish(&this.StubMailBox) {
		this.m_fsm = fsm_lease
		atomic.StoreInt32(&this.m_isRegister, 1)
		if this.m_ActorName != "" {
			actor.MGR.SendMsg(rpc.RpcHead{}, fmt.Sprintf("%s.Stub_On_Register", this.m_ActorName), this.StubMailBox.Id)
		}
		base.LOG.Printf("stub [%s]注册成功[%d]", this.StubMailBox.StubType.String(), this.StubMailBox.Id)
		time.Sleep(STUB_TTL_TIME / 3)
	} else if MGR.IsEnoughStub(this.StubMailBox.StubType) {
		this.m_fsm = fsm_idle
	}
}

func (this *Stub) idle() {
	if !MGR.IsEnoughStub(this.StubMailBox.StubType) {
		this.m_fsm = fsm_publish
	}
}

func (this *Stub) updateFsm() {
	for {
		switch this.m_fsm {
		case fsm_idle:
			this.idle()
		case fsm_publish:
			this.publish()
		case fsm_lease:
			this.lease()
		}
		time.Sleep(100 * time.Millisecond)
	}
}
