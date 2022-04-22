package actor

import (
	"gonet/rpc"
	"reflect"
	"sync"
)

//********************************************************
// actorpool 管理,不能动态分配
//********************************************************
type (
	ActorPool struct {
		m_MGR       IActor
		m_ActorList []IActor
		m_ActorSize int
	}
)

func (this *ActorPool) InitPool(pPool IActorPool, rType reflect.Type, num int) {
	this.m_ActorList = make([]IActor, num)
	this.m_ActorSize = num
	for i := 0; i < num; i++ {
		pActor := reflect.New(rType).Interface().(IActor)
		rType := reflect.TypeOf(pActor)
		op := Op{m_type: ACTOR_TYPE_POOL, m_name: rType.Name()}
		pActor.register(pActor, op)
		pActor.Init()
		this.m_ActorList[i] = pActor
	}
	this.m_MGR = reflect.New(rType).Interface().(IActor)
	MGR.RegisterActor(this.m_MGR, WithType(ACTOR_TYPE_POOL), withPool(pPool))
}

func (this *ActorPool) GetPoolSize() int {
	return this.m_ActorSize
}

func (this *ActorPool) SendAcotr(head rpc.RpcHead, packet rpc.Packet) bool {
	if this.m_MGR.HasRpc(packet.RpcPacket.FuncName) {
		switch head.SendType {
		case rpc.SEND_POINT:
			index := head.Id % int64(this.m_ActorSize)
			this.m_ActorList[index].Acotr().Send(head, packet)
		default:
			for i := 0; i < this.m_ActorSize; i++ {
				this.m_ActorList[i].Acotr().Send(head, packet)
			}
		}
		return true
	}
	return false
}

//********************************************************
// actorpool 管理,这里的actor可以动态添加
//********************************************************
type (
	VirtualActor struct {
		m_MGR       IActor
		m_ActorMap  map[int64]IActor
		m_ActorLock *sync.RWMutex
	}

	IVirtualActor interface {
		GetActor(Id int64) IActor //获取actor
		AddActor(pActor IActor)   //添加actor
		DelActor(Id int64)        //删除actor
		GetActorNum() int
		GetMgr() IActor
	}
)

func (this *VirtualActor) InitActor(pPool IActorPool, rType reflect.Type) {
	this.m_ActorMap = make(map[int64]IActor)
	this.m_ActorLock = &sync.RWMutex{}
	this.m_MGR = reflect.New(rType).Interface().(IActor)
	MGR.RegisterActor(this.m_MGR, WithType(ACTOR_TYPE_VIRTUAL), withPool(pPool))
}

func (this *VirtualActor) AddActor(pActor IActor) {
	rType := reflect.TypeOf(pActor)
	op := Op{m_type: ACTOR_TYPE_VIRTUAL, m_name: rType.Name()}
	pActor.register(pActor, op)
	this.m_ActorLock.Lock()
	this.m_ActorMap[pActor.GetId()] = pActor
	this.m_ActorLock.Unlock()
}

func (this *VirtualActor) DelActor(Id int64) {
	this.m_ActorLock.Lock()
	delete(this.m_ActorMap, Id)
	this.m_ActorLock.Unlock()
}

func (this *VirtualActor) GetActor(Id int64) IActor {
	this.m_ActorLock.RLock()
	pActor, bEx := this.m_ActorMap[Id]
	this.m_ActorLock.RUnlock()
	if bEx {
		return pActor
	}
	return nil
}

func (this *VirtualActor) GeActorrNum() int {
	nLen := 0
	this.m_ActorLock.RLock()
	nLen = len(this.m_ActorMap)
	this.m_ActorLock.RUnlock()
	return nLen
}

func (this *VirtualActor) GetMgr() IActor {
	return this.m_MGR
}

func (this *VirtualActor) SendAcotr(head rpc.RpcHead, packet rpc.Packet) bool {
	if this.m_MGR.HasRpc(packet.RpcPacket.FuncName) {
		if head.Id != 0 {
			pActor := this.GetActor(head.Id)
			if pActor != nil {
				pActor.Acotr().Send(head, packet)
			}
		}
		return true
	}
	return false
}
