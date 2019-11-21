package base

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"log"
	"math"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

const (
	INT_MAX = int(2147483647)
	TCP_END = "💞♡"						//解决tpc粘包半包,结束标志
	//TCP_END = "💞💞💞"				//解决tpc粘包半包,结束标志,-1
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

var(
	SEVERNAME string
	TCP_END_LENGTH = len([]byte(TCP_END)) //tcp结束标志长度
)

func Assert(x bool, y string) {
	if bool(x) == false {
		log.Printf("\nFatal :{%s}", y)
	}
}

func Abs(x float32) float32{
	return float32(math.Abs(float64(x)))
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
func IntToBytes(val int) []byte {
	tmp := uint32(val)
	buff := make([]byte, 4)
	binary.LittleEndian.PutUint32(buff, tmp)
	return buff
}

//字节转换成整形
func BytesToInt(data []byte) int {
	buff := make([]byte, 4)
	for i, v := range data{
		buff[i] = v
	}
	tmp := int32(binary.LittleEndian.Uint32(buff))
	return int(tmp)
}

//整形16转换成字节
func Int16ToBytes(val int16) []byte {
	tmp := uint16(val)
	buff := make([]byte, 2)
	binary.LittleEndian.PutUint16(buff, tmp)
	return buff
}

//字节转换成为int16
func BytesToInt16(data []byte) int16 {
	buff := make([]byte, 2)
	for i, v := range data{
		buff[i] = v
	}
	tmp := binary.LittleEndian.Uint16(buff)
	return int16(tmp)
}

//转化64位
func Int64ToBytes(val int64) []byte {
	tmp := uint64(val)
	buff := make([]byte, 8)
	binary.LittleEndian.PutUint64(buff, tmp)
	return buff
}

func BytesToInt64(data []byte) int64 {
	buff := make([]byte, 8)
	for i, v := range data{
		buff[i] = v
	}
	tmp := binary.LittleEndian.Uint64(buff)
	return int64(tmp)
}

//转化float
func Float32ToByte(val float32) []byte {
	tmp := math.Float32bits(val)
	buff := make([]byte, 4)
	binary.LittleEndian.PutUint32(buff, tmp)
	return buff
}

func BytesToFloat32(data []byte) float32 {
	buff := make([]byte, 4)
	for i, v := range data{
		buff[i] = v
	}
	tmp := binary.LittleEndian.Uint32(buff)
	return math.Float32frombits(tmp)
}

//转化float64
func Float64ToByte(val float64) []byte {
	tmp := math.Float64bits(val)
	buff := make([]byte, 8)
	binary.LittleEndian.PutUint64(buff, tmp)
	return buff
}

func BytesToFloat64(data []byte) float64 {
	buff := make([]byte, 8)
	for i, v := range data{
		buff[i] = v
	}
	tmp := binary.LittleEndian.Uint64(buff)
	return math.Float64frombits(tmp)
}

//[]int转[]int32
func IntToInt32(val []int) []int32 {
	tmp := []int32{}
	for _, v := range val{
		tmp = append(tmp, int32(v))
	}
	return tmp
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

func ChechErr(err error) {
	if err == nil {
		return
	}
	log.Panicf("错误：%s\n", err.Error())
}

func GetDBTime(strTime string) *time.Time{
	DefaultTimeLoc := time.Local
	loginTime, _ := time.ParseInLocation("2006-01-02 15:04:05", strTime, DefaultTimeLoc)
	return &loginTime
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
		sTypeName == "int" || sTypeName == "uint" ||
		sTypeName == "*bool" || sTypeName == "*float64" || sTypeName == "*float32" || sTypeName == "*int8" ||
		sTypeName == "*uint8" || sTypeName == "*int16" || sTypeName == "*uint16" || sTypeName == "*int32" ||
		sTypeName == "*uint32" || sTypeName == "*int64" || sTypeName == "*uint64" ||  sTypeName == "*string"||
		sTypeName == "*int" || sTypeName == "*uint"{
		return "[]" + sTypeName
	}else{
		if strings.Index(sTypeName, "*") != -1{
			return "[]*struct"
		}
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
		sTypeName == "int"  || sTypeName == "uint" ||
		sTypeName == "*bool" || sTypeName == "*float64" || sTypeName == "*float32" || sTypeName == "*int8" ||
		sTypeName == "*uint8" || sTypeName == "*int16" || sTypeName == "*uint16" || sTypeName == "*int32" ||
		sTypeName == "*uint32" || sTypeName == "*int64" || sTypeName == "*uint64" ||  sTypeName == "*string"||
		sTypeName == "*int" || sTypeName == "*uint"{
		return "[*]" + sTypeName
	}else{
		if strings.Index(sTypeName, "*") != -1{
			return "[*]*struct"
		}
		return "[*]struct"
	}

	return sTypeName
}

func GetSliceTypeStringEx(sTypeName string) string{
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

func GetArrayTypeStringEx(sTypeName string) string{
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

func ParseTag(sf reflect.StructField, tag string) map[string]string {
	setting := map[string]string{}
	tags := strings.Split(sf.Tag.Get(tag), ";")
	for _, value := range tags {
		v := strings.Split(value, ":")
		k := strings.TrimSpace(strings.ToLower(v[0]))
		if len(v) >= 2 {
			setting[k] = v[1]
		} else {
			setting[k] = k
		}
	}
	return setting
}

func GetClassName(param interface{}) string{
	sType := strings.ToLower(reflect.ValueOf(param).Type().String())
	index := strings.Index(sType, ".")
	if index!= -1{
		sType = sType[index+1:]
	}
	return sType
}

func GetPacketType(param interface{})string{
	sType := strings.ToLower(reflect.ValueOf(param).Type().String())
	index := strings.Index(sType, ".")
	if index!= -1{
		sType = sType[:index]
	}
	return sType
}

func GetTypeString(param interface{}) string{
	paramType := reflect.TypeOf(param)
	sType := ""
	if paramType.Kind() == reflect.Ptr{
		sType = "*" + paramType.Elem().Kind().String()
		paramType = paramType.Elem()
	}else if paramType.Kind() == reflect.Slice{
		sType = GetSliceTypeString(paramType.String())
	}else if paramType.Kind() == reflect.Array{
		sType = GetArrayTypeString(paramType.String())
	}else{
		sType = paramType.Kind().String()
	}

	if paramType.Kind() == reflect.Struct || paramType.Kind() == reflect.Map || sType == "[]*struct" ||
		sType == "[]struct" || sType == "[*]*struct" || sType == "[*]struct"{
		sType = "*gob"
	}
	return sType
}

func GetTypeStringEx(classField reflect.StructField, classVal reflect.Value) string{
	paramType := classField.Type
	sType := ""
	if paramType.Kind() == reflect.Ptr{
		sType = "*" + paramType.Elem().Kind().String()
	}else if paramType.Kind() == reflect.Slice{
		sType = GetSliceTypeStringEx(paramType.String())
	}else if paramType.Kind() == reflect.Array{
		sType = GetArrayTypeStringEx(paramType.String())
	} else{
		sType = classField.Type.Kind().String()
	}
	return sType
}

//copy name and type is right
func Copy(source interface{}, dest interface{}){
	defer func() {
		if err := recover(); err != nil{
			fmt.Printf("copy source to dest error")
		}
	}()
	getvaltype := func(val interface{}) (reflect.Value, reflect.Type){

		protoType := reflect.TypeOf(val)
		protoVal := reflect.ValueOf(val)
		for protoType.Kind() == reflect.Ptr {
			protoType = protoType.Elem()
			protoVal = protoVal.Elem()
		}

		return protoVal, protoType
	}

	val0, type0 := getvaltype(source)
	val1, type1 := getvaltype(dest)

	for i := 0; i < type0.NumField(); i++{
		if !val0.Field(i).CanSet(){//小写成员只有只读
			continue
		}

		for j := 0; j < type1.NumField(); j++{
			if val1.Field(j).Kind() == reflect.Struct{
				val := val1.Field(j).FieldByName(type0.Field(i).Name)
				if val.IsValid(){
					if val.Type() == type0.Field(i).Type{
						val.Set(val0.Field(i))
					}
				}
			}else{
				val := val1.FieldByName(type0.Field(i).Name)
				if val.IsValid(){
					if val.Type() == type0.Field(i).Type{
						val.Set(val0.Field(i))
					}
				}
			}
		}
	}
}

func ToLower(name string) string{
	return strings.ToLower(name)
}

func SetTcpEnd(buff []byte) []byte{
	buff = append(buff, []byte(TCP_END)...)
	return buff
}

//tcp粘包固定包头
/*func SetTcpEnd(buff []byte) []byte{
	buff = append(base.IntToBytes(len(buff)), buff...)
	return buff
}*/

func ToHash(str string) uint32{
	return crc32.ChecksumIEEE([]byte(str))
}

//-----------string strconv type-------------//
func Int(str string) int{
	n, _ := strconv.Atoi(str)
	return n
}

func Int64(str string) int64{
	n, _ := strconv.ParseInt(str, 0, 64)
	return n
}

func Float32(str string) float32{
	n, _ := strconv.ParseFloat(str, 32)
	return float32(n)
}

func Float64(str string) float64{
	n, _ := strconv.ParseFloat(str, 64)
	return n
}

func Bool(str string) bool{
	n, _ := strconv.ParseBool(str)
	return n
}

func Time(str string) int64 {
	return GetDBTime(str).Unix()
}
//--------------------------------------------//
