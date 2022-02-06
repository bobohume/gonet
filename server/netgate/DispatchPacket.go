package netgate

import (
	"bytes"
	"encoding/gob"
	"github.com/golang/protobuf/proto"
	"gonet/base"
	"gonet/rpc"
	"gonet/server/message"
	"reflect"
)

var(
	A_C_RegisterResponse = proto.MessageName(&message.A_C_RegisterResponse{})
	A_C_LoginResponse 	 = proto.MessageName(&message.A_C_LoginResponse{})
)

func SendToClient(socketId uint32, packet proto.Message){
	SERVER.GetServer().Send(rpc.RpcHead{SocketId:socketId}, rpc.Packet{Buff: message.Encode(packet)})
}

func DispatchPacket(packet rpc.Packet) bool{
	defer func(){
		if err := recover(); err != nil{
			base.TraceCode(err)
		}
	}()

	rpcPacket, head := rpc.Unmarshal(packet.Buff)
	switch head.DestServerType {
	case rpc.SERVICE_GATESERVER:
		messageName := ""
		buf := bytes.NewBuffer(rpcPacket.RpcBody)
		dec := gob.NewDecoder(buf)
		dec.Decode(&messageName)
		packet := reflect.New(proto.MessageType(messageName).Elem()).Interface().(proto.Message)
		dec.Decode(packet)
		buff := message.Encode(packet)
		if messageName== A_C_RegisterResponse || messageName == A_C_LoginResponse {
			SERVER.GetServer().Send(rpc.RpcHead{SocketId:head.SocketId}, rpc.Packet{Buff:buff})
		}else{
			socketId := SERVER.GetPlayerMgr().GetSocket(head.Id)
			SERVER.GetServer().Send(rpc.RpcHead{SocketId:socketId}, rpc.Packet{Buff:buff})
		}
	default:
		SERVER.GetCluster().Send(head, packet)
	}

	return true
}
