package rpc

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/json-iterator/go"
	"gonet/base"
	"gonet/message"
	"log"
	"reflect"
	"strings"
)

//rpc  Marshal
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
			bitstream.WriteFlag(param.(bool))
		case "float64":
			bitstream.WriteInt(RPC_FLOAT64, 8)
			bitstream.WriteFloat64(param.(float64))
		case "float32":
			bitstream.WriteInt(RPC_FLOAT32, 8)
			bitstream.WriteFloat(param.(float32))
		case "int8":
			bitstream.WriteInt(RPC_INT8, 8)
			bitstream.WriteInt(int(param.(int8)), 8)
		case "uint8":
			bitstream.WriteInt(RPC_UINT8, 8)
			bitstream.WriteInt(int(param.(uint8)),8)
		case "int16":
			bitstream.WriteInt(RPC_INT16, 8)
			bitstream.WriteInt(int(param.(int16)),16)
		case "uint16":
			bitstream.WriteInt(RPC_UINT16, 8)
			bitstream.WriteInt(int(param.(uint16)),16)
		case "int32":
			bitstream.WriteInt(RPC_INT32, 8)
			bitstream.WriteInt(int(param.(int32)),32)
		case "uint32":
			bitstream.WriteInt(RPC_UINT32, 8)
			bitstream.WriteInt(int(param.(uint32)),32)
		case "int64":
			bitstream.WriteInt(RPC_INT64, 8)
			bitstream.WriteInt64(param.(int64), 64)
		case "uint64":
			bitstream.WriteInt(RPC_UINT64, 8)
			bitstream.WriteInt64(int64(param.(uint64)), 64)
		case "string":
			bitstream.WriteInt(RPC_STRING, 8)
			bitstream.WriteString(param.(string))
		case "int":
			bitstream.WriteInt(RPC_INT, 8)
			bitstream.WriteInt(param.(int), 32)
		case "uint":
			bitstream.WriteInt(RPC_UINT, 8)
			bitstream.WriteInt(int(param.(uint)), 32)

		case "[]bool":
			bitstream.WriteInt(RPC_BOOL_SLICE, 8)
			nLen := len(param.([]bool))
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteFlag(param.([]bool)[i])
			}
		case "[]float64":
			bitstream.WriteInt(RPC_FLOAT64_SLICE, 8)
			nLen := len(param.([]float64))
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteFloat64(param.([]float64)[i])
			}
		case "[]float32":
			bitstream.WriteInt(RPC_FLOAT32_SLICE, 8)
			nLen := len(param.([]float32))
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteFloat(param.([]float32)[i])
			}
		case "[]int8":
			bitstream.WriteInt(RPC_INT8_SLICE, 8)
			nLen := len(param.([]int8))
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(param.([]int8)[i]), 8)
			}
		case "[]uint8":
			bitstream.WriteInt(RPC_UINT8_SLICE, 8)
			nLen := len(param.([]uint8))
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(param.([]uint8)[i]), 8)
			}
		case "[]int16":
			bitstream.WriteInt(RPC_INT16_SLICE, 8)
			nLen := len(param.([]int16))
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(param.([]int16)[i]), 16)
			}
		case "[]uint16":
			bitstream.WriteInt(RPC_UINT16_SLICE, 8)
			nLen := len(param.([]uint16))
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(param.([]uint16)[i]), 16)
			}
		case "[]int32":
			bitstream.WriteInt(RPC_INT32_SLICE, 8)
			nLen := len(param.([]int32))
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(param.([]int32)[i]), 32)
			}
		case "[]uint32":
			bitstream.WriteInt(RPC_UINT32_SLICE, 8)
			nLen := len(param.([]uint32))
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(param.([]uint32)[i]), 32)
			}
		case "[]int64":
			bitstream.WriteInt(RPC_INT64_SLICE, 8)
			nLen := len(param.([]int64))
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt64(param.([]int64)[i], 64)
			}
		case "[]uint64":
			bitstream.WriteInt(RPC_UINT64_SLICE, 8)
			nLen := len(param.([]uint64))
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt64(int64(param.([]uint64)[i]), 64)
			}
		case "[]string":
			bitstream.WriteInt(RPC_STRING_SLICE, 8)
			nLen := len(param.([]string))
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteString(param.([]string)[i])
			}
		case "[]int":
			bitstream.WriteInt(RPC_INT_SLICE, 8)
			nLen := len(param.([]int))
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(param.([]int)[i], 32)
			}
		case "[]uint":
			bitstream.WriteInt(RPC_UINT_SLICE, 8)
			nLen := len(param.([]uint))
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(param.([]uint)[i]), 32)
			}



		case "[*]bool":
			bitstream.WriteInt(RPC_BOOL_ARRAY, 8)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteFlag(val.Index(i).Bool())
			}
		case "[*]float64":
			bitstream.WriteInt(RPC_FLOAT64_ARRAY, 8)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteFloat64(val.Index(i).Float())
			}
		case "[*]float32":
			bitstream.WriteInt(RPC_FLOAT32_ARRAY, 8)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteFloat(float32(val.Index(i).Float()))
			}
		case "[*]int8":
			bitstream.WriteInt(RPC_INT8_ARRAY, 8)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(val.Index(i).Int()), 8)
			}
		case "[*]uint8":
			bitstream.WriteInt(RPC_UINT8_ARRAY, 8)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(val.Index(i).Uint()), 8)
			}
		case "[*]int16":
			bitstream.WriteInt(RPC_INT16_ARRAY, 8)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(val.Index(i).Int()), 16)
			}
		case "[*]uint16":
			bitstream.WriteInt(RPC_UINT16_ARRAY, 8)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(val.Index(i).Uint()), 16)
			}
		case "[*]int32":
			bitstream.WriteInt(RPC_INT32_ARRAY, 8)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(val.Index(i).Int()), 32)
			}
		case "[*]uint32":
			bitstream.WriteInt(RPC_UINT32_ARRAY, 8)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(val.Index(i).Uint()), 32)
			}
		case "[*]int64":
			bitstream.WriteInt(RPC_INT64_ARRAY, 8)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt64(val.Index(i).Int(), 64)
			}
		case "[*]uint64":
			bitstream.WriteInt(RPC_UINT64_ARRAY, 8)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt64(int64(val.Index(i).Uint()), 64)
			}
		case "[*]string":
			bitstream.WriteInt(RPC_STRING_ARRAY, 8)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteString(val.Index(i).String())
			}
		case "[*]int":
			bitstream.WriteInt(RPC_INT_ARRAY, 8)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(val.Index(i).Int()), 32)
			}
		case "[*]uint":
			bitstream.WriteInt(RPC_UINT_ARRAY, 8)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(val.Index(i).Uint()), 32)
			}



		case "*bool":
			bitstream.WriteInt(RPC_BOOL_PTR, 8)
			if param.(*bool) != nil{
				bitstream.WriteFlag(*param.(*bool))
			}else{
				bitstream.WriteFlag(false)
			}
		case "*float64":
			bitstream.WriteInt(RPC_FLOAT64_PTR, 8)
			if param.(*float64) != nil {
				bitstream.WriteFloat64(*param.(*float64))
			}else{
				bitstream.WriteFloat64(0)
			}
		case "*float32":
			bitstream.WriteInt(RPC_FLOAT32_PTR, 8)
			if param.(*float32) != nil {
				bitstream.WriteFloat(*param.(*float32))
			}else{
				bitstream.WriteFloat(0)
			}
		case "*int8":
			bitstream.WriteInt(RPC_INT8_PTR, 8)
			if param.(*int8) != nil {
				bitstream.WriteInt(int(*param.(*int8)), 8)
			}else{
				bitstream.WriteInt(0, 8)
			}
		case "*uint8":
			bitstream.WriteInt(RPC_UINT8_PTR, 8)
			if param.(*uint8) != nil {
				bitstream.WriteInt(int(*param.(*uint8)), 8)
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
		case "*uint16":
			bitstream.WriteInt(RPC_UINT16_PTR, 8)
			if param.(*uint16) != nil {
				bitstream.WriteInt(int(*param.(*uint16)), 16)
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
		case "*uint32":
			bitstream.WriteInt(RPC_UINT32_PTR, 8)
			if param.(*uint32) != nil {
				bitstream.WriteInt(int(*param.(*uint32)), 32)
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
		case "*uint64":
			bitstream.WriteInt(RPC_UINT64_PTR, 8)
			if param.(*uint64) != nil {
				bitstream.WriteInt64(int64(*param.(*uint64)), 64)
			}else{
				bitstream.WriteInt64(0, 64)
			}
		case "*string":
			bitstream.WriteInt(RPC_STRING_PTR, 8)
			if param.(*string) != nil {
				bitstream.WriteString(*param.(*string))
			}else{
				bitstream.WriteString("")
			}
		case "*int":
			bitstream.WriteInt(RPC_INT_PTR, 8)
			if param.(*int) != nil {
				bitstream.WriteInt(*param.(*int), 32)
			}else{
				bitstream.WriteInt(0, 32)
			}
		case "*uint":
			bitstream.WriteInt(RPC_UINT_PTR, 8)
			if param.(*uint) != nil {
				bitstream.WriteInt(int(*param.(*uint)), 32)
			}else{
				bitstream.WriteInt(0, 32)
			}



		case "[]*bool":
			bitstream.WriteInt(RPC_BOOL_PTR_SLICE, 8)
			nLen := len(param.([]*bool))
			bitstream.WriteInt(nLen, 16)
			for _, v := range param.([]*bool){
				if v != nil{
					bitstream.WriteFlag(*v)
				}else{
					bitstream.WriteFlag(false)
				}
			}
		case "[]*float64":
			bitstream.WriteInt(RPC_FLOAT64_PTR_SLICE, 8)
			nLen := len(param.([]float64))
			bitstream.WriteInt(nLen, 16)
			for _, v := range param.([]*float64){
				if v != nil{
					bitstream.WriteFloat64(*v)
				}else{
					bitstream.WriteFloat64(0)
				}
			}
		case "[]*float32":
			bitstream.WriteInt(RPC_FLOAT32_PTR_SLICE, 8)
			nLen := len(param.([]float32))
			bitstream.WriteInt(nLen, 16)
			for _, v := range param.([]*float32){
				if v != nil{
					bitstream.WriteFloat(*v)
				}else{
					bitstream.WriteFloat(0)
				}
			}
		case "[]*int8":
			bitstream.WriteInt(RPC_INT8_PTR_SLICE, 8)
			nLen := len(param.([]int8))
			bitstream.WriteInt(nLen, 16)
			for _, v := range param.([]*int8){
				if v != nil{
					bitstream.WriteInt(int(*v), 8)
				}else{
					bitstream.WriteInt(0, 8)
				}
			}
		case "[]*uint8":
			bitstream.WriteInt(RPC_UINT8_PTR_SLICE, 8)
			nLen := len(param.([]uint8))
			bitstream.WriteInt(nLen, 16)
			for _, v := range param.([]*uint8){
				if v != nil{
					bitstream.WriteInt(int(*v), 8)
				}else{
					bitstream.WriteInt(0, 8)
				}
			}
		case "[]*int16":
			bitstream.WriteInt(RPC_INT16_PTR_SLICE, 8)
			nLen := len(param.([]int16))
			bitstream.WriteInt(nLen, 16)
			for _, v := range param.([]*int16){
				if v != nil{
					bitstream.WriteInt(int(*v), 16)
				}else{
					bitstream.WriteInt(0, 16)
				}
			}
		case "[]*uint16":
			bitstream.WriteInt(RPC_UINT16_PTR_SLICE, 8)
			nLen := len(param.([]uint16))
			bitstream.WriteInt(nLen, 16)
			for _, v := range param.([]*uint16){
				if v != nil{
					bitstream.WriteInt(int(*v), 16)
				}else{
					bitstream.WriteInt(0, 16)
				}
			}
		case "[]*int32":
			bitstream.WriteInt(RPC_INT32_PTR_SLICE, 8)
			nLen := len(param.([]int32))
			bitstream.WriteInt(nLen, 16)
			for _, v := range param.([]*int32){
				if v != nil{
					bitstream.WriteInt(int(*v), 32)
				}else{
					bitstream.WriteInt(0, 32)
				}
			}
		case "[]*uint32":
			bitstream.WriteInt(RPC_UINT32_PTR_SLICE, 8)
			nLen := len(param.([]uint32))
			bitstream.WriteInt(nLen, 16)
			for _, v := range param.([]*uint32){
				if v != nil{
					bitstream.WriteInt(int(*v), 32)
				}else{
					bitstream.WriteInt(0, 32)
				}

			}
		case "[]*int64":
			bitstream.WriteInt(RPC_INT64_PTR_SLICE, 8)
			nLen := len(param.([]int64))
			bitstream.WriteInt(nLen, 16)
			for _, v := range param.([]*int64){
				if v != nil{
					bitstream.WriteInt64(*v, 64)
				}else{
					bitstream.WriteInt64(0, 64)
				}
			}
		case "[]*uint64":
			bitstream.WriteInt(RPC_UINT64_PTR_SLICE, 8)
			nLen := len(param.([]uint64))
			bitstream.WriteInt(nLen, 16)
			for _, v := range param.([]*uint64){
				if v != nil{
					bitstream.WriteInt64(int64(*v), 64)
				}else{
					bitstream.WriteInt64(0, 64)
				}
			}
		case "[]*string":
			bitstream.WriteInt(RPC_STRING_PTR_SLICE, 8)
			nLen := len(param.([]string))
			bitstream.WriteInt(nLen, 16)
			for _, v := range param.([]*string){
				if v != nil{
					bitstream.WriteString(*v)
				}else{
					bitstream.WriteString("")
				}
			}
		case "[]*int":
			bitstream.WriteInt(RPC_INT_PTR_SLICE, 8)
			nLen := len(param.([]*int))
			bitstream.WriteInt(nLen, 16)
			for _, v := range param.([]*int){
				if v != nil{
					bitstream.WriteInt(*v, 32)
				}else{
					bitstream.WriteInt(0, 32)
				}
			}
		case "[]*uint":
			bitstream.WriteInt(RPC_UINT_PTR_SLICE, 8)
			nLen := len(param.([]uint))
			bitstream.WriteInt(nLen, 16)
			for _, v := range param.([]*int){
				if v != nil{
					bitstream.WriteInt(int(*v), 32)
				}else{
					bitstream.WriteInt(0, 32)
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



		case "*gob":
			bitstream.WriteInt(RPC_GOB, 8)
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			buf, _ := json.Marshal(param)
			nLen := len(buf)
			bitstream.WriteInt(nLen, 32)
			bitstream.WriteBits(nLen << 3, buf)
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
func MarshalPB(packet proto.Message, bitstream *base.BitStream) bool {
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
}
