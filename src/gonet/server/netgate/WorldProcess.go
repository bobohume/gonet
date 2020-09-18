package netgate

import (
	"context"
	"gonet/actor"
	"gonet/base"
	"gonet/message"
	"gonet/rpc"
	"gonet/server/common"
)

type (
	WorldProcess struct {
		actor.Actor
		m_LostTimer *common.SimpleTimer

		m_ClusterId uint32
	}

	IWorldlProcess interface {
		actor.IActor

		SetClusterId(int)
	}
)

func (this * WorldProcess) SetClusterId(clusterId uint32){
	this.m_ClusterId = clusterId
}

func (this *WorldProcess) Init(num int) {
	this.Actor.Init(num)
	this.m_LostTimer = common.NewSimpleTimer(3)
	this.m_LostTimer.Start()
	this.m_ClusterId = 0
	this.RegisterTimer(1 * 1000 * 1000 * 1000, this.Update)
	this.RegisterCall("COMMON_RegisterRequest", func(ctx context.Context) {
		SERVER.GetWorldCluster().SendMsg(rpc.RpcHead{ClusterId:this.m_ClusterId},"COMMON_RegisterRequest", &message.ClusterInfo{Type:message.SERVICE_GATESERVER, Ip:UserNetIP, Port:int32(base.Int(UserNetPort))})
	})

	this.RegisterCall("COMMON_RegisterResponse", func(ctx context.Context) {
		//收到worldserver对自己注册的反馈
		this.m_LostTimer.Stop()
		SERVER.GetLog().Println("收到world对自己注册的反馈")
	})

	this.RegisterCall("STOP_ACTOR", func(ctx context.Context) {
		this.Stop()
	})

	this.RegisterCall("DISCONNECT", func(ctx context.Context, socketId uint32) {
		this.m_LostTimer.Start()
		SERVER.GetWorldCluster().Actor.SendMsg(rpc.RpcHead{},"DISCONNECT", this.m_ClusterId)
	})

	this.Actor.Start()
}

func (this* WorldProcess) Update(){
	if this.m_LostTimer.CheckTimer(){
		SERVER.GetWorldCluster().GetCluster(rpc.RpcHead{ClusterId:this.m_ClusterId}).Start()
	}
}