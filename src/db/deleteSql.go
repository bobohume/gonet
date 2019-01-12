package db

import (
	"base"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
)

const(
	delete_sql = "%s='%s',"
	delete_sqlarray = "%s%d='%s',"
)

func getDeleteSql(classField reflect.StructField, classVal reflect.Value) (bool,string,string) {
	classType := getSqlName(classField)
	/*defer func() {
		if err := recover(); err != nil {
			fmt.Println("getDeleteSql", classType,  err)
		}
	}()*/

	sType := base.GetTypeStringEx(classField, classVal)
	//fmt.Println(classVal, classType, sType, classVal.Type().String())
	var strsql *string
	primarysql := ""
	noramlsql := ""
	if isPrimary(classField){
		strsql = &primarysql
	}else{
		strsql = &noramlsql
	}
	switch sType {
	case "*float64":
		value := float64(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*float64)
		}
		*strsql += fmt.Sprintf(delete_sql, classType, strconv.FormatFloat(value, 'f', -1, 64))
	case "*float32":
		value := float32(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*float32)
		}
		*strsql += fmt.Sprintf(delete_sql, classType, strconv.FormatFloat(float64(value), 'f', -1, 32))
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
		*strsql += fmt.Sprintf(delete_sql, classType, strconv.FormatInt(int64(value),10))
	case "*uint8":
		value := uint8(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*uint8)
		}
		*strsql += fmt.Sprintf(delete_sql, classType, strconv.FormatUint(uint64(value),10))
	case "*int16":
		value := int16(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*int16)
		}
		*strsql += fmt.Sprintf(delete_sql, classType, strconv.FormatInt(int64(value),10))
	case "*uint16":
		value := uint16(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*uint16)
		}
		*strsql += fmt.Sprintf(delete_sql, classType, strconv.FormatUint(uint64(value),10))
	case "*int32":
		value := int32(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*int32)
		}
		*strsql += fmt.Sprintf(delete_sql, classType, strconv.FormatInt(int64(value),10))
	case "*uint32":
		value := uint32(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*uint32)
		}
		*strsql += fmt.Sprintf(delete_sql, classType, strconv.FormatUint(uint64(value),10))
	case "*int64":
		value := int64(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*int64)
		}
		if !isDatetime(classField){
			*strsql += fmt.Sprintf(delete_sql, classType, strconv.FormatInt(int64(value),10))
		}else{
			*strsql += fmt.Sprintf(delete_sql, classType, GetDBTimeString(int64(value)))
		}
	case "*uint64":
		value := uint64(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*uint64)
		}
		*strsql += fmt.Sprintf(delete_sql, classType, strconv.FormatUint(uint64(value),10))
	case "*string":
		value := string("")
		if !classVal.IsNil() {
			value = *classVal.Interface().(*string)
		}
		*strsql += fmt.Sprintf(delete_sql, classType, value)
	case "*int":
		value := int(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*int)
		}
		*strsql += fmt.Sprintf(delete_sql, classType, strconv.FormatInt(int64(value),10))
	case "*uint":
		value := uint(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*uint)
		}
		*strsql += fmt.Sprintf(delete_sql, classType, strconv.FormatUint(uint64(value),10))
	case "*struct":
		if !classVal.IsNil() {
			value := classVal.Elem().Interface()
			n, p := parseDeleteSql(value)
			noramlsql += n
			primarysql += p
		}
	case "float64":
		*strsql += fmt.Sprintf(delete_sql, classType, strconv.FormatFloat(classVal.Float(), 'f', -1, 64))
	case "float32":
		*strsql += fmt.Sprintf(delete_sql, classType, strconv.FormatFloat(classVal.Float(), 'f', -1, 32))
	case "bool":
		*strsql += fmt.Sprintf("%s=%t", classType, classVal.Bool())
	case "int8":
		*strsql += fmt.Sprintf(delete_sql, classType, strconv.FormatInt(classVal.Int(),10))
	case "uint8":
		*strsql += fmt.Sprintf(delete_sql, classType, strconv.FormatUint(classVal.Uint(),10))
	case "int16":
		*strsql += fmt.Sprintf(delete_sql, classType, strconv.FormatInt(classVal.Int(),10))
	case "uint16":
		*strsql += fmt.Sprintf(delete_sql, classType, strconv.FormatUint(classVal.Uint(),10))
	case "int32":
		*strsql += fmt.Sprintf(delete_sql, classType, strconv.FormatInt(classVal.Int(),10))
	case "uint32":
		*strsql += fmt.Sprintf(delete_sql, classType, strconv.FormatUint(classVal.Uint(), 10))
	case "int64":
		if !isDatetime(classField){
			*strsql += fmt.Sprintf(delete_sql, classType, strconv.FormatInt(classVal.Int(),10))
		}else{
			*strsql += fmt.Sprintf(delete_sql, classType, GetDBTimeString(classVal.Int()))
		}
	case "uint64":
		*strsql += fmt.Sprintf(delete_sql, classType, strconv.FormatUint(classVal.Uint(),10))
	case "string":
		*strsql += fmt.Sprintf(delete_sql, classType, classVal.String())
	case "int":
		*strsql += fmt.Sprintf(delete_sql, classType, strconv.FormatInt(classVal.Int(),10))
	case "uint":
		*strsql += fmt.Sprintf(delete_sql, classType, strconv.FormatUint(classVal.Uint(),10))
	case "struct":
		n, p := parseDeleteSql(classVal.Interface())
		noramlsql += n
		primarysql += p
	case "[]float64":
		value := []float64{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]float64)
		}
		for i,v := range value{
			*strsql += fmt.Sprintf(delete_sqlarray, classType, i, strconv.FormatFloat(v, 'f', -1, 64))
		}
	case "[]float32":
		value := []float32{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]float32)
		}
		for i,v := range value{
			*strsql += fmt.Sprintf(delete_sqlarray, classType, i, strconv.FormatFloat(float64(v), 'f', -1, 32))
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
			*strsql += fmt.Sprintf(delete_sqlarray, classType, i, strconv.FormatInt(int64(v), 10))
		}
	case "[]uint8":
		value := []uint8{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]uint8)
		}
		for i,v := range value{
			*strsql += fmt.Sprintf(delete_sqlarray, classType, i, strconv.FormatUint(uint64(v), 10))
		}
	case "[]int16":
		value := []int16{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]int16)
		}
		for i,v := range value{
			*strsql += fmt.Sprintf(delete_sqlarray, classType, i, strconv.FormatInt(int64(v), 10))
		}
	case "[]uint16":
		value := []uint16{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]uint16)
		}
		for i,v := range value{
			*strsql += fmt.Sprintf(delete_sqlarray, classType, i, strconv.FormatUint(uint64(v), 10))
		}
	case "[]int32":
		value := []int32{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]int32)
		}
		for i,v := range value{
			*strsql += fmt.Sprintf(delete_sqlarray, classType, i, strconv.FormatInt(int64(v), 10))
		}
	case "[]uint32":
		value := []uint32{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]uint32)
		}
		for i,v := range value{
			*strsql += fmt.Sprintf(delete_sqlarray, classType, i, strconv.FormatUint(uint64(v), 10))
		}
	case "[]int64":
		value := []int64{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]int64)
		}
		for i,v := range value{
			if !isDatetime(classField){
				*strsql += fmt.Sprintf(delete_sqlarray, classType, i, strconv.FormatInt(int64(v), 10))
			}else{
				*strsql += fmt.Sprintf(delete_sqlarray, classType, i, GetDBTimeString(v))
			}
		}
	case "[]uint64":
		value := []uint64{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]uint64)
		}
		for i,v := range value{
			*strsql += fmt.Sprintf(delete_sqlarray, classType, i, strconv.FormatUint(uint64(v), 10))
		}
	case "[]string":
		value := []string{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]string)
		}
		for i,v := range value{
			*strsql += fmt.Sprintf(delete_sqlarray, classType, i, v)
		}
	case "[]int":
		value := []int{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]int)
		}
		for i,v := range value{
			*strsql += fmt.Sprintf(delete_sqlarray, classType, i, strconv.FormatInt(int64(v), 10))
		}
	case "[]uint":
		value := []uint{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]uint)
		}
		for i,v := range value{
			*strsql += fmt.Sprintf(delete_sqlarray, classType, i, strconv.FormatUint(uint64(v), 10))
		}
	case "[]struct":
		for i := 0;  i < classVal.Len(); i++{
			n, p := parseDeleteSql(classVal.Index(i).Interface())
			noramlsql += n
			primarysql += p
		}
	case "[*]float64":
		for i:= 0; i < classVal.Len(); i++{
			*strsql += fmt.Sprintf(delete_sqlarray, classType, i, strconv.FormatFloat(classVal.Index(i).Float(), 'f', -1, 64))
		}
	case "[*]float32":
		for i:= 0; i < classVal.Len(); i++{
			*strsql += fmt.Sprintf(delete_sqlarray, classType, i, strconv.FormatFloat(classVal.Index(i).Float(), 'f', -1, 64))
		}
	case "[*]bool":
		for i:= 0; i < classVal.Len(); i++{
			*strsql += fmt.Sprintf("%s%d=%t,", classType, i, classVal.Index(i).Bool())
		}
	case "[*]int8":
		for i:= 0; i < classVal.Len(); i++{
			*strsql += fmt.Sprintf(delete_sqlarray, classType, i, strconv.FormatInt(classVal.Index(i).Int(), 10))
		}
	case "[*]uint8":
		for i:= 0; i < classVal.Len(); i++{
			*strsql += fmt.Sprintf(delete_sqlarray, classType, i, strconv.FormatUint(classVal.Index(i).Uint(), 10))
		}
	case "[*]int16":
		for i:= 0; i < classVal.Len(); i++{
			*strsql += fmt.Sprintf(delete_sqlarray, classType, i, strconv.FormatInt(classVal.Index(i).Int(), 10))
		}
	case "[*]uint16":
		for i:= 0; i < classVal.Len(); i++{
			*strsql += fmt.Sprintf(delete_sqlarray, classType, i, strconv.FormatUint(classVal.Index(i).Uint(), 10))
		}
	case "[*]int32":
		for i:= 0; i < classVal.Len(); i++{
			*strsql += fmt.Sprintf(delete_sqlarray, classType, i, strconv.FormatInt(classVal.Index(i).Int(), 10))
		}
	case "[*]uint32":
		for i:= 0; i < classVal.Len(); i++{
			*strsql += fmt.Sprintf(delete_sqlarray, classType, i, strconv.FormatUint(classVal.Index(i).Uint(), 10))
		}
	case "[*]int64":
		for i:= 0; i < classVal.Len(); i++{
			*strsql += fmt.Sprintf(delete_sqlarray, classType, i, strconv.FormatInt(classVal.Index(i).Int(), 10))
		}
	case "[*]uint64":
		for i:= 0; i < classVal.Len(); i++{
			*strsql += fmt.Sprintf(delete_sqlarray, classType, i, strconv.FormatUint(classVal.Index(i).Uint(), 10))
		}
	case "[*]string":
		value := []string{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]string)
		}
		for i,v := range value{
			*strsql += fmt.Sprintf(delete_sqlarray, classType, i, v)
		}
	case "[*]int":
		for i:= 0; i < classVal.Len(); i++{
			*strsql += fmt.Sprintf(delete_sqlarray, classType, i, strconv.FormatInt(classVal.Index(i).Int(), 10))
		}
	case "[*]uint":
		for i:= 0; i < classVal.Len(); i++{
			*strsql += fmt.Sprintf(delete_sqlarray, classType, i, strconv.FormatUint(classVal.Index(i).Uint(), 10))
		}
	case "[*]struct":
		for i := 0;  i < classVal.Len(); i++{
			n, p := parseDeleteSql(classVal.Index(i).Interface())
			noramlsql += n
			primarysql += p
		}
	default:
		fmt.Println("getDeleteSql type not supported", sType,  classField.Type)
		panic("getDeleteSql type not supported")
		return false, "", ""
		//}
	}
	return true, noramlsql, primarysql
}

func parseDeleteSql(obj interface{}) (string, string){
	var protoVal reflect.Value
	protoType := reflect.TypeOf(obj)
	if protoType.Kind() == reflect.Ptr {
		protoType = reflect.TypeOf(obj).Elem()
		protoVal = reflect.ValueOf(obj).Elem()
	}else if protoType.Kind() == reflect.Struct{
		protoVal = reflect.ValueOf(obj)
	}else{
		errorStr := fmt.Sprintf("parseDeleteSql no support ptr %s", protoType.Name())
		log.Println(errorStr)
		panic(errorStr)
		return "",""
	}

	str := ""
	primary := ""
	for i := 0; i < protoType.NumField(); i++{
		if !protoVal.Field(i).CanInterface(){
			continue
		}

		bRight, sqlstr, sqlprimary := getDeleteSql(protoType.Field(i), protoVal.Field(i))
		if !bRight{
			errorStr := fmt.Sprintf("parseDeleteSql type not supported %s", protoType.Name())
			panic(errorStr)
			return "",""//丢弃这个包
		}
		str += sqlstr
		primary += sqlprimary
	}
	return str,primary
}

func deleteSqlStr(sqltable string, str string, primary string) string{
	index := strings.LastIndex(str, ",")
	if index!= -1{
		str = str[:index]
	}

	index = strings.LastIndex(primary, ",")
	if index!= -1{
		primary = primary[:index]
	}
	return "delete from " + sqltable + " where " + primary
}

//--- struct to sql
func DeleteSql(obj interface{}, sqltable string)string{
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("DeleteSql", err)
		}
	}()

	str, primary := parseDeleteSql(obj)
	return  deleteSqlStr(sqltable, str, primary)
}

func DeleteSqlEx(obj interface{}, sqltable string, params ...string) string {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("DeleteSqlEx", err)
		}
	}()

	protoVal  := reflect.ValueOf(obj)
	protoType := reflect.TypeOf(obj)
	if protoType.Kind() == reflect.Ptr {
		protoType = reflect.TypeOf(obj).Elem()
		protoVal = reflect.ValueOf(obj).Elem()
	}

	str := ""
	primary := ""
	nameMap := make(map[string] string)
	for _,v := range params{
		v1 := strings.ToLower(v)
		nameMap[v1] = v1
	}
	for i := 0; i < protoType.NumField(); i++ {
		if !protoVal.Field(i).CanInterface() {
			continue
		}

		sf := protoType.Field(i)
		_, exist := nameMap[getSqlName(sf)]
		if exist{
			bRight, name, value := getDeleteSql(sf, protoVal.Field(i))
			if !bRight{
				errorStr := fmt.Sprintf("DeleteSqlEx error %s", reflect.TypeOf(obj).Name())
				panic(errorStr)
				return ""//丢弃这个包
			}
			str += name
			primary += value
		}
	}

	return deleteSqlStr(sqltable, str, primary)
}

