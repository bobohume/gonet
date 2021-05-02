package netgate

import (
	"context"
	"gonet/actor"
	"gonet/base"
	"gonet/common"
	"gonet/rpc"
	"strings"
)

var(
	A_C_RegisterResponse = strings.ToLower("A_C_RegisterResponse")
	A_C_LoginResponse 	 = strings.ToLower("A_C_LoginResponse")
)

type (
	AccountProcess struct {
		actor.Actor
		m_LostTimer *common.SimpleTimer

		m_ClusterId uint32
	}

	IAccountProcess interface {
		actor.IActor

		SetClusterId(uint32)
	}
)

func (this * AccountProcess) SetClusterId(clusterId uint32){
	this.m_ClusterId = clusterId
}

func (this *AccountProcess) Init(num int) {
	this.Actor.Init(num)
	this.m_LostTimer = common.NewSimpleTimer(3)
	this.m_LostTimer.Start()
	this.RegisterTimer(1 * 1000 * 1000 * 1000, this.Update)
	this.RegisterCall("COMMON_RegisterRequest", func(ctx context.Context) {
		SERVER.GetAccountCluster().SendMsg(rpc.RpcHead{ClusterId:this.m_ClusterId},"COMMON_RegisterRequest", &common.ClusterInfo{Type: rpc.SERVICE_GATESERVER, Ip:UserNetIP, Port:int32(base.Int(UserNetPort))})
	})

	this.RegisterCall("COMMON_RegisterResponse", func(ctx context.Context) {
		this.m_LostTimer.Stop()
	})

	this.RegisterCall("STOP_ACTOR", func(ctx context.Context) {
		this.Stop()
	})

	this.RegisterCall("DISCONNECT", func(ctx context.Context, socketId uint32) {
		this.m_LostTimer.Start()
		SERVER.GetAccountCluster().Actor.SendMsg(rpc.RpcHead{},"DISCONNECT", this.m_ClusterId)
	})

	this.RegisterCall("A_G_Account_Login", func(ctx context.Context, accountId int64, socketId uint32) {
		SERVER.GetPlayerMgr().SendMsg(rpc.RpcHead{},"ADD_ACCOUNT", accountId, socketId)
	})

	this.Actor.Start()
}

func (this* AccountProcess) Update(){
	if this.m_LostTimer.CheckTimer(){
		SERVER.GetAccountCluster().GetCluster(rpc.RpcHead{ClusterId:this.m_ClusterId}).Start()
	}
}