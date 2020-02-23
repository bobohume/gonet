package actor

import (
	"gonet/base"
	"gonet/rpc"
	"sync"
)

//********************************************************
// actorpool 管理,这里的actor可以动态添加,携程池
//********************************************************
type(
	ActorPool struct{
		Actor
		m_ActorMap  map[int64] IActor
		m_ActorLock *sync.RWMutex
		//m_ActorMap sync.Map
	}

	IActorPool interface {
		GetActor(Id int64) IActor//获取actor
		AddActor(Id int64, pActor IActor)//添加actor
		DelActor(Id int64)//删除actor
		BoardCast(funcName string, params ...interface{})//广播actor
		GetActorNum() int
	}
)

func (this *ActorPool) Init(chanNum int){
	this.m_ActorMap = make(map[int64] IActor)
	this.m_ActorLock = &sync.RWMutex{}
	this.Actor.Init(chanNum)
}

func (this *ActorPool) GetActor(Id int64) IActor{
	//v, bOk := this.m_ActorMap.Load(Id)
	this.m_ActorLock.RLock()
	pActor, bEx := this.m_ActorMap[Id]
	this.m_ActorLock.RUnlock()
	if bEx{
		return pActor
	}
	return nil
}

func (this *ActorPool) AddActor(Id int64, pActor IActor){
	this.m_ActorLock.Lock()
	this.m_ActorMap[Id] = pActor
	this.m_ActorLock.Unlock()
	//this.m_ActorMap.Store(Id, pActor)
}

func (this *ActorPool) DelActor(Id int64){
	this.m_ActorLock.Lock()
	delete(this.m_ActorMap, Id)
	this.m_ActorLock.Unlock()
	//this.m_ActorMap.Delete(Id)
}

func (this *ActorPool) GetActorNum() int{
	nLen := 0
	this.m_ActorLock.RLock()
	nLen = len(this.m_ActorMap)
	this.m_ActorLock.RUnlock()
	return nLen
}

func (this *ActorPool) BoardCast(funcName string, params ...interface{}){
	this.m_ActorLock.RLock()
	for _, v := range this.m_ActorMap{
		v.SendMsg(funcName, params...)
	}
	this.m_ActorLock.RUnlock()
}

func (this *ActorPool) SendMsg(funcName string, params ...interface{}) {
	var io CallIO
	io.ActorId = this.m_Id
	io.SocketId = 0
	buff, rpcPacket := rpc.MarshalEx(funcName, params...)
	io.Buff = buff
	if rpcPacket.RpcHead.Id != 0{
		pActor := this.GetActor(rpcPacket.RpcHead.Id)
		if pActor != nil && pActor.FindCall(funcName) != nil{
			pActor.Send(io)
			return
		}
	}
	this.Send(io)
}

//actor pool must rewrite PacketFunc
func (this *ActorPool) PacketFunc(id int, buff []byte) bool{
	defer func() {
		if err := recover(); err != nil {
			base.TraceCode(err)
		}
	}()

	var io CallIO
	io.Buff = buff
	io.SocketId = id

	rpcPacket := rpc.UnmarshalHead(io.Buff)
	if this.FindCall(rpcPacket.FuncName) != nil{
		this.Send(io)
		return true
	}else{
		pActor := this.GetActor(rpcPacket.RpcHead.Id)
		if pActor != nil && pActor.FindCall(rpcPacket.FuncName) != nil{
			pActor.Send(io)
			return true
		}
	}
	return false
}