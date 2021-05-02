package netgate

import (
	"github.com/golang/protobuf/proto"
	"gonet/base"
	"gonet/rpc"
	"gonet/server/message"
)

func SendToClient(socketId uint32, packet proto.Message){
	SERVER.GetServer().Send(rpc.RpcHead{SocketId:socketId}, base.SetTcpEnd(message.Encode(packet)))
}

func DispatchPacket(id uint32, buff []byte) bool{
	defer func(){
		if err := recover(); err != nil{
			base.TraceCode(err)
		}
	}()

	rpcPacket, head := rpc.Unmarshal(buff)
	switch head.DestServerType {
	case rpc.SERVICE_ACCOUNTSERVER:
		SERVER.GetAccountCluster().Send(head, base.SetTcpEnd(buff))
	case rpc.SERVICE_ZONESERVER:
		SERVER.GetZoneCluster().Send(head, base.SetTcpEnd(buff))
	case rpc.SERVICE_WORLDSERVER:
		SERVER.GetWorldCluster().Send(head, base.SetTcpEnd(buff))
	default:
		bitstream := base.NewBitStream(rpcPacket.RpcBody, len(rpcPacket.RpcBody))
		buff := message.EncodeEx(rpcPacket.FuncName, rpc.UnmarshalPB(bitstream))
		if rpcPacket.FuncName == A_C_RegisterResponse || rpcPacket.FuncName == A_C_LoginResponse {
			SERVER.GetServer().Send(rpc.RpcHead{SocketId:head.SocketId}, base.SetTcpEnd(buff))
		}else{
			socketId := SERVER.GetPlayerMgr().GetSocket(head.Id)
			SERVER.GetServer().Send(rpc.RpcHead{SocketId:socketId}, base.SetTcpEnd(buff))
		}
	}

	return true
}
