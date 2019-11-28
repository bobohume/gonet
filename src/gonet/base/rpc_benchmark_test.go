package base_test

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"github.com/golang/protobuf/proto"
	"github.com/json-iterator/go"
	"gonet/base"
	"gonet/message"
	"testing"
)

func Benchmark_TestMarshalJson(b *testing.B){
	b.StartTimer()
	data := &TopRank{}
	for i := 0; i < nArraySize; i++{
		data.Value = append(data.Value, nValue)
	}
	for i := 0; i < ntimes; i++{
		json.Marshal(data)
	}
	b.StopTimer()
}

func Benchmark_TestUMarshalJson(b *testing.B){
	b.StartTimer()
	data := &TopRank{}
	for i := 0; i < nArraySize; i++{
		data.Value = append(data.Value, nValue)
	}
	for i := 0; i < ntimes; i++{
		buff, _ := json.Marshal(data)
		json.Unmarshal(buff, &TopRank{})
	}
	b.StopTimer()
}

func Benchmark_TestMarshalJsonIter(b *testing.B){
	b.StartTimer()
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	data := &TopRank{}
	for i := 0; i < nArraySize; i++{
		data.Value = append(data.Value, nValue)
	}
	for i := 0; i < ntimes; i++{
		json.Marshal(data)
	}
	b.StopTimer()
}

func Benchmark_TestUMarshalJsonIter(b *testing.B){
	b.StartTimer()
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	data := &TopRank{}
	for i := 0; i < nArraySize; i++{
		data.Value = append(data.Value, nValue)
	}
	for i := 0; i < ntimes; i++{
		buff, _ := json.Marshal(data)
		json.Unmarshal(buff, &TopRank{})
	}
	b.StopTimer()
}

func Benchmark_TestMarshalPB(b *testing.B){
	b.StartTimer()
	aa := []int32{}
	for i := 0; i < nArraySize; i++{
		aa = append(aa, int32(nValue))
	}
	for i := 0; i < ntimes; i++{
		proto.Marshal(&message.W_C_Test{Recv:aa})
	}
	b.StopTimer()
}

func Benchmark_TestUMarshalPB(b *testing.B){
	b.StartTimer()
	aa := []int32{}
	for i := 0; i < nArraySize; i++{
		aa = append(aa, int32(nValue))
	}
	for i := 0; i < ntimes; i++{
		buff, _ := proto.Marshal(&message.W_C_Test{Recv:aa})
		proto.Unmarshal(buff, &message.W_C_Test{})
	}
	b.StopTimer()
}

func Benchmark_TestMarshalGob(b *testing.B){
	b.StartTimer()
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
	b.StopTimer()
}

func Benchmark_TestUMarshalGob(b *testing.B){
	b.StartTimer()
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
	b.StopTimer()
}

func Benchmark_TestMarshalRpc(b *testing.B){
	b.StartTimer()
	aa := []int32{}
	for i := 0; i < nArraySize; i++{
		aa = append(aa, int32(nValue))
	}
	for i := 0; i < ntimes; i++{
		base.GetPacket("test", aa)
	}
	b.StopTimer()
}

func Benchmark_TestUMarshalRpc(b *testing.B){
	b.StartTimer()
	aa := []int32{}
	for i := 0; i < nArraySize; i++{
		aa = append(aa, int32(nValue))
	}
	for i := 0; i < ntimes; i++{
		buff := base.GetPacket("test", aa)
		parse(buff)
	}
	b.StopTimer()
}