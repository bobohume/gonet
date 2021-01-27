package base

import (
	"encoding/binary"
	"fmt"
	"gonet/base/vector"
	"hash/crc32"
	"io/ioutil"
	"log"
	"math"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	INT_MAX = int(2147483647)
	TCP_HEAD_SIZE = 4			//解决tpc粘包半包,包头固定长度
	TCP_END = "💞♡"					//解决tpc粘包半包,特殊结束标志,pb采用Varint编码高位有特殊含义
)

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

func Max(a, b int) int{
	return int(math.Max(float64(a), float64(b)))
}

func Min(a, b int) int{
	return int(math.Min(float64(a), float64(b)))
}

func Clamp(val, low,  high int) int {
	return int(math.Max(math.Min(float64(val), float64(high)), float64(low)))
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
	copy(buff, data)
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
	copy(buff, data)
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
	copy(buff, data)
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
	copy(buff, data)
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
	copy(buff, data)
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

func ToLower(name string) string{
	return strings.ToLower(name)
}

func SetTcpEnd(buff []byte) []byte{
	buff = append(IntToBytes(len(buff)), buff...)
	return buff
}

//tcp粘包特殊结束标志
/*func SetTcpEnd(buff []byte) []byte{
	buff = append(buff, []byte(TCP_END)...)
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
// ToString converts a value to string.
func ToString(value interface{}) string {
	switch value.(type) {
	case string:
		return value.(string)
	case int8:
		return strconv.FormatInt(int64(value.(int8)), 10)
	case int16:
		return strconv.FormatInt(int64(value.(int16)), 10)
	case int32:
		return strconv.FormatInt(int64(value.(int32)), 10)
	case int64:
		return strconv.FormatInt(int64(value.(int64)), 10)
	case uint8:
		return strconv.FormatUint(uint64(value.(uint8)), 10)
	case uint16:
		return strconv.FormatUint(uint64(value.(uint16)), 10)
	case uint32:
		return strconv.FormatUint(uint64(value.(uint32)), 10)
	case uint64:
		return strconv.FormatUint(uint64(value.(uint64)), 10)
	case float32:
		return strconv.FormatFloat(float64(value.(float32)), 'g', -1, 64)
	case float64:
		return strconv.FormatFloat(float64(value.(float64)), 'g', -1, 64)
	case bool:
		return strconv.FormatBool(value.(bool))
	default:
		return fmt.Sprintf("%+v", value)
	}
}

//---------遍历子目录------------//
func WalkDir(dirpath string, filesVec *vector.Vector){
	files, err := ioutil.ReadDir(dirpath)//读取目录下文件
	if err != nil{
		return
	}
	for _, file := range files{
		if file.IsDir(){
			WalkDir(dirpath + "/" + file.Name(), filesVec)
			continue
		}else{
			filesVec.PushBack(file)
		}
	}
}