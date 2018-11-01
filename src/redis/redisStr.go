package redis

import (
	"base"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
)

const(
	redis_str = "%s"
	redis_strarray = "%s%d"
)

func getRedisName(sf reflect.StructField) string{
	tagMap := base.ParseTag(sf, "redis")
	if name, exist := tagMap["name"];exist{
		return name
	}

	return strings.ToLower(sf.Name)
}

func getRedisStr(classField reflect.StructField, classVal reflect.Value) (bool, []interface{}) {
	classType := getRedisName(classField)
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("getRedisSql", classType, classVal,  err)
		}
	}()

	sType := base.GetTypeStringEx(classField, classVal)
	//fmt.Println(classVal, classType, sType, classVal.Type().String())
	strarr := []interface{}{}
	appendStr := func(str string){
		strarr = append(strarr, str)
	}

	appendArray := func(arr []interface{}){
		strarr = append(strarr, arr...)
	}

	switch sType {
	case "*float64":
		value := float64(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*float64)
		}
		appendStr(classType)
		appendStr(strconv.FormatFloat(value, 'f', -1, 64))
	case "*float32":
		value := float32(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*float32)
		}
		appendStr(classType)
		appendStr(strconv.FormatFloat(float64(value), 'f', -1, 32))
	case "*bool":
		value := bool(false)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*bool)
		}
		appendStr(classType)
		appendStr( fmt.Sprintf("%t", value))
	case "*int8":
		value := int8(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*int8)
		}
		appendStr(classType)
		appendStr(strconv.FormatInt(int64(value),10))
	case "*uint8":
		value := uint8(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*uint8)
		}
		appendStr(classType)
		appendStr(strconv.FormatUint(uint64(value),10))
	case "*int16":
		value := int16(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*int16)
		}
		appendStr(classType)
		appendStr( strconv.FormatInt(int64(value),10))
	case "*uint16":
		value := uint16(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*uint16)
		}
		appendStr(classType)
		appendStr(strconv.FormatUint(uint64(value),10))
	case "*int32":
		value := int32(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*int32)
		}
		appendStr(classType)
		appendStr(strconv.FormatInt(int64(value),10))
	case "*uint32":
		value := uint32(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*uint32)
		}
		appendStr(classType)
		appendStr(strconv.FormatUint(uint64(value),10))
	case "*int64":
		value := int64(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*int64)
		}
		appendStr(classType)
		appendStr(strconv.FormatInt(int64(value),10))
	case "*uint64":
		value := uint64(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*uint64)
		}
		appendStr(classType)
		appendStr(strconv.FormatUint(uint64(value),10))
	case "*string":
		value := string("")
		if !classVal.IsNil() {
			value = *classVal.Interface().(*string)
		}
		appendStr(classType)
		appendStr(value)
	case "*int":
		value := int(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*int)
		}
		appendStr(classType)
		appendStr(strconv.FormatInt(int64(value),10))
	case "*uint":
		value := uint(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*uint)
		}
		appendStr(classType)
		appendStr(strconv.FormatUint(uint64(value),10))
	case "*struct":
		value := classVal.Elem().Interface()
		appendArray(parseRedisSql(value))
	case "float64":
		appendStr(classType)
		appendStr(strconv.FormatFloat(classVal.Float(), 'f', -1, 64))
	case "float32":
		appendStr(classType)
		appendStr(strconv.FormatFloat(classVal.Float(), 'f', -1, 32))
	case "bool":
		appendStr(classType)
		appendStr(fmt.Sprintf("%t", classVal.Bool()))
	case "int8":
		appendStr(classType)
		appendStr(strconv.FormatInt(classVal.Int(),10))
	case "uint8":
		appendStr(classType)
		appendStr(strconv.FormatUint(classVal.Uint(),10))
	case "int16":
		appendStr(classType)
		appendStr(strconv.FormatInt(classVal.Int(),10))
	case "uint16":
		appendStr(classType)
		appendStr(strconv.FormatUint(classVal.Uint(),10))
	case "int32":
		appendStr(classType)
		appendStr(strconv.FormatInt(classVal.Int(),10))
	case "uint32":
		appendStr(classType)
		appendStr(strconv.FormatUint(classVal.Uint(), 10))
	case "int64":
		appendStr(classType)
		appendStr(strconv.FormatInt(classVal.Int(),10))
	case "uint64":
		appendStr(classType)
		appendStr(strconv.FormatUint(classVal.Uint(),10))
	case "string":
		appendStr(classType)
		appendStr(classVal.String())
	case "int":
		appendStr(classType)
		appendStr(strconv.FormatInt(classVal.Int(),10))
	case "uint":
		appendStr(classType)
		appendStr(strconv.FormatUint(classVal.Uint(),10))
	case "struct":
		appendArray(parseRedisSql(classVal.Interface()))
	case "[]float64":
		value := []float64{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]float64)
		}
		for i,v := range value{
			appendStr(fmt.Sprintf(redis_strarray, classType, i))
			appendStr(strconv.FormatFloat(v, 'f', -1, 64))
		}
	case "[]float32":
		value := []float32{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]float32)
		}
		for i,v := range value{
			appendStr(fmt.Sprintf(redis_strarray, classType, i))
			appendStr(strconv.FormatFloat(float64(v), 'f', -1, 32))
		}
	case "[]bool":
		value := []bool{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]bool)
		}
		for i,v := range value{
			appendStr(fmt.Sprintf(redis_strarray, classType, i))
			appendStr(fmt.Sprintf("%t", v))
		}
	case "[]int8":
		value := []int8{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]int8)
		}
		for i,v := range value{
			appendStr(fmt.Sprintf(redis_strarray, classType, i))
			appendStr(strconv.FormatInt(int64(v), 10))
		}
	case "[]uint8":
		value := []uint8{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]uint8)
		}
		for i,v := range value{
			appendStr(fmt.Sprintf(redis_strarray, classType, i))
			appendStr(strconv.FormatUint(uint64(v), 10))
		}
	case "[]int16":
		value := []int16{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]int16)
		}
		for i,v := range value{
			appendStr(fmt.Sprintf(redis_strarray, classType, i))
			appendStr(strconv.FormatInt(int64(v), 10))
		}
	case "[]uint16":
		value := []uint16{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]uint16)
		}
		for i,v := range value{
			appendStr(fmt.Sprintf(redis_strarray, classType, i))
			appendStr(strconv.FormatUint(uint64(v), 10))
		}
	case "[]int32":
		value := []int32{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]int32)
		}
		for i,v := range value{
			appendStr(fmt.Sprintf(redis_strarray, classType, i))
			appendStr(strconv.FormatInt(int64(v), 10))
		}
	case "[]uint32":
		value := []uint32{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]uint32)
		}
		for i,v := range value{
			appendStr(fmt.Sprintf(redis_strarray, classType, i))
			appendStr(strconv.FormatUint(uint64(v), 10))
		}
	case "[]int64":
		value := []int64{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]int64)
		}
		for i,v := range value{
			appendStr(fmt.Sprintf(redis_strarray, classType, i))
			appendStr(strconv.FormatInt(int64(v), 10))
		}
	case "[]uint64":
		value := []uint64{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]uint64)
		}
		for i,v := range value{
			appendStr(fmt.Sprintf(redis_strarray, classType, i))
			appendStr(strconv.FormatUint(uint64(v), 10))
		}
	case "[]string":
		value := []string{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]string)
		}
		for i,v := range value{
			appendStr(fmt.Sprintf(redis_strarray, classType, i))
			appendStr(v)
		}
	case "[]int":
		value := []int{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]int)
		}
		for i,v := range value{
			appendStr(fmt.Sprintf(redis_strarray, classType, i))
			appendStr(strconv.FormatInt(int64(v), 10))
		}
	case "[]uint":
		value := []uint{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]uint)
		}
		for i,v := range value{
			appendStr(fmt.Sprintf(redis_strarray, classType, i))
			appendStr(strconv.FormatUint(uint64(v), 10))
		}
	case "[]struct"://no support
		for i := 0;  i < classVal.Len(); i++{
			appendArray(parseRedisSql(classVal.Index(i).Interface()))
		}
	case "[*]float64":
		for i := 0;  i < classVal.Len(); i++{
			appendStr(fmt.Sprintf(redis_strarray, classType, i))
			appendStr(strconv.FormatFloat(classVal.Index(i).Float(), 'f', -1, 64))
		}
	case "[*]float32":
		for i := 0;  i < classVal.Len(); i++{
			appendStr(fmt.Sprintf(redis_strarray, classType, i))
			appendStr(strconv.FormatFloat(classVal.Index(i).Float(), 'f', -1, 64))
		}
	case "[*]bool":
		for i := 0;  i < classVal.Len(); i++{
			appendStr(fmt.Sprintf(redis_strarray, classType, i))
			appendStr(fmt.Sprintf("%t", classVal.Index(i).Bool()))
		}
	case "[*]int8":
		for i := 0;  i < classVal.Len(); i++{
			appendStr(fmt.Sprintf(redis_strarray, classType, i))
			appendStr(strconv.FormatInt(classVal.Index(i).Int(), 10))
		}
	case "[*]uint8":
		for i := 0;  i < classVal.Len(); i++{
			appendStr(fmt.Sprintf(redis_strarray, classType, i))
			appendStr(strconv.FormatUint(classVal.Index(i).Uint(), 10))
		}
	case "[*]int16":
		for i := 0;  i < classVal.Len(); i++{
			appendStr(fmt.Sprintf(redis_strarray, classType, i))
			appendStr(strconv.FormatInt(classVal.Index(i).Int(), 10))
		}
	case "[*]uint16":
		for i := 0;  i < classVal.Len(); i++{
			appendStr(fmt.Sprintf(redis_strarray, classType, i))
			appendStr(strconv.FormatUint(classVal.Index(i).Uint(), 10))
		}
	case "[*]int32":
		for i := 0;  i < classVal.Len(); i++{
			appendStr(fmt.Sprintf(redis_strarray, classType, i))
			appendStr(strconv.FormatInt(classVal.Index(i).Int(), 10))
		}
	case "[*]uint32":
		for i := 0;  i < classVal.Len(); i++{
			appendStr(fmt.Sprintf(redis_strarray, classType, i))
			appendStr(strconv.FormatUint(classVal.Index(i).Uint(), 10))
		}
	case "[*]int64":
		for i := 0;  i < classVal.Len(); i++{
			appendStr(fmt.Sprintf(redis_strarray, classType, i))
			appendStr(strconv.FormatInt(classVal.Index(i).Int(), 10))
		}
	case "[*]uint64":
		for i := 0;  i < classVal.Len(); i++{
			appendStr(fmt.Sprintf(redis_strarray, classType, i))
			appendStr(strconv.FormatUint(classVal.Index(i).Uint(), 10))
		}
	case "[*]string":
		for i := 0;  i < classVal.Len(); i++{
			appendStr(fmt.Sprintf(redis_strarray, classType, i))
			appendStr( classVal.Index(i).String())
		}
	case "[*]int":
		for i := 0;  i < classVal.Len(); i++{
			appendStr(fmt.Sprintf(redis_strarray, classType, i))
			appendStr(strconv.FormatInt(classVal.Index(i).Int(), 10))
		}
	case "[*]uint":
		for i := 0;  i < classVal.Len(); i++{
			appendStr(fmt.Sprintf(redis_strarray, classType, i))
			appendStr(strconv.FormatUint(classVal.Index(i).Uint(), 10))
		}
	case "[*]struct"://no support
		for i := 0;  i < classVal.Len(); i++{
			appendArray(parseRedisSql(classVal.Index(i).Interface()))
		}
	default:
		/*if classVal.Kind() == reflect.Struct {
			n, p := parseUpdateSql(classVal.Interface())
			noramlsql += n
			primarysql += p
		}else if classVal.Kind() == reflect.Ptr && classVal.Elem().Kind() == reflect.Struct {
			n, p := parseUpdateSql(classVal.Elem().Interface())
			noramlsql += n
			primarysql += p
		} else{*/
		fmt.Println("getRedisStr type not supported", sType,  classField.Type)
		panic("getRedisStr type not supported")
		return false, strarr
		//}
	}
	return true, strarr
}

func parseRedisSql(obj interface{}) ([]interface{}){
	strarr := []interface{}{}
	var protoVal reflect.Value
	protoType := reflect.TypeOf(obj)
	if protoType.Kind() == reflect.Ptr {
		protoType = reflect.TypeOf(obj).Elem()
		protoVal = reflect.ValueOf(obj).Elem()
	}else if protoType.Kind() == reflect.Struct{
		protoVal = reflect.ValueOf(obj)
	}else{
		errorStr := fmt.Sprintf("parseRedisSql no support ptr %s", protoType.Name())
		log.Println(errorStr)
		panic(errorStr)
		return strarr
	}

	for i := 0; i < protoType.NumField(); i++{
		if !protoVal.Field(i).CanInterface(){
			continue
		}

		bRight, strarr1:= getRedisStr(protoType.Field(i), protoVal.Field(i))
		if !bRight{
			errorStr := fmt.Sprintf("parseRedisSql type not supported %s", protoType.Name())
			panic(errorStr)
			return strarr//丢弃这个包
		}
		strarr = append(strarr, strarr1...)
	}
	return strarr
}

//--- struct to sql
func RedisStr(obj interface{})[]interface{}{
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("RedisStr", err)
		}
	}()
	return  parseRedisSql(obj)
}

func RedisStrEx(obj interface{}, params ...string) []interface{} {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("RedisStrEx", err)
		}
	}()

	strarr := []interface{}{}
	protoVal  := reflect.ValueOf(obj)
	protoType := reflect.TypeOf(obj)
	if protoType.Kind() == reflect.Ptr {
		protoType = reflect.TypeOf(obj).Elem()
		protoVal = reflect.ValueOf(obj).Elem()
	}

	for _, v := range params{
		v1,_ := protoType.FieldByName(v)
		if !protoVal.FieldByName(v).CanInterface(){//private成员不能读取
			continue
		}

		bRight, strarr1 := getRedisStr(v1, protoVal.FieldByName(v))
		if !bRight{
			errorStr := fmt.Sprintf("RedisStrEx error %s", reflect.TypeOf(obj).Name())
			panic(errorStr)
			return strarr//丢弃这个包
		}
		strarr = append(strarr, strarr1...)
	}
	return strarr
}