package actor

import (
	"fmt"
	"gonet/base"
	"gonet/message"
	"strings"
	"sync"
)

//********************************************************
// actorpool 管理,这里的actor可以动态调节
//********************************************************
type(
	ActorPool struct{
		Actor
		m_ActorMap  map[int64] IActor
		m_ActorLock *sync.RWMutex
		m_Self IActorPool//类型c++的虚函数,由于go是组合继承
		//m_ActorMap sync.Map
	}

	IActorPool interface {
		IActor
		GetActor(Id int64) IActor//获取actor
		AddActor(Id int64, pActor IActor)//添加actor
		DelActor(Id int64)//删除actor
		SendActor(Id int64, io CallIO, funcName string) bool//发送到actor
		GetActorNum() int
	}
)

func (this *ActorPool) Init(num int){
	this.m_ActorMap = make(map[int64] IActor)
	this.m_ActorLock = &sync.RWMutex{}
	this.RegisterVirtual(this)
	this.Actor.Init(num)
}

//like c++ virtual func
func (this *ActorPool) RegisterVirtual(pActor IActorPool){
	if this.m_Self == nil{
		this.m_Self = pActor
	}
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
	this.m_ActorLock.Lock()
	nLen = len(this.m_ActorMap)
	this.m_ActorLock.Unlock()
	return nLen
}

func (this *ActorPool) SendActor(Id int64, io CallIO, funcName string) bool{
	pActor := this.GetActor(Id)
	if pActor != nil && pActor.FindCall(funcName) != nil{
		pActor.Send(io)
		return true
	}
	return false
}

func (this *ActorPool) PacketFunc(id int, buff []byte) bool{
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("PlayerMgr PacketFunc", err)
		}
	}()

	var io CallIO
	io.Buff = buff
	io.SocketId = id

	bitstream := base.NewBitStream(io.Buff, len(io.Buff))
	funcName := bitstream.ReadString()
	funcName = strings.ToLower(funcName)
	pFunc := this.FindCall(funcName)
	if pFunc != nil{
		this.Send(io)
		return true
	}else{
		bitstream.ReadInt(base.Bit8)
		nType := bitstream.ReadInt(base.Bit8)
		if (nType == base.RPC_Int64 || nType == base.RPC_UInt64 || nType == base.RPC_PInt64 || nType == base.RPC_PUInt64){
			nId := bitstream.ReadInt64(base.Bit64)
			return this.m_Self.(IActorPool).SendActor(nId, io, funcName)
		}else if (nType == base.RPC_MESSAGE){
			packet := message.GetPakcetByName(funcName)
			nLen := bitstream.ReadInt(base.Bit32)
			packetBuf := bitstream.ReadBits(nLen << 3)
			message.UnmarshalText(packet, packetBuf)
			packetHead := message.GetPakcetHead(packet)
			nId := int64(*packetHead.Id)
			return this.m_Self.(IActorPool).SendActor(nId, io, funcName)
		}
	}

	return false
}