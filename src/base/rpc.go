package base

import (
	"fmt"
	"reflect"
)

const(
	RPC_Int64 		= 10
	RPC_UInt64 		= 11
	RPC_PInt64 		= 70
	RPC_PUInt64 	= 71
	RPC_Message 	= 120
)

func GetPacket(funcName string, params ...interface{})[]byte {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("GetPacket", err)
		}
	}()

	msg := make([]byte, 1024)
	bitstream := NewBitStream(msg, 1024)
	bitstream.WriteString(funcName)
	bitstream.WriteInt(len(params), 8)
	for _, param := range params {
		sType := GetTypeString(param)
		switch sType {
		case "bool":
			bitstream.WriteInt(1, 8)
			bitstream.WriteFlag(param.(bool))
		case "float64":
			bitstream.WriteInt(2, 8)
			bitstream.WriteFloat64(param.(float64))
		case "float32":
			bitstream.WriteInt(3, 8)
			bitstream.WriteFloat(param.(float32))
		case "int8":
			bitstream.WriteInt(4, 8)
			bitstream.WriteInt(int(param.(int8)), 8)
		case "uint8":
			bitstream.WriteInt(5, 8)
			bitstream.WriteInt(int(param.(uint8)),8)
		case "int16":
			bitstream.WriteInt(6, 8)
			bitstream.WriteInt(int(param.(int16)),16)
		case "uint16":
			bitstream.WriteInt(7, 8)
			bitstream.WriteInt(int(param.(uint16)),16)
		case "int32":
			bitstream.WriteInt(8, 8)
			bitstream.WriteInt(int(param.(int32)),32)
		case "uint32":
			bitstream.WriteInt(9, 8)
			bitstream.WriteInt(int(param.(uint32)),32)
		case "int64":
			bitstream.WriteInt(10, 8)
			bitstream.WriteInt64(param.(int64), 64)
		case "uint64":
			bitstream.WriteInt(11, 8)
			bitstream.WriteInt64(int64(param.(uint64)), 64)
		case "string":
			bitstream.WriteInt(12, 8)
			bitstream.WriteString(param.(string))
		case "int":
			bitstream.WriteInt(13, 8)
			bitstream.WriteInt(param.(int), 32)
		case "uint":
			bitstream.WriteInt(14, 8)
			bitstream.WriteInt(int(param.(uint)), 32)


		case "[]bool":
			bitstream.WriteInt(21, 8)
			nLen := len(param.([]bool))
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteFlag(param.([]bool)[i])
			}
		case "[]float64":
			bitstream.WriteInt(22, 8)
			nLen := len(param.([]float64))
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteFloat64(param.([]float64)[i])
			}
		case "[]float32":
			bitstream.WriteInt(23, 8)
			nLen := len(param.([]float32))
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteFloat(param.([]float32)[i])
			}
		case "[]int8":
			bitstream.WriteInt(24, 8)
			nLen := len(param.([]int8))
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(param.([]int8)[i]), 8)
			}
		case "[]uint8":
			bitstream.WriteInt(25, 8)
			nLen := len(param.([]uint8))
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(param.([]uint8)[i]), 8)
			}
		case "[]int16":
			bitstream.WriteInt(26, 8)
			nLen := len(param.([]int16))
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(param.([]int16)[i]), 16)
			}
		case "[]uint16":
			bitstream.WriteInt(27, 8)
			nLen := len(param.([]uint16))
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(param.([]uint16)[i]), 16)
			}
		case "[]int32":
			bitstream.WriteInt(28, 8)
			nLen := len(param.([]int32))
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(param.([]int32)[i]), 32)
			}
		case "[]uint32":
			bitstream.WriteInt(29, 8)
			nLen := len(param.([]uint32))
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(param.([]uint32)[i]), 32)
			}
		case "[]int64":
			bitstream.WriteInt(30, 8)
			nLen := len(param.([]int64))
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt64(param.([]int64)[i], 64)
			}
		case "[]uint64":
			bitstream.WriteInt(31, 8)
			nLen := len(param.([]uint64))
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt64(int64(param.([]uint64)[i]), 64)
			}
		case "[]string":
			bitstream.WriteInt(32, 8)
			nLen := len(param.([]string))
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteString(param.([]string)[i])
			}
		case "[]int":
			bitstream.WriteInt(33, 8)
			nLen := len(param.([]int))
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(param.([]int)[i], 32)
			}
		case "[]uint":
			bitstream.WriteInt(34, 8)
			nLen := len(param.([]uint))
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(param.([]uint)[i]), 32)
			}
		case "[]struct"://结构体必须重写WriteData and ReadData
			bitstream.WriteInt(35, 8)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteString(getMessageName(val.Index(i).Addr().Interface()))
				WriteData(val.Index(i).Addr().Interface(), bitstream)
				//val.Index(i).Addr().Interface().(Message).WriteData(bitstream)
				//bitstream.WriteString(getMessageName(val.Index(i).Addr().Interface().(Message)))
				//val.Index(i).Addr().Interface().(Message).WriteData(bitstream)
			}


		case "[*]bool":
			bitstream.WriteInt(41, 8)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteFlag(val.Index(i).Bool())
			}
		case "[*]float64":
			bitstream.WriteInt(42, 8)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteFloat64(val.Index(i).Float())
			}
		case "[*]float32":
			bitstream.WriteInt(43, 8)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteFloat(float32(val.Index(i).Float()))
			}
		case "[*]int8":
			bitstream.WriteInt(44, 8)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(val.Index(i).Int()), 8)
			}
		case "[*]uint8":
			bitstream.WriteInt(45, 8)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(val.Index(i).Uint()), 8)
			}
		case "[*]int16":
			bitstream.WriteInt(46, 8)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(val.Index(i).Int()), 16)
			}
		case "[*]uint16":
			bitstream.WriteInt(47, 8)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(val.Index(i).Uint()), 16)
			}
		case "[*]int32":
			bitstream.WriteInt(48, 8)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(val.Index(i).Int()), 32)
			}
		case "[*]uint32":
			bitstream.WriteInt(49, 8)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(val.Index(i).Uint()), 32)
			}
		case "[*]int64":
			bitstream.WriteInt(50, 8)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt64(val.Index(i).Int(), 64)
			}
		case "[*]uint64":
			bitstream.WriteInt(51, 8)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt64(int64(val.Index(i).Uint()), 64)
			}
		case "[*]string":
			bitstream.WriteInt(52, 8)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteString(val.Index(i).String())
			}
		case "[*]int":
			bitstream.WriteInt(53, 8)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(val.Index(i).Int()), 32)
			}
		case "[*]uint":
			bitstream.WriteInt(54, 8)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(val.Index(i).Uint()), 32)
			}
		/*case "[*]struct"://结构体必须重写WriteData and ReadData
			bitstream.WriteInt(55, 8)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteString(getMessageName(val.Index(i).Interface()))
				WriteData(val.Index(i), bitstream)
			}*/


		case "*bool":
			if param.(*bool) != nil{
				bitstream.WriteInt(61, 8)
				bitstream.WriteFlag(*param.(*bool))
			}else{
				bitstream.WriteInt(61, 8)
				bitstream.WriteFlag(false)
			}
		case "*float64":
			if param.(*float64) != nil {
				bitstream.WriteInt(62, 8)
				bitstream.WriteFloat64(*param.(*float64))
			}else{
				bitstream.WriteInt(62, 8)
				bitstream.WriteFloat64(0)
			}
		case "*float32":
			if param.(*float32) != nil {
				bitstream.WriteInt(63, 8)
				bitstream.WriteFloat(*param.(*float32))
			}else{
				bitstream.WriteInt(63, 8)
				bitstream.WriteFloat(0)
			}
		case "*int8":
			if param.(*int8) != nil {
				bitstream.WriteInt(64, 8)
				bitstream.WriteInt(int(*param.(*int8)), 8)
			}else{
				bitstream.WriteInt(64, 8)
				bitstream.WriteInt(0, 8)
			}
		case "*uint8":
			if param.(*uint8) != nil {
				bitstream.WriteInt(65, 8)
				bitstream.WriteInt(int(*param.(*uint8)), 8)
			}else{
				bitstream.WriteInt(65, 8)
				bitstream.WriteInt(0, 8)
			}
		case "*int16":
			if param.(*int16) != nil {
				bitstream.WriteInt(66, 8)
				bitstream.WriteInt(int(*param.(*int16)), 16)
			}else{
				bitstream.WriteInt(66, 8)
				bitstream.WriteInt(0, 16)
			}
		case "*uint16":
			if param.(*uint16) != nil {
				bitstream.WriteInt(67, 8)
				bitstream.WriteInt(int(*param.(*uint16)), 16)
			}else{
				bitstream.WriteInt(67, 8)
				bitstream.WriteInt(0, 16)
			}
		case "*int32":
			if param.(*int32) != nil {
				bitstream.WriteInt(68, 8)
				bitstream.WriteInt(int(*param.(*int32)), 32)
			}else{
				bitstream.WriteInt(68, 8)
				bitstream.WriteInt(0, 32)
			}
		case "*uint32":
			if param.(*uint32) != nil {
				bitstream.WriteInt(69, 8)
				bitstream.WriteInt(int(*param.(*uint32)), 32)
			}else{
				bitstream.WriteInt(69, 8)
				bitstream.WriteInt(0, 32)
			}
		case "*int64":
			if param.(*int64) != nil {
				bitstream.WriteInt(70, 8)
				bitstream.WriteInt64(*param.(*int64), 64)
			}else{
				bitstream.WriteInt(70, 8)
				bitstream.WriteInt64(0, 64)
			}
		case "*uint64":
			if param.(*uint64) != nil {
				bitstream.WriteInt(71, 8)
				bitstream.WriteInt64(int64(*param.(*uint64)), 64)
			}else{
				bitstream.WriteInt(71, 8)
				bitstream.WriteInt64(0, 64)
			}
		case "*string":
			if param.(*string) != nil {
				bitstream.WriteInt(72, 8)
				bitstream.WriteString(*param.(*string))
			}else{
				bitstream.WriteInt(72, 8)
				bitstream.WriteString("")
			}
		case "*int":
			if param.(*int) != nil {
				bitstream.WriteInt(73, 8)
				bitstream.WriteInt(*param.(*int), 32)
			}else{
				bitstream.WriteInt(73, 8)
				bitstream.WriteInt(0, 32)
			}
		case "*uint":
			if param.(*uint) != nil {
				bitstream.WriteInt(74, 8)
				bitstream.WriteInt(int(*param.(*uint)), 32)
			}else{
				bitstream.WriteInt(74, 8)
				bitstream.WriteInt(0, 32)
			}
		case "*struct"://结构体必须重写WriteData and ReadData
			bitstream.WriteInt(75, 8)
			bitstream.WriteString(getMessageName(param))
			WriteData(param, bitstream)



		case "[]*bool":
			bitstream.WriteInt(81, 8)
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
			bitstream.WriteInt(82, 8)
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
			bitstream.WriteInt(83, 8)
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
			bitstream.WriteInt(84, 8)
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
			bitstream.WriteInt(85, 8)
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
			bitstream.WriteInt(86, 8)
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
			bitstream.WriteInt(87, 8)
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
			bitstream.WriteInt(88, 8)
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
			bitstream.WriteInt(89, 8)
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
			bitstream.WriteInt(90, 8)
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
			bitstream.WriteInt(91, 8)
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
			bitstream.WriteInt(92, 8)
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
			bitstream.WriteInt(93, 8)
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
			bitstream.WriteInt(94, 8)
			nLen := len(param.([]uint))
			bitstream.WriteInt(nLen, 16)
			for _, v := range param.([]*int){
				if v != nil{
					bitstream.WriteInt(int(*v), 32)
				}else{
					bitstream.WriteInt(0, 32)
				}
			}
		case "[]*struct"://结构体必须重写WriteData and ReadData
			bitstream.WriteInt(95, 8)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteString(getMessageName(val.Index(i).Interface()))
				WriteData(val.Index(i).Interface(), bitstream)
			}


		case "[*]*bool":
			bitstream.WriteInt(101, 8)
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
			bitstream.WriteInt(102, 8)
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
			bitstream.WriteInt(103, 8)
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
			bitstream.WriteInt(104, 8)
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
			bitstream.WriteInt(105, 8)
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
			bitstream.WriteInt(106, 8)
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
			bitstream.WriteInt(107, 8)
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
			bitstream.WriteInt(108, 8)
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
			bitstream.WriteInt(109, 8)
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
			bitstream.WriteInt(110, 8)
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
			bitstream.WriteInt(111, 8)
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
			bitstream.WriteInt(112, 8)
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
			bitstream.WriteInt(113, 8)
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
			bitstream.WriteInt(114, 8)
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
		case "[*]*struct"://结构体必须重写WriteData and ReadData
			bitstream.WriteInt(115, 8)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteString(getMessageName(val.Index(i).Interface()))
				WriteData(val.Index(i).Interface(), bitstream)
			}

		default:
			fmt.Println("params type not supported", sType,  reflect.TypeOf(param))
			panic("params type not supported")
		}
	}

	return bitstream.GetBuffer()
}
