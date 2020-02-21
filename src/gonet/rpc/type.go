package rpc

import (
	"github.com/golang/protobuf/proto"
	"gonet/base"
	"reflect"
	"strings"
)

const(
	RPC_BOOL = iota
	RPC_STRING
	RPC_FLOAT32
	RPC_FLOAT64
	RPC_INT
	RPC_INT8
	RPC_INT16
	RPC_INT32
	RPC_INT64
	RPC_UINT
	RPC_UINT16
	RPC_UINT8
	RPC_UINT32
	RPC_UINT64

	RPC_BOOL_SLICE
	RPC_STRING_SLICE
	RPC_FLOAT32_SLICE
	RPC_FLOAT64_SLICE
	RPC_INT_SLICE
	RPC_INT8_SLICE
	RPC_INT16_SLICE
	RPC_INT32_SLICE
	RPC_INT64_SLICE
	RPC_UINT_SLICE
	RPC_UINT8_SLICE
	RPC_UINT16_SLICE
	RPC_UINT32_SLICE
	RPC_UINT64_SLICE

	RPC_BOOL_ARRAY
	RPC_STRING_ARRAY
	RPC_FLOAT32_ARRAY
	RPC_FLOAT64_ARRAY
	RPC_INT_ARRAY
	RPC_INT8_ARRAY
	RPC_INT16_ARRAY
	RPC_INT32_ARRAY
	RPC_INT64_ARRAY
	RPC_UINT_ARRAY
	RPC_UINT8_ARRAY
	RPC_UINT16_ARRAY
	RPC_UINT32_ARRAY
	RPC_UINT64_ARRAY

	RPC_BOOL_PTR
	RPC_STRING_PTR
	RPC_FLOAT32_PTR
	RPC_FLOAT64_PTR
	RPC_INT_PTR
	RPC_INT8_PTR
	RPC_INT16_PTR
	RPC_INT32_PTR
	RPC_INT64_PTR
	RPC_UINT_PTR
	RPC_UINT8_PTR
	RPC_UINT16_PTR
	RPC_UINT32_PTR
	RPC_UINT64_PTR

	RPC_BOOL_PTR_SLICE
	RPC_STRING_PTR_SLICE
	RPC_FLOAT32_PTR_SLICE
	RPC_FLOAT64_PTR_SLICE
	RPC_INT_PTR_SLICE
	RPC_INT8_PTR_SLICE
	RPC_INT16_PTR_SLICE
	RPC_INT32_PTR_SLICE
	RPC_INT64_PTR_SLICE
	RPC_UINT_PTR_SLICE
	RPC_UINT8_PTR_SLICE
	RPC_UINT16_PTR_SLICE
	RPC_UINT32_PTR_SLICE
	RPC_UINT64_PTR_SLICE

	RPC_BOOL_PTR_ARRAY
	RPC_STRING_PTR_ARRAY
	RPC_FLOAT32_PTR_ARRAY
	RPC_FLOAT64_PTR_ARRAY
	RPC_INT_PTR_ARRAY
	RPC_INT8_PTR_ARRAY
	RPC_INT16_PTR_ARRAY
	RPC_INT32_PTR_ARRAY
	RPC_INT64_PTR_ARRAY
	RPC_UINT_PTR_ARRAY
	RPC_UINT8_PTR_ARRAY
	RPC_UINT16_PTR_ARRAY
	RPC_UINT32_PTR_ARRAY
	RPC_UINT64_PTR_ARRAY

	RPC_BOOL_SLICE_PTR
	RPC_STRING_SLICE_PTR
	RPC_FLOAT32_SLICE_PTR
	RPC_FLOAT64_SLICE_PTR
	RPC_INT_SLICE_PTR
	RPC_INT8_SLICE_PTR
	RPC_INT16_SLICE_PTR
	RPC_INT32_SLICE_PTR
	RPC_INT64_SLICE_PTR
	RPC_UINT_SLICE_PTR
	RPC_UINT8_SLICE_PTR
	RPC_UINT16_SLICE_PTR
	RPC_UINT32_SLICE_PTR
	RPC_UINT64_SLICE_PTR

	RPC_BOOL_PTR_SLICE_PTR
	RPC_STRING_PTR_SLICE_PTR
	RPC_FLOAT32_PTR_SLICE_PTR
	RPC_FLOAT64_PTR_SLICE_PTR
	RPC_INT_PTR_SLICE_PTR
	RPC_INT8_PTR_SLICE_PTR
	RPC_INT16_PTR_SLICE_PTR
	RPC_INT32_PTR_SLICE_PTR
	RPC_INT64_PTR_SLICE_PTR
	RPC_UINT_PTR_SLICE_PTR
	RPC_UINT8_PTR_SLICE_PTR
	RPC_UINT16_PTR_SLICE_PTR
	RPC_UINT32_PTR_SLICE_PTR
	RPC_UINT64_PTR_SLICE_PTR

	RPC_BOOL_ARRAY_PTR
	RPC_STRING_ARRAY_PTR
	RPC_FLOAT32_ARRAY_PTR
	RPC_FLOAT64_ARRAY_PTR
	RPC_INT_ARRAY_PTR
	RPC_INT8_ARRAY_PTR
	RPC_INT16_ARRAY_PTR
	RPC_INT32_ARRAY_PTR
	RPC_INT64_ARRAY_PTR
	RPC_UINT_ARRAY_PTR
	RPC_UINT8_ARRAY_PTR
	RPC_UINT16_ARRAY_PTR
	RPC_UINT32_ARRAY_PTR
	RPC_UINT64_ARRAY_PTR

	RPC_BOOL_PTR_ARRAY_PTR
	RPC_STRING_PTR_ARRAY_PTR
	RPC_FLOAT32_PTR_ARRAY_PTR
	RPC_FLOAT64_PTR_ARRAY_PTR
	RPC_INT_PTR_ARRAY_PTR
	RPC_INT8_PTR_ARRAY_PTR
	RPC_INT16_PTR_ARRAY_PTR
	RPC_INT32_PTR_ARRAY_PTR
	RPC_INT64_PTR_ARRAY_PTR
	RPC_UINT_PTR_ARRAY_PTR
	RPC_UINT8_PTR_ARRAY_PTR
	RPC_UINT16_PTR_ARRAY_PTR
	RPC_UINT32_PTR_ARRAY_PTR
	RPC_UINT64_PTR_ARRAY_PTR

	RPC_MESSAGE
	RPC_GOB//暂时用json,gob包头解析小包太慢
)


func getSliceTypeString(sTypeName string, bPtr bool) string{
	index := strings.Index(sTypeName, "]")
	if index != -1{
		sTypeName = sTypeName[index+1:]
	}

	switch sTypeName {
	case "*bool", "*float64", "*float32", "*int8", "*uint8", "*int16", "*uint16",
		"*int32", "*uint32", "*int64", "*uint64", "*string", "*int", "*uint",
		"bool", "float64", "float32", "int8", "uint8", "int16", "uint16",
		"int32", "uint32", "int64", "uint64", "string", "int", "uint":
		if !bPtr{
			return "[]" + sTypeName
		}else{
			return "*[]" + sTypeName
		}
	}
	return "*gob"
}

func getArrayTypeString(sTypeName string, bPtr bool) string{
	index := strings.Index(sTypeName, "]")
	if index != -1{
		sTypeName = sTypeName[index+1:]
	}

	switch sTypeName {
	case "*bool", "*float64", "*float32", "*int8", "*uint8", "*int16", "*uint16",
		"*int32", "*uint32", "*int64", "*uint64", "*string", "*int", "*uint",
		"bool", "float64", "float32", "int8", "uint8", "int16", "uint16",
		"int32", "uint32", "int64", "uint64", "string", "int", "uint":
		if !bPtr{
			return "[*]" + sTypeName
		}else{
			return "*[*]" + sTypeName
		}
	}
	return "*gob"
}

func getTypeString(param interface{}) string{
	paramType := reflect.TypeOf(param)
	sType := ""

	if paramType.Kind() == reflect.Ptr{
		switch paramType.String() {
		case "*bool", "*float64", "*float32", "*int8", "*uint8", "*int16", "*uint16",
			"*int32", "*uint32", "*int64", "*uint64", "*string", "*int", "*uint":
			sType = paramType.String()
		default:
			if paramType.Elem().Kind() == reflect.Array{
				sType = getArrayTypeString(paramType.String(), true)
			}else if paramType.Elem().Kind() == reflect.Slice{
				sType = getSliceTypeString(paramType.String(), true)
			}else if strings.Index(paramType.String(), "*message")!= -1{
				sType = "*message"
			}else{
				sType = "*gob"
			}
		}
	}else if paramType.Kind() == reflect.Slice{
		sType = getSliceTypeString(paramType.String(), false)
	}else if paramType.Kind() == reflect.Array{
		sType = getArrayTypeString(paramType.String(), false)
	}else{
		switch paramType.String() {
		case "bool", "float64", "float32", "int8", "uint8", "int16", "uint16",
			"int32", "uint32", "int64", "uint64", "string", "int", "uint":
			sType = paramType.String()
		default:
			sType = "*gob"
		}
	}

	return sType
}

//rpc  MarshalPB
func MarshalPB(bitstream *base.BitStream, packet proto.Message) {
	bitstream.WriteString(proto.MessageName(packet))
	buf, _ :=proto.Marshal(packet)
	nLen := len(buf)
	bitstream.WriteInt(nLen, 32)
	bitstream.WriteBits(buf, nLen << 3)
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