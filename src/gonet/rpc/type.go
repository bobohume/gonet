package rpc

import (
	"reflect"
	"strings"
	"unsafe"
)

const(
	RPC_BOOL = iota
	RPC_FLOAT64
	RPC_FLOAT32
	RPC_INT8
	RPC_UINT8
	RPC_INT16
	RPC_UINT16
	RPC_INT32
	RPC_UINT32
	RPC_INT64
	RPC_UINT64
	RPC_STRING
	RPC_INT
	RPC_UINT

	RPC_BOOL_SLICE
	RPC_FLOAT64_SLICE
	RPC_FLOAT32_SLICE
	RPC_INT8_SLICE
	RPC_UINT8_SLICE
	RPC_INT16_SLICE
	RPC_UINT16_SLICE
	RPC_INT32_SLICE
	RPC_UINT32_SLICE
	RPC_INT64_SLICE
	RPC_UINT64_SLICE
	RPC_STRING_SLICE
	RPC_INT_SLICE
	RPC_UINT_SLICE

	RPC_BOOL_ARRAY
	RPC_FLOAT64_ARRAY
	RPC_FLOAT32_ARRAY
	RPC_INT8_ARRAY
	RPC_UINT8_ARRAY
	RPC_INT16_ARRAY
	RPC_UINT16_ARRAY
	RPC_INT32_ARRAY
	RPC_UINT32_ARRAY
	RPC_INT64_ARRAY
	RPC_UINT64_ARRAY
	RPC_STRING_ARRAY
	RPC_INT_ARRAY
	RPC_UINT_ARRAY

	RPC_BOOL_PTR
	RPC_FLOAT64_PTR
	RPC_FLOAT32_PTR
	RPC_INT8_PTR
	RPC_UINT8_PTR
	RPC_INT16_PTR
	RPC_UINT16_PTR
	RPC_INT32_PTR
	RPC_UINT32_PTR
	RPC_INT64_PTR
	RPC_UINT64_PTR
	RPC_STRING_PTR
	RPC_INT_PTR
	RPC_UINT_PTR

	RPC_BOOL_PTR_SLICE
	RPC_FLOAT64_PTR_SLICE
	RPC_FLOAT32_PTR_SLICE
	RPC_INT8_PTR_SLICE
	RPC_UINT8_PTR_SLICE
	RPC_INT16_PTR_SLICE
	RPC_UINT16_PTR_SLICE
	RPC_INT32_PTR_SLICE
	RPC_UINT32_PTR_SLICE
	RPC_INT64_PTR_SLICE
	RPC_UINT64_PTR_SLICE
	RPC_STRING_PTR_SLICE
	RPC_INT_PTR_SLICE
	RPC_UINT_PTR_SLICE

	RPC_BOOL_PTR_ARRAY
	RPC_FLOAT64_PTR_ARRAY
	RPC_FLOAT32_PTR_ARRAY
	RPC_INT8_PTR_ARRAY
	RPC_UINT8_PTR_ARRAY
	RPC_INT16_PTR_ARRAY
	RPC_UINT16_PTR_ARRAY
	RPC_INT32_PTR_ARRAY
	RPC_UINT32_PTR_ARRAY
	RPC_INT64_PTR_ARRAY
	RPC_UINT64_PTR_ARRAY
	RPC_STRING_PTR_ARRAY
	RPC_INT_PTR_ARRAY
	RPC_UINT_PTR_ARRAY

	RPC_MESSAGE 	= 120
	RPC_GOB			= 121//暂时用json,gob包头解析小包太慢
)

const(
	SIZE_BOOL = unsafe.Sizeof(bool(false))
	SIZE_INT = unsafe.Sizeof(int(0))
	SIZE_INT8 = unsafe.Sizeof(int8(0))
	SIZE_INT16 = unsafe.Sizeof(int16(0))
	SIZE_INT32 = unsafe.Sizeof(int32(0))
	SIZE_INT64 = unsafe.Sizeof(int64(0))
	SIZE_UINT = unsafe.Sizeof(uint(0))
	SIZE_UINT8 = unsafe.Sizeof(uint8(0))
	SIZE_UINT16 = unsafe.Sizeof(uint16(0))
	SIZE_UINT32 = unsafe.Sizeof(uint32(0))
	SIZE_UINT64 = unsafe.Sizeof(uint64(0))
	SIZE_FLOAT32 = unsafe.Sizeof(float32(0))
	SIZE_FLOAT64 = unsafe.Sizeof(float64(0))
	SIZE_STRING = unsafe.Sizeof(string(0))
	SIZE_PTR 	= unsafe.Sizeof(uintptr(0))
)//packet size


func getSliceTypeString(sTypeName string) string{
	index := strings.Index(sTypeName, "]")
	if index != -1{
		sTypeName = sTypeName[index+1:]
	}

	switch sTypeName {
	case "*bool", "*float64", "*float32", "*int8", "*uint8", "*int16", "*uint16",
		"*int32", "*uint32", "*int64", "*uint64", "*string", "*int", "*uint",
		"bool", "float64", "float32", "int8", "uint8", "int16", "uint16",
		"int32", "uint32", "int64", "uint64", "string", "int", "uint":
		return "[]" + sTypeName
	}
	return "*gob"
}

func getArrayTypeString(sTypeName string) string{
	index := strings.Index(sTypeName, "]")
	if index != -1{
		sTypeName = sTypeName[index+1:]
	}

	switch sTypeName {
	case "*bool", "*float64", "*float32", "*int8", "*uint8", "*int16", "*uint16",
		"*int32", "*uint32", "*int64", "*uint64", "*string", "*int", "*uint",
		"bool", "float64", "float32", "int8", "uint8", "int16", "uint16",
		"int32", "uint32", "int64", "uint64", "string", "int", "uint":
		return "[*]" + sTypeName
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
			sType = "*gob"
		}
	}else if paramType.Kind() == reflect.Slice{
		sType = getSliceTypeString(paramType.String())
	}else if paramType.Kind() == reflect.Array{
		sType = getArrayTypeString(paramType.String())
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