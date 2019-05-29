package db

import (
	"fmt"
	"gonet/base"
	"reflect"
	"strconv"
	"strings"
)

func deletesql(sqlData *SqlData, p *Properties, val string){
	if p.IsPrimary(){
		sqlData.SqlName += fmt.Sprintf("`%s`='%s',", p.Name, val)
	}else{
		//sqlData.SqlValue += fmt.Sprintf("`%s`='%s',", p.Name, val)
	}
}

func deletesqlarray(sqlData *SqlData, p *Properties, val string, i int){
	if p.IsPrimary() {
		sqlData.SqlName += fmt.Sprintf("`%s%d`='%s',", p.Name, i, val)
	}else{
		//sqlData.SqlValue += fmt.Sprintf("`%s%d`='%s',", p.Name, i, val)
	}
}

func getDeleteSql(classField reflect.StructField, classVal reflect.Value, sqlData *SqlData) (bool) {
	p := getProperties(classField)

	sType := base.GetTypeStringEx(classField, classVal)
	//过略json
	if p.IsJson(){
		return true
	} else if p.IsBlob(){
		return true
	} else if p.IsIgnore(){
		return true
	}

	switch sType {
	case "*float64":
		value := float64(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*float64)
		}
		deletesql(sqlData, p, strconv.FormatFloat(value, 'f', -1, 64))
	case "*float32":
		value := float32(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*float32)
		}
		deletesql(sqlData, p, strconv.FormatFloat(float64(value), 'f', -1, 32))
	case "*bool":
		value := bool(false)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*bool)
		}
		deletesql(sqlData, p, strconv.FormatBool(value))
	case "*int8":
		value := int8(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*int8)
		}
		deletesql(sqlData, p, strconv.FormatInt(int64(value),10))
	case "*uint8":
		value := uint8(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*uint8)
		}
		deletesql(sqlData, p, strconv.FormatUint(uint64(value),10))
	case "*int16":
		value := int16(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*int16)
		}
		deletesql(sqlData, p, strconv.FormatInt(int64(value),10))
	case "*uint16":
		value := uint16(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*uint16)
		}
		deletesql(sqlData, p, strconv.FormatUint(uint64(value),10))
	case "*int32":
		value := int32(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*int32)
		}
		deletesql(sqlData, p, strconv.FormatInt(int64(value),10))
	case "*uint32":
		value := uint32(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*uint32)
		}
		deletesql(sqlData, p, strconv.FormatUint(uint64(value),10))
	case "*int64":
		value := int64(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*int64)
		}
		if !p.IsDatetime(){
			deletesql(sqlData, p, strconv.FormatInt(int64(value),10))
		}else{
			deletesql(sqlData, p, GetDBTimeString(int64(value)))
		}
	case "*uint64":
		value := uint64(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*uint64)
		}
		deletesql(sqlData, p, strconv.FormatUint(uint64(value),10))
	case "*string":
		value := string("")
		if !classVal.IsNil() {
			value = *classVal.Interface().(*string)
		}
		deletesql(sqlData, p, value)
	case "*int":
		value := int(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*int)
		}
		deletesql(sqlData, p, strconv.FormatInt(int64(value),10))
	case "*uint":
		value := uint(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*uint)
		}
		deletesql(sqlData, p, strconv.FormatUint(uint64(value),10))
	case "*struct":
		if !classVal.IsNil() {
			value := classVal.Elem().Interface()
			parseDeleteSql(value, sqlData)
		}
	case "float64":
		deletesql(sqlData, p, strconv.FormatFloat(classVal.Float(), 'f', -1, 64))
	case "float32":
		deletesql(sqlData, p, strconv.FormatFloat(classVal.Float(), 'f', -1, 32))
	case "bool":
		deletesql(sqlData, p, strconv.FormatBool(classVal.Bool()))
	case "int8":
		deletesql(sqlData, p, strconv.FormatInt(classVal.Int(),10))
	case "uint8":
		deletesql(sqlData, p, strconv.FormatUint(classVal.Uint(),10))
	case "int16":
		deletesql(sqlData, p, strconv.FormatInt(classVal.Int(),10))
	case "uint16":
		deletesql(sqlData, p, strconv.FormatUint(classVal.Uint(),10))
	case "int32":
		deletesql(sqlData, p, strconv.FormatInt(classVal.Int(),10))
	case "uint32":
		deletesql(sqlData, p, strconv.FormatUint(classVal.Uint(), 10))
	case "int64":
		if !p.IsDatetime(){
			deletesql(sqlData, p, strconv.FormatInt(classVal.Int(),10))
		}else{
			deletesql(sqlData, p, GetDBTimeString(classVal.Int()))
		}
	case "uint64":
		deletesql(sqlData, p, strconv.FormatUint(classVal.Uint(),10))
	case "string":
		deletesql(sqlData, p, classVal.String())
	case "int":
		deletesql(sqlData, p, strconv.FormatInt(classVal.Int(),10))
	case "uint":
		deletesql(sqlData, p, strconv.FormatUint(classVal.Uint(),10))
	case "struct":
		parseDeleteSql(classVal.Interface(), sqlData)
	case "[]float64":
		value := []float64{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]float64)
		}
		for i,v := range value{
			deletesqlarray(sqlData, p, strconv.FormatFloat(v, 'f', -1, 64), i)
		}
	case "[]float32":
		value := []float32{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]float32)
		}
		for i,v := range value{
			deletesqlarray(sqlData, p, strconv.FormatFloat(float64(v), 'f', -1, 32), i)
		}
	case "[]bool":
		value := []bool{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]bool)
		}
		for i,v := range value{
			deletesqlarray(sqlData, p, strconv.FormatBool(v), i)
		}
	case "[]int8":
		value := []int8{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]int8)
		}
		for i,v := range value{
			deletesqlarray(sqlData, p, strconv.FormatInt(int64(v), 10), i)
		}
	case "[]uint8":
		value := []uint8{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]uint8)
		}
		for i,v := range value{
			deletesqlarray(sqlData, p, strconv.FormatUint(uint64(v), 10), i)
		}
	case "[]int16":
		value := []int16{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]int16)
		}
		for i,v := range value{
			deletesqlarray(sqlData, p, strconv.FormatInt(int64(v), 10), i)
		}
	case "[]uint16":
		value := []uint16{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]uint16)
		}
		for i,v := range value{
			deletesqlarray(sqlData, p, strconv.FormatUint(uint64(v), 10), i)
		}
	case "[]int32":
		value := []int32{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]int32)
		}
		for i,v := range value{
			deletesqlarray(sqlData, p, strconv.FormatInt(int64(v), 10), i)
		}
	case "[]uint32":
		value := []uint32{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]uint32)
		}
		for i,v := range value{
			deletesqlarray(sqlData, p, strconv.FormatUint(uint64(v), 10), i)
		}
	case "[]int64":
		value := []int64{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]int64)
		}
		for i,v := range value{
			if !p.IsDatetime(){
				deletesqlarray(sqlData, p, strconv.FormatInt(int64(v), 10), i)
			}else{
				deletesqlarray(sqlData, p, GetDBTimeString(v), i)
			}
		}
	case "[]uint64":
		value := []uint64{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]uint64)
		}
		for i,v := range value{
			deletesqlarray(sqlData, p, strconv.FormatUint(uint64(v), 10), i)
		}
	case "[]string":
		value := []string{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]string)
		}
		for i,v := range value{
			deletesqlarray(sqlData, p, v, i)
		}
	case "[]int":
		value := []int{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]int)
		}
		for i,v := range value{
			deletesqlarray(sqlData, p, strconv.FormatInt(int64(v), 10), i)
		}
	case "[]uint":
		value := []uint{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]uint)
		}
		for i,v := range value{
			deletesqlarray(sqlData, p, strconv.FormatUint(uint64(v), 10), i)
		}
	case "[]struct":
		for i := 0;  i < classVal.Len(); i++{
			parseDeleteSql(classVal.Index(i).Interface(), sqlData)
		}
	case "[*]float64":
		for i:= 0; i < classVal.Len(); i++{
			deletesqlarray(sqlData, p, strconv.FormatFloat(classVal.Index(i).Float(), 'f', -1, 64), i)
		}
	case "[*]float32":
		for i:= 0; i < classVal.Len(); i++{
			deletesqlarray(sqlData, p, strconv.FormatFloat(classVal.Index(i).Float(), 'f', -1, 64), i)
		}
	case "[*]bool":
		for i:= 0; i < classVal.Len(); i++{
			deletesqlarray(sqlData, p, strconv.FormatBool(classVal.Index(i).Bool()), i)
		}
	case "[*]int8":
		for i:= 0; i < classVal.Len(); i++{
			deletesqlarray(sqlData, p, strconv.FormatInt(classVal.Index(i).Int(), 10), i)
		}
	case "[*]uint8":
		for i:= 0; i < classVal.Len(); i++{
			deletesqlarray(sqlData, p, strconv.FormatUint(classVal.Index(i).Uint(), 10), i)
		}
	case "[*]int16":
		for i:= 0; i < classVal.Len(); i++{
			deletesqlarray(sqlData, p, strconv.FormatInt(classVal.Index(i).Int(), 10), i)
		}
	case "[*]uint16":
		for i:= 0; i < classVal.Len(); i++{
			deletesqlarray(sqlData, p, strconv.FormatUint(classVal.Index(i).Uint(), 10), i)
		}
	case "[*]int32":
		for i:= 0; i < classVal.Len(); i++{
			deletesqlarray(sqlData, p, strconv.FormatInt(classVal.Index(i).Int(), 10), i)
		}
	case "[*]uint32":
		for i:= 0; i < classVal.Len(); i++{
			deletesqlarray(sqlData, p, strconv.FormatUint(classVal.Index(i).Uint(), 10), i)
		}
	case "[*]int64":
		for i:= 0; i < classVal.Len(); i++{
			deletesqlarray(sqlData, p, strconv.FormatInt(classVal.Index(i).Int(), 10), i)
		}
	case "[*]uint64":
		for i:= 0; i < classVal.Len(); i++{
			deletesqlarray(sqlData, p, strconv.FormatUint(classVal.Index(i).Uint(), 10), i)
		}
	case "[*]string":
		value := []string{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]string)
		}
		for i,v := range value{
			deletesqlarray(sqlData, p, v, i)
		}
	case "[*]int":
		for i:= 0; i < classVal.Len(); i++{
			deletesqlarray(sqlData, p, strconv.FormatInt(classVal.Index(i).Int(), 10), i)
		}
	case "[*]uint":
		for i:= 0; i < classVal.Len(); i++{
			deletesqlarray(sqlData, p, strconv.FormatUint(classVal.Index(i).Uint(), 10), i)
		}
	case "[*]struct":
		for i := 0;  i < classVal.Len(); i++{
			parseDeleteSql(classVal.Index(i).Interface(), sqlData)
		}
	default:
		fmt.Println("getDeleteSql type not supported", sType,  classField.Type)
		panic("getDeleteSql type not supported")
		return false
		//}
	}
	return true
}

func parseDeleteSql(obj interface{}, sqlData *SqlData){
	classVal := reflect.ValueOf(obj)
	for classVal.Kind() == reflect.Ptr {
		classVal = classVal.Elem()
	}
	classType := classVal.Type()

	for i := 0; i < classType.NumField(); i++{
		if !classVal.Field(i).CanInterface(){
			continue
		}

		bRight := getDeleteSql(classType.Field(i), classVal.Field(i), sqlData)
		if !bRight{
			errorStr := fmt.Sprintf("parseDeleteSql type not supported %s", classType.Name())
			panic(errorStr)
			return//丢弃这个包
		}
	}
	return
}

func deleteSqlStr(sqltable string, sqlData *SqlData) string{
	str := sqlData.SqlValue
	primary := sqlData.SqlName
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

	sqlData := &SqlData{}
	parseDeleteSql(obj, sqlData)
	return  deleteSqlStr(sqltable, sqlData)
}

func DeleteSqlEx(obj interface{}, sqltable string, params ...string) string {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("DeleteSqlEx", err)
		}
	}()

	classVal := reflect.ValueOf(obj)
	for classVal.Kind() == reflect.Ptr {
		classVal = classVal.Elem()
	}
	classType := classVal.Type()

	sqlData := &SqlData{}
	nameMap := make(map[string] string)
	for _,v := range params{
		v1 := strings.ToLower(v)
		nameMap[v1] = v1
	}
	for i := 0; i < classType.NumField(); i++ {
		if !classVal.Field(i).CanInterface() {
			continue
		}

		sf := classType.Field(i)
		_, exist := nameMap[getProperties(sf).Name]
		if exist{
			bRight := getDeleteSql(sf, classVal.Field(i), sqlData)
			if !bRight{
				errorStr := fmt.Sprintf("DeleteSqlEx error %s", reflect.TypeOf(obj).Name())
				panic(errorStr)
				return ""//丢弃这个包
			}
		}
	}

	return deleteSqlStr(sqltable, sqlData)
}

