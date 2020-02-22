package rpc_test

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"github.com/golang/protobuf/proto"
	"github.com/json-iterator/go"
	"gonet/message"
	"gonet/rpc"
	"testing"
)

type(
	TopRank struct{
		Value []int `sql:"name:value"`
	}
)

var(
	ntimes = 1000
	nArraySize = 2000
	nValue = 0x7fffffff
)

func TestMarshalJson(t *testing.T){
	data := &TopRank{}
	for i := 0; i < nArraySize; i++{
		data.Value = append(data.Value, nValue)
	}
	for i := 0; i < ntimes; i++{
		json.Marshal(data)
	}
}

func TestUMarshalJson(t *testing.T){
	data := &TopRank{}
	for i := 0; i < nArraySize; i++{
		data.Value = append(data.Value, nValue)
	}

	for i := 0; i < ntimes; i++{
		buff, _ := json.Marshal(data)
		json.Unmarshal(buff, &TopRank{})
	}
}

func TestMarshalJsonIter(t *testing.T){
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	data := &TopRank{}
	for i := 0; i < nArraySize; i++{
		data.Value = append(data.Value, nValue)
	}
	for i := 0; i < ntimes; i++{
		json.Marshal(data)
	}
}

func TestUMarshalJsonIter(t *testing.T){
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	data := &TopRank{}
	for i := 0; i < nArraySize; i++{
		data.Value = append(data.Value, nValue)
	}

	for i := 0; i < ntimes; i++{
		buff, _ := json.Marshal(data)
		json.Unmarshal(buff, &TopRank{})
	}
}

func TestMarshalPB(t *testing.T){
	aa := []int32{}
	for i := 0; i < nArraySize; i++{
		aa = append(aa, int32(nValue))
	}
	for i := 0; i < ntimes; i++{
		proto.Marshal(&message.W_C_Test{Recv:aa})
	}
}

func TestUMarshalPB(t *testing.T){
	aa := []int32{}
	for i := 0; i < nArraySize; i++{
		aa = append(aa, int32(nValue))
	}
	for i := 0; i < ntimes; i++{
		buff, _ := proto.Marshal(&message.W_C_Test{Recv:aa})
		proto.Unmarshal(buff, &message.W_C_Test{})
	}
}

func TestMarshalGob(t *testing.T){
	data := &TopRank{}
	for i := 0; i < nArraySize; i++{
		data.Value = append(data.Value, nValue)
	}
	for i := 0; i < ntimes; i++{
		//enc.Encode(int(0))
		buf := &bytes.Buffer{}
		enc := gob.NewEncoder(buf)
		enc.Encode(data)
	}
}

func TestUMarshalGob(t *testing.T){
	data := &TopRank{}
	for i := 0; i < nArraySize; i++{
		data.Value = append(data.Value, nValue)
	}

	//fmt.Println(buf.Bytes(), len(buf.Bytes()))
	for i := 0; i < ntimes; i++{
		buf := bytes.NewBuffer([]byte{})
		enc := gob.NewEncoder(buf)
		dec := gob.NewDecoder(buf)
		enc.Encode(data)
		aa1 := &TopRank{}
		dec.Decode(aa1)
	}
}

func TestMarshalRpc(t *testing.T){
	aa := []int32{}
	for i := 0; i < nArraySize; i++{
		aa = append(aa, int32(nValue))
	}
	for i := 0; i < ntimes; i++{
		rpc.Marshal("test", aa)
	}
}

func TestUMarshalRpc(t *testing.T){
	aa := []int32{}
	for i := 0; i < nArraySize; i++{
		aa = append(aa, int32(nValue))
	}
	for i := 0; i < ntimes; i++{
		buff := rpc.Marshal("test", aa)
		parse(buff)
	}
}

func parse (buff []byte) {
	rpcPacket := rpc.UnmarshalHead(buff)
	rpc.UnmarshalBody(rpcPacket, nil)
}
