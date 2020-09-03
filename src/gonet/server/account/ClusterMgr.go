package account

import (
	"github.com/golang/protobuf/proto"
	"gonet/actor"
	"gonet/base"
	"gonet/message"
	"gonet/rpc"
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
	this.InitService(message.SERVICE_ACCOUNTSERVER, UserNetIP, base.Int(UserNetPort), EtcdEndpoints)
	this.RegisterClusterCall()

	this.Actor.Start()
}

func (this *ClusterManager) KickWorldPlayer(accountId int64){
	BoardCastToWorld("G_ClientLost", accountId)
}

//发送world
func SendToWorld(ClusterId uint32, funcName string, params  ...interface{}){
	head := rpc.RpcHead{ClusterId:ClusterId, DestServerType:message.SERVICE_WORLDSERVER, SrcClusterId:SERVER.GetClusterMgr().Id()}
	SERVER.GetClusterMgr().SendMsg(head, funcName, params...)
}

//广播world
func BoardCastToWorld(funcName string, params  ...interface{}){
	head := rpc.RpcHead{DestServerType:message.SERVICE_WORLDSERVER, SendType:message.SEND_BOARD_CAST, SrcClusterId:SERVER.GetClusterMgr().Id()}
	SERVER.GetClusterMgr().SendMsg(head, funcName, params...)
}

//发送到客户端
func SendToClient(head rpc.RpcHead, packet proto.Message){
	buff := message.Encode(packet)
	pakcetHead := packet.(message.Packet).GetPacketHead()
	head.Id = pakcetHead.Id
	head.DestServerType = message.SERVICE_GATESERVER
	rpcPacket := &message.RpcPacket{FuncName:message.GetMessageName(packet), ArgLen:1, RpcHead:(*message.RpcHead)(&head), RpcBody:buff}
	data, _ := proto.Marshal(rpcPacket)
	SERVER.GetServer().Send(head, base.SetTcpEnd(data))
}