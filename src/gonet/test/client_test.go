package main_test

import (
	"encoding/json"
	"github.com/golang/protobuf/proto"
	"gonet/base"
	"gonet/message"
	"reflect"
	"strings"
	"testing"
	"unsafe"
)

type(
	TopRank struct{
		Value []int `sql:"name:value"				json:"value"	json:"value"`
	}
)

var(
	ntimes = 1000000
	nArraySize = 25
	nValue = 0x7fffffff
)

func TestJson(t *testing.T){
	data := &TopRank{}
	for i := 0; i < nArraySize; i++{
		data.Value = append(data.Value, nValue)
	}
	for i := 0; i < ntimes; i++{
		json.Marshal(data)
	}
}

func TestUJson(t *testing.T){
	data := &TopRank{}
	for i := 0; i < nArraySize; i++{
		data.Value = append(data.Value, nValue)
	}
	buff, _ := json.Marshal(data)
	for i := 0; i < ntimes; i++{
		json.Unmarshal(buff, &TopRank{})
	}
}

func TestPB(t *testing.T){
	aa := []int32{}
	for i := 0; i < nArraySize; i++{
		aa = append(aa, int32(nValue))
	}
	for i := 0; i < ntimes; i++{
		proto.Marshal(&message.W_C_Test{Recv:aa})
	}
}

func TestUPB(t *testing.T){
	aa := []int32{}
	for i := 0; i < nArraySize; i++{
		aa = append(aa, int32(nValue))
	}
	buff, _ := proto.Marshal(&message.W_C_Test{Recv:aa})
	for i := 0; i < ntimes; i++{
		proto.Unmarshal(buff, &message.W_C_Test{})
	}
}

func TestRpc(t *testing.T){
	aa := []int{}
	for i := 0; i < nArraySize; i++{
		aa = append(aa, nValue)
	}
	for i := 0; i < ntimes; i++{
		base.GetPacket("test", aa)
	}
}

func TestURpc(t *testing.T){
	aa := []int{}
	for i := 0; i < nArraySize; i++{
		aa = append(aa, nValue)
	}
	buff := base.GetPacket("test", aa)
	for i := 0; i < ntimes; i++{
		parse(buff)
	}
}

func parse (buff []byte) {
	funcName := ""
	bitstream := base.NewBitStream(buff, len(buff))
	funcName = bitstream.ReadString()
	funcName = strings.ToLower(funcName)
	if 1 != 0 {
		nCurLen := bitstream.ReadInt(8)
		params := make([]interface{}, nCurLen)
		for i := 0; i < nCurLen; i++  {
			switch bitstream.ReadInt(8) {
			case 1:
				params[i] = bitstream.ReadFlag()
			case 2:
				params[i] = bitstream.ReadFloat64()
			case 3:
				params[i] = bitstream.ReadFloat()
			case 4:
				params[i] = int8(bitstream.ReadInt(8))
			case 5:
				params[i] = uint8(bitstream.ReadInt(8))
			case 6:
				params[i] = int16(bitstream.ReadInt(16))
			case 7:
				params[i] = uint16(bitstream.ReadInt(16))
			case 8:
				params[i] = int32(bitstream.ReadInt(32))
			case 9:
				params[i] = uint32(bitstream.ReadInt(32))
			case 10:
				params[i] = int64(bitstream.ReadInt64(64))
			case 11:
				params[i] = uint64(bitstream.ReadInt64(64))
			case 12:
				params[i] = bitstream.ReadString()
			case 13:
				params[i] = bitstream.ReadInt(32)
			case 14:
				params[i] = uint(bitstream.ReadInt(32))


			case 21:
				nLen := bitstream.ReadInt(16)
				val := make([]bool, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = bitstream.ReadFlag()
				}
				params[i] = val
			case 22:
				nLen := bitstream.ReadInt(16)
				val := make([]float64, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = bitstream.ReadFloat64()
				}
				params[i] = val
			case 23:
				nLen := bitstream.ReadInt(16)
				val := make([]float32, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = bitstream.ReadFloat()
				}
				params[i] = val
			case 24:
				nLen := bitstream.ReadInt(16)
				val := make([]int8, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = int8(bitstream.ReadInt(8))
				}
				params[i] = val
			case 25:
				nLen := bitstream.ReadInt(16)
				val := make([]uint8, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = uint8(bitstream.ReadInt(8))
				}
				params[i] = val
			case 26:
				nLen := bitstream.ReadInt(16)
				val := make([]int16, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = int16(bitstream.ReadInt(16))
				}
				params[i] = val
			case 27:
				nLen := bitstream.ReadInt(16)
				val := make([]uint16, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = uint16(bitstream.ReadInt(16))
				}
				params[i] = val
			case 28:
				nLen := bitstream.ReadInt(16)
				val := make([]int32, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = int32(bitstream.ReadInt(32))
				}
				params[i] = val
			case 29:
				nLen := bitstream.ReadInt(16)
				val := make([]uint32, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = uint32(bitstream.ReadInt(32))
				}
				params[i] = val
			case 30:
				nLen := bitstream.ReadInt(16)
				val := make([]int64, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = int64(bitstream.ReadInt64(64))
				}
				params[i] = val
			case 31:
				nLen := bitstream.ReadInt(16)
				val := make([]uint64, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = uint64(bitstream.ReadInt64(64))
				}
				params[i] = val
			case 32:
				nLen := bitstream.ReadInt(16)
				val := make([]string, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = bitstream.ReadString()
				}
				params[i] = val
			case 33:
				nLen := bitstream.ReadInt(16)
				val := make([]int, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = bitstream.ReadInt(32)
				}
				params[i] = val
			case 34:
				nLen := bitstream.ReadInt(16)
				val := make([]uint, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = uint(bitstream.ReadInt(32))
				}
				params[i] = val


			case 41:
				nLen := bitstream.ReadInt(16)
				aa := bool(false)
				tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
				val := reflect.New(tVal).Elem()
				for i := 0; i < nLen; i++ {
					val.Index(i).SetBool(bitstream.ReadFlag())
				}
				params[i] = val.Interface()
			case 42:
				nLen := bitstream.ReadInt(16)
				aa := float64(0)
				tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
				val := reflect.New(tVal).Elem()
				for i := 0; i < nLen; i++ {
					val.Index(i).SetFloat(bitstream.ReadFloat64())
				}
				params[i] = val.Interface()
			case 43:
				nLen := bitstream.ReadInt(16)
				aa := float32(0)
				tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
				val := reflect.New(tVal).Elem()
				for i := 0; i < nLen; i++ {
					val.Index(i).SetFloat(float64(bitstream.ReadFloat()))
				}
				params[i] = val.Interface()
			case 44:
				nLen := bitstream.ReadInt(16)
				aa := int8(0)
				tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
				val := reflect.New(tVal).Elem()
				for i := 0; i < nLen; i++ {
					val.Index(i).SetInt(int64(bitstream.ReadInt(8)))
				}
				params[i] = val.Interface()
			case 45:
				nLen := bitstream.ReadInt(16)
				aa := uint8(0)
				tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
				val := reflect.New(tVal).Elem()
				for i := 0; i < nLen; i++ {
					val.Index(i).SetInt(int64(bitstream.ReadInt(8)))
				}
				params[i] = val.Interface()
			case 46:
				nLen := bitstream.ReadInt(16)
				aa := int16(0)
				tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
				val := reflect.New(tVal).Elem()
				for i := 0; i < nLen; i++ {
					val.Index(i).SetInt(int64(bitstream.ReadInt(16)))
				}
				params[i] = val.Interface()
			case 47:
				nLen := bitstream.ReadInt(16)
				aa := uint16(0)
				tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
				val := reflect.New(tVal).Elem()
				for i := 0; i < nLen; i++ {
					val.Index(i).SetInt(int64(bitstream.ReadInt(16)))
				}
				params[i] = val.Interface()
			case 48:
				nLen := bitstream.ReadInt(16)
				aa := int32(0)
				tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
				val := reflect.New(tVal).Elem()
				for i := 0; i < nLen; i++ {
					val.Index(i).SetInt(int64(bitstream.ReadInt(32)))
				}
				params[i] = val.Interface()
			case 49:
				nLen := bitstream.ReadInt(16)
				aa := uint32(0)
				tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
				val := reflect.New(tVal).Elem()
				for i := 0; i < nLen; i++ {
					val.Index(i).SetInt(int64(bitstream.ReadInt(32)))
				}
				params[i] = val.Interface()
			case 50:
				nLen := bitstream.ReadInt(16)
				aa := int64(0)
				tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
				val := reflect.New(tVal).Elem()
				for i := 0; i < nLen; i++ {
					val.Index(i).SetInt(int64(bitstream.ReadInt64(64)))
				}
				params[i] = val.Interface()
			case 51:
				nLen := bitstream.ReadInt(16)
				aa := uint64(0)
				tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
				val := reflect.New(tVal).Elem()
				for i := 0; i < nLen; i++ {
					val.Index(i).SetInt(int64(bitstream.ReadInt64(64)))
				}
				params[i] = val.Interface()
			case 52:
				nLen := bitstream.ReadInt(16)
				aa := string("")
				tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
				val := reflect.New(tVal).Elem()
				for i := 0; i < nLen; i++ {
					val.Index(i).SetString(bitstream.ReadString())
				}
				params[i] = val.Interface()
			case 53:
				nLen := bitstream.ReadInt(16)
				aa := int(0)
				tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
				val := reflect.New(tVal).Elem()
				for i := 0; i < nLen; i++ {
					val.Index(i).SetInt(int64(bitstream.ReadInt(32)))
				}
				params[i] = val.Interface()
			case 54:
				nLen := bitstream.ReadInt(16)
				aa := uint(0)
				tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
				val := reflect.New(tVal).Elem()
				for i := 0; i < nLen; i++ {
					val.Index(i).SetInt(int64(bitstream.ReadInt(32)))
				}
				params[i] = val.Interface()
			/*case 55://[*]struct
				if k.In(i).Kind() != reflect.Array{
					log.Printf("func [%s] params no fit, func params [%s], params [%v]", funcName, strParams, params)
					return
				}
				nLen := bitstream.ReadInt(16)
				tVal := reflect.ArrayOf(nLen, reflect.TypeOf(k.In(i)))
				val := reflect.New(tVal).Elem()
				arrayPtr := uintptr(unsafe.Pointer(val.Addr().Pointer()))
				for i := 0; i < nLen; i++ {
					value :=  unsafe.Pointer(unsafe.Pointer(arrayPtr))
					packet:= base.GetMessage(bitstream.ReadString())
					base.ReadData(packet, bitstream)
					arrayPtr = arrayPtr + unsafe.Sizeof(packet)
					*(*unsafe.Pointer)(value) = unsafe.Pointer(reflect.ValueOf(packet).Pointer())
				}

				params[i] = val.Interface()*/

			case 61:
				val := new(bool)
				*val = bitstream.ReadFlag()
				params[i] = val
			case 62:
				val := new(float64)
				*val = bitstream.ReadFloat64()
				params[i] = val
			case 63:
				val := new(float32)
				*val = bitstream.ReadFloat()
				params[i] = val
			case 64:
				val := new(int8)
				*val = int8(bitstream.ReadInt(8))
				params[i] = val
			case 65:
				val := new(uint8)
				*val = uint8(bitstream.ReadInt(8))
				params[i] = val
			case 66:
				val := new(int16)
				*val = int16(bitstream.ReadInt(16))
				params[i] = val
			case 67:
				val := new(uint16)
				*val = uint16(bitstream.ReadInt(16))
				params[i] = val
			case 68:
				val := new(int32)
				*val = int32(bitstream.ReadInt(32))
				params[i] = val
			case 69:
				val := new(uint32)
				*val = uint32(bitstream.ReadInt(32))
				params[i] = val
			case 70:
				val := new(int64)
				*val = int64(bitstream.ReadInt64(64))
				params[i] = val
			case 71:
				val := new(uint64)
				*val = uint64(bitstream.ReadInt64(64))
				params[i] = val
			case 72:
				val := new(string)
				*val = bitstream.ReadString()
				params[i] = val
			case 73:
				val := new(int)
				*val = bitstream.ReadInt(32)
				params[i] = val
			case 74:
				val := new(uint)
				*val = uint(bitstream.ReadInt(32))
				params[i] = val
			case 75://*struct
				packet := base.GetMessage(bitstream.ReadString())
				base.ReadData(packet, bitstream)
				params[i] = packet



			case 81:
				nLen := bitstream.ReadInt(16)
				val := make([]*bool, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = new(bool)
					*val[i] = bitstream.ReadFlag()
				}
				params[i] = val
			case 82:
				nLen := bitstream.ReadInt(16)
				val := make([]*float64, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = new(float64)
					*val[i] = bitstream.ReadFloat64()
				}
				params[i] = val
			case 83:
				nLen := bitstream.ReadInt(16)
				val := make([]*float32, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = new(float32)
					*val[i] = bitstream.ReadFloat()
				}
				params[i] = val
			case 84:
				nLen := bitstream.ReadInt(16)
				val := make([]*int8, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = new(int8)
					*val[i] = int8(bitstream.ReadInt(8))
				}
				params[i] = val
			case 85:
				nLen := bitstream.ReadInt(16)
				val := make([]*uint8, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = new(uint8)
					*val[i] = uint8(bitstream.ReadInt(8))
				}
				params[i] = val
			case 86:
				nLen := bitstream.ReadInt(16)
				val := make([]*int16, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = new(int16)
					*val[i] = int16(bitstream.ReadInt(16))
				}
				params[i] = val
			case 87:
				nLen := bitstream.ReadInt(16)
				val := make([]*uint16, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = new(uint16)
					*val[i] = uint16(bitstream.ReadInt(16))
				}
				params[i] = val
			case 88:
				nLen := bitstream.ReadInt(16)
				val := make([]*int32, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = new(int32)
					*val[i] = int32(bitstream.ReadInt(32))
				}
				params[i] = val
			case 89:
				nLen := bitstream.ReadInt(16)
				val := make([]*uint32, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = new(uint32)
					*val[i] = uint32(bitstream.ReadInt(32))
				}
				params[i] = val
			case 90:
				nLen := bitstream.ReadInt(16)
				val := make([]*int64, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = new(int64)
					*val[i] = int64(bitstream.ReadInt64(64))
				}
				params[i] = val
			case 91:
				nLen := bitstream.ReadInt(16)
				val := make([]*uint64, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = new(uint64)
					*val[i] = uint64(bitstream.ReadInt64(64))
				}
				params[i] = val
			case 92:
				nLen := bitstream.ReadInt(16)
				val := make([]*string, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = new(string)
					*val[i] = bitstream.ReadString()
				}
				params[i] = val
			case 93:
				nLen := bitstream.ReadInt(16)
				val := make([]*int, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = new(int)
					*val[i] = bitstream.ReadInt(32)
				}
				params[i] = val
			case 94:
				nLen := bitstream.ReadInt(16)
				val := make([]*uint, nLen)
				for i := 0; i < nLen; i++ {
					val[i] = new(uint)
					*val[i] = uint(bitstream.ReadInt(32))
				}
				params[i] = val
			case 101:
				nLen := bitstream.ReadInt(16)
				aa := bool(false)
				tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
				val := reflect.New(tVal).Elem()
				arrayPtr := uintptr(unsafe.Pointer(val.Addr().Pointer()))
				for i := 0; i < nLen; i++ {
					value :=  (**bool)(unsafe.Pointer(arrayPtr))
					arrayPtr = arrayPtr + base.SIZE_BOOL
					val1 := bitstream.ReadFlag()
					*value = &val1
				}
				params[i] = val.Interface()
			case 102:
				nLen := bitstream.ReadInt(16)
				aa := float64(0)
				tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
				val := reflect.New(tVal).Elem()
				arrayPtr := uintptr(unsafe.Pointer(val.Addr().Pointer()))
				for i := 0; i < nLen; i++ {
					value :=  (**float64)(unsafe.Pointer(arrayPtr))
					arrayPtr = arrayPtr + base.SIZE_FLOAT64
					val1 := bitstream.ReadFloat64()
					*value = &val1
				}
				params[i] = val.Interface()
			case 103:
				nLen := bitstream.ReadInt(16)
				aa := float32(0)
				tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
				val := reflect.New(tVal).Elem()
				arrayPtr := uintptr(unsafe.Pointer(val.Addr().Pointer()))
				for i := 0; i < nLen; i++ {
					value :=  (**float32)(unsafe.Pointer(arrayPtr))
					arrayPtr = arrayPtr + base.SIZE_FLOAT32
					val1 := float32(bitstream.ReadFloat64())
					*value =  &val1
				}
				params[i] = val.Interface()
			case 104:
				nLen := bitstream.ReadInt(16)
				aa := int8(0)
				tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
				val := reflect.New(tVal).Elem()
				arrayPtr := uintptr(unsafe.Pointer(val.Addr().Pointer()))
				for i := 0; i < nLen; i++ {
					value :=  (**int8)(unsafe.Pointer(arrayPtr))
					arrayPtr = arrayPtr + base.SIZE_INT8
					val1 := int8(bitstream.ReadInt(8))
					*value =  &val1
				}
				params[i] = val.Interface()
			case 105:
				nLen := bitstream.ReadInt(16)
				aa := uint8(0)
				tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
				val := reflect.New(tVal).Elem()
				arrayPtr := uintptr(unsafe.Pointer(val.Addr().Pointer()))
				for i := 0; i < nLen; i++ {
					value :=  (**uint8)(unsafe.Pointer(arrayPtr))
					arrayPtr = arrayPtr + base.SIZE_UINT8
					val1 := uint8(bitstream.ReadInt(8))
					*value = &val1
				}
				params[i] = val.Interface()
			case 106:
				nLen := bitstream.ReadInt(16)
				aa := int16(0)
				tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
				val := reflect.New(tVal).Elem()
				arrayPtr := uintptr(unsafe.Pointer(val.Addr().Pointer()))
				for i := 0; i < nLen; i++ {
					value :=  (**int16)(unsafe.Pointer(arrayPtr))
					arrayPtr = arrayPtr + base.SIZE_INT16
					val1 := int16(bitstream.ReadInt(16))
					*value =&val1
				}
				params[i] = val.Interface()
			case 107:
				nLen := bitstream.ReadInt(16)
				aa := uint16(0)
				tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
				val := reflect.New(tVal).Elem()
				arrayPtr := uintptr(unsafe.Pointer(val.Addr().Pointer()))
				for i := 0; i < nLen; i++ {
					value :=  (**uint16)(unsafe.Pointer(arrayPtr))
					arrayPtr = arrayPtr + base.SIZE_UINT16
					val1 := uint16(bitstream.ReadInt(16))
					*value = &val1
				}
				params[i] = val.Interface()
			case 108:
				nLen := bitstream.ReadInt(16)
				aa := int32(0)
				tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
				val := reflect.New(tVal).Elem()
				arrayPtr := uintptr(unsafe.Pointer(val.Addr().Pointer()))
				for i := 0; i < nLen; i++ {
					value :=  (**int32)(unsafe.Pointer(arrayPtr))
					arrayPtr = arrayPtr + base.SIZE_INT32
					val1 := int32(bitstream.ReadInt(32))
					*value = &val1
				}
				params[i] = val.Interface()
			case 109:
				nLen := bitstream.ReadInt(16)
				aa := uint32(0)
				tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
				val := reflect.New(tVal).Elem()
				arrayPtr := uintptr(unsafe.Pointer(val.Addr().Pointer()))
				for i := 0; i < nLen; i++ {
					value :=  (**uint32)(unsafe.Pointer(arrayPtr))
					arrayPtr = arrayPtr + base.SIZE_UINT32
					val1 := uint32(bitstream.ReadInt(32))
					*value = &val1
				}
				params[i] = val.Interface()
			case 110:
				nLen := bitstream.ReadInt(16)
				aa := int64(0)
				tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
				val := reflect.New(tVal).Elem()
				arrayPtr := uintptr(unsafe.Pointer(val.Addr().Pointer()))
				for i := 0; i < nLen; i++ {
					value :=  (**int64)(unsafe.Pointer(arrayPtr))
					arrayPtr = arrayPtr + base.SIZE_INT64
					val1 := int64(bitstream.ReadInt64(64))
					*value =  &val1
				}
				params[i] = val.Interface()
			case 111:
				nLen := bitstream.ReadInt(16)
				aa := uint64(0)
				tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
				val := reflect.New(tVal).Elem()
				arrayPtr := uintptr(unsafe.Pointer(val.Addr().Pointer()))
				for i := 0; i < nLen; i++ {
					value :=  (**uint64)(unsafe.Pointer(arrayPtr))
					arrayPtr = arrayPtr + base.SIZE_UINT64
					val1 := uint64(bitstream.ReadInt64(64))
					*value = &val1
				}
				params[i] = val.Interface()
			case 112:
				nLen := bitstream.ReadInt(16)
				aa := string("")
				tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
				val := reflect.New(tVal).Elem()
				arrayPtr := uintptr(unsafe.Pointer(val.Addr().Pointer()))
				for i := 0; i < nLen; i++ {
					value :=  (**string)(unsafe.Pointer(arrayPtr))
					arrayPtr = arrayPtr + base.SIZE_STRING
					val1 := string(bitstream.ReadString())
					*value = &val1
				}
				params[i] = val.Interface()
			case 113:
				nLen := bitstream.ReadInt(16)
				aa := int(0)
				tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
				val := reflect.New(tVal).Elem()
				arrayPtr := uintptr(unsafe.Pointer(val.Addr().Pointer()))
				for i := 0; i < nLen; i++ {
					value :=  (**int)(unsafe.Pointer(arrayPtr))
					arrayPtr = arrayPtr + base.SIZE_INT
					val1 := bitstream.ReadInt(32)
					*value = &val1
				}
				params[i] = val.Interface()
			case 114:
				nLen := bitstream.ReadInt(16)
				aa := uint(0)
				tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
				val := reflect.New(tVal).Elem()
				arrayPtr := uintptr(unsafe.Pointer(val.Addr().Pointer()))
				for i := 0; i < nLen; i++ {
					value :=  (**uint)(unsafe.Pointer(arrayPtr))
					arrayPtr = arrayPtr + base.SIZE_UINT
					val1 := uint(bitstream.ReadInt(32))
					*value = &val1
				}
				params[i] = val.Interface()

			case base.RPC_MESSAGE://protobuf
				packet := message.GetPakcetByName(funcName)
				nLen := bitstream.ReadInt(base.Bit32)
				packetBuf := bitstream.ReadBits(nLen << 3)
				message.UnmarshalText(packet, packetBuf)
				params[i] = packet

			default:
				panic("func [%s] params type not supported")
			}
		}
	}
}
