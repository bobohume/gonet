package rpc

import (
	"bytes"
	"encoding/gob"
	"github.com/golang/protobuf/proto"
	"gonet/base"
)

//rpc  Marshal
func Marshal(head* RpcHead, funcName *string, params ...interface{}) Packet {
	return marshal(head, funcName, params...)
}

//rpc  marshal
func marshal(head *RpcHead, funcName *string, params ...interface{}) Packet {
	defer func() {
		if err := recover(); err != nil {
			base.TraceCode(err)
		}
	}()

	*funcName = Route(head, *funcName)
	rpcPacket := &RpcPacket{FuncName: *funcName, ArgLen: int32(len(params)), RpcHead: (*RpcHead)(head)}
	buf := bytes.NewBuffer([]byte{})
	enc := gob.NewEncoder(buf)
	for _, param := range params {
		enc.Encode(param)
	}
	rpcPacket.RpcBody = buf.Bytes()
	dat, _ := proto.Marshal(rpcPacket)
	return Packet{Buff: dat, RpcPacket: rpcPacket}
}

//rpc  MarshalPB
func marshalPB(bitstream *base.BitStream, packet proto.Message) {
	bitstream.WriteString(proto.MessageName(packet))
	buf, _ := proto.Marshal(packet)
	nLen := len(buf)
	bitstream.WriteInt(nLen, 32)
	bitstream.WriteBits(buf, nLen<<3)
}
