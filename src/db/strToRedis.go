package db

import (
	"reflect"
	"fmt"
	"strings"
	"strconv"
	"log"
)

const(
	redis_str = "%s=%s,"
	redis_strarray = "%s%d=%s,"
)

func getRedisToStr(classField reflect.StructField, classVal reflect.Value) (bool,string) {
	classType := classField.Name
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("getRedisToStr", classType, classVal,  err)
		}
	}()

	sType := getTypeString(classField, classVal)
	//fmt.Println(classVal, classType, sType, classVal.Type().String())
	var strsql *string
	noramlsql := ""
	strsql = &noramlsql

	switch sType {
	case "*float64":
		value := float64(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*float64)
		}
		*strsql += fmt.Sprintf(redis_str, classType, strconv.FormatFloat(value, 'f', -1, 64))
	case "*float32":
		value := float32(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*float32)
		}
		*strsql += fmt.Sprintf(redis_str, classType, strconv.FormatFloat(float64(value), 'f', -1, 32))
	case "*bool":
		value := bool(false)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*bool)
		}
		*strsql += fmt.Sprintf("%s=%t", classType, value)
	case "*int8":
		value := int8(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*int8)
		}
		*strsql += fmt.Sprintf(redis_str, classType, strconv.FormatInt(int64(value),10))
	case "*uint8":
		value := uint8(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*uint8)
		}
		*strsql += fmt.Sprintf(redis_str, classType, strconv.FormatUint(uint64(value),10))
	case "*int16":
		value := int16(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*int16)
		}
		*strsql += fmt.Sprintf(redis_str, classType, strconv.FormatInt(int64(value),10))
	case "*uint16":
		value := uint16(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*uint16)
		}
		*strsql += fmt.Sprintf(redis_str, classType, strconv.FormatUint(uint64(value),10))
	case "*int32":
		value := int32(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*int32)
		}
		*strsql += fmt.Sprintf(redis_str, classType, strconv.FormatInt(int64(value),10))
	case "*uint32":
		value := uint32(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*uint32)
		}
		*strsql += fmt.Sprintf(redis_str, classType, strconv.FormatUint(uint64(value),10))
	case "*int64":
		value := int64(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*int64)
		}
		*strsql += fmt.Sprintf(redis_str, classType, strconv.FormatInt(int64(value),10))
	case "*uint64":
		value := uint64(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*uint64)
		}
		*strsql += fmt.Sprintf(redis_str, classType, strconv.FormatUint(uint64(value),10))
	case "*string":
		value := string("")
		if !classVal.IsNil() {
			value = *classVal.Interface().(*string)
		}
		*strsql += fmt.Sprintf(redis_str, classType, value)
	case "*int":
		value := int(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*int)
		}
		*strsql += fmt.Sprintf(redis_str, classType, strconv.FormatInt(int64(value),10))
	case "*uint":
		value := uint(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*uint)
		}
		*strsql += fmt.Sprintf(redis_str, classType, strconv.FormatUint(uint64(value),10))
	case "*struct":
		value := classVal.Elem().Interface()
		n := parseRedisToStr(value)
		noramlsql += n
	case "float64":
		*strsql += fmt.Sprintf(redis_str, classType, strconv.FormatFloat(classVal.Float(), 'f', -1, 64))
	case "float32":
		*strsql += fmt.Sprintf(redis_str, classType, strconv.FormatFloat(classVal.Float(), 'f', -1, 32))
	case "bool":
		*strsql += fmt.Sprintf("%s=%t", classType, classVal.Bool())
	case "int8":
		*strsql += fmt.Sprintf(redis_str, classType, strconv.FormatInt(classVal.Int(),10))
	case "uint8":
		*strsql += fmt.Sprintf(redis_str, classType, strconv.FormatUint(classVal.Uint(),10))
	case "int16":
		*strsql += fmt.Sprintf(redis_str, classType, strconv.FormatInt(classVal.Int(),10))
	case "uint16":
		*strsql += fmt.Sprintf(redis_str, classType, strconv.FormatUint(classVal.Uint(),10))
	case "int32":
		*strsql += fmt.Sprintf(redis_str, classType, strconv.FormatInt(classVal.Int(),10))
	case "uint32":
		*strsql += fmt.Sprintf(redis_str, classType, strconv.FormatUint(classVal.Uint(), 10))
	case "int64":
		*strsql += fmt.Sprintf(redis_str, classType, strconv.FormatInt(classVal.Int(),10))
	case "uint64":
		*strsql += fmt.Sprintf(redis_str, classType, strconv.FormatUint(classVal.Uint(),10))
	case "string":
		*strsql += fmt.Sprintf(redis_str, classType, classVal.String())
	case "int":
		*strsql += fmt.Sprintf(redis_str, classType, strconv.FormatInt(classVal.Int(),10))
	case "uint":
		*strsql += fmt.Sprintf(redis_str, classType, strconv.FormatUint(classVal.Uint(),10))
	case "struct":
		n := parseRedisToStr(classVal.Interface())
		noramlsql += n
	case "[]float64":
		value := []float64{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]float64)
		}
		for i,v := range value{
			*strsql += fmt.Sprintf(redis_strarray, classType, i, strconv.FormatFloat(v, 'f', -1, 64))
		}
	case "[]float32":
		value := []float32{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]float32)
		}
		for i,v := range value{
			*strsql += fmt.Sprintf(redis_strarray, classType, i, strconv.FormatFloat(float64(v), 'f', -1, 32))
		}
	case "[]bool":
		value := []bool{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]bool)
		}
		for i,v := range value{
			*strsql += fmt.Sprintf("%s%d=%t,", classType, i, v)
		}
	case "[]int8":
		value := []int8{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]int8)
		}
		for i,v := range value{
			*strsql += fmt.Sprintf(redis_strarray, classType, i, strconv.FormatInt(int64(v), 10))
		}
	case "[]uint8":
		value := []uint8{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]uint8)
		}
		for i,v := range value{
			*strsql += fmt.Sprintf(redis_strarray, classType, i, strconv.FormatUint(uint64(v), 10))
		}
	case "[]int16":
		value := []int16{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]int16)
		}
		for i,v := range value{
			*strsql += fmt.Sprintf(redis_strarray, classType, i, strconv.FormatInt(int64(v), 10))
		}
	case "[]uint16":
		value := []uint16{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]uint16)
		}
		for i,v := range value{
			*strsql += fmt.Sprintf(redis_strarray, classType, i, strconv.FormatUint(uint64(v), 10))
		}
	case "[]int32":
		value := []int32{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]int32)
		}
		for i,v := range value{
			*strsql += fmt.Sprintf(redis_strarray, classType, i, strconv.FormatInt(int64(v), 10))
		}
	case "[]uint32":
		value := []uint32{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]uint32)
		}
		for i,v := range value{
			*strsql += fmt.Sprintf(redis_strarray, classType, i, strconv.FormatUint(uint64(v), 10))
		}
	case "[]int64":
		value := []int64{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]int64)
		}
		for i,v := range value{
			*strsql += fmt.Sprintf(redis_strarray, classType, i, strconv.FormatInt(int64(v), 10))
		}
	case "[]uint64":
		value := []uint64{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]uint64)
		}
		for i,v := range value{
			*strsql += fmt.Sprintf(redis_strarray, classType, i, strconv.FormatUint(uint64(v), 10))
		}
	case "[]string":
		value := []string{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]string)
		}
		for i,v := range value{
			*strsql += fmt.Sprintf(redis_strarray, classType, i, v)
		}
	case "[]int":
		value := []int{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]int)
		}
		for i,v := range value{
			*strsql += fmt.Sprintf(redis_strarray, classType, i, strconv.FormatInt(int64(v), 10))
		}
	case "[]uint":
		value := []uint{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]uint)
		}
		for i,v := range value{
			*strsql += fmt.Sprintf(redis_strarray, classType, i, strconv.FormatUint(uint64(v), 10))
		}
	case "[]struct"://no support
		for i := 0;  i < classVal.Len(); i++{
			n := parseRedisToStr(classVal.Index(i).Interface())
			noramlsql += n
		}
	case "[*]float64":
		for i := 0;  i < classVal.Len(); i++{
			*strsql += fmt.Sprintf(redis_strarray, classType, i, strconv.FormatFloat(classVal.Index(i).Float(), 'f', -1, 64))
		}
	case "[*]float32":
		for i := 0;  i < classVal.Len(); i++{
			*strsql += fmt.Sprintf(redis_strarray, classType, i, strconv.FormatFloat(classVal.Index(i).Float(), 'f', -1, 64))
		}
	case "[*]bool":
		for i := 0;  i < classVal.Len(); i++{
			*strsql += fmt.Sprintf("%s%d=%t,", classType, i, classVal.Index(i).Bool())
		}
	case "[*]int8":
		for i := 0;  i < classVal.Len(); i++{
			*strsql += fmt.Sprintf(redis_strarray, classType, i, strconv.FormatInt(classVal.Index(i).Int(), 10))
		}
	case "[*]uint8":
		for i := 0;  i < classVal.Len(); i++{
			*strsql += fmt.Sprintf(redis_strarray, classType, i, strconv.FormatUint(classVal.Index(i).Uint(), 10))
		}
	case "[*]int16":
		for i := 0;  i < classVal.Len(); i++{
			*strsql += fmt.Sprintf(redis_strarray, classType, i, strconv.FormatInt(classVal.Index(i).Int(), 10))
		}
	case "[*]uint16":
		for i := 0;  i < classVal.Len(); i++{
			*strsql += fmt.Sprintf(redis_strarray, classType, i, strconv.FormatUint(classVal.Index(i).Uint(), 10))
		}
	case "[*]int32":
		for i := 0;  i < classVal.Len(); i++{
			*strsql += fmt.Sprintf(redis_strarray, classType, i, strconv.FormatInt(classVal.Index(i).Int(), 10))
		}
	case "[*]uint32":
		for i := 0;  i < classVal.Len(); i++{
			*strsql += fmt.Sprintf(redis_strarray, classType, i, strconv.FormatUint(classVal.Index(i).Uint(), 10))
		}
	case "[*]int64":
		for i := 0;  i < classVal.Len(); i++{
			*strsql += fmt.Sprintf(redis_strarray, classType, i, strconv.FormatInt(classVal.Index(i).Int(), 10))
		}
	case "[*]uint64":
		for i := 0;  i < classVal.Len(); i++{
			*strsql += fmt.Sprintf(redis_strarray, classType, i, strconv.FormatUint(classVal.Index(i).Uint(), 10))
		}
	case "[*]string":
		for i := 0;  i < classVal.Len(); i++{
			*strsql += fmt.Sprintf(redis_strarray, classType, i, classVal.Index(i).String())
		}
	case "[*]int":
		for i := 0;  i < classVal.Len(); i++{
			*strsql += fmt.Sprintf(redis_strarray, classType, i, strconv.FormatInt(classVal.Index(i).Int(), 10))
		}
	case "[*]uint":
		for i := 0;  i < classVal.Len(); i++{
			*strsql += fmt.Sprintf(redis_strarray, classType, i, strconv.FormatUint(classVal.Index(i).Uint(), 10))
		}
	case "[*]struct"://no support
		for i := 0;  i < classVal.Len(); i++{
			n := parseRedisToStr(classVal.Index(i).Interface())
			noramlsql += n
		}
	default:
		fmt.Println("getRedisToStr type not supported", sType,  classField.Type)
		panic("getRedisToStr type not supported")
		return false, ""
	}
	return true, noramlsql
}

func parseRedisToStr(obj interface{}) (string){
	var protoVal reflect.Value
	protoType := reflect.TypeOf(obj)
	if protoType.Kind() == reflect.Ptr {
		protoType = reflect.TypeOf(obj).Elem()
		protoVal = reflect.ValueOf(obj).Elem()
	}else if protoType.Kind() == reflect.Struct{
		protoVal = reflect.ValueOf(obj)
	}else{
		errorStr := fmt.Sprintf("parseRedisStr no support ptr %s", protoType.Name())
		log.Println(errorStr)
		panic(errorStr)
		return ""
	}

	str := ""
	for i := 0; i < protoType.NumField(); i++{
		if !protoVal.Field(i).CanInterface(){
			continue
		}

		bRight, sqlstr := getRedisToStr(protoType.Field(i), protoVal.Field(i))
		if !bRight{
			errorStr := fmt.Sprintf("parseRedisStr type not supported %s", protoType.Name())
			panic(errorStr)
			return ""//丢弃这个包
		}
		str += sqlstr
	}
	return str
}

func redistoStr(str string) string{
	index := strings.LastIndex(str, ",")
	if index!= -1{
		str = str[:index]
	}

	return str
}

//--- struct to str
func RedisStr(obj interface{})string{
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("RedisSql", err)
		}
	}()
	str:= parseRedisToStr(obj)
	return  redistoStr(str)
}
/*
type Sqltest1 struct{
	MM int8
	MM1 uint8
}

type sqltest struct{
	I uint8 `primary`
	J int8 `primary`
	K string
	I2 []uint
	J2 []int
	*Sqltest1
	K2 []string
	T int64 `datetime`
}
	var1 :=sqltest{1, 2, "test", []uint{1, 2}, []int{3,4}, &Sqltest1{1, 1}, []string{"tes21", "tes31"}, time.Now().Unix()}
	fmt.Println(db.RedisStr(var1))
 */