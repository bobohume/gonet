package base

import (
	"bytes"
	"encoding/binary"
	"log"
	"math"
	"os"
	"fmt"
	"reflect"
	"strings"
)

const (
	INT_MAX = int(2147483647)
	TCP_END = "#@"						//解决tpc粘包半包,结束标志
)

var(
	SEVERNAME string
)

func Assert(x bool, y string) {
	if bool(x) == false {
		log.Printf("\nFatal :{%s}", y)
	}
}

func IFAssert(x bool, y string) {
	if bool(x) == false {
		log.Fatalf("\nFatal :{%s}", y)
	}
}

func BIT(x interface{}) interface{}{
	return (1 << x.(uint32))
}

func BIT64(x interface{}) interface{}{
	return (1 << x.(uint64))
}

//整形转换成字节
func IntToBytes(n int) []byte {
	tmp := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.LittleEndian, tmp)
	return bytesBuffer.Bytes()
}

//字节转换成整形
func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)
	var tmp int32
	binary.Read(bytesBuffer, binary.LittleEndian, &tmp)
	return int(tmp)
}

//字节转换成为int16
func BytesToInt16(b []byte) int16 {
	bytesBuffer := bytes.NewBuffer(b)
	var tmp int16
	binary.Read(bytesBuffer, binary.LittleEndian, &tmp)
	return int16(tmp)
}

//转化64位
func Int64ToBytes(n int64) []byte {
	tmp := uint64(n)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, tmp)
	return bytes
}

func BytesToInt64(b []byte) int64 {
	var tmp uint64
	tmp = binary.LittleEndian.Uint64(b)
	return int64(tmp)
}

//转化float
func Float32ToByte(float float32) []byte {
	bits := math.Float32bits(float)
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, bits)
	return bytes
}

func ByteToFloat32(b []byte) float32 {
	bits := binary.LittleEndian.Uint32(b)
	return math.Float32frombits(bits)
}

//转化float64
func Float64ToByte(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)
	return bytes
}

func ByteToFloat64(b []byte) float64 {
	bits := binary.LittleEndian.Uint64(b)
	return math.Float64frombits(bits)
}

func Htons(n uint16) []byte{
	bytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(bytes, n)
	return bytes
}

func Htonl(n uint64) []byte{
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, n)
	return bytes
}

func PathExists(path string) (bool) {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func GetSliceTypeString(sTypeName string) string{
	index := strings.Index(sTypeName, "]")
	if index != -1{
		sTypeName = sTypeName[index+1:]
	}

	if sTypeName == "bool" || sTypeName == "float64" || sTypeName == "float32" || sTypeName == "int8" ||
		sTypeName == "uint8" || sTypeName == "int16" || sTypeName == "uint16" || sTypeName == "int32" ||
		sTypeName == "uint32" || sTypeName == "int64" || sTypeName == "uint64" ||  sTypeName == "string"||
		sTypeName == "int" || sTypeName == "uint"{
		return "[]" + sTypeName
	}else{
		return "[]struct"
	}

	return sTypeName
}

func GetArrayTypeString(sTypeName string) string{
	index := strings.Index(sTypeName, "]")
	if index != -1{
		sTypeName = sTypeName[index+1:]
	}

	if sTypeName == "bool" || sTypeName == "float64" || sTypeName == "float32" || sTypeName == "int8" ||
		sTypeName == "uint8" || sTypeName == "int16" || sTypeName == "uint16" || sTypeName == "int32" ||
		sTypeName == "uint32" || sTypeName == "int64" || sTypeName == "uint64" ||  sTypeName == "string"||
		sTypeName == "int"  || sTypeName == "uint"{
		return "[*]" + sTypeName
	}else{
		return "[*]struct"
	}

	return sTypeName
}

func GetTypeString(param interface{}) string{
	paramType := reflect.TypeOf(param)
	sType := ""
	if paramType.Kind() == reflect.Ptr{
		sType = "*" + paramType.Elem().Kind().String()
	}else if paramType.Kind() == reflect.Slice{
		sType = GetSliceTypeString(paramType.String())
	}else if paramType.Kind() == reflect.Array{
		sType = GetArrayTypeString(paramType.String())
	}else{
		sType = paramType.Kind().String()
	}
	return sType
}

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
			bitstream.WriteInt(15, 8)
			nLen := len(param.([]bool))
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteFlag(param.([]bool)[i])
			}
		case "[]float64":
			bitstream.WriteInt(16, 8)
			nLen := len(param.([]float64))
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteFloat64(param.([]float64)[i])
			}
		case "[]float32":
			bitstream.WriteInt(17, 8)
			nLen := len(param.([]float32))
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteFloat(param.([]float32)[i])
			}
		case "[]int8":
			bitstream.WriteInt(18, 8)
			nLen := len(param.([]int8))
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(param.([]int8)[i]), 8)
			}
		case "[]uint8":
			bitstream.WriteInt(19, 8)
			nLen := len(param.([]uint8))
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(param.([]uint8)[i]), 8)
			}
		case "[]int16":
			bitstream.WriteInt(20, 8)
			nLen := len(param.([]int16))
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(param.([]int16)[i]), 16)
			}
		case "[]uint16":
			bitstream.WriteInt(21, 8)
			nLen := len(param.([]uint16))
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(param.([]uint16)[i]), 16)
			}
		case "[]int32":
			bitstream.WriteInt(22, 8)
			nLen := len(param.([]int32))
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(param.([]int32)[i]), 32)
			}
		case "[]uint32":
			bitstream.WriteInt(23, 8)
			nLen := len(param.([]uint32))
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(param.([]uint32)[i]), 32)
			}
		case "[]int64":
			bitstream.WriteInt(24, 8)
			nLen := len(param.([]int64))
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt64(param.([]int64)[i], 64)
			}
		case "[]uint64":
			bitstream.WriteInt(25, 8)
			nLen := len(param.([]uint64))
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt64(int64(param.([]uint64)[i]), 64)
			}
		case "[]string":
			bitstream.WriteInt(26, 8)
			nLen := len(param.([]string))
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteString(param.([]string)[i])
			}
		case "[]int":
			bitstream.WriteInt(27, 8)
			nLen := len(param.([]int))
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(param.([]int)[i], 32)
			}
		case "[]uint":
			bitstream.WriteInt(28, 8)
			nLen := len(param.([]uint))
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(param.([]uint)[i]), 32)
			}
		case "*struct"://结构体必须重写WriteData and ReadData
			bitstream.WriteInt(29, 8)
			bitstream.WriteString(getMessageName(param.(Message)))
			param.(Message).WriteData(bitstream)
		case "[]struct"://结构体必须重写WriteData and ReadData
			bitstream.WriteInt(30, 8)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteString(getMessageName(val.Index(i).Interface().(Message)))
				val.Index(i).Interface().(Message).WriteData(bitstream)
			}
		default:
			fmt.Println("params type not supported", sType,  reflect.TypeOf(param))
			panic("params type not supported")
		}
	}

	return bitstream.GetBuffer()
}

func ToLower(name string) string{
	return strings.ToLower(name)
}

func SetTcpEnd(buff []byte) []byte{
	buff = append(buff, []byte(TCP_END)...)
	return buff
}