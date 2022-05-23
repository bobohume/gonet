package gate

import (
	"bytes"
	"encoding/gob"
	"gonet/base"
	"gonet/rpc"
	"gonet/server/message"
	"reflect"

	"github.com/golang/protobuf/proto"
)

var (
	LoginAccountResponse = proto.MessageName(&message.LoginAccountResponse{})
	SelectPlayerResponse = proto.MessageName(&message.SelectPlayerResponse{})
)

func SendToClient(socketId uint32, packet proto.Message) {
	SERVER.GetServer().Send(rpc.RpcHead{SocketId: socketId}, rpc.Packet{Buff: message.Encode(packet)})
}

func DispatchPacket(packet rpc.Packet) bool {
	defer func() {
		if err := recover(); err != nil {
			base.TraceCode(err)
		}
	}()

	rpcPacket, head := rpc.Unmarshal(packet.Buff)
	switch head.DestServerType {
	case rpc.SERVICE_GATE:
		messageName := ""
		buf := bytes.NewBuffer(rpcPacket.RpcBody)
		dec := gob.NewDecoder(buf)
		dec.Decode(&messageName)
		packet := reflect.New(proto.MessageType(messageName).Elem()).Interface().(proto.Message)
		dec.Decode(packet)
		buff := message.Encode(packet)
		switch messageName {
		case LoginAccountResponse, SelectPlayerResponse:
			SERVER.GetServer().Send(rpc.RpcHead{SocketId: head.SocketId}, rpc.Packet{Buff: buff})
		default:
			socketId := SERVER.GetPlayerMgr().GetSocket(head.Id)
			SERVER.GetServer().Send(rpc.RpcHead{SocketId: socketId}, rpc.Packet{Buff: buff})
		}

	default:
		SERVER.GetCluster().Send(head, packet)
	}

	return true
}
