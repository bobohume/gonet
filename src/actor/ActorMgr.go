package actor

import (
	"base"
	"log"
	"network"
	"strings"
)

type (
	ActorMgr struct{
		m_ActorMap 	map[string] IActor
	}

	IActorMgr interface {
		Init()
		AddActor(IActor, ...string)
		InitActorHandle(network.ISocket)
		SendMsg(string, string, ...interface{})
	}

	ActorChan struct {
		pActor IActor
		state int
	}
)

var(
	pAcotrMgr *ActorMgr
)

func (this *ActorMgr) Init() {
	this.m_ActorMap = make(map[string] IActor)
}

func (this *ActorMgr) AddActor(pActor IActor,  names ...string) {
	name := ""
	if len(names) == 0 {
		name = base.GetClassName(pActor)
		_, exist := this.m_ActorMap[name]
		if exist{
			log.Printf("Register an existed GobalActor")
			return
		}
	}else{
		name = names[0]
	}

	this.m_ActorMap[name] = pActor
}

func (this *ActorMgr) InitActorHandle(pServer network.ISocket){
	for _,v := range this.m_ActorMap{
		pServer.BindPacketFunc(v.PacketFunc)
	}
}

func (this *ActorMgr) SendMsg(name, funcName string, params  ...interface{}){
	name = strings.ToLower(name)
	pActor, exist := this.m_ActorMap[name]
	if exist{
		pActor.SendMsg(funcName, params...)
	}
}

func SendMsg(name, funcName string, params  ...interface{}){
	MGR().SendMsg(name, funcName, params...)
}

func MGR() IActorMgr{
	if pAcotrMgr == nil{
		pAcotrMgr = &ActorMgr{}
		pAcotrMgr.Init()
	}

	return pAcotrMgr
}