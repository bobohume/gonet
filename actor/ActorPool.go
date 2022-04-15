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
		op := Op{m_type: ACTOR_TYPE_PLAYER, m_name: rType.Name()}
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
			this.m_ActorList[index].GetAcotr().Send(head, packet)
		default:
			for i := 0; i < this.m_ActorSize; i++ {
				this.m_ActorList[i].GetAcotr().Send(head, packet)
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
	ActorPlayer struct {
		m_MGR        IActor
		m_PlayerMap  map[int64]IActor
		m_PlayerLock *sync.RWMutex
	}

	IActorPlayer interface {
		GetPlayer(Id int64) IActor //获取actor
		AddPlayer(pActor IActor)   //添加actor
		DelPlayer(Id int64)        //删除actor
		GetPlayerNum() int
		GetMgr() IActor
	}
)

func (this *ActorPlayer) InitPlayer(pPool IActorPool, rType reflect.Type) {
	this.m_PlayerMap = make(map[int64]IActor)
	this.m_PlayerLock = &sync.RWMutex{}
	this.m_MGR = reflect.New(rType).Interface().(IActor)
	MGR.RegisterActor(this.m_MGR, WithType(ACTOR_TYPE_PLAYER), withPool(pPool))
}

func (this *ActorPlayer) AddPlayer(pActor IActor) {
	rType := reflect.TypeOf(pActor)
	op := Op{m_type: ACTOR_TYPE_PLAYER, m_name: rType.Name()}
	pActor.register(pActor, op)
	this.m_PlayerLock.Lock()
	this.m_PlayerMap[pActor.GetId()] = pActor
	this.m_PlayerLock.Unlock()
}

func (this *ActorPlayer) DelPlayer(Id int64) {
	this.m_PlayerLock.Lock()
	delete(this.m_PlayerMap, Id)
	this.m_PlayerLock.Unlock()
}

func (this *ActorPlayer) GetPlayer(Id int64) IActor {
	this.m_PlayerLock.RLock()
	pActor, bEx := this.m_PlayerMap[Id]
	this.m_PlayerLock.RUnlock()
	if bEx {
		return pActor
	}
	return nil
}

func (this *ActorPlayer) GetPlayerNum() int {
	nLen := 0
	this.m_PlayerLock.RLock()
	nLen = len(this.m_PlayerMap)
	this.m_PlayerLock.RUnlock()
	return nLen
}

func (this *ActorPlayer) GetMgr() IActor {
	return this.m_MGR
}

func (this *ActorPlayer) SendAcotr(head rpc.RpcHead, packet rpc.Packet) bool {
	if this.m_MGR.HasRpc(packet.RpcPacket.FuncName) {
		if head.Id != 0 {
			pActor := this.GetPlayer(head.Id)
			if pActor != nil {
				pActor.GetAcotr().Send(head, packet)
				return true
			}
		}
	}
	return false
}
