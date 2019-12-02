package account

import (
	"gonet/actor"
	"gonet/base"
	"gonet/message"
	"gonet/server/common/cluster"
)

type (
	ClusterManager struct{
		cluster.ClusterServer
	}

	IClusterManager interface {
		actor.IActor
		KickWorldPlayer(accountId int64)
	}
)

func (this *ClusterManager) Init(num int){
	this.Actor.Init(num)
	//注册account集群
	this.InitService(int(message.SERVICE_ACCOUNTSERVER), UserNetIP, base.Int(UserNetPort), EtcdEndpoints)
	this.RegisterClusterCall()

	this.Actor.Start()
}


func (this *ClusterManager) KickWorldPlayer(accountId int64){
	this.BoardCastMsg(int(message.SERVICE_WORLDSERVER), "G_ClientLost", accountId)
}