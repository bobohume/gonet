package rpc

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/json-iterator/go"
	"gonet/base"
	"gonet/message"
	"reflect"
)

//rpc  Marshal
func writeBool(bitstream base.IBitStream, val bool)(){
	dat, _ := proto.Marshal(&message.Bool{Val:val})
	nLen := len(dat)
	bitstream.WriteInt(nLen, 8)
	bitstream.WriteBits(dat, nLen << 3)
}

func writeString(bitstream base.IBitStream, val string)(){
	dat, _ := proto.Marshal(&message.String{Val:val})
	nLen := len(dat)
	bitstream.WriteInt(nLen, 8)
	bitstream.WriteBits(dat, nLen << 3)
}

func writeFloat32(bitstream base.IBitStream, val float32)(){
	dat, _ := proto.Marshal(&message.Float{Val:val})
	nLen := len(dat)
	bitstream.WriteInt(nLen, 8)
	bitstream.WriteBits(dat, nLen << 3)
}

func writeFloat64(bitstream base.IBitStream, val float64)(){
	dat, _ := proto.Marshal(&message.Double{Val:val})
	nLen := len(dat)
	bitstream.WriteInt(nLen, 8)
	bitstream.WriteBits(dat, nLen << 3)
}

func writeInt32(bitstream base.IBitStream, val int32)(){
	dat, _ := proto.Marshal(&message.Int{Val:val})
	nLen := len(dat)
	bitstream.WriteInt(nLen, 8)
	bitstream.WriteBits(dat, nLen << 3)
}

func writeInt64(bitstream base.IBitStream, val int64)(){
	dat, _ := proto.Marshal(&message.Int64{Val:val})
	nLen := len(dat)
	bitstream.WriteInt(nLen, 8)
	bitstream.WriteBits(dat, nLen << 3)
}

func writeUInt32(bitstream base.IBitStream, val uint32)(){
	dat, _ := proto.Marshal(&message.UInt{Val:val})
	nLen := len(dat)
	bitstream.WriteInt(nLen, 8)
	bitstream.WriteBits(dat, nLen << 3)
}

func writeUInt64(bitstream base.IBitStream, val uint64)(){
	dat, _ := proto.Marshal(&message.UInt64{Val:val})
	nLen := len(dat)
	bitstream.WriteInt(nLen, 8)
	bitstream.WriteBits(dat, nLen << 3)
}

func writeBoolSlice(bitstream base.IBitStream, val []bool)(){
	dat, _ := proto.Marshal(&message.BoolSlice{Val:val})
	nLen := len(dat)
	bitstream.WriteInt(nLen, 16)
	bitstream.WriteBits(dat, nLen << 3)
}

func writeStringSlice(bitstream base.IBitStream, val []string)(){
	dat, _ := proto.Marshal(&message.StringSlice{Val:val})
	nLen := len(dat)
	bitstream.WriteInt(nLen, 16)
	bitstream.WriteBits(dat,  nLen << 3)
}

func writeFloat32Slice(bitstream base.IBitStream, val []float32)(){
	dat, _ := proto.Marshal(&message.FloatSlice{Val:val})
	nLen := len(dat)
	bitstream.WriteInt(nLen, 16)
	bitstream.WriteBits(dat,  nLen << 3)
}

func writeFloat64Slice(bitstream base.IBitStream, val []float64)(){
	dat, _ := proto.Marshal(&message.DoubleSlice{Val:val})
	nLen := len(dat)
	bitstream.WriteInt(nLen, 16)
	bitstream.WriteBits(dat,  nLen << 3)
}

func writeIntSlice(bitstream base.IBitStream, val []int)(){
	val0 := make([]int32, len(val))
	for i, v := range val{
		val0[i] = int32(v)
	}
	writeInt32Slice(bitstream, val0)
}

func writeInt8Slice(bitstream base.IBitStream, val []int8)(){
	val0 := make([]int32, len(val))
	for i, v := range val{
		val0[i] = int32(v)
	}
	writeInt32Slice(bitstream, val0)
}

func writeInt16Slice(bitstream base.IBitStream, val []int16)(){
	val0 := make([]int32, len(val))
	for i, v := range val{
		val0[i] = int32(v)
	}
	writeInt32Slice(bitstream, val0)
}

func writeInt32Slice(bitstream base.IBitStream, val []int32)(){
	dat, _ := proto.Marshal(&message.IntSlice{Val:val})
	nLen := len(dat)
	bitstream.WriteInt(nLen, 16)
	bitstream.WriteBits(dat,  nLen << 3)
}

func writeInt64Slice(bitstream base.IBitStream, val []int64)(){
	dat, _ := proto.Marshal(&message.Int64Slice{Val:val})
	nLen := len(dat)
	bitstream.WriteInt(nLen, 16)
	bitstream.WriteBits(dat,  nLen << 3)
}

func writeUIntSlice(bitstream base.IBitStream, val []uint)(){
	val0 := make([]uint32, len(val))
	for i, v := range val{
		val0[i] = uint32(v)
	}
	writeUInt32Slice(bitstream, val0)
}

func writeUInt8Slice(bitstream base.IBitStream, val []uint8)(){
	val0 := make([]uint32, len(val))
	for i, v := range val{
		val0[i] = uint32(v)
	}
	writeUInt32Slice(bitstream, val0)
}

func writeUInt16Slice(bitstream base.IBitStream, val []uint16)(){
	val0 := make([]uint32, len(val))
	for i, v := range val{
		val0[i] = uint32(v)
	}
	writeUInt32Slice(bitstream, val0)
}

func writeUInt32Slice(bitstream base.IBitStream, val []uint32)(){
	dat, _ := proto.Marshal(&message.UIntSlice{Val:val})
	nLen := len(dat)
	bitstream.WriteInt(nLen, 16)
	bitstream.WriteBits(dat, nLen << 3)
}

func writeUInt64Slice(bitstream base.IBitStream, val []uint64)(){
	dat, _ := proto.Marshal(&message.UInt64Slice{Val:val})
	nLen := len(dat)
	bitstream.WriteInt(nLen, 16)
	bitstream.WriteBits(dat, nLen << 3)
}

func writeBoolPtrSlice(bitstream base.IBitStream, val []*bool)(){
	val0 := make([]bool, len(val))
	for i, v := range val{
		if v != nil{
			val0[i] = *v
		}
	}
	writeBoolSlice(bitstream, val0)
}

func writeStringPtrSlice(bitstream base.IBitStream, val []*string)(){
	val0 := make([]string, len(val))
	for i, v := range val{
		if v != nil{
			val0[i] = *v
		}
	}
	writeStringSlice(bitstream, val0)
}

func writeFloat32PtrSlice(bitstream base.IBitStream, val []*float32)(){
	val0 := make([]float32, len(val))
	for i, v := range val{
		if v != nil{
			val0[i] = *v
		}
	}
	writeFloat32Slice(bitstream, val0)
}

func writeFloat64PtrSlice(bitstream base.IBitStream, val []*float64)(){
	val0 := make([]float64, len(val))
	for i, v := range val{
		if v != nil{
			val0[i] = *v
		}
	}
	writeFloat64Slice(bitstream, val0)
}

func writeIntPtrSlice(bitstream base.IBitStream, val []*int)(){
	val0 := make([]int, len(val))
	for i, v := range val{
		if v != nil{
			val0[i] = *v
		}
	}
	writeIntSlice(bitstream, val0)
}

func writeInt8PtrSlice(bitstream base.IBitStream, val []*int8)(){
	val0 := make([]int8, len(val))
	for i, v := range val{
		if v != nil{
			val0[i] = *v
		}
	}
	writeInt8Slice(bitstream, val0)
}

func writeInt16PtrSlice(bitstream base.IBitStream, val []*int16)(){
	val0 := make([]int16, len(val))
	for i, v := range val{
		if v != nil{
			val0[i] = *v
		}
	}
	writeInt16Slice(bitstream, val0)
}

func writeInt32PtrSlice(bitstream base.IBitStream, val []*int32)(){
	val0 := make([]int32, len(val))
	for i, v := range val{
		if v != nil{
			val0[i] = *v
		}
	}
	writeInt32Slice(bitstream, val0)
}

func writeInt64PtrSlice(bitstream base.IBitStream, val []*int64)(){
	val0 := make([]int64, len(val))
	for i, v := range val{
		if v != nil{
			val0[i] = *v
		}
	}
	writeInt64Slice(bitstream, val0)
}

func writeUIntPtrSlice(bitstream base.IBitStream, val []*uint)(){
	val0 := make([]uint, len(val))
	for i, v := range val{
		if v != nil{
			val0[i] = *v
		}
	}
	writeUIntSlice(bitstream, val0)
}

func writeUInt8PtrSlice(bitstream base.IBitStream, val []*uint8)(){
	val0 := make([]uint8, len(val))
	for i, v := range val{
		if v != nil{
			val0[i] = *v
		}
	}
	writeUInt8Slice(bitstream, val0)
}

func writeUInt16PtrSlice(bitstream base.IBitStream, val []*uint16)(){
	val0 := make([]uint16, len(val))
	for i, v := range val{
		if v != nil{
			val0[i] = *v
		}
	}
	writeUInt16Slice(bitstream, val0)
}

func writeUInt32PtrSlice(bitstream base.IBitStream, val []*uint32)(){
	val0 := make([]uint32, len(val))
	for i, v := range val{
		if v != nil{
			val0[i] = *v
		}
	}
	writeUInt32Slice(bitstream, val0)
}

func writeUInt64PtrSlice(bitstream base.IBitStream, val []*uint64)(){
	val0 := make([]uint64, len(val))
	for i, v := range val{
		if v != nil{
			val0[i] = *v
		}
	}
	writeUInt64Slice(bitstream, val0)
}

func Marshal(funcName string, params ...interface{})[]byte {
	defer func() {
		if err := recover(); err != nil {
			base.TraceCode(err)
		}
	}()

	msg := make([]byte, 1024)
	bitstream := base.NewBitStream(msg, 1024)
	bitstream.WriteString(funcName)
	bitstream.WriteInt(len(params), 8)
	for _, param := range params {
		sType := getTypeString(param)

		switch sType {
		case "bool":
			bitstream.WriteInt(RPC_BOOL, 8)
			writeBool(bitstream, param.(bool))
		case "string":
			bitstream.WriteInt(RPC_STRING, 8)
			writeString(bitstream, string(param.(string)))
		case "float32":
			bitstream.WriteInt(RPC_FLOAT32, 8)
			writeFloat32(bitstream, param.(float32))
		case "float64":
			bitstream.WriteInt(RPC_FLOAT64, 8)
			writeFloat64(bitstream, param.(float64))
		case "int":
			bitstream.WriteInt(RPC_INT, 8)
			writeInt32(bitstream, int32(param.(int)))
		case "int8":
			bitstream.WriteInt(RPC_INT8, 8)
			writeInt32(bitstream, int32(param.(int8)))
		case "int16":
			bitstream.WriteInt(RPC_INT16, 8)
			writeInt32(bitstream, int32(param.(int16)))
		case "int32":
			bitstream.WriteInt(RPC_INT32, 8)
			writeInt32(bitstream, int32(param.(int32)))
		case "int64":
			bitstream.WriteInt(RPC_INT64, 8)
			writeInt64(bitstream, int64(param.(int64)))
		case "uint":
			bitstream.WriteInt(RPC_UINT, 8)
			writeUInt32(bitstream, uint32(param.(uint)))
		case "uint8":
			bitstream.WriteInt(RPC_UINT8, 8)
			writeUInt32(bitstream, uint32(param.(uint8)))
		case "uint16":
			bitstream.WriteInt(RPC_UINT16, 8)
			writeUInt32(bitstream, uint32(param.(uint16)))
		case "uint32":
			bitstream.WriteInt(RPC_UINT32, 8)
			writeUInt32(bitstream, uint32(param.(uint32)))
		case "uint64":
			bitstream.WriteInt(RPC_UINT64, 8)
			writeUInt64(bitstream, uint64(param.(uint64)))



		case "[]bool":
			bitstream.WriteInt(RPC_BOOL_SLICE, 8)
			writeBoolSlice(bitstream, param.([]bool))
		case "[]string":
			bitstream.WriteInt(RPC_STRING_SLICE, 8)
			writeStringSlice(bitstream, param.([]string))
		case "[]float32":
			bitstream.WriteInt(RPC_FLOAT32_SLICE, 8)
			writeFloat32Slice(bitstream, param.([]float32))
		case "[]float64":
			bitstream.WriteInt(RPC_FLOAT64_SLICE, 8)
			writeFloat64Slice(bitstream, param.([]float64))
		case "[]int":
			bitstream.WriteInt(RPC_INT_SLICE, 8)
			writeIntSlice(bitstream, param.([]int))
		case "[]int8":
			bitstream.WriteInt(RPC_INT8_SLICE, 8)
			writeInt8Slice(bitstream, param.([]int8))
		case "[]int16":
			bitstream.WriteInt(RPC_INT16_SLICE, 8)
			writeInt16Slice(bitstream, param.([]int16))
		case "[]int32":
			bitstream.WriteInt(RPC_INT32_SLICE, 8)
			writeInt32Slice(bitstream, param.([]int32))
		case "[]int64":
			bitstream.WriteInt(RPC_INT64_SLICE, 8)
			writeInt64Slice(bitstream, param.([]int64))
		case "[]uint":
			bitstream.WriteInt(RPC_UINT8_SLICE, 8)
			writeUIntSlice(bitstream, param.([]uint))
		case "[]uint8":
			bitstream.WriteInt(RPC_UINT8_SLICE, 8)
			writeUInt8Slice(bitstream, param.([]uint8))
		case "[]uint16":
			bitstream.WriteInt(RPC_UINT8_SLICE, 8)
			writeUInt16Slice(bitstream, param.([]uint16))
		case "[]uint32":
			bitstream.WriteInt(RPC_UINT32_SLICE, 8)
			writeUInt32Slice(bitstream, param.([]uint32))
		case "[]uint64":
			bitstream.WriteInt(RPC_UINT64_SLICE, 8)
			writeUInt64Slice(bitstream, param.([]uint64))



		case "[*]bool":
			bitstream.WriteInt(RPC_BOOL_ARRAY, 8)
			val0 := reflect.ValueOf(param)
			val1 := reflect.MakeSlice(reflect.TypeOf([]bool{}), val0.Len(), val0.Len())
			reflect.Copy(val1, val0)
			writeBoolSlice(bitstream, val1.Interface().([]bool))
		case "[*]string":
			bitstream.WriteInt(RPC_STRING_ARRAY, 8)
			val0 := reflect.ValueOf(param)
			val1 := reflect.MakeSlice(reflect.TypeOf([]string{}), val0.Len(), val0.Len())
			reflect.Copy(val1, val0)
			writeStringSlice(bitstream, val1.Interface().([]string))
		case "[*]float32":
			bitstream.WriteInt(RPC_FLOAT32_ARRAY, 8)
			val0 := reflect.ValueOf(param)
			val1 := reflect.MakeSlice(reflect.TypeOf([]float32{}), val0.Len(), val0.Len())
			reflect.Copy(val1, val0)
			writeFloat32Slice(bitstream, val1.Interface().([]float32))
		case "[*]float64":
			bitstream.WriteInt(RPC_FLOAT64_ARRAY, 8)
			val0 := reflect.ValueOf(param)
			val1 := reflect.MakeSlice(reflect.TypeOf([]float64{}), val0.Len(), val0.Len())
			reflect.Copy(val1, val0)
			writeFloat64Slice(bitstream, val1.Interface().([]float64))
		case "[*]int":
			bitstream.WriteInt(RPC_INT_ARRAY, 8)
			val0 := reflect.ValueOf(param)
			val1 := reflect.MakeSlice(reflect.TypeOf([]int{}), val0.Len(), val0.Len())
			reflect.Copy(val1, val0)
			writeIntSlice(bitstream, val1.Interface().([]int))
		case "[*]int8":
			bitstream.WriteInt(RPC_INT8_ARRAY, 8)
			val0 := reflect.ValueOf(param)
			val1 := reflect.MakeSlice(reflect.TypeOf([]int8{}), val0.Len(), val0.Len())
			reflect.Copy(val1, val0)
			writeInt8Slice(bitstream, val1.Interface().([]int8))
		case "[*]int16":
			bitstream.WriteInt(RPC_INT16_ARRAY, 8)
			val0 := reflect.ValueOf(param)
			val1 := reflect.MakeSlice(reflect.TypeOf([]int16{}), val0.Len(), val0.Len())
			reflect.Copy(val1, val0)
			writeInt16Slice(bitstream, val1.Interface().([]int16))
		case "[*]int32":
			bitstream.WriteInt(RPC_INT32_ARRAY, 8)
			val0 := reflect.ValueOf(param)
			val1 := reflect.MakeSlice(reflect.TypeOf([]int32{}), val0.Len(), val0.Len())
			reflect.Copy(val1, val0)
			writeInt32Slice(bitstream, val1.Interface().([]int32))
		case "[*]int64":
			bitstream.WriteInt(RPC_INT64_ARRAY, 8)
			val0 := reflect.ValueOf(param)
			val1 := reflect.MakeSlice(reflect.TypeOf([]int64{}), val0.Len(), val0.Len())
			reflect.Copy(val1, val0)
			writeInt64Slice(bitstream, val1.Interface().([]int64))
		case "[*]uint":
			bitstream.WriteInt(RPC_UINT_ARRAY, 8)
			val0 := reflect.ValueOf(param)
			val1 := reflect.MakeSlice(reflect.TypeOf([]uint{}), val0.Len(), val0.Len())
			reflect.Copy(val1, val0)
			writeUIntSlice(bitstream, val1.Interface().([]uint))
		case "[*]uint8":
			bitstream.WriteInt(RPC_UINT8_ARRAY, 8)
			val0 := reflect.ValueOf(param)
			val1 := reflect.MakeSlice(reflect.TypeOf([]uint8{}), val0.Len(), val0.Len())
			reflect.Copy(val1, val0)
			writeUInt8Slice(bitstream, val1.Interface().([]uint8))
		case "[*]uint16":
			bitstream.WriteInt(RPC_UINT16_ARRAY, 8)
			val0 := reflect.ValueOf(param)
			val1 := reflect.MakeSlice(reflect.TypeOf([]uint16{}), val0.Len(), val0.Len())
			reflect.Copy(val1, val0)
			writeUInt16Slice(bitstream, val1.Interface().([]uint16))
		case "[*]uint32":
			bitstream.WriteInt(RPC_UINT32_ARRAY, 8)
			val0 := reflect.ValueOf(param)
			val1 := reflect.MakeSlice(reflect.TypeOf([]uint32{}), val0.Len(), val0.Len())
			reflect.Copy(val1, val0)
			writeUInt32Slice(bitstream, val1.Interface().([]uint32))
		case "[*]uint64":
			bitstream.WriteInt(RPC_UINT64_ARRAY, 8)
			val0 := reflect.ValueOf(param)
			val1 := reflect.MakeSlice(reflect.TypeOf([]uint64{}), val0.Len(), val0.Len())
			reflect.Copy(val1, val0)
			writeUInt64Slice(bitstream, val1.Interface().([]uint64))



		case "*bool":
			bitstream.WriteInt(RPC_BOOL_PTR, 8)
			if param.(*bool) != nil{
				writeBool(bitstream, *param.(*bool))
			}else{
				writeBool(bitstream, false)
			}
		case "*string":
			bitstream.WriteInt(RPC_STRING_PTR, 8)
			if param.(*string) != nil {
				writeString(bitstream, *param.(*string))
			}else{
				writeString(bitstream, "")
			}
		case "*float32":
			bitstream.WriteInt(RPC_FLOAT32_PTR, 8)
			if param.(*float32) != nil {
				writeFloat32(bitstream, *param.(*float32))
			}else{
				writeFloat32(bitstream, 0)
			}
		case "*float64":
			bitstream.WriteInt(RPC_FLOAT64_PTR, 8)
			if param.(*float64) != nil {
				writeFloat64(bitstream, *param.(*float64))
			}else{
				writeFloat64(bitstream, 0)
			}
		case "*int":
			bitstream.WriteInt(RPC_INT_PTR, 8)
			if param.(*int) != nil {
				writeInt32(bitstream, int32(*param.(*int)))
			}else{
				writeInt32(bitstream, 0)
			}
		case "*int8":
			bitstream.WriteInt(RPC_INT8_PTR, 8)
			if param.(*int8) != nil {
				writeInt32(bitstream, int32(*param.(*int8)))
			}else{
				writeInt32(bitstream, 0)
			}
		case "*int16":
			bitstream.WriteInt(RPC_INT16_PTR, 8)
			if param.(*int16) != nil {
				writeInt32(bitstream, int32(*param.(*int16)))
			}else{
				writeInt32(bitstream, 0)
			}
		case "*int32":
			bitstream.WriteInt(RPC_INT32_PTR, 8)
			if param.(*int32) != nil {
				writeInt32(bitstream, int32(*param.(*int32)))
			}else{
				writeInt32(bitstream, 0)
			}
		case "*int64":
			bitstream.WriteInt(RPC_INT64_PTR, 8)
			if param.(*int64) != nil {
				writeInt64(bitstream, int64(*param.(*int64)))
			}else{
				writeInt64(bitstream, 0)
			}
		case "*uint":
			bitstream.WriteInt(RPC_UINT_PTR, 8)
			if param.(*uint) != nil {
				writeUInt32(bitstream, uint32(*param.(*uint)))
			}else{
				writeUInt32(bitstream, 0)
			}
		case "*uint8":
			bitstream.WriteInt(RPC_UINT8_PTR, 8)
			if param.(*uint8) != nil {
				writeUInt32(bitstream, uint32(*param.(*uint8)))
			}else{
				writeUInt32(bitstream, 0)
			}
		case "*uint16":
			bitstream.WriteInt(RPC_UINT16_PTR, 8)
			if param.(*uint16) != nil {
				writeUInt32(bitstream, uint32(*param.(*uint16)))
			}else{
				writeUInt32(bitstream, 0)
			}
		case "*uint32":
			bitstream.WriteInt(RPC_UINT32_PTR, 8)
			if param.(*uint32) != nil {
				writeUInt32(bitstream, uint32(*param.(*uint32)))
			}else{
				writeUInt32(bitstream, 0)
			}
		case "*uint64":
			bitstream.WriteInt(RPC_UINT64_PTR, 8)
			if param.(*uint64) != nil {
				writeUInt64(bitstream, uint64(*param.(*uint64)))
			}else{
				writeUInt64(bitstream, 0)
			}



		case "[]*bool":
			bitstream.WriteInt(RPC_BOOL_PTR_SLICE, 8)
			writeBoolPtrSlice(bitstream, param.([]*bool))
		case "[]*string":
			bitstream.WriteInt(RPC_STRING_PTR_SLICE, 8)
			writeStringPtrSlice(bitstream, param.([]*string))
		case "[]*float32":
			bitstream.WriteInt(RPC_FLOAT32_PTR_SLICE, 8)
			writeFloat32PtrSlice(bitstream, param.([]*float32))
		case "[]*float64":
			bitstream.WriteInt(RPC_FLOAT64_PTR_SLICE, 8)
			writeFloat64PtrSlice(bitstream, param.([]*float64))
		case "[]*int":
			bitstream.WriteInt(RPC_INT_PTR_SLICE, 8)
			writeIntPtrSlice(bitstream, param.([]*int))
		case "[]*int8":
			bitstream.WriteInt(RPC_INT8_PTR_SLICE, 8)
			writeInt8PtrSlice(bitstream, param.([]*int8))
		case "[]*int16":
			bitstream.WriteInt(RPC_INT16_PTR_SLICE, 8)
			writeInt16PtrSlice(bitstream, param.([]*int16))
		case "[]*int32":
			bitstream.WriteInt(RPC_INT32_PTR_SLICE, 8)
			writeInt32PtrSlice(bitstream, param.([]*int32))
		case "[]*int64":
			bitstream.WriteInt(RPC_INT64_PTR_SLICE, 8)
			writeInt64PtrSlice(bitstream, param.([]*int64))
		case "[]*uint":
			bitstream.WriteInt(RPC_UINT_PTR_SLICE, 8)
			writeUIntPtrSlice(bitstream, param.([]*uint))
		case "[]*uint8":
			bitstream.WriteInt(RPC_UINT8_PTR_SLICE, 8)
			writeUInt8PtrSlice(bitstream, param.([]*uint8))
		case "[]*uint16":
			bitstream.WriteInt(RPC_UINT16_PTR_SLICE, 8)
			writeUInt16PtrSlice(bitstream, param.([]*uint16))
		case "[]*uint32":
			bitstream.WriteInt(RPC_UINT32_PTR_SLICE, 8)
			writeUInt32PtrSlice(bitstream, param.([]*uint32))
		case "[]*uint64":
			bitstream.WriteInt(RPC_UINT64_PTR_SLICE, 8)
			writeUInt64PtrSlice(bitstream, param.([]*uint64))



		case "[*]*bool":
			bitstream.WriteInt(RPC_BOOL_PTR_ARRAY, 8)
			val0 := reflect.ValueOf(param)
			val1 := reflect.MakeSlice(reflect.TypeOf([]*bool{}), val0.Len(), val0.Len())
			reflect.Copy(val1, val0)
			writeBoolPtrSlice(bitstream, val1.Interface().([]*bool))
		case "[*]*string":
			bitstream.WriteInt(RPC_STRING_PTR_ARRAY, 8)
			val0 := reflect.ValueOf(param)
			val1 := reflect.MakeSlice(reflect.TypeOf([]*string{}), val0.Len(), val0.Len())
			reflect.Copy(val1, val0)
			writeStringPtrSlice(bitstream, val1.Interface().([]*string))
		case "[*]*float32":
			bitstream.WriteInt(RPC_FLOAT32_PTR_ARRAY, 8)
			val0 := reflect.ValueOf(param)
			val1 := reflect.MakeSlice(reflect.TypeOf([]*float32{}), val0.Len(), val0.Len())
			reflect.Copy(val1, val0)
			writeFloat32PtrSlice(bitstream, val1.Interface().([]*float32))
		case "[*]*float64":
			bitstream.WriteInt(RPC_FLOAT64_PTR_ARRAY, 8)
			val0 := reflect.ValueOf(param)
			val1 := reflect.MakeSlice(reflect.TypeOf([]*float64{}), val0.Len(), val0.Len())
			reflect.Copy(val1, val0)
			writeFloat64PtrSlice(bitstream, val1.Interface().([]*float64))
		case "[*]*int":
			bitstream.WriteInt(RPC_INT_PTR_ARRAY, 8)
			val0 := reflect.ValueOf(param)
			val1 := reflect.MakeSlice(reflect.TypeOf([]*int{}), val0.Len(), val0.Len())
			reflect.Copy(val1, val0)
			writeIntPtrSlice(bitstream, val1.Interface().([]*int))
		case "[*]*int8":
			bitstream.WriteInt(RPC_INT8_PTR_ARRAY, 8)
			val0 := reflect.ValueOf(param)
			val1 := reflect.MakeSlice(reflect.TypeOf([]*int8{}), val0.Len(), val0.Len())
			reflect.Copy(val1, val0)
			writeInt8PtrSlice(bitstream, val1.Interface().([]*int8))
		case "[*]*int16":
			bitstream.WriteInt(RPC_INT16_PTR_ARRAY, 8)
			val0 := reflect.ValueOf(param)
			val1 := reflect.MakeSlice(reflect.TypeOf([]*int16{}), val0.Len(), val0.Len())
			reflect.Copy(val1, val0)
			writeInt16PtrSlice(bitstream, val1.Interface().([]*int16))
		case "[*]*int32":
			bitstream.WriteInt(RPC_INT32_PTR_ARRAY, 8)
			val0 := reflect.ValueOf(param)
			val1 := reflect.MakeSlice(reflect.TypeOf([]*int32{}), val0.Len(), val0.Len())
			reflect.Copy(val1, val0)
			writeInt32PtrSlice(bitstream, val1.Interface().([]*int32))
		case "[*]*int64":
			bitstream.WriteInt(RPC_INT64_PTR_ARRAY, 8)
			val0 := reflect.ValueOf(param)
			val1 := reflect.MakeSlice(reflect.TypeOf([]*int64{}), val0.Len(), val0.Len())
			reflect.Copy(val1, val0)
			writeInt64PtrSlice(bitstream, val1.Interface().([]*int64))
		case "[*]*uint":
			bitstream.WriteInt(RPC_UINT_PTR_ARRAY, 8)
			val0 := reflect.ValueOf(param)
			val1 := reflect.MakeSlice(reflect.TypeOf([]*uint{}), val0.Len(), val0.Len())
			reflect.Copy(val1, val0)
			writeUIntPtrSlice(bitstream, val1.Interface().([]*uint))
		case "[*]*uint8":
			bitstream.WriteInt(RPC_UINT8_PTR_ARRAY, 8)
			val0 := reflect.ValueOf(param)
			val1 := reflect.MakeSlice(reflect.TypeOf([]*uint8{}), val0.Len(), val0.Len())
			reflect.Copy(val1, val0)
			writeUInt8PtrSlice(bitstream, val1.Interface().([]*uint8))
		case "[*]*uint16":
			bitstream.WriteInt(RPC_UINT16_PTR_ARRAY, 8)
			val0 := reflect.ValueOf(param)
			val1 := reflect.MakeSlice(reflect.TypeOf([]*uint16{}), val0.Len(), val0.Len())
			reflect.Copy(val1, val0)
			writeUInt16PtrSlice(bitstream, val1.Interface().([]*uint16))
		case "[*]*uint32":
			bitstream.WriteInt(RPC_UINT32_PTR_ARRAY, 8)
			val0 := reflect.ValueOf(param)
			val1 := reflect.MakeSlice(reflect.TypeOf([]*uint32{}), val0.Len(), val0.Len())
			reflect.Copy(val1, val0)
			writeUInt32PtrSlice(bitstream, val1.Interface().([]*uint32))
		case "[*]*uint64":
			bitstream.WriteInt(RPC_UINT64_PTR_ARRAY, 8)
			val0 := reflect.ValueOf(param)
			val1 := reflect.MakeSlice(reflect.TypeOf([]*uint64{}), val0.Len(), val0.Len())
			reflect.Copy(val1, val0)
			writeUInt64PtrSlice(bitstream, val1.Interface().([]*uint64))



		case "*[]bool":
			bitstream.WriteInt(RPC_BOOL_SLICE_PTR, 8)
			writeBoolSlice(bitstream, *param.(*[]bool))
		case "*[]string":
			bitstream.WriteInt(RPC_STRING_SLICE_PTR, 8)
			writeStringSlice(bitstream, *param.(*[]string))
		case "*[]float32":
			bitstream.WriteInt(RPC_FLOAT32_SLICE_PTR, 8)
			writeFloat32Slice(bitstream, *param.(*[]float32))
		case "*[]float64":
			bitstream.WriteInt(RPC_FLOAT64_SLICE_PTR, 8)
			writeFloat64Slice(bitstream, *param.(*[]float64))
		case "*[]int":
			bitstream.WriteInt(RPC_INT_SLICE_PTR, 8)
			writeIntSlice(bitstream, *param.(*[]int))
		case "*[]int8":
			bitstream.WriteInt(RPC_INT8_SLICE_PTR, 8)
			writeInt8Slice(bitstream, *param.(*[]int8))
		case "*[]int16":
			bitstream.WriteInt(RPC_INT16_SLICE_PTR, 8)
			writeInt16Slice(bitstream, *param.(*[]int16))
		case "*[]int32":
			bitstream.WriteInt(RPC_INT32_SLICE_PTR, 8)
			writeInt32Slice(bitstream, *param.(*[]int32))
		case "*[]int64":
			bitstream.WriteInt(RPC_INT64_SLICE_PTR, 8)
			writeInt64Slice(bitstream, *param.(*[]int64))
		case "*[]uint":
			bitstream.WriteInt(RPC_UINT8_SLICE_PTR, 8)
			writeUIntSlice(bitstream, *param.(*[]uint))
		case "*[]uint8":
			bitstream.WriteInt(RPC_UINT8_SLICE_PTR, 8)
			writeUInt8Slice(bitstream, *param.(*[]uint8))
		case "*[]uint16":
			bitstream.WriteInt(RPC_UINT8_SLICE_PTR, 8)
			writeUInt16Slice(bitstream, *param.(*[]uint16))
		case "*[]uint32":
			bitstream.WriteInt(RPC_UINT32_SLICE_PTR, 8)
			writeUInt32Slice(bitstream, *param.(*[]uint32))
		case "*[]uint64":
			bitstream.WriteInt(RPC_UINT64_SLICE_PTR, 8)
			writeUInt64Slice(bitstream, *param.(*[]uint64))



		case "*[]*bool":
			bitstream.WriteInt(RPC_BOOL_PTR_SLICE_PTR, 8)
			writeBoolPtrSlice(bitstream, *param.(*[]*bool))
		case "*[]*string":
			bitstream.WriteInt(RPC_STRING_PTR_SLICE_PTR, 8)
			writeStringPtrSlice(bitstream, *param.(*[]*string))
		case "*[]*float32":
			bitstream.WriteInt(RPC_FLOAT32_PTR_SLICE_PTR, 8)
			writeFloat32PtrSlice(bitstream, *param.(*[]*float32))
		case "*[]*float64":
			bitstream.WriteInt(RPC_FLOAT64_PTR_SLICE_PTR, 8)
			writeFloat64PtrSlice(bitstream, *param.(*[]*float64))
		case "*[]*int":
			bitstream.WriteInt(RPC_INT_PTR_SLICE_PTR, 8)
			writeIntPtrSlice(bitstream, *param.(*[]*int))
		case "*[]*int8":
			bitstream.WriteInt(RPC_INT8_PTR_SLICE_PTR, 8)
			writeInt8PtrSlice(bitstream, *param.(*[]*int8))
		case "*[]*int16":
			bitstream.WriteInt(RPC_INT16_PTR_SLICE_PTR, 8)
			writeInt16PtrSlice(bitstream, *param.(*[]*int16))
		case "*[]*int32":
			bitstream.WriteInt(RPC_INT32_PTR_SLICE_PTR, 8)
			writeInt32PtrSlice(bitstream, *param.(*[]*int32))
		case "*[]*int64":
			bitstream.WriteInt(RPC_INT64_PTR_SLICE_PTR, 8)
			writeInt64PtrSlice(bitstream, *param.(*[]*int64))
		case "*[]*uint":
			bitstream.WriteInt(RPC_UINT_PTR_SLICE_PTR, 8)
			writeUIntPtrSlice(bitstream, *param.(*[]*uint))
		case "*[]*uint8":
			bitstream.WriteInt(RPC_UINT8_PTR_SLICE_PTR, 8)
			writeUInt8PtrSlice(bitstream, *param.(*[]*uint8))
		case "*[]*uint16":
			bitstream.WriteInt(RPC_UINT16_PTR_SLICE_PTR, 8)
			writeUInt16PtrSlice(bitstream, *param.(*[]*uint16))
		case "*[]*uint32":
			bitstream.WriteInt(RPC_UINT32_PTR_SLICE_PTR, 8)
			writeUInt32PtrSlice(bitstream, *param.(*[]*uint32))
		case "*[]*uint64":
			bitstream.WriteInt(RPC_UINT64_PTR_SLICE_PTR, 8)
			writeUInt64PtrSlice(bitstream, *param.(*[]*uint64))



		case "*[*]bool":
			bitstream.WriteInt(RPC_BOOL_ARRAY_PTR, 8)
			val0 := reflect.ValueOf(param)
			val1 := val0.Elem().Slice(0, val0.Elem().Len()).Interface().([]bool)
			writeBoolSlice(bitstream, val1)
		case "*[*]string":
			bitstream.WriteInt(RPC_STRING_ARRAY_PTR, 8)
			val0 := reflect.ValueOf(param)
			val1 := val0.Elem().Slice(0, val0.Elem().Len()).Interface().([]string)
			writeStringSlice(bitstream, val1)
		case "*[*]float32":
			bitstream.WriteInt(RPC_FLOAT32_ARRAY_PTR, 8)
			val0 := reflect.ValueOf(param)
			val1 := val0.Elem().Slice(0, val0.Elem().Len()).Interface().([]float32)
			writeFloat32Slice(bitstream, val1)
		case "*[*]float64":
			bitstream.WriteInt(RPC_FLOAT64_ARRAY_PTR, 8)
			val0 := reflect.ValueOf(param)
			val1 := val0.Elem().Slice(0, val0.Elem().Len()).Interface().([]float64)
			writeFloat64Slice(bitstream, val1)
		case "*[*]int":
			bitstream.WriteInt(RPC_INT_ARRAY_PTR, 8)
			val0 := reflect.ValueOf(param)
			val1 := val0.Elem().Slice(0, val0.Elem().Len()).Interface().([]int)
			writeIntSlice(bitstream, val1)
		case "*[*]int8":
			bitstream.WriteInt(RPC_INT8_ARRAY_PTR, 8)
			val0 := reflect.ValueOf(param)
			val1 := val0.Elem().Slice(0, val0.Elem().Len()).Interface().([]int8)
			writeInt8Slice(bitstream, val1)
		case "*[*]int16":
			bitstream.WriteInt(RPC_INT16_ARRAY_PTR, 8)
			val0 := reflect.ValueOf(param)
			val1 := val0.Elem().Slice(0, val0.Elem().Len()).Interface().([]int16)
			writeInt16Slice(bitstream, val1)
		case "*[*]int32":
			bitstream.WriteInt(RPC_INT32_ARRAY_PTR, 8)
			val0 := reflect.ValueOf(param)
			val1 := val0.Elem().Slice(0, val0.Elem().Len()).Interface().([]int32)
			writeInt32Slice(bitstream, val1)
		case "*[*]int64":
			bitstream.WriteInt(RPC_INT64_ARRAY_PTR, 8)
			val0 := reflect.ValueOf(param)
			val1 := val0.Elem().Slice(0, val0.Elem().Len()).Interface().([]int64)
			writeInt64Slice(bitstream, val1)
		case "*[*]uint":
			bitstream.WriteInt(RPC_UINT_ARRAY_PTR, 8)
			val0 := reflect.ValueOf(param)
			val1 := val0.Elem().Slice(0, val0.Elem().Len()).Interface().([]uint)
			writeUIntSlice(bitstream, val1)
		case "*[*]uint8":
			bitstream.WriteInt(RPC_UINT8_ARRAY_PTR, 8)
			val0 := reflect.ValueOf(param)
			val1 := val0.Elem().Slice(0, val0.Elem().Len()).Interface().([]uint8)
			writeUInt8Slice(bitstream, val1)
		case "*[*]uint16":
			bitstream.WriteInt(RPC_UINT16_ARRAY_PTR, 8)
			val0 := reflect.ValueOf(param)
			val1 := val0.Elem().Slice(0, val0.Elem().Len()).Interface().([]uint16)
			writeUInt16Slice(bitstream, val1)
		case "*[*]uint32":
			bitstream.WriteInt(RPC_UINT32_ARRAY_PTR, 8)
			val0 := reflect.ValueOf(param)
			val1 := val0.Elem().Slice(0, val0.Elem().Len()).Interface().([]uint32)
			writeUInt32Slice(bitstream, val1)
		case "*[*]uint64":
			bitstream.WriteInt(RPC_UINT64_ARRAY_PTR, 8)
			val0 := reflect.ValueOf(param)
			val1 := val0.Elem().Slice(0, val0.Elem().Len()).Interface().([]uint64)
			writeUInt64Slice(bitstream, val1)



		case "*[*]*bool":
			bitstream.WriteInt(RPC_BOOL_PTR_ARRAY_PTR, 8)
			val0 := reflect.ValueOf(param)
			val1 := val0.Elem().Slice(0, val0.Elem().Len()).Interface().([]*bool)
			writeBoolPtrSlice(bitstream, val1)
		case "*[*]*string":
			bitstream.WriteInt(RPC_STRING_PTR_ARRAY_PTR, 8)
			val0 := reflect.ValueOf(param)
			val1 := val0.Elem().Slice(0, val0.Elem().Len()).Interface().([]*string)
			writeStringPtrSlice(bitstream, val1)
		case "*[*]*float32":
			bitstream.WriteInt(RPC_FLOAT32_PTR_ARRAY_PTR, 8)
			val0 := reflect.ValueOf(param)
			val1 := val0.Elem().Slice(0, val0.Elem().Len()).Interface().([]*float32)
			writeFloat32PtrSlice(bitstream, val1)
		case "*[*]*float64":
			bitstream.WriteInt(RPC_FLOAT64_PTR_ARRAY_PTR, 8)
			val0 := reflect.ValueOf(param)
			val1 := val0.Elem().Slice(0, val0.Elem().Len()).Interface().([]*float64)
			writeFloat64PtrSlice(bitstream, val1)
		case "*[*]*int":
			bitstream.WriteInt(RPC_INT_PTR_ARRAY_PTR, 8)
			val0 := reflect.ValueOf(param)
			val1 := val0.Elem().Slice(0, val0.Elem().Len()).Interface().([]*int)
			writeIntPtrSlice(bitstream, val1)
		case "*[*]*int8":
			bitstream.WriteInt(RPC_INT8_PTR_ARRAY_PTR, 8)
			val0 := reflect.ValueOf(param)
			val1 := val0.Elem().Slice(0, val0.Elem().Len()).Interface().([]*int8)
			writeInt8PtrSlice(bitstream, val1)
		case "*[*]*int16":
			bitstream.WriteInt(RPC_INT16_PTR_ARRAY_PTR, 8)
			val0 := reflect.ValueOf(param)
			val1 := val0.Elem().Slice(0, val0.Elem().Len()).Interface().([]*int16)
			writeInt16PtrSlice(bitstream, val1)
		case "*[*]*int32":
			bitstream.WriteInt(RPC_INT32_PTR_ARRAY_PTR, 8)
			val0 := reflect.ValueOf(param)
			val1 := val0.Elem().Slice(0, val0.Elem().Len()).Interface().([]*int32)
			writeInt32PtrSlice(bitstream, val1)
		case "*[*]*int64":
			bitstream.WriteInt(RPC_INT64_PTR_ARRAY_PTR, 8)
			val0 := reflect.ValueOf(param)
			val1 := val0.Elem().Slice(0, val0.Elem().Len()).Interface().([]*int64)
			writeInt64PtrSlice(bitstream, val1)
		case "*[*]*uint":
			bitstream.WriteInt(RPC_UINT_PTR_ARRAY_PTR, 8)
			val0 := reflect.ValueOf(param)
			val1 := val0.Elem().Slice(0, val0.Elem().Len()).Interface().([]*uint)
			writeUIntPtrSlice(bitstream, val1)
		case "*[*]*uint8":
			bitstream.WriteInt(RPC_UINT8_PTR_ARRAY_PTR, 8)
			val0 := reflect.ValueOf(param)
			val1 := val0.Elem().Slice(0, val0.Elem().Len()).Interface().([]*uint8)
			writeUInt8PtrSlice(bitstream, val1)
		case "*[*]*uint16":
			bitstream.WriteInt(RPC_UINT16_PTR_ARRAY_PTR, 8)
			val0 := reflect.ValueOf(param)
			val1 := val0.Elem().Slice(0, val0.Elem().Len()).Interface().([]*uint16)
			writeUInt16PtrSlice(bitstream, val1)
		case "*[*]*uint32":
			bitstream.WriteInt(RPC_UINT32_PTR_ARRAY_PTR, 8)
			val0 := reflect.ValueOf(param)
			val1 := val0.Elem().Slice(0, val0.Elem().Len()).Interface().([]*uint32)
			writeUInt32PtrSlice(bitstream, val1)
		case "*[*]*uint64":
			bitstream.WriteInt(RPC_UINT64_PTR_ARRAY_PTR, 8)
			val0 := reflect.ValueOf(param)
			val1 := val0.Elem().Slice(0, val0.Elem().Len()).Interface().([]*uint64)
			writeUInt64PtrSlice(bitstream, val1)



		case "*message":
			bitstream.WriteInt(RPC_MESSAGE, 8)
			bitstream.WriteString(proto.MessageName(param.(proto.Message)))
			buf, _ :=proto.Marshal(param.(proto.Message))
			nLen := len(buf)
			bitstream.WriteInt(nLen, 32)
			bitstream.WriteBits(buf, nLen << 3)



		case "*gob":
			bitstream.WriteInt(RPC_GOB, 8)
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			buf, _ := json.Marshal(param)
			nLen := len(buf)
			bitstream.WriteInt(nLen, 32)
			bitstream.WriteBits(buf, nLen << 3)
			/*buf := &bytes.Buffer{}
			enc := gob.NewEncoder(buf)
			enc.Encode(param)
			nLen := buf.Len()
			bitstream.WriteInt(nLen, Bit32)
			bitstream.WriteBits(nLen << 3, buf.Bytes())*/

		default:
			fmt.Println("params type not supported", sType,  reflect.TypeOf(param))
			panic("params type not supported")
		}
	}

	return bitstream.GetBuffer()
}

//rpc  MarshalPB
/*func MarshalPB(packet proto.Message, bitstream *base.BitStream) bool {
	bitstream.WriteString(message.GetMessageName(packet))
	bitstream.WriteInt(1, 8)
	{
		sType := strings.ToLower(reflect.ValueOf(packet).Type().String())
		index := strings.Index(sType, ".")
		if index!= -1{
			sType = sType[:index]
		}
		switch sType {
		case "*message":
			bitstream.WriteInt(RPC_MESSAGE, 8)
			buf, _ :=proto.Marshal(packet)
			nLen := len(buf)
			bitstream.WriteInt(nLen, 32)
			bitstream.WriteBits(nLen << 3, buf)
		default:
			log.Panicln("packet params type not supported", packet, sType)
			return false
		}
	}
	return true
}*/