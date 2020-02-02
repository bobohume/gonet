package rpc

import (
	"github.com/golang/protobuf/proto"
	"github.com/json-iterator/go"
	"gonet/base"
	"gonet/message"
	"reflect"
)

func readBool(bitstream base.IBitStream)(bool){
	nLen := bitstream.ReadInt(8)
	val1 := &message.Bool{}
	proto.Unmarshal(bitstream.ReadBits(nLen << 3), val1)
	return val1.Val
}

func readString(bitstream base.IBitStream)(string){
	nLen := bitstream.ReadInt(8)
	val1 := &message.String{}
	proto.Unmarshal(bitstream.ReadBits(nLen << 3), val1)
	return val1.Val
}

func readFloat32(bitstream base.IBitStream)(float32){
	nLen := bitstream.ReadInt(8)
	val1 := &message.Float{}
	proto.Unmarshal(bitstream.ReadBits(nLen << 3), val1)
	return val1.Val
}

func readFloat64(bitstream base.IBitStream)(float64){
	nLen := bitstream.ReadInt(8)
	val1 := &message.Double{}
	proto.Unmarshal(bitstream.ReadBits(nLen << 3), val1)
	return val1.Val
}

func readInt32(bitstream base.IBitStream)(int32){
	nLen := bitstream.ReadInt(8)
	val1 := &message.Int{}
	proto.Unmarshal(bitstream.ReadBits(nLen << 3), val1)
	return val1.Val
}

func ReadInt64(bitstream base.IBitStream)(int64){
	return readInt64(bitstream)
}

func readInt64(bitstream base.IBitStream)(int64){
	nLen := bitstream.ReadInt(8)
	val1 := &message.Int64{}
	proto.Unmarshal(bitstream.ReadBits(nLen << 3), val1)
	return val1.Val
}

func readUInt32(bitstream base.IBitStream)(uint32){
	nLen := bitstream.ReadInt(8)
	val1 := &message.UInt{}
	proto.Unmarshal(bitstream.ReadBits(nLen << 3), val1)
	return val1.Val
}

func readUInt64(bitstream base.IBitStream)(uint64){
	nLen := bitstream.ReadInt(8)
	val1 := &message.UInt64{}
	proto.Unmarshal(bitstream.ReadBits(nLen << 3), val1)
	return val1.Val
}

func readBoolSlice(bitstream base.IBitStream)([]bool){
	nLen := bitstream.ReadInt(16)
	val1 := &message.BoolSlice{}
	proto.Unmarshal(bitstream.ReadBits(nLen << 3), val1)
	return val1.Val
}

func readStringSlice(bitstream base.IBitStream)([]string){
	nLen := bitstream.ReadInt(16)
	val1 := &message.StringSlice{}
	proto.Unmarshal(bitstream.ReadBits(nLen << 3), val1)
	return val1.Val
}

func readFloat32Slice(bitstream base.IBitStream)([]float32){
	nLen := bitstream.ReadInt(16)
	val1 := &message.FloatSlice{}
	proto.Unmarshal(bitstream.ReadBits(nLen << 3), val1)
	return val1.Val
}

func readFloat64Slice(bitstream base.IBitStream)([]float64){
	nLen := bitstream.ReadInt(16)
	val1 := &message.DoubleSlice{}
	proto.Unmarshal(bitstream.ReadBits(nLen << 3), val1)
	return val1.Val
}

func readIntSlice(bitstream base.IBitStream)([]int){
	nLen := bitstream.ReadInt(16)
	val1 := &message.IntSlice{}
	proto.Unmarshal(bitstream.ReadBits(nLen << 3), val1)
	val0 := make([]int, len(val1.Val))
	for i, v := range val1.Val{
		val0[i] = int(v)
	}
	return val0
}

func readInt8Slice(bitstream base.IBitStream)([]int8){
	nLen := bitstream.ReadInt(16)
	val1 := &message.IntSlice{}
	proto.Unmarshal(bitstream.ReadBits(nLen << 3), val1)
	val0 := make([]int8, len(val1.Val))
	for i, v := range val1.Val{
		val0[i] = int8(v)
	}
	return val0
}

func readInt16Slice(bitstream base.IBitStream)([]int16){
	nLen := bitstream.ReadInt(16)
	val1 := &message.IntSlice{}
	proto.Unmarshal(bitstream.ReadBits(nLen << 3), val1)
	val0 := make([]int16, len(val1.Val))
	for i, v := range val1.Val{
		val0[i] = int16(v)
	}
	return val0
}

func readInt32Slice(bitstream base.IBitStream)([]int32){
	nLen := bitstream.ReadInt(16)
	val1 := &message.IntSlice{}
	proto.Unmarshal(bitstream.ReadBits(nLen << 3), val1)
	return val1.Val
}

func readInt64Slice(bitstream base.IBitStream)([]int64){
	nLen := bitstream.ReadInt(16)
	val1 := &message.Int64Slice{}
	proto.Unmarshal(bitstream.ReadBits(nLen << 3), val1)
	return val1.Val
}

func readUIntSlice(bitstream base.IBitStream)([]uint){
	nLen := bitstream.ReadInt(16)
	val1 := &message.UIntSlice{}
	proto.Unmarshal(bitstream.ReadBits(nLen << 3), val1)
	val0 := make([]uint, len(val1.Val))
	for i, v := range val1.Val{
		val0[i] = uint(v)
	}
	return val0
}

func readUInt8Slice(bitstream base.IBitStream)([]uint8){
	nLen := bitstream.ReadInt(16)
	val1 := &message.UIntSlice{}
	proto.Unmarshal(bitstream.ReadBits(nLen << 3), val1)
	val0 := make([]uint8, len(val1.Val))
	for i, v := range val1.Val{
		val0[i] = uint8(v)
	}
	return val0
}

func readUInt16Slice(bitstream base.IBitStream)([]uint16){
	nLen := bitstream.ReadInt(16)
	val1 := &message.UIntSlice{}
	proto.Unmarshal(bitstream.ReadBits(nLen << 3), val1)
	val0 := make([]uint16, len(val1.Val))
	for i, v := range val1.Val{
		val0[i] = uint16(v)
	}
	return val0
}

func readUInt32Slice(bitstream base.IBitStream)([]uint32){
	nLen := bitstream.ReadInt(16)
	val1 := &message.UIntSlice{}
	proto.Unmarshal(bitstream.ReadBits(nLen << 3), val1)
	val0 := make([]uint32, len(val1.Val))
	for i, v := range val1.Val{
		val0[i] = uint32(v)
	}
	return val0
}

func readUInt64Slice(bitstream base.IBitStream)([]uint64){
	nLen := bitstream.ReadInt(16)
	val1 := &message.UInt64Slice{}
	proto.Unmarshal(bitstream.ReadBits(nLen << 3), val1)
	return val1.Val
}

func readBoolPtrSlice(bitstream base.IBitStream)([]*bool){
	nLen := bitstream.ReadInt(16)
	val1 := &message.BoolSlice{}
	proto.Unmarshal(bitstream.ReadBits(nLen << 3), val1)
	val0 := make([]*bool, len(val1.Val))
	for i, v := range val1.Val{
		val0[i] = &v
	}
	return val0
}

func readStringPtrSlice(bitstream base.IBitStream)([]*string){
	nLen := bitstream.ReadInt(16)
	val1 := &message.StringSlice{}
	proto.Unmarshal(bitstream.ReadBits(nLen << 3), val1)
	val0 := make([]*string, len(val1.Val))
	for i, v := range val1.Val{
		val0[i] = &v
	}
	return val0
}

func readFloat32PtrSlice(bitstream base.IBitStream)([]*float32){
	nLen := bitstream.ReadInt(16)
	val1 := &message.FloatSlice{}
	proto.Unmarshal(bitstream.ReadBits(nLen << 3), val1)
	val0 := make([]*float32, len(val1.Val))
	for i, v := range val1.Val{
		val0[i] = &v
	}
	return val0
}

func readFloat64PtrSlice(bitstream base.IBitStream)([]*float64){
	nLen := bitstream.ReadInt(16)
	val1 := &message.DoubleSlice{}
	proto.Unmarshal(bitstream.ReadBits(nLen << 3), val1)
	val0 := make([]*float64, len(val1.Val))
	for i, v := range val1.Val{
		val0[i] = &v
	}
	return val0
}

func readIntPtrSlice(bitstream base.IBitStream)([]*int){
	nLen := bitstream.ReadInt(16)
	val1 := &message.IntSlice{}
	proto.Unmarshal(bitstream.ReadBits(nLen << 3), val1)
	val0 := make([]*int, len(val1.Val))
	for i, v := range val1.Val{
		v1 := int(v)
		val0[i] = &v1
	}
	return val0
}

func readInt8PtrSlice(bitstream base.IBitStream)([]*int8){
	nLen := bitstream.ReadInt(16)
	val1 := &message.IntSlice{}
	proto.Unmarshal(bitstream.ReadBits(nLen << 3), val1)
	val0 := make([]*int8, len(val1.Val))
	for i, v := range val1.Val{
		v1 := int8(v)
		val0[i] = &v1
	}
	return val0
}

func readInt16PtrSlice(bitstream base.IBitStream)([]*int16){
	nLen := bitstream.ReadInt(16)
	val1 := &message.IntSlice{}
	proto.Unmarshal(bitstream.ReadBits(nLen << 3), val1)
	val0 := make([]*int16, len(val1.Val))
	for i, v := range val1.Val{
		v1 := int16(v)
		val0[i] = &v1
	}
	return val0
}

func readInt32PtrSlice(bitstream base.IBitStream)([]*int32){
	nLen := bitstream.ReadInt(16)
	val1 := &message.IntSlice{}
	proto.Unmarshal(bitstream.ReadBits(nLen << 3), val1)
	val0 := make([]*int32, len(val1.Val))
	for i, v := range val1.Val{
		v1 := int32(v)
		val0[i] = &v1
	}
	return val0
}

func readInt64PtrSlice(bitstream base.IBitStream)([]*int64){
	nLen := bitstream.ReadInt(16)
	val1 := &message.Int64Slice{}
	proto.Unmarshal(bitstream.ReadBits(nLen << 3), val1)
	val0 := make([]*int64, len(val1.Val))
	for i, v := range val1.Val{
		val0[i] = &v
	}
	return val0
}

func readUIntPtrSlice(bitstream base.IBitStream)([]*uint){
	nLen := bitstream.ReadInt(16)
	val1 := &message.UIntSlice{}
	proto.Unmarshal(bitstream.ReadBits(nLen << 3), val1)
	val0 := make([]*uint, len(val1.Val))
	for i, v := range val1.Val{
		v1 := uint(v)
		val0[i] = &v1
	}
	return val0
}

func readUInt8PtrSlice(bitstream base.IBitStream)([]*uint8){
	nLen := bitstream.ReadInt(16)
	val1 := &message.UIntSlice{}
	proto.Unmarshal(bitstream.ReadBits(nLen << 3), val1)
	val0 := make([]*uint8, len(val1.Val))
	for i, v := range val1.Val{
		v1 := uint8(v)
		val0[i] = &v1
	}
	return val0
}

func readUInt16PtrSlice(bitstream base.IBitStream)([]*uint16){
	nLen := bitstream.ReadInt(16)
	val1 := &message.UIntSlice{}
	proto.Unmarshal(bitstream.ReadBits(nLen << 3), val1)
	val0 := make([]*uint16, len(val1.Val))
	for i, v := range val1.Val{
		v1 := uint16(v)
		val0[i] = &v1
	}
	return val0
}

func readUInt32PtrSlice(bitstream base.IBitStream)([]*uint32){
	nLen := bitstream.ReadInt(16)
	val1 := &message.UIntSlice{}
	proto.Unmarshal(bitstream.ReadBits(nLen << 3), val1)
	val0 := make([]*uint32, len(val1.Val))
	for i, v := range val1.Val{
		v1 := uint32(v)
		val0[i] = &v1
	}
	return val0
}

func readUInt64PtrSlice(bitstream base.IBitStream)([]*uint64){
	nLen := bitstream.ReadInt(16)
	val1 := &message.UInt64Slice{}
	proto.Unmarshal(bitstream.ReadBits(nLen << 3), val1)
	val0 := make([]*uint64, len(val1.Val))
	for i, v := range val1.Val{
		val0[i] = &v
	}
	return val0
}

//rpc Unmarshal
//pFuncType for RegisterCall func
func Unmarshal(bitstream *base.BitStream, funcName string, pFuncType reflect.Type) []interface{}{
	nCurLen := bitstream.ReadInt(8)
	params := make([]interface{}, nCurLen)
	for i := 0; i < nCurLen; i++  {
		switch bitstream.ReadInt(8) {
		case RPC_BOOL:
			params[i] = readBool(bitstream)
		case RPC_STRING:
			params[i] = readString(bitstream)
		case RPC_FLOAT32:
			params[i] = readFloat32(bitstream)
		case RPC_FLOAT64:
			params[i] = readFloat64(bitstream)
		case RPC_INT:
			params[i] = int(readInt32(bitstream))
		case RPC_INT8:
			params[i] = int8(readInt32(bitstream))
		case RPC_INT16:
			params[i] = int16(readInt32(bitstream))
		case RPC_INT32:
			params[i] = int32(readInt32(bitstream))
		case RPC_INT64:
			params[i] = int64(readInt64(bitstream))
		case RPC_UINT:
			params[i] = uint(readUInt32(bitstream))
		case RPC_UINT8:
			params[i] = uint8(readUInt32(bitstream))
		case RPC_UINT16:
			params[i] = uint16(readUInt32(bitstream))
		case RPC_UINT32:
			params[i] = uint32(readUInt32(bitstream))
		case RPC_UINT64:
			params[i] = uint64(readUInt64(bitstream))



		case RPC_BOOL_SLICE:
			params[i] = readBoolSlice(bitstream)
		case RPC_STRING_SLICE:
			params[i] = readStringSlice(bitstream)
		case RPC_FLOAT32_SLICE:
			params[i] = readFloat32Slice(bitstream)
		case RPC_FLOAT64_SLICE:
			params[i] = readFloat64Slice(bitstream)
		case RPC_INT_SLICE:
			params[i] = readIntSlice(bitstream)
		case RPC_INT8_SLICE:
			params[i] = readInt8Slice(bitstream)
		case RPC_INT16_SLICE:
			params[i] = readInt16Slice(bitstream)
		case RPC_INT32_SLICE:
			params[i] = readInt32Slice(bitstream)
		case RPC_INT64_SLICE:
			params[i] = readInt64Slice(bitstream)
		case RPC_UINT_SLICE:
			params[i] = readUIntSlice(bitstream)
		case RPC_UINT8_SLICE:
			params[i] = readUInt8Slice(bitstream)
		case RPC_UINT16_SLICE:
			params[i] = readUInt16Slice(bitstream)
		case RPC_UINT32_SLICE:
			params[i] = readUInt32Slice(bitstream)
		case RPC_UINT64_SLICE:
			params[i] = readUInt64Slice(bitstream)



		case RPC_BOOL_ARRAY:
			val0 := readBoolSlice(bitstream)
			val1 := reflect.New(reflect.ArrayOf(len(val0), reflect.TypeOf(bool(false)))).Elem()
			reflect.Copy(val1, reflect.ValueOf(val0))
			params[i] = val1.Interface()
		case RPC_STRING_ARRAY:
			val0 := readStringSlice(bitstream)
			val1 := reflect.New(reflect.ArrayOf(len(val0), reflect.TypeOf(string("")))).Elem()
			reflect.Copy(val1, reflect.ValueOf(val0))
			params[i] = val1.Interface()
		case RPC_FLOAT32_ARRAY:
			val0 := readFloat32Slice(bitstream)
			val1 := reflect.New(reflect.ArrayOf(len(val0), reflect.TypeOf(float32(0)))).Elem()
			reflect.Copy(val1, reflect.ValueOf(val0))
			params[i] = val1.Interface()
		case RPC_FLOAT64_ARRAY:
			val0 := readFloat64Slice(bitstream)
			val1 := reflect.New(reflect.ArrayOf(len(val0), reflect.TypeOf(float64(0)))).Elem()
			reflect.Copy(val1, reflect.ValueOf(val0))
			params[i] = val1.Interface()
		case RPC_INT_ARRAY:
			val0 := readIntSlice(bitstream)
			val1 := reflect.New(reflect.ArrayOf(len(val0), reflect.TypeOf(int(0)))).Elem()
			reflect.Copy(val1, reflect.ValueOf(val0))
			params[i] = val1.Interface()
		case RPC_INT8_ARRAY:
			val0 := readInt8Slice(bitstream)
			val1 := reflect.New(reflect.ArrayOf(len(val0), reflect.TypeOf(int8(0)))).Elem()
			reflect.Copy(val1, reflect.ValueOf(val0))
			params[i] = val1.Interface()
		case RPC_INT16_ARRAY:
			val0 := readInt16Slice(bitstream)
			val1 := reflect.New(reflect.ArrayOf(len(val0), reflect.TypeOf(int16(0)))).Elem()
			reflect.Copy(val1, reflect.ValueOf(val0))
			params[i] = val1.Interface()
		case RPC_INT32_ARRAY:
			val0 := readInt32Slice(bitstream)
			val1 := reflect.New(reflect.ArrayOf(len(val0), reflect.TypeOf(int32(0)))).Elem()
			reflect.Copy(val1, reflect.ValueOf(val0))
			params[i] = val1.Interface()
		case RPC_INT64_ARRAY:
			val0 := readInt64Slice(bitstream)
			val1 := reflect.New(reflect.ArrayOf(len(val0), reflect.TypeOf(int64(0)))).Elem()
			reflect.Copy(val1, reflect.ValueOf(val0))
			params[i] = val1.Interface()
		case RPC_UINT_ARRAY:
			val0 := readUIntSlice(bitstream)
			val1 := reflect.New(reflect.ArrayOf(len(val0), reflect.TypeOf(uint(0)))).Elem()
			reflect.Copy(val1, reflect.ValueOf(val0))
			params[i] = val1.Interface()
		case RPC_UINT8_ARRAY:
			val0 := readUInt8Slice(bitstream)
			val1 := reflect.New(reflect.ArrayOf(len(val0), reflect.TypeOf(uint8(0)))).Elem()
			reflect.Copy(val1, reflect.ValueOf(val0))
			params[i] = val1.Interface()
		case RPC_UINT16_ARRAY:
			val0 := readUInt16Slice(bitstream)
			val1 := reflect.New(reflect.ArrayOf(len(val0), reflect.TypeOf(uint16(0)))).Elem()
			reflect.Copy(val1, reflect.ValueOf(val0))
			params[i] = val1.Interface()
		case RPC_UINT32_ARRAY:
			val0 := readUInt32Slice(bitstream)
			val1 := reflect.New(reflect.ArrayOf(len(val0), reflect.TypeOf(uint32(0)))).Elem()
			reflect.Copy(val1, reflect.ValueOf(val0))
			params[i] = val1.Interface()
		case RPC_UINT64_ARRAY:
			val0 := readUInt64Slice(bitstream)
			val1 := reflect.New(reflect.ArrayOf(len(val0), reflect.TypeOf(uint64(0)))).Elem()
			reflect.Copy(val1, reflect.ValueOf(val0))
			params[i] = val1.Interface()



		case RPC_BOOL_PTR:
			val := new(bool)
			*val = readBool(bitstream)
			params[i] = val
		case RPC_STRING_PTR:
			val := new(string)
			*val = readString(bitstream)
			params[i] = val
		case RPC_FLOAT32_PTR:
			val := new(float32)
			*val = readFloat32(bitstream)
			params[i] = val
		case RPC_FLOAT64_PTR:
			val := new(float64)
			*val = readFloat64(bitstream)
			params[i] = val
		case RPC_INT_PTR:
			val := new(int)
			*val = int(readInt32(bitstream))
			params[i] = val
		case RPC_INT8_PTR:
			val := new(int8)
			*val = int8(readInt32(bitstream))
			params[i] = val
		case RPC_INT16_PTR:
			val := new(int16)
			*val = int16(readInt32(bitstream))
			params[i] = val
		case RPC_INT32_PTR:
			val := new(int32)
			*val = int32(readInt32(bitstream))
			params[i] = val
		case RPC_INT64_PTR:
			val := new(int64)
			*val = int64(readInt64(bitstream))
			params[i] = val
		case RPC_UINT_PTR:
			val := new(uint)
			*val = uint(readUInt32(bitstream))
			params[i] = val
		case RPC_UINT8_PTR:
			val := new(uint8)
			*val = uint8(readUInt32(bitstream))
			params[i] = val
		case RPC_UINT16_PTR:
			val := new(uint16)
			*val = uint16(readUInt32(bitstream))
			params[i] = val
		case RPC_UINT32_PTR:
			val := new(uint32)
			*val = uint32(readUInt32(bitstream))
			params[i] = val
		case RPC_UINT64_PTR:
			val := new(uint64)
			*val = uint64(readUInt64(bitstream))
			params[i] = val



		case RPC_BOOL_PTR_SLICE:
			params[i] = readBoolPtrSlice(bitstream)
		case RPC_STRING_PTR_SLICE:
			params[i] = readStringPtrSlice(bitstream)
		case RPC_FLOAT32_PTR_SLICE:
			params[i] = readFloat32PtrSlice(bitstream)
		case RPC_FLOAT64_PTR_SLICE:
			params[i] = readFloat64PtrSlice(bitstream)
		case RPC_INT_PTR_SLICE:
			params[i] = readIntPtrSlice(bitstream)
		case RPC_INT8_PTR_SLICE:
			params[i] = readInt8PtrSlice(bitstream)
		case RPC_INT16_PTR_SLICE:
			params[i] = readInt16PtrSlice(bitstream)
		case RPC_INT32_PTR_SLICE:
			params[i] = readInt32PtrSlice(bitstream)
		case RPC_INT64_PTR_SLICE:
			params[i] = readInt64PtrSlice(bitstream)
		case RPC_UINT_PTR_SLICE:
			params[i] = readUIntPtrSlice(bitstream)
		case RPC_UINT8_PTR_SLICE:
			params[i] = readUInt8PtrSlice(bitstream)
		case RPC_UINT16_PTR_SLICE:
			params[i] = readUInt16PtrSlice(bitstream)
		case RPC_UINT32_PTR_SLICE:
			params[i] = readUInt32PtrSlice(bitstream)
		case RPC_UINT64_PTR_SLICE:
			params[i] = readUInt64PtrSlice(bitstream)



		case RPC_BOOL_PTR_ARRAY:
			val0 := readBoolPtrSlice(bitstream)
			val2 := bool(false)
			val1 := reflect.New(reflect.ArrayOf(len(val0), reflect.TypeOf(&val2))).Elem()
			reflect.Copy(val1, reflect.ValueOf(val0))
			params[i] = val1.Interface()
		case RPC_STRING_PTR_ARRAY:
			val0 := readStringPtrSlice(bitstream)
			val2 := string("")
			val1 := reflect.New(reflect.ArrayOf(len(val0), reflect.TypeOf(&val2))).Elem()
			reflect.Copy(val1, reflect.ValueOf(val0))
			params[i] = val1.Interface()
		case RPC_FLOAT32_PTR_ARRAY:
			val0 := readFloat32PtrSlice(bitstream)
			val2 := float32(0)
			val1 := reflect.New(reflect.ArrayOf(len(val0), reflect.TypeOf(&val2))).Elem()
			reflect.Copy(val1, reflect.ValueOf(val0))
			params[i] = val1.Interface()
		case RPC_FLOAT64_PTR_ARRAY:
			val0 := readFloat64PtrSlice(bitstream)
			val2 := float64(0)
			val1 := reflect.New(reflect.ArrayOf(len(val0), reflect.TypeOf(&val2))).Elem()
			reflect.Copy(val1, reflect.ValueOf(val0))
			params[i] = val1.Interface()
		case RPC_INT_PTR_ARRAY:
			val0 := readIntPtrSlice(bitstream)
			val2 := int(0)
			val1 := reflect.New(reflect.ArrayOf(len(val0), reflect.TypeOf(&val2))).Elem()
			reflect.Copy(val1, reflect.ValueOf(val0))
			params[i] = val1.Interface()
		case RPC_INT8_PTR_ARRAY:
			val0 := readInt8PtrSlice(bitstream)
			val2 := int8(0)
			val1 := reflect.New(reflect.ArrayOf(len(val0), reflect.TypeOf(&val2))).Elem()
			reflect.Copy(val1, reflect.ValueOf(val0))
			params[i] = val1.Interface()
		case RPC_INT16_PTR_ARRAY:
			val0 := readInt16PtrSlice(bitstream)
			val2 := int16(0)
			val1 := reflect.New(reflect.ArrayOf(len(val0), reflect.TypeOf(&val2))).Elem()
			reflect.Copy(val1, reflect.ValueOf(val0))
			params[i] = val1.Interface()
		case RPC_INT32_PTR_ARRAY:
			val0 := readInt32PtrSlice(bitstream)
			val2 := int32(0)
			val1 := reflect.New(reflect.ArrayOf(len(val0), reflect.TypeOf(&val2))).Elem()
			reflect.Copy(val1, reflect.ValueOf(val0))
			params[i] = val1.Interface()
		case RPC_INT64_PTR_ARRAY:
			val0 := readInt64PtrSlice(bitstream)
			val2 := int64(0)
			val1 := reflect.New(reflect.ArrayOf(len(val0), reflect.TypeOf(&val2))).Elem()
			reflect.Copy(val1, reflect.ValueOf(val0))
			params[i] = val1.Interface()
		case RPC_UINT_PTR_ARRAY:
			val0 := readUIntPtrSlice(bitstream)
			val2 := uint(0)
			val1 := reflect.New(reflect.ArrayOf(len(val0), reflect.TypeOf(&val2))).Elem()
			reflect.Copy(val1, reflect.ValueOf(val0))
			params[i] = val1.Interface()
		case RPC_UINT8_PTR_ARRAY:
			val0 := readUInt8PtrSlice(bitstream)
			val2 := uint8(0)
			val1 := reflect.New(reflect.ArrayOf(len(val0), reflect.TypeOf(&val2))).Elem()
			reflect.Copy(val1, reflect.ValueOf(val0))
			params[i] = val1.Interface()
		case RPC_UINT16_PTR_ARRAY:
			val0 := readUInt16PtrSlice(bitstream)
			val2 := uint16(0)
			val1 := reflect.New(reflect.ArrayOf(len(val0), reflect.TypeOf(&val2))).Elem()
			reflect.Copy(val1, reflect.ValueOf(val0))
			params[i] = val1.Interface()
		case RPC_UINT32_PTR_ARRAY:
			val0 := readUInt32PtrSlice(bitstream)
			val2 := uint32(0)
			val1 := reflect.New(reflect.ArrayOf(len(val0), reflect.TypeOf(&val2))).Elem()
			reflect.Copy(val1, reflect.ValueOf(val0))
			params[i] = val1.Interface()
		case RPC_UINT64_PTR_ARRAY:
			val0 := readUInt64PtrSlice(bitstream)
			val2 := uint64(0)
			val1 := reflect.New(reflect.ArrayOf(len(val0), reflect.TypeOf(&val2))).Elem()
			reflect.Copy(val1, reflect.ValueOf(val0))
			params[i] = val1.Interface()



		case RPC_BOOL_SLICE_PTR:
			val := readBoolSlice(bitstream)
			params[i] = &val
		case RPC_STRING_SLICE_PTR:
			val := readStringSlice(bitstream)
			params[i] = &val
		case RPC_FLOAT32_SLICE_PTR:
			val := readFloat32Slice(bitstream)
			params[i] = &val
		case RPC_FLOAT64_SLICE_PTR:
			val := readFloat64Slice(bitstream)
			params[i] = &val
		case RPC_INT_SLICE_PTR:
			val := readIntSlice(bitstream)
			params[i] = &val
		case RPC_INT8_SLICE_PTR:
			val := readInt8Slice(bitstream)
			params[i] = &val
		case RPC_INT16_SLICE_PTR:
			val := readInt16Slice(bitstream)
			params[i] = &val
		case RPC_INT32_SLICE_PTR:
			val := readInt32Slice(bitstream)
			params[i] = &val
		case RPC_INT64_SLICE_PTR:
			val := readInt64Slice(bitstream)
			params[i] = &val
		case RPC_UINT_SLICE_PTR:
			val := readUIntSlice(bitstream)
			params[i] = &val
		case RPC_UINT8_SLICE_PTR:
			val := readUInt8Slice(bitstream)
			params[i] = &val
		case RPC_UINT16_SLICE_PTR:
			val := readUInt16Slice(bitstream)
			params[i] = &val
		case RPC_UINT32_SLICE_PTR:
			val := readUInt32Slice(bitstream)
			params[i] = &val
		case RPC_UINT64_SLICE_PTR:
			val := readUInt64Slice(bitstream)
			params[i] = &val


		case RPC_BOOL_PTR_SLICE_PTR:
			val := readBoolPtrSlice(bitstream)
			params[i] = &val
		case RPC_STRING_PTR_SLICE_PTR:
			val := readStringPtrSlice(bitstream)
			params[i] = &val
		case RPC_FLOAT32_PTR_SLICE_PTR:
			val := readFloat32PtrSlice(bitstream)
			params[i] = &val
		case RPC_FLOAT64_PTR_SLICE_PTR:
			val := readFloat64PtrSlice(bitstream)
			params[i] = &val
		case RPC_INT_PTR_SLICE_PTR:
			val := readIntPtrSlice(bitstream)
			params[i] = &val
		case RPC_INT8_PTR_SLICE_PTR:
			val := readInt8PtrSlice(bitstream)
			params[i] = &val
		case RPC_INT16_PTR_SLICE_PTR:
			val := readInt16PtrSlice(bitstream)
			params[i] = &val
		case RPC_INT32_PTR_SLICE_PTR:
			val := readInt32PtrSlice(bitstream)
			params[i] = &val
		case RPC_INT64_PTR_SLICE_PTR:
			val := readInt64PtrSlice(bitstream)
			params[i] = &val
		case RPC_UINT_PTR_SLICE_PTR:
			val := readUIntPtrSlice(bitstream)
			params[i] = &val
		case RPC_UINT8_PTR_SLICE_PTR:
			val := readUInt8PtrSlice(bitstream)
			params[i] = &val
		case RPC_UINT16_PTR_SLICE_PTR:
			val := readUInt16PtrSlice(bitstream)
			params[i] = &val
		case RPC_UINT32_PTR_SLICE_PTR:
			val := readUInt32PtrSlice(bitstream)
			params[i] = &val
		case RPC_UINT64_PTR_SLICE_PTR:
			val := readUInt64PtrSlice(bitstream)
			params[i] = &val



		case RPC_BOOL_ARRAY_PTR:
			val0 := readBoolSlice(bitstream)
			val1 := reflect.New(reflect.ArrayOf(len(val0), reflect.TypeOf(bool(false))))
			reflect.Copy(val1.Elem(), reflect.ValueOf(val0))
			params[i] = val1.Interface()
		case RPC_STRING_ARRAY_PTR:
			val0 := readStringSlice(bitstream)
			val1 := reflect.New(reflect.ArrayOf(len(val0), reflect.TypeOf(string(""))))
			reflect.Copy(val1.Elem(), reflect.ValueOf(val0))
			params[i] = val1.Interface()
		case RPC_FLOAT32_ARRAY_PTR:
			val0 := readFloat32Slice(bitstream)
			val1 := reflect.New(reflect.ArrayOf(len(val0), reflect.TypeOf(float32(0))))
			reflect.Copy(val1.Elem(), reflect.ValueOf(val0))
			params[i] = val1.Interface()
		case RPC_FLOAT64_ARRAY_PTR:
			val0 := readFloat64Slice(bitstream)
			val1 := reflect.New(reflect.ArrayOf(len(val0), reflect.TypeOf(float64(0))))
			reflect.Copy(val1.Elem(), reflect.ValueOf(val0))
			params[i] = val1.Interface()
		case RPC_INT_ARRAY_PTR:
			val0 := readIntSlice(bitstream)
			val1 := reflect.New(reflect.ArrayOf(len(val0), reflect.TypeOf(int(0))))
			reflect.Copy(val1.Elem(), reflect.ValueOf(val0))
			params[i] = val1.Interface()
		case RPC_INT8_ARRAY_PTR:
			val0 := readInt8Slice(bitstream)
			val1 := reflect.New(reflect.ArrayOf(len(val0), reflect.TypeOf(int8(0))))
			reflect.Copy(val1.Elem(), reflect.ValueOf(val0))
			params[i] = val1.Interface()
		case RPC_INT16_ARRAY_PTR:
			val0 := readInt16Slice(bitstream)
			val1 := reflect.New(reflect.ArrayOf(len(val0), reflect.TypeOf(int16(0))))
			reflect.Copy(val1.Elem(), reflect.ValueOf(val0))
			params[i] = val1.Interface()
		case RPC_INT32_ARRAY_PTR:
			val0 := readInt32Slice(bitstream)
			val1 := reflect.New(reflect.ArrayOf(len(val0), reflect.TypeOf(int32(0))))
			reflect.Copy(val1.Elem(), reflect.ValueOf(val0))
			params[i] = val1.Interface()
		case RPC_INT64_ARRAY_PTR:
			val0 := readInt64Slice(bitstream)
			val1 := reflect.New(reflect.ArrayOf(len(val0), reflect.TypeOf(int64(0))))
			reflect.Copy(val1.Elem(), reflect.ValueOf(val0))
			params[i] = val1.Interface()
		case RPC_UINT_ARRAY_PTR:
			val0 := readUIntSlice(bitstream)
			val1 := reflect.New(reflect.ArrayOf(len(val0), reflect.TypeOf(uint(0))))
			reflect.Copy(val1.Elem(), reflect.ValueOf(val0))
			params[i] = val1.Interface()
		case RPC_UINT8_ARRAY_PTR:
			val0 := readUInt8Slice(bitstream)
			val1 := reflect.New(reflect.ArrayOf(len(val0), reflect.TypeOf(uint8(0))))
			reflect.Copy(val1.Elem(), reflect.ValueOf(val0))
			params[i] = val1.Interface()
		case RPC_UINT16_ARRAY_PTR:
			val0 := readUInt16Slice(bitstream)
			val1 := reflect.New(reflect.ArrayOf(len(val0), reflect.TypeOf(uint16(0))))
			reflect.Copy(val1.Elem(), reflect.ValueOf(val0))
			params[i] = val1.Interface()
		case RPC_UINT32_ARRAY_PTR:
			val0 := readUInt32Slice(bitstream)
			val1 := reflect.New(reflect.ArrayOf(len(val0), reflect.TypeOf(uint32(0))))
			reflect.Copy(val1.Elem(), reflect.ValueOf(val0))
			params[i] = val1.Interface()
		case RPC_UINT64_ARRAY_PTR:
			val0 := readUInt64Slice(bitstream)
			val1 := reflect.New(reflect.ArrayOf(len(val0), reflect.TypeOf(uint64(0))))
			reflect.Copy(val1.Elem(), reflect.ValueOf(val0))
			params[i] = val1.Interface()



		case RPC_BOOL_PTR_ARRAY_PTR:
			val0 := readBoolPtrSlice(bitstream)
			val2 := bool(false)
			val1 := reflect.New(reflect.ArrayOf(len(val0), reflect.TypeOf(&val2)))
			reflect.Copy(val1.Elem(), reflect.ValueOf(val0))
			params[i] = val1.Interface()
		case RPC_STRING_PTR_ARRAY_PTR:
			val0 := readStringPtrSlice(bitstream)
			val2 := string("")
			val1 := reflect.New(reflect.ArrayOf(len(val0), reflect.TypeOf(&val2)))
			reflect.Copy(val1.Elem(), reflect.ValueOf(val0))
			params[i] = val1.Interface()
		case RPC_FLOAT32_PTR_ARRAY_PTR:
			val0 := readFloat32PtrSlice(bitstream)
			val2 := float32(0)
			val1 := reflect.New(reflect.ArrayOf(len(val0), reflect.TypeOf(&val2)))
			reflect.Copy(val1.Elem(), reflect.ValueOf(val0))
			params[i] = val1.Interface()
		case RPC_FLOAT64_PTR_ARRAY_PTR:
			val0 := readFloat64PtrSlice(bitstream)
			val2 := float64(0)
			val1 := reflect.New(reflect.ArrayOf(len(val0), reflect.TypeOf(&val2)))
			reflect.Copy(val1.Elem(), reflect.ValueOf(val0))
			params[i] = val1.Interface()
		case RPC_INT_PTR_ARRAY_PTR:
			val0 := readIntPtrSlice(bitstream)
			val2 := int(0)
			val1 := reflect.New(reflect.ArrayOf(len(val0), reflect.TypeOf(&val2)))
			reflect.Copy(val1.Elem(), reflect.ValueOf(val0))
			params[i] = val1.Interface()
		case RPC_INT8_PTR_ARRAY_PTR:
			val0 := readInt8PtrSlice(bitstream)
			val2 := int8(0)
			val1 := reflect.New(reflect.ArrayOf(len(val0), reflect.TypeOf(&val2)))
			reflect.Copy(val1.Elem(), reflect.ValueOf(val0))
			params[i] = val1.Interface()
		case RPC_INT16_PTR_ARRAY_PTR:
			val0 := readInt16PtrSlice(bitstream)
			val2 := int16(0)
			val1 := reflect.New(reflect.ArrayOf(len(val0), reflect.TypeOf(&val2)))
			reflect.Copy(val1.Elem(), reflect.ValueOf(val0))
			params[i] = val1.Interface()
		case RPC_INT32_PTR_ARRAY_PTR:
			val0 := readInt32PtrSlice(bitstream)
			val2 := int32(0)
			val1 := reflect.New(reflect.ArrayOf(len(val0), reflect.TypeOf(&val2)))
			reflect.Copy(val1.Elem(), reflect.ValueOf(val0))
			params[i] = val1.Interface()
		case RPC_INT64_PTR_ARRAY_PTR:
			val0 := readInt64PtrSlice(bitstream)
			val2 := int64(0)
			val1 := reflect.New(reflect.ArrayOf(len(val0), reflect.TypeOf(&val2)))
			reflect.Copy(val1.Elem(), reflect.ValueOf(val0))
			params[i] = val1.Interface()
		case RPC_UINT_PTR_ARRAY_PTR:
			val0 := readUIntPtrSlice(bitstream)
			val2 := uint(0)
			val1 := reflect.New(reflect.ArrayOf(len(val0), reflect.TypeOf(&val2)))
			reflect.Copy(val1.Elem(), reflect.ValueOf(val0))
			params[i] = val1.Interface()
		case RPC_UINT8_PTR_ARRAY_PTR:
			val0 := readUInt8PtrSlice(bitstream)
			val2 := uint8(0)
			val1 := reflect.New(reflect.ArrayOf(len(val0), reflect.TypeOf(&val2)))
			reflect.Copy(val1.Elem(), reflect.ValueOf(val0))
			params[i] = val1.Interface()
		case RPC_UINT16_PTR_ARRAY_PTR:
			val0 := readUInt16PtrSlice(bitstream)
			val2 := uint16(0)
			val1 := reflect.New(reflect.ArrayOf(len(val0), reflect.TypeOf(&val2)))
			reflect.Copy(val1.Elem(), reflect.ValueOf(val0))
			params[i] = val1.Interface()
		case RPC_UINT32_PTR_ARRAY_PTR:
			val0 := readUInt32PtrSlice(bitstream)
			val2 := uint32(0)
			val1 := reflect.New(reflect.ArrayOf(len(val0), reflect.TypeOf(&val2)))
			reflect.Copy(val1.Elem(), reflect.ValueOf(val0))
			params[i] = val1.Interface()
		case RPC_UINT64_PTR_ARRAY_PTR:
			val0 := readUInt64PtrSlice(bitstream)
			val2 := uint64(0)
			val1 := reflect.New(reflect.ArrayOf(len(val0), reflect.TypeOf(&val2)))
			reflect.Copy(val1.Elem(), reflect.ValueOf(val0))
			params[i] = val1.Interface()



		case RPC_MESSAGE://protobuf
			packet, err := UnmarshalPB(bitstream)
			if err == nil{
				params[i] = packet
			}
			/*nLen := bitstream.ReadInt(32)
			packetBuf := bitstream.ReadBits(nLen << 3)
			if pFuncType != nil{
				if i < pFuncType.NumIn() {
					val := reflect.New(pFuncType.In(i).Elem())
					err := proto.Unmarshal(packetBuf, val.Interface().(proto.Message))
					//packet := message.GetPakcetByName(funcName)
					//err := message.UnmarshalText(packet, packetBuf)
					if err == nil{
						params[i] = val.Interface()
					}
				}
			}*/



		case RPC_GOB://gob
			nLen := bitstream.ReadInt(32)
			packetBuf := bitstream.ReadBits(nLen << 3)

			if pFuncType != nil{
				if i < pFuncType.NumIn() {
					val := reflect.New(pFuncType.In(i))
					json := jsoniter.ConfigCompatibleWithStandardLibrary
					err := json.Unmarshal(packetBuf, val.Interface())
					/*buf := bytes.NewBuffer(packetBuf)
					enc := gob.NewDecoder(buf)
					err := enc.DecodeValue(val)*/
					if err == nil{
						params[i] = val.Elem().Interface()
					}
				}
			}

		default:
			panic("func [%s] params type not supported")
		}
	}
	return params
}

//rpc  UnmarshalPB
func UnmarshalPB(bitstream *base.BitStream) (proto.Message, error) {
	packetName := bitstream.ReadString()
	nLen := bitstream.ReadInt(32)
	packetBuf := bitstream.ReadBits(nLen << 3)
	packet := reflect.New(proto.MessageType(packetName).Elem()).Interface().(proto.Message)
	err := proto.Unmarshal(packetBuf, packet)
	return  packet, err
}

