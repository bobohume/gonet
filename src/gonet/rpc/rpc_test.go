package rpc_test

import (
	"github.com/golang/protobuf/proto"
	"gonet/base"
	"gonet/message"
	"gonet/rpc"
	"strings"
	"testing"
)

type(
	TopRank struct{
		Value []int `sql:"name:value"`
	}
)

var(
	ntimes = 100000
	nArraySize = 2000
	nValue = 0x7fffffff
)

func TestMarshalPB(t *testing.T){
	for i := 0; i < ntimes; i++{
		proto.Marshal(&message.W_C_Test1{Recv:uint32(nValue)})
	}
}

func TestUMarshalPB(t *testing.T){
	for i := 0; i < ntimes; i++{
		buff, _ := proto.Marshal(&message.W_C_Test1{Recv:uint32(nValue)})
		proto.Unmarshal(buff, &message.W_C_Test{})
	}
}

func TestMarshalRpc(t *testing.T){
	for i := 0; i < ntimes; i++{
		rpc.Marshal("", int32(1), int32(1), int32(1), int32(1), int32(1), int32(1))
	}
}

func TestUMarshalRpc(t *testing.T){
	for i := 0; i < ntimes; i++{
		buff := rpc.Marshal("", int32(1), int32(1), int32(1), int32(1), int32(1), int32(1))
		parse(buff)
	}
}

func TestMarshalRpcStream(t *testing.T){
	for i := 0; i < ntimes; i++{
		rpc.MarshalSteam("", uint32(1), uint32(1), uint32(1), uint32(1), uint32(1), uint32(1))
	}
}

func TestUMarshalRpcStream(t *testing.T){
	for i := 0; i < ntimes; i++{
		buff := rpc.MarshalSteam("", uint32(1), uint32(1), uint32(1), uint32(1), uint32(1), uint32(1))
		parseStream(buff)
	}
}

func parseStream (buff []byte) {
	funcName := ""
	bitstream := base.NewBitStream(buff, len(buff))
	funcName = bitstream.ReadString()
	funcName = strings.ToLower(funcName)
	rpc.UnmarshalStream(bitstream, funcName, nil)
}

func parse (buff []byte) {
	funcName := ""
	bitstream := base.NewBitStream(buff, len(buff))
	funcName = bitstream.ReadString()
	funcName = strings.ToLower(funcName)
	rpc.Unmarshal(bitstream, funcName, nil)
}
