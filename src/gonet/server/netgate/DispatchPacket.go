package netgate

import (
	"github.com/golang/protobuf/proto"
	"gonet/base"
	"gonet/message"
	"gonet/rpc"
)

func SendToClient(socketId uint32, packet proto.Message){
	buff, err := proto.Marshal(packet)
	if err == nil {
		SERVER.GetServer().Send(rpc.RpcHead{SocketId:socketId}, base.SetTcpEnd(buff))
	}
}

func DispatchPacket(id uint32, buff []byte) bool{
	defer func(){
		if err := recover(); err != nil{
			base.TraceCode(err)
		}
	}()

	rpcPacket, head := rpc.UnmarshalHead(buff)
	switch head.DestServerType {
	case message.SERVICE_ACCOUNTSERVER:
		SERVER.GetAccountCluster().Send(head, base.SetTcpEnd(buff))
	case message.SERVICE_ZONESERVER:
		SERVER.GetZoneCluster().Send(head, base.SetTcpEnd(buff))
	case message.SERVICE_WORLDSERVER:
		SERVER.GetWorldCluster().Send(head, base.SetTcpEnd(buff))
	default:
		if rpcPacket.FuncName == A_C_RegisterResponse || rpcPacket.FuncName == A_C_LoginResponse {
			SERVER.GetServer().Send(rpc.RpcHead{SocketId:head.ClusterId}, base.SetTcpEnd(rpcPacket.RpcBody))
		}else{
			socketId := SERVER.GetPlayerMgr().GetSocket(head.Id)
			SERVER.GetServer().Send(rpc.RpcHead{SocketId:socketId}, base.SetTcpEnd(rpcPacket.RpcBody))
		}
	}

	return true
}
