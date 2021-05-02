package world

import (
	"github.com/golang/protobuf/proto"
	"gonet/actor"
	"gonet/base"
	"gonet/common"
	"gonet/common/cluster"
	"gonet/rpc"
	"gonet/server/message"
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
	this.InitService(&common.ClusterInfo{Type: rpc.SERVICE_WORLDSERVER, Ip:UserNetIP, Port:int32(base.Int(UserNetPort))}, EtcdEndpoints)
	this.RegisterClusterCall()

	this.Actor.Start()
}

//发送account
func SendToAccount(funcName string, params  ...interface{}){
	head := rpc.RpcHead{DestServerType:rpc.SERVICE_ACCOUNTSERVER, SendType:rpc.SEND_BALANCE, SrcClusterId:SERVER.GetClusterMgr().Id()}
	SERVER.GetClusterMgr().SendMsg(head,  funcName, params...)
}

//发送给客户端
func SendToClient(clusterId uint32, packet proto.Message){
	pakcetHead := packet.(message.Packet).GetPacketHead()
	if pakcetHead != nil {
		SERVER.GetClusterMgr().SendMsg(rpc.RpcHead{DestServerType:rpc.SERVICE_GATESERVER, ClusterId:clusterId, Id:pakcetHead.Id}, message.GetMessageName(packet), packet)
	}
}

//--------------发送给地图----------------------//
func SendToZone(Id int64, ClusterId uint32, funcName string, params  ...interface{}){
	head := rpc.RpcHead{Id:Id, ClusterId:ClusterId, DestServerType:rpc.SERVICE_ZONESERVER, SrcClusterId:SERVER.GetClusterMgr().Id()}
	SERVER.GetClusterMgr().SendMsg(head, funcName, params...)
}