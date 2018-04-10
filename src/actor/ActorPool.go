package actor

import (
	"log"
	"strings"
	"network"
)

type (
	ActorList struct{
		m_ActorMap map[string] IActor
	}

	IActorList interface {
		RegisterGActorList(string, IActor)
		SendGActor(string, int, string, ...interface{})
		GetGActor(string) *IActor
		InitGActorListHandle(network.ISocket)
	}
)

var(
	g_pActorList *ActorList
	g_bActorListInit bool
)

func (this *ActorList) RegisterGActorList(name string, pActor IActor){
	name = strings.ToLower(name)
	_, exist := this.m_ActorMap[name]
	if exist{
		log.Printf("Register an existed GobalActor")
		return
	}

	this.m_ActorMap[name] = pActor
}

func (this *ActorList) SendGActor(name string, sokcetId int, funcName string, params ...interface{}){
	name = strings.ToLower(name)
	pActor := this.GetGActor(name)
	if pActor != nil{
		pActor.SendMsg(sokcetId, funcName, params...)
	}
}

func (this *ActorList) GetGActor(name string) IActor{
	name = strings.ToLower(name)
	pActor, exist := this.m_ActorMap[name]
	if exist{
		return pActor
	}

	return nil
}

func (this *ActorList) InitGActorListHandle(pServer network.ISocket){
	for _,v := range this.m_ActorMap{
		pServer.BindPacketFunc(v.PacketFunc)
	}
}

func GetGActorList() *ActorList{
	if !g_bActorListInit{
		g_pActorList = &ActorList{}
		g_pActorList.m_ActorMap = make(map[string] IActor)
		g_bActorListInit = true
	}

	return g_pActorList
}