package actor

import (
	"sync"
)

//********************************************************
// actorpool 管理,这里的actor可以动态调节
//********************************************************
type(
	ActorPool struct{
		m_ActorMap  map[int64] IActor
		m_ActorLock *sync.RWMutex
		//m_ActorMap sync.Map
	}

	IActorPool interface {
		GetActor(Id int64) IActor//获取actor
		AddActor(Id int64, pActor IActor)//添加actor
		DelActor(Id int64)//删除actor
		BoardCast(funcName string, params ...interface{})//广播actor
		Send(Id int64, funcName string, io CallIO) bool//发送到actor
		SendMsg(Id int64, funcName string, params  ...interface{}) bool//发送到actor
		GetActorNum() int
	}
)

func (this *ActorPool) Init(){
	this.m_ActorMap = make(map[int64] IActor)
	this.m_ActorLock = &sync.RWMutex{}
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

func (this *ActorPool) Send(Id int64, funcName string, io CallIO) bool{
	pActor := this.GetActor(Id)
	if pActor != nil && pActor.FindCall(funcName) != nil{
		pActor.Send(io)
		return true
	}
	return false
}

func (this *ActorPool) SendMsg(Id int64, funcName string, params  ...interface{}) bool{
	pActor := this.GetActor(Id)
	if pActor != nil && pActor.FindCall(funcName) != nil{
		pActor.SendMsg(funcName, params...)
		return true
	}
	return false
}

//actor pool must rewrite PacketFunc
/*func (this *ActorPool) PacketFunc(id int, buff []byte) bool{
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
			nId := packetHead.Id
			return this.m_Self.(IActorPool).SendActor(nId, io, funcName)
		}
	}

	return false
}*/