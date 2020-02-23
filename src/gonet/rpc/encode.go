package rpc

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/json-iterator/go"
	"gonet/base"
	"gonet/message"
	"reflect"
	"strings"
)

//rpc  Marshal
//rpc  特定rpc头部设置需求，params[0]传入RpcHead
func Marshal(funcName string, params ...interface{})[]byte {
	data, _ := MarshalEx(funcName, params...)
	return data
}

//rpc  MarshalEx
//rpc  特定rpc头部设置需求，params[0]传入RpcHead
func MarshalEx(funcName string, params ...interface{})([]byte, *message.RpcPacket) {
	defer func() {
		if err := recover(); err != nil {
			base.TraceCode(err)
		}
	}()

	rpcPacket := &message.RpcPacket{FuncName:strings.ToLower(funcName), ArgLen:int32(len(params)), RpcHead:&message.RpcHead{}}
	msg := make([]byte, 1024)
	bitstream := base.NewBitStream(msg, 1024)
	for i, param := range params {
		sType := getTypeString(param)
		switch sType {
		case "bool":
			bitstream.WriteInt(RPC_BOOL, 8)
			bitstream.WriteFlag(param.(bool))
		case "string":
			bitstream.WriteInt(RPC_STRING, 8)
			bitstream.WriteString(param.(string))
		case "float32":
			bitstream.WriteInt(RPC_FLOAT32, 8)
			bitstream.WriteFloat(param.(float32))
		case "float64":
			bitstream.WriteInt(RPC_FLOAT64, 8)
			bitstream.WriteFloat64(param.(float64))
		case "int":
			bitstream.WriteInt(RPC_INT, 8)
			bitstream.WriteInt(param.(int), 32)
		case "int8":
			bitstream.WriteInt(RPC_INT8, 8)
			bitstream.WriteInt(int(param.(int8)), 8)
		case "int16":
			bitstream.WriteInt(RPC_INT16, 8)
			bitstream.WriteInt(int(param.(int16)),16)
		case "int32":
			bitstream.WriteInt(RPC_INT32, 8)
			bitstream.WriteInt(int(param.(int32)),32)
		case "int64":
			bitstream.WriteInt(RPC_INT64, 8)
			bitstream.WriteInt64(param.(int64), 64)
		case "uint":
			bitstream.WriteInt(RPC_UINT, 8)
			bitstream.WriteInt(int(param.(uint)), 32)
		case "uint8":
			bitstream.WriteInt(RPC_UINT8, 8)
			bitstream.WriteInt(int(param.(uint8)),8)
		case "uint16":
			bitstream.WriteInt(RPC_UINT16, 8)
			bitstream.WriteInt(int(param.(uint16)),16)
		case "uint32":
			bitstream.WriteInt(RPC_UINT32, 8)
			bitstream.WriteInt(int(param.(uint32)),32)
		case "uint64":
			bitstream.WriteInt(RPC_UINT64, 8)
			bitstream.WriteInt64(int64(param.(uint64)), 64)



		case "[]bool":
			bitstream.WriteInt(RPC_BOOL_SLICE, 8)
			val := param.([]bool)
			nLen := len(val)
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteFlag(val[i])
			}
		case "[]string":
			bitstream.WriteInt(RPC_STRING_SLICE, 8)
			val := param.([]string)
			nLen := len(val)
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteString(val[i])
			}
		case "[]float32":
			bitstream.WriteInt(RPC_FLOAT32_SLICE, 8)
			val := param.([]float32)
			nLen := len(val)
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteFloat(val[i])
			}
		case "[]float64":
			bitstream.WriteInt(RPC_FLOAT64_SLICE, 8)
			val := param.([]float64)
			nLen := len(val)
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteFloat64(val[i])
			}
		case "[]int":
			bitstream.WriteInt(RPC_INT_SLICE, 8)
			val := param.([]int)
			nLen := len(val)
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(val[i], 32)
			}
		case "[]int8":
			bitstream.WriteInt(RPC_INT8_SLICE, 8)
			val := param.([]int8)
			nLen := len(val)
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(val[i]), 8)
			}
		case "[]int16":
			bitstream.WriteInt(RPC_INT16_SLICE, 8)
			val := param.([]int16)
			nLen := len(val)
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(val[i]), 16)
			}
		case "[]int32":
			bitstream.WriteInt(RPC_INT32_SLICE, 8)
			val := param.([]int32)
			nLen := len(val)
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(val[i]), 32)
			}
		case "[]int64":
			bitstream.WriteInt(RPC_INT64_SLICE, 8)
			val := param.([]int64)
			nLen := len(val)
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt64(val[i], 64)
			}
		case "[]uint":
			bitstream.WriteInt(RPC_UINT_SLICE, 8)
			val := param.([]uint)
			nLen := len(val)
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(val[i]), 32)
			}
		case "[]uint8":
			bitstream.WriteInt(RPC_UINT8_SLICE, 8)
			val := param.([]uint8)
			nLen := len(val)
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(val[i]), 8)
			}
		case "[]uint16":
			bitstream.WriteInt(RPC_UINT16_SLICE, 8)
			val := param.([]uint16)
			nLen := len(val)
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(val[i]), 16)
			}
		case "[]uint32":
			bitstream.WriteInt(RPC_UINT32_SLICE, 8)
			val := param.([]uint32)
			nLen := len(val)
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(val[i]), 32)
			}
		case "[]uint64":
			bitstream.WriteInt(RPC_UINT64_SLICE, 8)
			val := param.([]uint64)
			nLen := len(val)
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt64(int64(val[i]), 64)
			}



		case "[*]bool":
			bitstream.WriteInt(RPC_BOOL_ARRAY, 8)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteFlag(val.Index(i).Bool())
			}
		case "[*]string":
			bitstream.WriteInt(RPC_STRING_ARRAY, 8)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteString(val.Index(i).String())
			}
		case "[*]float32":
			bitstream.WriteInt(RPC_FLOAT32_ARRAY, 8)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteFloat(float32(val.Index(i).Float()))
			}
		case "[*]float64":
			bitstream.WriteInt(RPC_FLOAT64_ARRAY, 8)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteFloat64(val.Index(i).Float())
			}
		case "[*]int":
			bitstream.WriteInt(RPC_INT_ARRAY, 8)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(val.Index(i).Int()), 32)
			}
		case "[*]int8":
			bitstream.WriteInt(RPC_INT8_ARRAY, 8)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(val.Index(i).Int()), 8)
			}
		case "[*]int16":
			bitstream.WriteInt(RPC_INT16_ARRAY, 8)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(val.Index(i).Int()), 16)
			}
		case "[*]int32":
			bitstream.WriteInt(RPC_INT32_ARRAY, 8)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(val.Index(i).Int()), 32)
			}
		case "[*]int64":
			bitstream.WriteInt(RPC_INT64_ARRAY, 8)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt64(val.Index(i).Int(), 64)
			}
		case "[*]uint":
			bitstream.WriteInt(RPC_UINT_ARRAY, 8)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(val.Index(i).Uint()), 32)
			}
		case "[*]uint8":
			bitstream.WriteInt(RPC_UINT8_ARRAY, 8)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(val.Index(i).Uint()), 8)
			}
		case "[*]uint16":
			bitstream.WriteInt(RPC_UINT16_ARRAY, 8)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(val.Index(i).Uint()), 16)
			}
		case "[*]uint32":
			bitstream.WriteInt(RPC_UINT32_ARRAY, 8)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(val.Index(i).Uint()), 32)
			}
		case "[*]uint64":
			bitstream.WriteInt(RPC_UINT64_ARRAY, 8)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt64(int64(val.Index(i).Uint()), 64)
			}


		case "*bool":
			bitstream.WriteInt(RPC_BOOL_PTR, 8)
			if param.(*bool) != nil{
				bitstream.WriteFlag(*param.(*bool))
			}else{
				bitstream.WriteFlag(false)
			}
		case "*string":
			bitstream.WriteInt(RPC_STRING_PTR, 8)
			if param.(*string) != nil {
				bitstream.WriteString(*param.(*string))
			}else{
				bitstream.WriteString("")
			}
		case "*float32":
			bitstream.WriteInt(RPC_FLOAT32_PTR, 8)
			if param.(*float32) != nil {
				bitstream.WriteFloat(*param.(*float32))
			}else{
				bitstream.WriteFloat(0)
			}
		case "*float64":
			bitstream.WriteInt(RPC_FLOAT64_PTR, 8)
			if param.(*float64) != nil {
				bitstream.WriteFloat64(*param.(*float64))
			}else{
				bitstream.WriteFloat64(0)
			}
		case "*int":
			bitstream.WriteInt(RPC_INT_PTR, 8)
			if param.(*int) != nil {
				bitstream.WriteInt(*param.(*int), 32)
			}else{
				bitstream.WriteInt(0, 32)
			}
		case "*int8":
			bitstream.WriteInt(RPC_INT8_PTR, 8)
			if param.(*int8) != nil {
				bitstream.WriteInt(int(*param.(*int8)), 8)
			}else{
				bitstream.WriteInt(0, 8)
			}
		case "*int16":
			bitstream.WriteInt(RPC_INT16_PTR, 8)
			if param.(*int16) != nil {
				bitstream.WriteInt(int(*param.(*int16)), 16)
			}else{
				bitstream.WriteInt(0, 16)
			}
		case "*int32":
			bitstream.WriteInt(RPC_INT32_PTR, 8)
			if param.(*int32) != nil {
				bitstream.WriteInt(int(*param.(*int32)), 32)
			}else{
				bitstream.WriteInt(0, 32)
			}
		case "*int64":
			bitstream.WriteInt(RPC_INT64_PTR, 8)
			if param.(*int64) != nil {
				bitstream.WriteInt64(*param.(*int64), 64)
			}else{
				bitstream.WriteInt64(0, 64)
			}
		case "*uint":
			bitstream.WriteInt(RPC_UINT_PTR, 8)
			if param.(*uint) != nil {
				bitstream.WriteInt(int(*param.(*uint)), 32)
			}else{
				bitstream.WriteInt(0, 32)
			}
		case "*uint8":
			bitstream.WriteInt(RPC_UINT8_PTR, 8)
			if param.(*uint8) != nil {
				bitstream.WriteInt(int(*param.(*uint8)), 8)
			}else{
				bitstream.WriteInt(0, 8)
			}
		case "*uint16":
			bitstream.WriteInt(RPC_UINT16_PTR, 8)
			if param.(*uint16) != nil {
				bitstream.WriteInt(int(*param.(*uint16)), 16)
			}else{
				bitstream.WriteInt(0, 16)
			}
		case "*uint32":
			bitstream.WriteInt(RPC_UINT32_PTR, 8)
			if param.(*uint32) != nil {
				bitstream.WriteInt(int(*param.(*uint32)), 32)
			}else{
				bitstream.WriteInt(0, 32)
			}
		case "*uint64":
			bitstream.WriteInt(RPC_UINT64_PTR, 8)
			if param.(*uint64) != nil {
				bitstream.WriteInt64(int64(*param.(*uint64)), 64)
			}else{
				bitstream.WriteInt64(0, 64)
			}



		case "[]*bool":
			bitstream.WriteInt(RPC_BOOL_PTR_SLICE, 8)
			val := param.([]*bool)
			nLen := len(val)
			bitstream.WriteInt(nLen, 16)
			for _, v := range val{
				if v != nil{
					bitstream.WriteFlag(*v)
				}else{
					bitstream.WriteFlag(false)
				}
			}
		case "[]*string":
			bitstream.WriteInt(RPC_STRING_PTR_SLICE, 8)
			val := param.([]*string)
			nLen := len(val)
			bitstream.WriteInt(nLen, 16)
			for _, v := range val{
				if v != nil{
					bitstream.WriteString(*v)
				}else{
					bitstream.WriteString("")
				}
			}
		case "[]*float32":
			bitstream.WriteInt(RPC_FLOAT32_PTR_SLICE, 8)
			val := param.([]*float32)
			nLen := len(val)
			bitstream.WriteInt(nLen, 16)
			for _, v := range val{
				if v != nil{
					bitstream.WriteFloat(*v)
				}else{
					bitstream.WriteFloat(0)
				}
			}
		case "[]*float64":
			bitstream.WriteInt(RPC_FLOAT64_PTR_SLICE, 8)
			val := param.([]*float64)
			nLen := len(val)
			bitstream.WriteInt(nLen, 16)
			for _, v := range val{
				if v != nil{
					bitstream.WriteFloat64(*v)
				}else{
					bitstream.WriteFloat64(0)
				}
			}
		case "[]*int":
			bitstream.WriteInt(RPC_INT_PTR_SLICE, 8)
			val := param.([]*int)
			nLen := len(val)
			bitstream.WriteInt(nLen, 16)
			for _, v := range val{
				if v != nil{
					bitstream.WriteInt(*v, 32)
				}else{
					bitstream.WriteInt(0, 32)
				}
			}
		case "[]*int8":
			bitstream.WriteInt(RPC_INT8_PTR_SLICE, 8)
			val := param.([]*int8)
			nLen := len(val)
			bitstream.WriteInt(nLen, 16)
			for _, v := range val{
				if v != nil{
					bitstream.WriteInt(int(*v), 8)
				}else{
					bitstream.WriteInt(0, 8)
				}
			}
		case "[]*int16":
			bitstream.WriteInt(RPC_INT16_PTR_SLICE, 8)
			val := param.([]*int16)
			nLen := len(val)
			bitstream.WriteInt(nLen, 16)
			for _, v := range val{
				if v != nil{
					bitstream.WriteInt(int(*v), 16)
				}else{
					bitstream.WriteInt(0, 16)
				}
			}
		case "[]*int32":
			bitstream.WriteInt(RPC_INT32_PTR_SLICE, 8)
			val := param.([]*int32)
			nLen := len(val)
			bitstream.WriteInt(nLen, 16)
			for _, v := range val{
				if v != nil{
					bitstream.WriteInt(int(*v), 32)
				}else{
					bitstream.WriteInt(0, 32)
				}
			}
		case "[]*int64":
			bitstream.WriteInt(RPC_INT64_PTR_SLICE, 8)
			val := param.([]*int64)
			nLen := len(val)
			bitstream.WriteInt(nLen, 16)
			for _, v := range val{
				if v != nil{
					bitstream.WriteInt64(*v, 64)
				}else{
					bitstream.WriteInt64(0, 64)
				}
			}
		case "[]*uint":
			bitstream.WriteInt(RPC_UINT_PTR_SLICE, 8)
			val := param.([]*uint)
			nLen := len(val)
			bitstream.WriteInt(nLen, 16)
			for _, v := range val{
				if v != nil{
					bitstream.WriteInt(int(*v), 32)
				}else{
					bitstream.WriteInt(0, 32)
				}
			}
		case "[]*uint8":
			bitstream.WriteInt(RPC_UINT8_PTR_SLICE, 8)
			val := param.([]*uint8)
			nLen := len(val)
			bitstream.WriteInt(nLen, 16)
			for _, v := range val{
				if v != nil{
					bitstream.WriteInt(int(*v), 8)
				}else{
					bitstream.WriteInt(0, 8)
				}
			}
		case "[]*uint16":
			bitstream.WriteInt(RPC_UINT16_PTR_SLICE, 8)
			val := param.([]*uint16)
			nLen := len(val)
			bitstream.WriteInt(nLen, 16)
			for _, v := range val{
				if v != nil{
					bitstream.WriteInt(int(*v), 16)
				}else{
					bitstream.WriteInt(0, 16)
				}
			}
		case "[]*uint32":
			bitstream.WriteInt(RPC_UINT32_PTR_SLICE, 8)
			val := param.([]*uint32)
			nLen := len(val)
			bitstream.WriteInt(nLen, 16)
			for _, v := range val{
				if v != nil{
					bitstream.WriteInt(int(*v), 32)
				}else{
					bitstream.WriteInt(0, 32)
				}
			}
		case "[]*uint64":
			bitstream.WriteInt(RPC_UINT64_PTR_SLICE, 8)
			val := param.([]*uint64)
			nLen := len(val)
			bitstream.WriteInt(nLen, 16)
			for _, v := range val{
				if v != nil{
					bitstream.WriteInt64(int64(*v), 64)
				}else{
					bitstream.WriteInt64(0, 64)
				}
			}



		case "[*]*bool":
			bitstream.WriteInt(RPC_BOOL_PTR_ARRAY, 8)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				if !val.Index(i).IsNil(){
					bitstream.WriteFlag(val.Index(i).Elem().Bool())
				}else{
					bitstream.WriteFlag(false)
				}
			}
		case "[*]*string":
			bitstream.WriteInt(RPC_STRING_PTR_ARRAY, 8)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				if !val.Index(i).IsNil() {
					bitstream.WriteString(val.Index(i).Elem().String())
				}else{
					bitstream.WriteString("")
				}
			}
		case "[*]*float32":
			bitstream.WriteInt(RPC_FLOAT32_PTR_ARRAY, 8)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				if !val.Index(i).IsNil() {
					bitstream.WriteFloat(float32(val.Index(i).Elem().Float()))
				}else{
					bitstream.WriteFloat(0)
				}
			}
		case "[*]*float64":
			bitstream.WriteInt(RPC_FLOAT64_PTR_ARRAY, 8)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				if !val.Index(i).IsNil() {
					bitstream.WriteFloat64(val.Index(i).Elem().Float())
				}else{
					bitstream.WriteFloat64(0)
				}
			}
		case "[*]*int":
			bitstream.WriteInt(RPC_INT_PTR_ARRAY, 8)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				if !val.Index(i).IsNil() {
					bitstream.WriteInt(int(val.Index(i).Elem().Int()), 32)
				}else{
					bitstream.WriteInt(0, 32)
				}
			}
		case "[*]*int8":
			bitstream.WriteInt(RPC_INT8_PTR_ARRAY, 8)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				if !val.Index(i).IsNil() {
					bitstream.WriteInt(int(val.Index(i).Elem().Int()), 8)
				}else{
					bitstream.WriteInt(0, 8)
				}
			}
		case "[*]*int16":
			bitstream.WriteInt(RPC_INT16_PTR_ARRAY, 8)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				if !val.Index(i).IsNil() {
					bitstream.WriteInt(int(val.Index(i).Elem().Int()), 16)
				}else{
					bitstream.WriteInt(0, 16)
				}
			}
		case "[*]*int32":
			bitstream.WriteInt(RPC_INT32_PTR_ARRAY, 8)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				if !val.Index(i).IsNil() {
					bitstream.WriteInt(int(val.Index(i).Elem().Int()), 32)
				}else{
					bitstream.WriteInt(0, 32)
				}
			}
		case "[*]*int64":
			bitstream.WriteInt(RPC_INT64_PTR_ARRAY, 8)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				if !val.Index(i).IsNil() {
					bitstream.WriteInt64(val.Index(i).Elem().Int(), 64)
				}else{
					bitstream.WriteInt64(0, 64)
				}
			}
		case "[*]*uint":
			bitstream.WriteInt(RPC_UINT_PTR_ARRAY, 8)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				if !val.Index(i).IsNil() {
					bitstream.WriteInt(int(val.Index(i).Elem().Uint()), 32)
				}else{
					bitstream.WriteInt(0, 32)
				}
			}
		case "[*]*uint8":
			bitstream.WriteInt(RPC_UINT8_PTR_ARRAY, 8)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				if !val.Index(i).IsNil() {
					bitstream.WriteInt(int(val.Index(i).Elem().Uint()), 8)
				}else{
					bitstream.WriteInt(0, 8)
				}
			}
		case "[*]*uint16":
			bitstream.WriteInt(RPC_UINT16_PTR_ARRAY, 8)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				if !val.Index(i).IsNil() {
					bitstream.WriteInt(int(val.Index(i).Elem().Uint()), 16)
				}else{
					bitstream.WriteInt(0, 16)
				}
			}
		case "[*]*uint32":
			bitstream.WriteInt(RPC_UINT32_PTR_ARRAY, 8)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				if !val.Index(i).IsNil() {
					bitstream.WriteInt(int(val.Index(i).Elem().Uint()), 32)
				}else{
					bitstream.WriteInt(0, 32)
				}
			}
		case "[*]*uint64":
			bitstream.WriteInt(RPC_UINT64_PTR_ARRAY, 8)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				if !val.Index(i).IsNil() {
					bitstream.WriteInt64(int64(val.Index(i).Elem().Uint()), 64)
				}else{
					bitstream.WriteInt64(0, 64)
				}
			}



		case "*[]bool":
			bitstream.WriteInt(RPC_BOOL_SLICE_PTR, 8)
			val := *param.(*[]bool)
			nLen := len(val)
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteFlag(val[i])
			}
		case "*[]string":
			bitstream.WriteInt(RPC_STRING_SLICE_PTR, 8)
			val := *param.(*[]string)
			nLen := len(val)
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteString(val[i])
			}
		case "*[]float32":
			bitstream.WriteInt(RPC_FLOAT32_SLICE_PTR, 8)
			val := *param.(*[]float32)
			nLen := len(val)
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteFloat(val[i])
			}
		case "*[]float64":
			bitstream.WriteInt(RPC_FLOAT64_SLICE_PTR, 8)
			val := *param.(*[]float64)
			nLen := len(val)
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteFloat64(val[i])
			}
		case "*[]int":
			bitstream.WriteInt(RPC_INT_SLICE_PTR, 8)
			val := *param.(*[]int)
			nLen := len(val)
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(val[i], 32)
			}
		case "*[]int8":
			bitstream.WriteInt(RPC_INT8_SLICE_PTR, 8)
			val := *param.(*[]int8)
			nLen := len(val)
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(val[i]), 8)
			}
		case "*[]int16":
			bitstream.WriteInt(RPC_INT16_SLICE_PTR, 8)
			val := *param.(*[]int16)
			nLen := len(val)
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(val[i]), 16)
			}
		case "*[]int32":
			bitstream.WriteInt(RPC_INT32_SLICE_PTR, 8)
			val := *param.(*[]int32)
			nLen := len(val)
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(val[i]), 32)
			}
		case "*[]int64":
			bitstream.WriteInt(RPC_INT64_SLICE_PTR, 8)
			val := *param.(*[]int64)
			nLen := len(val)
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt64(val[i], 64)
			}
		case "*[]uint":
			bitstream.WriteInt(RPC_UINT_SLICE_PTR, 8)
			val := *param.(*[]uint)
			nLen := len(val)
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(val[i]), 32)
			}
		case "*[]uint8":
			bitstream.WriteInt(RPC_UINT8_SLICE_PTR, 8)
			val := *param.(*[]uint8)
			nLen := len(val)
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(val[i]), 8)
			}
		case "*[]uint16":
			bitstream.WriteInt(RPC_UINT16_SLICE_PTR, 8)
			val := *param.(*[]uint16)
			nLen := len(val)
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(val[i]), 16)
			}
		case "*[]uint32":
			bitstream.WriteInt(RPC_UINT32_SLICE_PTR, 8)
			val := *param.(*[]uint32)
			nLen := len(val)
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(val[i]), 32)
			}
		case "*[]uint64":
			bitstream.WriteInt(RPC_UINT64_SLICE_PTR, 8)
			val := *param.(*[]uint64)
			nLen := len(val)
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt64(int64(val[i]), 64)
			}




		case "*[]*bool":
			bitstream.WriteInt(RPC_BOOL_PTR_SLICE_PTR, 8)
			val := *param.(*[]*bool)
			nLen := len(val)
			bitstream.WriteInt(nLen, 16)
			for _, v := range val{
				if v != nil{
					bitstream.WriteFlag(*v)
				}else{
					bitstream.WriteFlag(false)
				}
			}
		case "*[]*string":
			bitstream.WriteInt(RPC_STRING_PTR_SLICE_PTR, 8)
			val := *param.(*[]*string)
			nLen := len(val)
			bitstream.WriteInt(nLen, 16)
			for _, v := range val{
				if v != nil{
					bitstream.WriteString(*v)
				}else{
					bitstream.WriteString("")
				}
			}
		case "*[]*float32":
			bitstream.WriteInt(RPC_FLOAT32_PTR_SLICE_PTR, 8)
			val := *param.(*[]*float32)
			nLen := len(val)
			bitstream.WriteInt(nLen, 16)
			for _, v := range val{
				if v != nil{
					bitstream.WriteFloat(*v)
				}else{
					bitstream.WriteFloat(0)
				}
			}
		case "*[]*float64":
			bitstream.WriteInt(RPC_FLOAT64_PTR_SLICE_PTR, 8)
			val := *param.(*[]*float64)
			nLen := len(val)
			bitstream.WriteInt(nLen, 16)
			for _, v := range val{
				if v != nil{
					bitstream.WriteFloat64(*v)
				}else{
					bitstream.WriteFloat64(0)
				}
			}
		case "*[]*int":
			bitstream.WriteInt(RPC_INT_PTR_SLICE_PTR, 8)
			val := *param.(*[]*int)
			nLen := len(val)
			bitstream.WriteInt(nLen, 16)
			for _, v := range val{
				if v != nil{
					bitstream.WriteInt(*v, 32)
				}else{
					bitstream.WriteInt(0, 32)
				}
			}
		case "*[]*int8":
			bitstream.WriteInt(RPC_INT8_PTR_SLICE_PTR, 8)
			val := *param.(*[]*int8)
			nLen := len(val)
			bitstream.WriteInt(nLen, 16)
			for _, v := range val{
				if v != nil{
					bitstream.WriteInt(int(*v), 8)
				}else{
					bitstream.WriteInt(0, 8)
				}
			}
		case "*[]*int16":
			bitstream.WriteInt(RPC_INT16_PTR_SLICE_PTR, 8)
			val := *param.(*[]*int16)
			nLen := len(val)
			bitstream.WriteInt(nLen, 16)
			for _, v := range val{
				if v != nil{
					bitstream.WriteInt(int(*v), 16)
				}else{
					bitstream.WriteInt(0, 16)
				}
			}
		case "*[]*int32":
			bitstream.WriteInt(RPC_INT32_PTR_SLICE_PTR, 8)
			val := *param.(*[]*int32)
			nLen := len(val)
			bitstream.WriteInt(nLen, 16)
			for _, v := range val{
				if v != nil{
					bitstream.WriteInt(int(*v), 32)
				}else{
					bitstream.WriteInt(0, 32)
				}
			}
		case "*[]*int64":
			bitstream.WriteInt(RPC_INT64_PTR_SLICE_PTR, 8)
			val := *param.(*[]*int64)
			nLen := len(val)
			bitstream.WriteInt(nLen, 16)
			for _, v := range val{
				if v != nil{
					bitstream.WriteInt64(*v, 64)
				}else{
					bitstream.WriteInt64(0, 64)
				}
			}
		case "*[]*uint":
			bitstream.WriteInt(RPC_UINT_PTR_SLICE_PTR, 8)
			val := *param.(*[]*uint)
			nLen := len(val)
			bitstream.WriteInt(nLen, 16)
			for _, v := range val{
				if v != nil{
					bitstream.WriteInt(int(*v), 32)
				}else{
					bitstream.WriteInt(0, 32)
				}
			}
		case "*[]*uint8":
			bitstream.WriteInt(RPC_UINT8_PTR_SLICE_PTR, 8)
			val := *param.(*[]*uint8)
			nLen := len(val)
			bitstream.WriteInt(nLen, 16)
			for _, v := range val{
				if v != nil{
					bitstream.WriteInt(int(*v), 8)
				}else{
					bitstream.WriteInt(0, 8)
				}
			}
		case "*[]*uint16":
			bitstream.WriteInt(RPC_UINT16_PTR_SLICE_PTR, 8)
			val := *param.(*[]*uint16)
			nLen := len(val)
			bitstream.WriteInt(nLen, 16)
			for _, v := range val{
				if v != nil{
					bitstream.WriteInt(int(*v), 16)
				}else{
					bitstream.WriteInt(0, 16)
				}
			}
		case "*[]*uint32":
			bitstream.WriteInt(RPC_UINT32_PTR_SLICE_PTR, 8)
			val := *param.(*[]*uint32)
			nLen := len(val)
			bitstream.WriteInt(nLen, 16)
			for _, v := range val{
				if v != nil{
					bitstream.WriteInt(int(*v), 32)
				}else{
					bitstream.WriteInt(0, 32)
				}
			}
		case "*[]*uint64":
			bitstream.WriteInt(RPC_UINT64_PTR_SLICE_PTR, 8)
			val := *param.(*[]*uint64)
			nLen := len(val)
			bitstream.WriteInt(nLen, 16)
			for _, v := range val{
				if v != nil{
					bitstream.WriteInt64(int64(*v), 64)
				}else{
					bitstream.WriteInt64(0, 64)
				}
			}



		case "*[*]bool":
			bitstream.WriteInt(RPC_BOOL_ARRAY_PTR, 8)
			val := reflect.ValueOf(param).Elem()
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteFlag(val.Index(i).Bool())
			}
		case "*[*]string":
			bitstream.WriteInt(RPC_STRING_ARRAY_PTR, 8)
			val := reflect.ValueOf(param).Elem()
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteString(val.Index(i).String())
			}
		case "*[*]float32":
			bitstream.WriteInt(RPC_FLOAT32_ARRAY_PTR, 8)
			val := reflect.ValueOf(param).Elem()
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteFloat(float32(val.Index(i).Float()))
			}
		case "*[*]float64":
			bitstream.WriteInt(RPC_FLOAT64_ARRAY_PTR, 8)
			val := reflect.ValueOf(param).Elem()
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteFloat64(val.Index(i).Float())
			}
		case "*[*]int":
			bitstream.WriteInt(RPC_INT_ARRAY_PTR, 8)
			val := reflect.ValueOf(param).Elem()
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(val.Index(i).Int()), 32)
			}
		case "*[*]int8":
			bitstream.WriteInt(RPC_INT8_ARRAY_PTR, 8)
			val := reflect.ValueOf(param).Elem()
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(val.Index(i).Int()), 8)
			}
		case "*[*]int16":
			bitstream.WriteInt(RPC_INT16_ARRAY_PTR, 8)
			val := reflect.ValueOf(param).Elem()
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(val.Index(i).Int()), 16)
			}
		case "*[*]int32":
			bitstream.WriteInt(RPC_INT32_ARRAY_PTR, 8)
			val := reflect.ValueOf(param).Elem()
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(val.Index(i).Int()), 32)
			}
		case "*[*]int64":
			bitstream.WriteInt(RPC_INT64_ARRAY_PTR, 8)
			val := reflect.ValueOf(param).Elem()
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt64(val.Index(i).Int(), 64)
			}
		case "*[*]uint":
			bitstream.WriteInt(RPC_UINT_ARRAY_PTR, 8)
			val := reflect.ValueOf(param).Elem()
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(val.Index(i).Uint()), 32)
			}
		case "*[*]uint8":
			bitstream.WriteInt(RPC_UINT8_ARRAY_PTR, 8)
			val := reflect.ValueOf(param).Elem()
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(val.Index(i).Uint()), 8)
			}
		case "*[*]uint16":
			bitstream.WriteInt(RPC_UINT16_ARRAY_PTR, 8)
			val := reflect.ValueOf(param).Elem()
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(val.Index(i).Uint()), 16)
			}
		case "*[*]uint32":
			bitstream.WriteInt(RPC_UINT32_ARRAY_PTR, 8)
			val := reflect.ValueOf(param).Elem()
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(val.Index(i).Uint()), 32)
			}
		case "*[*]uint64":
			bitstream.WriteInt(RPC_UINT64_ARRAY_PTR, 8)
			val := reflect.ValueOf(param).Elem()
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt64(int64(val.Index(i).Uint()), 64)
			}



		case "*[*]*bool":
			bitstream.WriteInt(RPC_BOOL_PTR_ARRAY_PTR, 8)
			val := reflect.ValueOf(param).Elem()
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				if !val.Index(i).IsNil(){
					bitstream.WriteFlag(val.Index(i).Elem().Bool())
				}else{
					bitstream.WriteFlag(false)
				}
			}
		case "*[*]*string":
			bitstream.WriteInt(RPC_STRING_PTR_ARRAY_PTR, 8)
			val := reflect.ValueOf(param).Elem()
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				if !val.Index(i).IsNil() {
					bitstream.WriteString(val.Index(i).Elem().String())
				}else{
					bitstream.WriteString("")
				}
			}
		case "*[*]*float32":
			bitstream.WriteInt(RPC_FLOAT32_PTR_ARRAY_PTR, 8)
			val := reflect.ValueOf(param).Elem()
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				if !val.Index(i).IsNil() {
					bitstream.WriteFloat(float32(val.Index(i).Elem().Float()))
				}else{
					bitstream.WriteFloat(0)
				}
			}
		case "*[*]*float64":
			bitstream.WriteInt(RPC_FLOAT64_PTR_ARRAY_PTR, 8)
			val := reflect.ValueOf(param).Elem()
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				if !val.Index(i).IsNil() {
					bitstream.WriteFloat64(val.Index(i).Elem().Float())
				}else{
					bitstream.WriteFloat64(0)
				}
			}
		case "*[*]*int":
			bitstream.WriteInt(RPC_INT_PTR_ARRAY_PTR, 8)
			val := reflect.ValueOf(param).Elem()
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				if !val.Index(i).IsNil() {
					bitstream.WriteInt(int(val.Index(i).Elem().Int()), 32)
				}else{
					bitstream.WriteInt(0, 32)
				}
			}
		case "*[*]*int8":
			bitstream.WriteInt(RPC_INT8_PTR_ARRAY_PTR, 8)
			val := reflect.ValueOf(param).Elem()
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				if !val.Index(i).IsNil() {
					bitstream.WriteInt(int(val.Index(i).Elem().Int()), 8)
				}else{
					bitstream.WriteInt(0, 8)
				}
			}
		case "*[*]*int16":
			bitstream.WriteInt(RPC_INT16_PTR_ARRAY_PTR, 8)
			val := reflect.ValueOf(param).Elem()
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				if !val.Index(i).IsNil() {
					bitstream.WriteInt(int(val.Index(i).Elem().Int()), 16)
				}else{
					bitstream.WriteInt(0, 16)
				}
			}
		case "*[*]*int32":
			bitstream.WriteInt(RPC_INT32_PTR_ARRAY_PTR, 8)
			val := reflect.ValueOf(param).Elem()
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				if !val.Index(i).IsNil() {
					bitstream.WriteInt(int(val.Index(i).Elem().Int()), 32)
				}else{
					bitstream.WriteInt(0, 32)
				}
			}
		case "*[*]*int64":
			bitstream.WriteInt(RPC_INT64_PTR_ARRAY_PTR, 8)
			val := reflect.ValueOf(param).Elem()
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				if !val.Index(i).IsNil() {
					bitstream.WriteInt64(val.Index(i).Elem().Int(), 64)
				}else{
					bitstream.WriteInt64(0, 64)
				}
			}
		case "*[*]*uint":
			bitstream.WriteInt(RPC_UINT_PTR_ARRAY_PTR, 8)
			val := reflect.ValueOf(param).Elem()
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				if !val.Index(i).IsNil() {
					bitstream.WriteInt(int(val.Index(i).Elem().Uint()), 32)
				}else{
					bitstream.WriteInt(0, 32)
				}
			}
		case "*[*]*uint8":
			bitstream.WriteInt(RPC_UINT8_PTR_ARRAY_PTR, 8)
			val := reflect.ValueOf(param).Elem()
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				if !val.Index(i).IsNil() {
					bitstream.WriteInt(int(val.Index(i).Elem().Uint()), 8)
				}else{
					bitstream.WriteInt(0, 8)
				}
			}
		case "*[*]*uint16":
			bitstream.WriteInt(RPC_UINT16_PTR_ARRAY_PTR, 8)
			val := reflect.ValueOf(param).Elem()
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				if !val.Index(i).IsNil() {
					bitstream.WriteInt(int(val.Index(i).Elem().Uint()), 16)
				}else{
					bitstream.WriteInt(0, 16)
				}
			}
		case "*[*]*uint32":
			bitstream.WriteInt(RPC_UINT32_PTR_ARRAY_PTR, 8)
			val := reflect.ValueOf(param).Elem()
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				if !val.Index(i).IsNil() {
					bitstream.WriteInt(int(val.Index(i).Elem().Uint()), 32)
				}else{
					bitstream.WriteInt(0, 32)
				}
			}
		case "*[*]*uint64":
			bitstream.WriteInt(RPC_UINT64_PTR_ARRAY_PTR, 8)
			val := reflect.ValueOf(param).Elem()
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				if !val.Index(i).IsNil() {
					bitstream.WriteInt64(int64(val.Index(i).Elem().Uint()), 64)
				}else{
					bitstream.WriteInt64(0, 64)
				}
			}




		case "*message":
			//rpc  特定rpc头部设置需求，params[0]传入RpcHead
			if i == 0{
				rpcHead, bOk := params[0].(*message.RpcHead)
				if bOk{
					rpcPacket.ArgLen--
					rpcPacket.RpcHead = rpcHead
					continue
				}
			}
			bitstream.WriteInt(RPC_MESSAGE, 8)
			marshalPB(bitstream, param.(proto.Message))



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

	rpcPacket.RpcBody = bitstream.GetBuffer()
	buf, _ := proto.Marshal(rpcPacket)
	return buf, rpcPacket
}

//rpc  MarshalPB
func marshalPB(bitstream *base.BitStream, packet proto.Message) {
	bitstream.WriteString(proto.MessageName(packet))
	buf, _ :=proto.Marshal(packet)
	nLen := len(buf)
	bitstream.WriteInt(nLen, 32)
	bitstream.WriteBits(buf, nLen << 3)
}