package netgate

import (
	"github.com/golang/protobuf/proto"
	"gonet/base"
	"gonet/rpc"
	"gonet/server/message"
	"strings"
)

var(
	A_C_RegisterResponse = strings.ToLower("A_C_RegisterResponse")
	A_C_LoginResponse 	 = strings.ToLower("A_C_LoginResponse")
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
	case rpc.SERVICE_GATESERVER:
		bitstream := base.NewBitStream(rpcPacket.RpcBody, len(rpcPacket.RpcBody))
		buff := message.EncodeEx(rpcPacket.FuncName, rpc.UnmarshalPB(bitstream))
		if rpcPacket.FuncName == A_C_RegisterResponse || rpcPacket.FuncName == A_C_LoginResponse {
			SERVER.GetServer().Send(rpc.RpcHead{SocketId:head.SocketId}, base.SetTcpEnd(buff))
		}else{
			socketId := SERVER.GetPlayerMgr().GetSocket(head.Id)
			SERVER.GetServer().Send(rpc.RpcHead{SocketId:socketId}, base.SetTcpEnd(buff))
		}
	default:
		SERVER.GetCluster().Send(head, buff)
	}

	return true
}
