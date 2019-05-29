package actor

import (
	"gonet/base"
	"log"
	"gonet/network"
	"strings"
)

//一些全局的actor,不可删除的,不用锁考虑性能
//不是全局的actor,请使用actor pool
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
	MGR.SendMsg(name, funcName, params...)
}

var(
	MGR *ActorMgr
)

func init(){
	MGR = &ActorMgr{}
	MGR.Init()
}