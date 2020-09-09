package world

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
	}
)

func (this *ClusterManager) Init(num int){
	this.Actor.Init(num)
	//注册到集群
	this.InitService(message.SERVICE_WORLDSERVER, UserNetIP, base.Int(UserNetPort), EtcdEndpoints)
	this.RegisterClusterCall()

	this.Actor.Start()
}

//发送account
func SendToAccount(funcName string, params  ...interface{}){
	head := rpc.RpcHead{DestServerType:message.SERVICE_ACCOUNTSERVER, SendType:message.SEND_BALANCE, SrcClusterId:SERVER.GetClusterMgr().Id()}
	SERVER.GetClusterMgr().SendMsg(head,  funcName, params...)
}

//发送给客户端
func SendToClient(clusterId uint32, packet proto.Message){
	buff := message.Encode(packet)
	pakcetHead := packet.(message.Packet).GetPacketHead()
	if pakcetHead != nil {
		rpcPacket := &message.RpcPacket{FuncName:message.GetMessageName(packet), ArgLen:1, RpcHead:&message.RpcHead{Id:pakcetHead.Id}, RpcBody:buff}
		data, _ := proto.Marshal(rpcPacket)
		SERVER.GetClusterMgr().Send(rpc.RpcHead{DestServerType:message.SERVICE_GATESERVER, ClusterId:clusterId}, base.SetTcpEnd(data))
	}
}

func SendToClientBySocketId(socketId uint32, packet proto.Message) {
	buff := message.Encode(packet)
	pakcetHead := packet.(message.Packet).GetPacketHead()
	if pakcetHead != nil {
		rpcPacket := &message.RpcPacket{FuncName: message.GetMessageName(packet), ArgLen: 1, RpcHead: &message.RpcHead{Id: pakcetHead.Id}, RpcBody: buff}
		data, _ := proto.Marshal(rpcPacket)
		SERVER.GetClusterMgr().Send(rpc.RpcHead{DestServerType:message.SERVICE_GATESERVER, SocketId:socketId}, base.SetTcpEnd(data))
	}
}

//--------------发送给地图----------------------//
func SendToZone(Id int64, ClusterId uint32, funcName string, params  ...interface{}){
	head := rpc.RpcHead{Id:Id, ClusterId:ClusterId, DestServerType:message.SERVICE_ZONESERVER, SrcClusterId:SERVER.GetClusterMgr().Id()}
	SERVER.GetClusterMgr().SendMsg(head, funcName, params...)
}