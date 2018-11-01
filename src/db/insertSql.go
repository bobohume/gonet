package db
import (
	"base"
	"reflect"
"fmt"
"strings"
"strconv"
"log"
)

const(
	insert_sql = "'%s',"
	insert_sqlarray = "'%s',"
	insert_sqlname = "%s,"
	insert_sqlarrayname = "%s%d,"
)

func getInsertSql(classField reflect.StructField, classVal reflect.Value) (bool,string,string) {
	classType := getSqlName(classField)
	/*defer func() {
		if err := recover(); err != nil {
			fmt.Println("getInsertSql", classType,  err)
		}
	}()*/

	sType := base.GetTypeStringEx(classField, classVal)
	//fmt.Println(classVal, classType, sType, classVal.Type().String())
	sqlname := ""
	sqlvalue := ""
	switch sType {
	case "*float64":
		value := float64(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*float64)
		}
		sqlvalue += fmt.Sprintf(insert_sql, strconv.FormatFloat(value, 'f', -1, 64))
		sqlname += fmt.Sprintf(insert_sqlname, classType)
	case "*float32":
		value := float32(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*float32)
		}
		sqlvalue += fmt.Sprintf(insert_sql, strconv.FormatFloat(float64(value), 'f', -1, 32))
		sqlname += fmt.Sprintf(insert_sqlname, classType)
	case "*bool":
		value := bool(false)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*bool)
		}
		sqlvalue += fmt.Sprintf("%t", value)
		sqlname += fmt.Sprintf(insert_sqlname, classType)
	case "*int8":
		value := int8(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*int8)
		}
		sqlvalue += fmt.Sprintf(insert_sql, strconv.FormatInt(int64(value),10))
		sqlname += fmt.Sprintf(insert_sqlname, classType)
	case "*uint8":
		value := uint8(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*uint8)
		}
		sqlvalue += fmt.Sprintf(insert_sql, strconv.FormatUint(uint64(value),10))
		sqlname += fmt.Sprintf(insert_sqlname, classType)
	case "*int16":
		value := int16(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*int16)
		}
		sqlvalue += fmt.Sprintf(insert_sql, strconv.FormatInt(int64(value),10))
		sqlname += fmt.Sprintf(insert_sqlname, classType)
	case "*uint16":
		value := uint16(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*uint16)
		}
		sqlvalue += fmt.Sprintf(insert_sql, strconv.FormatUint(uint64(value),10))
		sqlname += fmt.Sprintf(insert_sqlname, classType)
	case "*int32":
		value := int32(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*int32)
		}
		sqlvalue += fmt.Sprintf(insert_sql, strconv.FormatInt(int64(value),10))
		sqlname += fmt.Sprintf(insert_sqlname, classType)
	case "*uint32":
		value := uint32(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*uint32)
		}
		sqlvalue += fmt.Sprintf(insert_sql, strconv.FormatUint(uint64(value),10))
		sqlname += fmt.Sprintf(insert_sqlname, classType)
	case "*int64":
		value := int64(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*int64)
		}
		if !isDatetime(classField){
			sqlvalue += fmt.Sprintf(insert_sql, strconv.FormatInt(int64(value),10))
		}else{
			sqlvalue += fmt.Sprintf(insert_sql, GetDBTimeString(int64(value)))
		}
		sqlname += fmt.Sprintf(insert_sqlname, classType)
	case "*uint64":
		value := uint64(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*uint64)
		}
		sqlvalue += fmt.Sprintf(insert_sql, strconv.FormatUint(uint64(value),10))
		sqlname += fmt.Sprintf(insert_sqlname, classType)
	case "*string":
		value := string("")
		if !classVal.IsNil() {
			value = *classVal.Interface().(*string)
		}
		sqlvalue += fmt.Sprintf(insert_sql, value)
		sqlname += fmt.Sprintf(insert_sqlname, classType)
	case "*int":
		value := int(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*int)
		}
		sqlvalue += fmt.Sprintf(insert_sql, strconv.FormatInt(int64(value),10))
		sqlname += fmt.Sprintf(insert_sqlname, classType)
	case "*uint":
		value := uint(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*uint)
		}
		sqlvalue += fmt.Sprintf(insert_sql, strconv.FormatUint(uint64(value),10))
		sqlname += fmt.Sprintf(insert_sqlname, classType)
	case "*struct":
		value := classVal.Elem().Interface()
		n,p := parseInserSql(value)
		sqlname += n
		sqlvalue += p
	case "float64":
		sqlvalue += fmt.Sprintf(insert_sql, strconv.FormatFloat(classVal.Float(), 'f', -1, 64))
		sqlname += fmt.Sprintf(insert_sqlname, classType)
	case "float32":
		sqlvalue += fmt.Sprintf(insert_sql, strconv.FormatFloat(classVal.Float(), 'f', -1, 32))
		sqlname += fmt.Sprintf(insert_sqlname, classType)
	case "bool":
		sqlvalue += fmt.Sprintf("%t", classVal.Bool())
		sqlname += fmt.Sprintf(insert_sqlname, classType)
	case "int8":
		sqlvalue += fmt.Sprintf(insert_sql, strconv.FormatInt(classVal.Int(),10))
		sqlname += fmt.Sprintf(insert_sqlname, classType)
	case "uint8":
		sqlvalue += fmt.Sprintf(insert_sql, strconv.FormatUint(classVal.Uint(),10))
		sqlname += fmt.Sprintf(insert_sqlname, classType)
	case "int16":
		sqlvalue += fmt.Sprintf(insert_sql, strconv.FormatInt(classVal.Int(),10))
		sqlname += fmt.Sprintf(insert_sqlname, classType)
	case "uint16":
		sqlvalue += fmt.Sprintf(insert_sql, strconv.FormatUint(classVal.Uint(),10))
		sqlname += fmt.Sprintf(insert_sqlname, classType)
	case "int32":
		sqlvalue += fmt.Sprintf(insert_sql, strconv.FormatInt(classVal.Int(),10))
		sqlname += fmt.Sprintf(insert_sqlname, classType)
	case "uint32":
		sqlvalue += fmt.Sprintf(insert_sql, strconv.FormatUint(classVal.Uint(), 10))
		sqlname += fmt.Sprintf(insert_sqlname, classType)
	case "int64":
		if !isDatetime(classField){
			sqlvalue += fmt.Sprintf(insert_sql, strconv.FormatInt(classVal.Int(),10))
		}else{
			sqlvalue += fmt.Sprintf(insert_sql, GetDBTimeString(classVal.Int()))
		}
		sqlname += fmt.Sprintf(insert_sqlname, classType)
	case "uint64":
		sqlvalue += fmt.Sprintf(insert_sql, strconv.FormatUint(classVal.Uint(),10))
		sqlname += fmt.Sprintf(insert_sqlname, classType)
	case "string":
		sqlvalue += fmt.Sprintf(insert_sql, classVal.String())
		sqlname += fmt.Sprintf(insert_sqlname, classType)
	case "int":
		sqlvalue += fmt.Sprintf(insert_sql, strconv.FormatInt(classVal.Int(),10))
		sqlname += fmt.Sprintf(insert_sqlname, classType)
	case "uint":
		sqlvalue += fmt.Sprintf(insert_sql, strconv.FormatUint(classVal.Uint(),10))
		sqlname += fmt.Sprintf(insert_sqlname, classType)
	case "struct":
		n, p := parseInserSql(classVal.Interface())
		sqlname += n
		sqlvalue += p
	case "[]float64":
		value := []float64{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]float64)
		}
		for i,v := range value{
			sqlvalue += fmt.Sprintf(insert_sqlarray, strconv.FormatFloat(v, 'f', -1, 64))
			sqlname += fmt.Sprintf(insert_sqlarrayname, classType, i)
		}
	case "[]float32":
		value := []float32{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]float32)
		}
		for i,v := range value{
			sqlvalue += fmt.Sprintf(insert_sqlarray, strconv.FormatFloat(float64(v), 'f', -1, 32))
			sqlname += fmt.Sprintf(insert_sqlarrayname, classType, i)
		}
	case "[]bool":
		value := []bool{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]bool)
		}
		for i,v := range value{
			sqlvalue += fmt.Sprintf("%t,", v)
			sqlname += fmt.Sprintf(insert_sqlarrayname, classType, i)
		}
	case "[]int8":
		value := []int8{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]int8)
		}
		for i,v := range value{
			sqlvalue += fmt.Sprintf(insert_sqlarray, strconv.FormatInt(int64(v), 10))
			sqlname += fmt.Sprintf(insert_sqlarrayname, classType, i)
		}
	case "[]uint8":
		value := []uint8{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]uint8)
		}
		for i,v := range value{
			sqlvalue += fmt.Sprintf(insert_sqlarray, strconv.FormatUint(uint64(v), 10))
			sqlname += fmt.Sprintf(insert_sqlarrayname, classType, i)
		}
	case "[]int16":
		value := []int16{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]int16)
		}
		for i,v := range value{
			sqlvalue += fmt.Sprintf(insert_sqlarray, strconv.FormatInt(int64(v), 10))
			sqlname += fmt.Sprintf(insert_sqlarrayname, classType, i)
		}
	case "[]uint16":
		value := []uint16{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]uint16)
		}
		for i,v := range value{
			sqlvalue += fmt.Sprintf(insert_sqlarray, strconv.FormatUint(uint64(v), 10))
			sqlname += fmt.Sprintf(insert_sqlarrayname, classType, i)
		}
	case "[]int32":
		value := []int32{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]int32)
		}
		for i,v := range value{
			sqlvalue += fmt.Sprintf(insert_sqlarray, strconv.FormatInt(int64(v), 10))
			sqlname += fmt.Sprintf(insert_sqlarrayname, classType, i)
		}
	case "[]uint32":
		value := []uint32{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]uint32)
		}
		for i,v := range value{
			sqlvalue += fmt.Sprintf(insert_sqlarray, strconv.FormatUint(uint64(v), 10))
			sqlname += fmt.Sprintf(insert_sqlarrayname, classType, i)
		}
	case "[]int64":
		value := []int64{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]int64)
		}
		for i,v := range value{
			if !isDatetime(classField){
				sqlvalue += fmt.Sprintf(insert_sqlarray, strconv.FormatInt(int64(v), 10))
			}else{
				sqlvalue += fmt.Sprintf(insert_sqlarray, GetDBTimeString(v))
			}
			sqlname += fmt.Sprintf(insert_sqlarrayname, classType, i)
		}
	case "[]uint64":
		value := []uint64{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]uint64)
		}
		for i,v := range value{
			sqlvalue += fmt.Sprintf(insert_sqlarray, strconv.FormatUint(uint64(v), 10))
			sqlname += fmt.Sprintf(insert_sqlarrayname, classType, i)
		}
	case "[]string":
		value := []string{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]string)
		}
		for i,v := range value{
			sqlvalue += fmt.Sprintf(insert_sqlarray, v)
			sqlname += fmt.Sprintf(insert_sqlarrayname, classType, i)
		}
	case "[]int":
		value := []int{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]int)
		}
		for i,v := range value{
			sqlvalue += fmt.Sprintf(insert_sqlarray, strconv.FormatInt(int64(v), 10))
			sqlname += fmt.Sprintf(insert_sqlarrayname, classType, i)
		}
	case "[]uint":
		value := []uint{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]uint)
		}
		for i,v := range value{
			sqlvalue += fmt.Sprintf(insert_sqlarray, strconv.FormatUint(uint64(v), 10))
			sqlname += fmt.Sprintf(insert_sqlarrayname, classType, i)
		}
	case "[]struct":
		for i := 0;  i < classVal.Len(); i++{
			n, p := parseInserSql(classVal.Index(i).Interface())
			sqlname += n
			sqlvalue += p
		}
	case "[*]float64":
		for i := 0; i < classVal.Len(); i++{
			sqlvalue += fmt.Sprintf(insert_sqlarray, strconv.FormatFloat(classVal.Index(i).Float(), 'f', -1, 64))
			sqlname += fmt.Sprintf(insert_sqlarrayname, classType, i)
		}
	case "[*]float32":
		for i := 0; i < classVal.Len(); i++{
			sqlvalue += fmt.Sprintf(insert_sqlarray, strconv.FormatFloat(classVal.Index(i).Float(), 'f', -1, 64))
			sqlname += fmt.Sprintf(insert_sqlarrayname, classType, i)
		}
	case "[*]bool":
		for i := 0; i < classVal.Len(); i++{
			sqlvalue += fmt.Sprintf("%t,", classVal.Index(i).Bool())
			sqlname += fmt.Sprintf(insert_sqlarrayname, classType, i)
		}
	case "[*]int8":
		for i := 0; i < classVal.Len(); i++{
			sqlvalue += fmt.Sprintf(insert_sqlarray, strconv.FormatInt(classVal.Index(i).Int(), 10))
			sqlname += fmt.Sprintf(insert_sqlarrayname, classType, i)
		}
	case "[*]uint8":
		for i := 0; i < classVal.Len(); i++{
			sqlvalue += fmt.Sprintf(insert_sqlarray, strconv.FormatUint(classVal.Index(i).Uint(), 10))
			sqlname += fmt.Sprintf(insert_sqlarrayname, classType, i)
		}
	case "[*]int16":
		for i := 0; i < classVal.Len(); i++{
			sqlvalue += fmt.Sprintf(insert_sqlarray, strconv.FormatInt(classVal.Index(i).Int(), 10))
			sqlname += fmt.Sprintf(insert_sqlarrayname, classType, i)
		}
	case "[*]uint16":
		for i := 0; i < classVal.Len(); i++{
			sqlvalue += fmt.Sprintf(insert_sqlarray, strconv.FormatUint(classVal.Index(i).Uint(), 10))
			sqlname += fmt.Sprintf(insert_sqlarrayname, classType, i)
		}
	case "[*]int32":
		for i := 0; i < classVal.Len(); i++{
			sqlvalue += fmt.Sprintf(insert_sqlarray, strconv.FormatInt(classVal.Index(i).Int(), 10))
			sqlname += fmt.Sprintf(insert_sqlarrayname, classType, i)
		}
	case "[*]uint32":
		for i := 0; i < classVal.Len(); i++{
			sqlvalue += fmt.Sprintf(insert_sqlarray, strconv.FormatUint(classVal.Index(i).Uint(), 10))
			sqlname += fmt.Sprintf(insert_sqlarrayname, classType, i)
		}
	case "[*]int64":
		for i := 0; i < classVal.Len(); i++{
			sqlvalue += fmt.Sprintf(insert_sqlarray, strconv.FormatInt(classVal.Index(i).Int(), 10))
			sqlname += fmt.Sprintf(insert_sqlarrayname, classType, i)
		}
	case "[*]uint64":
		for i := 0; i < classVal.Len(); i++{
			sqlvalue += fmt.Sprintf(insert_sqlarray, strconv.FormatUint(classVal.Index(i).Uint(), 10))
			sqlname += fmt.Sprintf(insert_sqlarrayname, classType, i)
		}
	case "[*]string":
		for i := 0; i < classVal.Len(); i++{
			sqlvalue += fmt.Sprintf(insert_sqlarray, classVal.Index(i).String())
			sqlname += fmt.Sprintf(insert_sqlarrayname, classType, i)
		}
	case "[*]int":
		for i := 0; i < classVal.Len(); i++{
			sqlvalue += fmt.Sprintf(insert_sqlarray, strconv.FormatInt(classVal.Index(i).Int(), 10))
			sqlname += fmt.Sprintf(insert_sqlarrayname, classType, i)
		}
	case "[*]uint":
		for i := 0; i < classVal.Len(); i++{
			sqlvalue += fmt.Sprintf(insert_sqlarray, strconv.FormatUint(classVal.Index(i).Uint(), 10))
			sqlname += fmt.Sprintf(insert_sqlarrayname, classType, i)
		}
	case "[*]struct":
		for i := 0;  i < classVal.Len(); i++{
			n, p := parseInserSql(classVal.Index(i).Interface())
			sqlname += n
			sqlvalue += p
		}
	default:
		fmt.Println("getInsertSql type not supported", sType,  classField.Type)
		panic("getInsertSql type not supported")
		return false, "", ""
		//}
	}
	return true, sqlname, sqlvalue
}

func parseInserSql(obj interface{}) (string, string){
	var protoVal reflect.Value
	protoType := reflect.TypeOf(obj)
	if protoType.Kind() == reflect.Ptr {
		protoType = reflect.TypeOf(obj).Elem()
		protoVal = reflect.ValueOf(obj).Elem()
	}else if protoType.Kind() == reflect.Struct{
		protoVal = reflect.ValueOf(obj)
	}else{
		errorStr := fmt.Sprintf("parseInserSql no support %s", protoType.Name())
		log.Println(errorStr)
		panic(errorStr)
		return "",""
	}

	sqlname := ""
	sqlvalue := ""
	for i := 0; i < protoType.NumField(); i++{
		if !protoVal.Field(i).CanInterface(){
			continue
		}

		bRight, name, value := getInsertSql(protoType.Field(i), protoVal.Field(i))
		if !bRight{
			errorStr := fmt.Sprintf("parseInserSql type not supported %s", protoType.Name())
			panic(errorStr)
			return "",""//丢弃这个包
		}
		sqlname += name
		sqlvalue += value
	}
	return sqlname,sqlvalue
}

func insertSqlStr(sqltable string, sqlname string, sqlvalue string) string{
	index := strings.LastIndex(sqlname, ",")
	if index!= -1{
		sqlname = sqlname[:index]
	}

	index = strings.LastIndex(sqlvalue, ",")
	if index!= -1{
		sqlvalue = sqlvalue[:index]
	}
	return "insert into "+ sqltable + " (" + sqlname+") VALUES (" + sqlvalue + ")"
}

//--- struct to sql
func InsertSql(obj interface{}, sqltable string,)string{
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("InsertSql", err)
		}
	}()

	sqlname, sqlvalue := parseInserSql(obj)
	return  insertSqlStr(sqltable, sqlname, sqlvalue)
}

func InsertSqlEx(obj interface{}, sqltable string, params ...string) string {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("InsertSqlEx", err)
		}
	}()

	protoVal  := reflect.ValueOf(obj)
	protoType := reflect.TypeOf(obj)
	if protoType.Kind() == reflect.Ptr {
		protoType = reflect.TypeOf(obj).Elem()
		protoVal = reflect.ValueOf(obj).Elem()
	}

	sqlname := ""
	sqlvalue := ""
	nameMap := make(map[string] string)
	for _,v := range params{
		v1 := strings.ToLower(v)
		nameMap[v1] = v1
	}
	for i := 0; i < protoType.NumField(); i++ {
		if !protoVal.Field(i).CanInterface() {//private成员不能读取
			continue
		}

		sf := protoType.Field(i)
		_, exist := nameMap[getSqlName(sf)]
		if exist{
			bRight, name, value := getInsertSql(sf, protoVal.Field(i))
			if !bRight{
				errorStr := fmt.Sprintf("InsertSqlEx error %s", reflect.TypeOf(obj).Name())
				panic(errorStr)
				return ""//丢弃这个包
			}
			sqlname += name
			sqlvalue += value
		}
	}
	return insertSqlStr(sqltable, sqlname, sqlvalue)
}



