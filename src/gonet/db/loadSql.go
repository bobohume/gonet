package db

import (
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"gonet/base"
	"math"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

func loadsql(sqlData *SqlData, p *Properties, val string){
	//sqlData.SqlValue += fmt.Sprintf("'%s',", val)
	sqlData.SqlName += fmt.Sprintf("`%s`,", p.Name)
}

func loadsqlarray(sqlData *SqlData, p *Properties, val string, i int){
	//sqlData.SqlValue += fmt.Sprintf("'%s',", val)
	if sqlData.bitMap == nil || !sqlData.bitMap.Test(i){
		return
	}
	sqlData.SqlName += fmt.Sprintf("`%s%d`,", p.Name, i)
}

func getLoadSql(classField reflect.StructField, classVal reflect.Value, sqlData *SqlData) (bool) {
	p := getProperties(classField)
	sType := base.GetTypeStringEx(classField, classVal)
	if p.IsJson(){
		data, _ := json.Marshal(classVal.Interface())
		loadsql(sqlData, p, string(data))
		return true
	}else if p.IsBlob(){
		for classVal.Kind() == reflect.Ptr {
			classVal = classVal.Elem()
		}
		data, _ := proto.Marshal(classVal.Addr().Interface().(proto.Message))
		loadsql(sqlData, p, string(data))
		return true
	}else if p.IsIgnore(){
		return true
	}

	switch sType {
	case "*float64":
		value := float64(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*float64)
		}
		loadsql(sqlData, p, strconv.FormatFloat(value, 'f', -1, 64))
	case "*float32":
		value := float32(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*float32)
		}
		loadsql(sqlData, p, strconv.FormatFloat(float64(value), 'f', -1, 32))
	case "*bool":
		value := bool(false)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*bool)
		}
		loadsql(sqlData, p, strconv.FormatBool(value))
	case "*int8":
		value := int8(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*int8)
		}
		loadsql(sqlData, p, strconv.FormatInt(int64(value),10))
	case "*uint8":
		value := uint8(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*uint8)
		}
		loadsql(sqlData, p, strconv.FormatUint(uint64(value),10))
	case "*int16":
		value := int16(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*int16)
		}
		loadsql(sqlData, p, strconv.FormatInt(int64(value),10))
	case "*uint16":
		value := uint16(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*uint16)
		}
		loadsql(sqlData, p, strconv.FormatUint(uint64(value),10))
	case "*int32":
		value := int32(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*int32)
		}
		loadsql(sqlData, p, strconv.FormatInt(int64(value),10))
	case "*uint32":
		value := uint32(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*uint32)
		}
		loadsql(sqlData, p, strconv.FormatUint(uint64(value),10))
	case "*int64":
		value := int64(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*int64)
		}
		if !p.IsDatetime(){
			loadsql(sqlData, p, strconv.FormatInt(int64(value),10))
		}else{
			loadsql(sqlData, p, GetDBTimeString(int64(value)))
		}
	case "*uint64":
		value := uint64(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*uint64)
		}
		loadsql(sqlData, p, strconv.FormatUint(uint64(value),10))
	case "*string":
		value := string("")
		if !classVal.IsNil() {
			value = *classVal.Interface().(*string)
		}
		loadsql(sqlData, p, value)
	case "*int":
		value := int(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*int)
		}
		loadsql(sqlData, p, strconv.FormatInt(int64(value),10))
	case "*uint":
		value := uint(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*uint)
		}
		loadsql(sqlData, p, strconv.FormatUint(uint64(value),10))
	case "*struct":
		if !classVal.IsNil() {
			value := classVal.Elem().Interface()
			parseLoadSql(value, sqlData)
		}
	case "float64":
		loadsql(sqlData, p, strconv.FormatFloat(classVal.Float(), 'f', -1, 64))
	case "float32":
		loadsql(sqlData, p, strconv.FormatFloat(classVal.Float(), 'f', -1, 32))
	case "bool":
		loadsql(sqlData, p, strconv.FormatBool(classVal.Bool()))
	case "int8":
		loadsql(sqlData, p, strconv.FormatInt(classVal.Int(),10))
	case "uint8":
		loadsql(sqlData, p, strconv.FormatUint(classVal.Uint(),10))
	case "int16":
		loadsql(sqlData, p, strconv.FormatInt(classVal.Int(),10))
	case "uint16":
		loadsql(sqlData, p, strconv.FormatUint(classVal.Uint(),10))
	case "int32":
		loadsql(sqlData, p, strconv.FormatInt(classVal.Int(),10))
	case "uint32":
		loadsql(sqlData, p, strconv.FormatUint(classVal.Uint(), 10))
	case "int64":
		if !p.IsDatetime(){
			loadsql(sqlData, p, strconv.FormatInt(classVal.Int(),10))
		}else{
			loadsql(sqlData, p, GetDBTimeString(classVal.Int()))
		}
	case "uint64":
		loadsql(sqlData, p, strconv.FormatUint(classVal.Uint(),10))
	case "string":
		loadsql(sqlData, p, classVal.String())
	case "int":
		loadsql(sqlData, p, strconv.FormatInt(classVal.Int(),10))
	case "uint":
		loadsql(sqlData, p, strconv.FormatUint(classVal.Uint(),10))
	case "struct":
		parseLoadSql(classVal.Interface(), sqlData)
	case "[]float64":
		value := []float64{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]float64)
		}
		for i,v := range value{
			loadsqlarray(sqlData, p, strconv.FormatFloat(v, 'f', -1, 64), i)
		}
	case "[]float32":
		value := []float32{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]float32)
		}
		for i,v := range value{
			loadsqlarray(sqlData, p, strconv.FormatFloat(float64(v), 'f', -1, 32), i)
		}
	case "[]bool":
		value := []bool{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]bool)
		}
		for i,v := range value{
			loadsqlarray(sqlData, p, strconv.FormatBool(v), i)
		}
	case "[]int8":
		value := []int8{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]int8)
		}
		for i,v := range value{
			loadsqlarray(sqlData, p, strconv.FormatInt(int64(v), 10), i)
		}
	case "[]uint8":
		value := []uint8{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]uint8)
		}
		for i,v := range value{
			loadsqlarray(sqlData, p, strconv.FormatUint(uint64(v), 10), i)
		}
	case "[]int16":
		value := []int16{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]int16)
		}
		for i,v := range value{
			loadsqlarray(sqlData, p, strconv.FormatInt(int64(v), 10), i)
		}
	case "[]uint16":
		value := []uint16{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]uint16)
		}
		for i,v := range value{
			loadsqlarray(sqlData, p, strconv.FormatUint(uint64(v), 10), i)
		}
	case "[]int32":
		value := []int32{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]int32)
		}
		for i,v := range value{
			loadsqlarray(sqlData, p, strconv.FormatInt(int64(v), 10), i)
		}
	case "[]uint32":
		value := []uint32{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]uint32)
		}
		for i,v := range value{
			loadsqlarray(sqlData, p, strconv.FormatUint(uint64(v), 10), i)
		}
	case "[]int64":
		value := []int64{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]int64)
		}
		for i,v := range value{
			if !p.IsDatetime(){
				loadsqlarray(sqlData, p, strconv.FormatInt(int64(v), 10), i)
			}else{
				loadsqlarray(sqlData, p, GetDBTimeString(v), i)
			}
		}
	case "[]uint64":
		value := []uint64{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]uint64)
		}
		for i,v := range value{
			loadsqlarray(sqlData, p, strconv.FormatUint(uint64(v), 10), i)
		}
	case "[]string":
		value := []string{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]string)
		}
		for i,v := range value{
			loadsqlarray(sqlData, p, v, i)
		}
	case "[]int":
		value := []int{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]int)
		}
		for i,v := range value{
			loadsqlarray(sqlData, p, strconv.FormatInt(int64(v), 10), i)
		}
	case "[]uint":
		value := []uint{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]uint)
		}
		for i,v := range value{
			loadsqlarray(sqlData, p, strconv.FormatUint(uint64(v), 10), i)
		}
	case "[]struct":
		for i := 0;  i < classVal.Len(); i++{
			parseLoadSql(classVal.Index(i).Interface(), sqlData)
		}
	case "[*]float64":
		for i := 0; i < classVal.Len(); i++{
			loadsqlarray(sqlData, p, strconv.FormatFloat(classVal.Index(i).Float(), 'f', -1, 64), i)
		}
	case "[*]float32":
		for i := 0; i < classVal.Len(); i++{
			loadsqlarray(sqlData, p, strconv.FormatFloat(classVal.Index(i).Float(), 'f', -1, 64), i)
		}
	case "[*]bool":
		for i := 0; i < classVal.Len(); i++{
			loadsqlarray(sqlData, p, strconv.FormatBool(classVal.Index(i).Bool()), i)
		}
	case "[*]int8":
		for i := 0; i < classVal.Len(); i++{
			loadsqlarray(sqlData, p, strconv.FormatInt(classVal.Index(i).Int(), 10), i)
		}
	case "[*]uint8":
		for i := 0; i < classVal.Len(); i++{
			loadsqlarray(sqlData, p, strconv.FormatUint(classVal.Index(i).Uint(), 10), i)
		}
	case "[*]int16":
		for i := 0; i < classVal.Len(); i++{
			loadsqlarray(sqlData, p, strconv.FormatInt(classVal.Index(i).Int(), 10), i)
		}
	case "[*]uint16":
		for i := 0; i < classVal.Len(); i++{
			loadsqlarray(sqlData, p, strconv.FormatUint(classVal.Index(i).Uint(), 10), i)
		}
	case "[*]int32":
		for i := 0; i < classVal.Len(); i++{
			loadsqlarray(sqlData, p, strconv.FormatInt(classVal.Index(i).Int(), 10), i)
		}
	case "[*]uint32":
		for i := 0; i < classVal.Len(); i++{
			loadsqlarray(sqlData, p, strconv.FormatUint(classVal.Index(i).Uint(), 10), i)
		}
	case "[*]int64":
		for i := 0; i < classVal.Len(); i++{
			loadsqlarray(sqlData, p, strconv.FormatInt(classVal.Index(i).Int(), 10), i)
		}
	case "[*]uint64":
		for i := 0; i < classVal.Len(); i++{
			loadsqlarray(sqlData, p, strconv.FormatUint(classVal.Index(i).Uint(), 10), i)
		}
	case "[*]string":
		for i := 0; i < classVal.Len(); i++{
			loadsqlarray(sqlData, p, classVal.Index(i).String(), i)
		}
	case "[*]int":
		for i := 0; i < classVal.Len(); i++{
			loadsqlarray(sqlData, p, strconv.FormatInt(classVal.Index(i).Int(), 10), i)
		}
	case "[*]uint":
		for i := 0; i < classVal.Len(); i++{
			loadsqlarray(sqlData, p, strconv.FormatUint(classVal.Index(i).Uint(), 10), i)
		}
	case "[*]struct":
		for i := 0;  i < classVal.Len(); i++{
			parseLoadSql(classVal.Index(i).Interface(), sqlData)
		}
	default:
		fmt.Println("getLoadSql type not supported", sType,  classField.Type)
		panic("getLoadSql type not supported")
		return false
		//}
	}
	return true
}

func parseLoadSql(obj interface{}, sqlData *SqlData) (){
	classVal := reflect.ValueOf(obj)
	for classVal.Kind() == reflect.Ptr {
		classVal = classVal.Elem()
	}
	classType := classVal.Type()

	for i := 0; i < classType.NumField(); i++{
		if !classVal.Field(i).CanInterface(){
			continue
		}

		bRight := getLoadSql(classType.Field(i), classVal.Field(i), sqlData)
		if !bRight{
			errorStr := fmt.Sprintf("parseLoadSql type not supported %s", classType.Name())
			panic(errorStr)
			return//丢弃这个包
		}
	}
	return
}

func loadSqlStr(sqltable string, sqlData *SqlData) string{
	sqlname := sqlData.SqlName
	sqlvalue := sqlData.SqlValue
	index := strings.LastIndex(sqlname, ",")
	if index!= -1{
		sqlname = sqlname[:index]
	}

	index = strings.LastIndex(sqlvalue, ",")
	if index!= -1{
		sqlvalue = sqlvalue[:index]
	}
	return "select " + sqlname + " from " + sqltable
}

//--- struct to sql
func LoadSql(obj interface{}, sqltable string, key string)string{
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("LoadSql", err)
		}
	}()

	if key != ""{
		key = "where " + key;
	}

	sqlData := &SqlData{}
	parseLoadSql(obj, sqlData)
	return  loadSqlStr(sqltable, sqlData) + " " +  key
}

func LoadSqlEx(obj interface{}, sqltable string, key string, params ...string) string {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("LoadSqlEx", err)
		}
	}()

	classVal := reflect.ValueOf(obj)
	for classVal.Kind() == reflect.Ptr {
		classVal = classVal.Elem()
	}
	classType := classVal.Type()

	if key != ""{
		key = "where " + key;
	}

	sqlData := &SqlData{}
	nameMap := make(map[string] *base.BitMap)//name index[for array]
	for _,v := range params{
		nIndex, i := 0, 0
		v1 := strings.ToLower(v)
		v2 := strings.TrimRightFunc(v, func(r rune) bool {
			if unicode.IsNumber(r){
				nIndex = int(r - '0') * int(math.Pow(10, float64(i))) + nIndex
				i++
				return true
			}
			return false
		})
		if v1 != v2{
			bitMap, bOk := nameMap[v2]
			if !bOk{
				bitMap = base.NewBitMap(MAX_ARRAY_LENGTH)
				nameMap[v2] = bitMap
			}
			bitMap.Set(nIndex)
		}else{
			nameMap[v1] = nil
		}
	}
	for i := 0; i < classType.NumField(); i++ {
		if !classVal.Field(i).CanInterface() {//private成员不能读取
			continue
		}

		sf := classType.Field(i)
		p := getProperties(sf)
		bitMap, exist := nameMap[p.Name]
		if exist{
			sqlData.bitMap = bitMap
			bRight := getLoadSql(sf, classVal.Field(i), sqlData)
			if !bRight{
				errorStr := fmt.Sprintf("LoadSqlEx error %s", reflect.TypeOf(obj).Name())
				panic(errorStr)
				return ""//丢弃这个包
			}
		}
	}
	return loadSqlStr(sqltable, sqlData) + " " + key
}



