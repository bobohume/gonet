package actor

import (
	"gonet/base"
	"gonet/rpc"
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
		GetActor(string) IActor
		InitActorHandle(ICluster)
		SendMsg(rpc.RpcHead, string, ...interface{})
	}

	ICluster interface{
		BindPacketFunc(network.HandleFunc)
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

func (this *ActorMgr) GetActor(name string) IActor{
	name = strings.ToLower(name)
	pActor, exist := this.m_ActorMap[name]
	if exist{
		return pActor
	}
	return nil
}

func (this *ActorMgr) InitActorHandle(pCluster ICluster){
	for _,v := range this.m_ActorMap{
		pCluster.BindPacketFunc(v.PacketFunc)
	}
}

func (this *ActorMgr) SendMsg(head rpc.RpcHead, funcName string, params  ...interface{}){
	name := strings.ToLower(head.ActorName)
	pActor, exist := this.m_ActorMap[name]
	if exist{
		pActor.SendMsg(head, funcName, params...)
	}
}

var(
	MGR *ActorMgr
)

func init(){
	MGR = &ActorMgr{}
	MGR.Init()
}