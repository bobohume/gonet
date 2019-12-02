package world

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
	}
)

func (this *ClusterManager) Init(num int){
	this.Actor.Init(num)
	//注册到集群
	this.InitService(int(message.SERVICE_WORLDSERVER), UserNetIP, base.Int(UserNetPort), EtcdEndpoints)
	this.RegisterClusterCall()

	this.Actor.Start()
}