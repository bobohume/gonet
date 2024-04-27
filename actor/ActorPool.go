package actor

import (
	"gonet/base"
	"gonet/rpc"
	"reflect"
	"sync"
)

// ********************************************************
// actorpool 管理,不能动态分配
// ********************************************************
type (
	ActorPool struct {
		MGR       IActor
		actorList []IActor
		actorSize int
	}
)

func (a *ActorPool) InitPool(pPool IActorPool, rType reflect.Type, num int) {
	a.actorList = make([]IActor, num)
	a.actorSize = num
	for i := 0; i < num; i++ {
		ac := reflect.New(rType).Interface().(IActor)
		rType := reflect.TypeOf(ac)
		op := Op{actorType: ACTOR_TYPE_POOL, name: base.GetClassName(rType)}
		ac.register(ac, op)
		ac.Init()
		a.actorList[i] = ac
	}
	a.MGR = reflect.New(rType).Interface().(IActor)
	MGR.RegisterActor(a.MGR, WithType(ACTOR_TYPE_POOL), withPool(pPool))
}

func (a *ActorPool) GetPoolSize() int {
	return a.actorSize
}

func (a *ActorPool) SendAcotr(head rpc.RpcHead, packet rpc.Packet) bool {
	if a.MGR.HasRpc(packet.RpcPacket.FuncName) {
		switch head.SendType {
		case rpc.SEND_POINT, rpc.SEND_LOCAL:
			index := base.UUIDHASH(head.Id) % int64(a.actorSize)
			a.actorList[index].Acotr().Send(head, packet)
		default:
			for i := 0; i < a.actorSize; i++ {
				a.actorList[i].Acotr().Send(head, packet)
			}
		}
		return true
	}
	return false
}

func (a *ActorPool) GetActor(Id int64) (IActor, bool) {
	index := base.UUIDHASH(Id) % int64(a.actorSize)
	return a.actorList[index], true
}

// ********************************************************
// actorpool 管理,这里的actor可以动态添加
// ********************************************************
type (
	VirtualActor struct {
		MGR       IActor
		actorMap  map[int64]IActor
		actorLock *sync.RWMutex
	}

	IVirtualActor interface {
		GetActor(Id int64) IActor //获取actor
		AddActor(ac IActor)       //添加actor
		DelActor(Id int64)        //删除actor
		GetActorNum() int
		GetMgr() IActor
	}
)

func (a *VirtualActor) InitActor(pPool IActorPool, rType reflect.Type) {
	a.actorMap = make(map[int64]IActor)
	a.actorLock = &sync.RWMutex{}
	a.MGR = reflect.New(rType).Interface().(IActor)
	MGR.RegisterActor(a.MGR, WithType(ACTOR_TYPE_VIRTUAL), withPool(pPool))
}

func (a *VirtualActor) AddActor(ac IActor) {
	rType := reflect.TypeOf(ac)
	op := Op{actorType: ACTOR_TYPE_VIRTUAL, name: base.GetClassName(rType)}
	ac.register(ac, op)
	a.actorLock.Lock()
	a.actorMap[ac.GetId()] = ac
	a.actorLock.Unlock()
}

func (a *VirtualActor) DelActor(Id int64) {
	a.actorLock.Lock()
	delete(a.actorMap, Id)
	a.actorLock.Unlock()
}

func (a *VirtualActor) GetActor(Id int64) IActor {
	a.actorLock.RLock()
	ac, bEx := a.actorMap[Id]
	a.actorLock.RUnlock()
	if bEx {
		return ac
	}
	return nil
}

func (a *VirtualActor) GeActorrNum() int {
	nLen := 0
	a.actorLock.RLock()
	nLen = len(a.actorMap)
	a.actorLock.RUnlock()
	return nLen
}

func (a *VirtualActor) GetMgr() IActor {
	return a.MGR
}

func (a *VirtualActor) SendAcotr(head rpc.RpcHead, packet rpc.Packet) bool {
	if a.MGR.HasRpc(packet.RpcPacket.FuncName) {
		if head.Id != 0 {
			ac := a.GetActor(head.Id)
			if ac != nil {
				ac.Acotr().Send(head, packet)
			}
		}
		return true
	}
	return false
}
