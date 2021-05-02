package account

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
		KickWorldPlayer(accountId int64)
	}
)

func (this *ClusterManager) Init(num int){
	this.Actor.Init(num)
	//注册account集群
	this.InitService(&common.ClusterInfo{Type: rpc.SERVICE_ACCOUNTSERVER, Ip:UserNetIP, Port:int32(base.Int(UserNetPort))}, EtcdEndpoints)
	this.RegisterClusterCall()

	this.Actor.Start()
}

func (this *ClusterManager) KickWorldPlayer(accountId int64){
	BoardCastToWorld("G_ClientLost", accountId)
}

//发送world
func SendToWorld(ClusterId uint32, funcName string, params  ...interface{}){
	head := rpc.RpcHead{ClusterId:ClusterId, DestServerType:rpc.SERVICE_WORLDSERVER, SrcClusterId:SERVER.GetClusterMgr().Id()}
	SERVER.GetClusterMgr().SendMsg(head, funcName, params...)
}

//广播world
func BoardCastToWorld(funcName string, params  ...interface{}){
	head := rpc.RpcHead{DestServerType:rpc.SERVICE_WORLDSERVER, SendType:rpc.SEND_BOARD_CAST, SrcClusterId:SERVER.GetClusterMgr().Id()}
	SERVER.GetClusterMgr().SendMsg(head, funcName, params...)
}

//发送到客户端
func SendToClient(head rpc.RpcHead, packet proto.Message){
	pakcetHead := packet.(message.Packet).GetPacketHead()
	SERVER.GetClusterMgr().SendMsg(rpc.RpcHead{DestServerType:rpc.SERVICE_GATESERVER, Id:pakcetHead.Id}, message.GetMessageName(packet), packet)
}