package db

import (
	"gonet/base"
	"reflect"
	"fmt"
	"strings"
	"strconv"
	"log"
)

const(
	update_sql = "`%s`='%s',"
	update_sqlarray = "`%s%d`='%s',"
)

func getUpdateSql(classField reflect.StructField, classVal reflect.Value) (bool,string,string) {
	classType := getSqlName(classField)
	/*defer func() {
		if err := recover(); err != nil {
			fmt.Println("getUpdateSql", classType, classVal,  err)
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
		*strsql += fmt.Sprintf(update_sql, classType, strconv.FormatFloat(value, 'f', -1, 64))
	case "*float32":
		value := float32(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*float32)
		}
		*strsql += fmt.Sprintf(update_sql, classType, strconv.FormatFloat(float64(value), 'f', -1, 32))
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
		*strsql += fmt.Sprintf(update_sql, classType, strconv.FormatInt(int64(value),10))
	case "*uint8":
		value := uint8(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*uint8)
		}
		*strsql += fmt.Sprintf(update_sql, classType, strconv.FormatUint(uint64(value),10))
	case "*int16":
		value := int16(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*int16)
		}
		*strsql += fmt.Sprintf(update_sql, classType, strconv.FormatInt(int64(value),10))
	case "*uint16":
		value := uint16(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*uint16)
		}
		*strsql += fmt.Sprintf(update_sql, classType, strconv.FormatUint(uint64(value),10))
	case "*int32":
		value := int32(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*int32)
		}
		*strsql += fmt.Sprintf(update_sql, classType, strconv.FormatInt(int64(value),10))
	case "*uint32":
		value := uint32(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*uint32)
		}
		*strsql += fmt.Sprintf(update_sql, classType, strconv.FormatUint(uint64(value),10))
	case "*int64":
		value := int64(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*int64)
		}
		if !isDatetime(classField){
			*strsql += fmt.Sprintf(update_sql, classType, strconv.FormatInt(int64(value),10))
		}else{
			*strsql += fmt.Sprintf(update_sql, classType, GetDBTimeString(int64(value)))
		}
	case "*uint64":
		value := uint64(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*uint64)
		}
		*strsql += fmt.Sprintf(update_sql, classType, strconv.FormatUint(uint64(value),10))
	case "*string":
		value := string("")
		if !classVal.IsNil() {
			value = *classVal.Interface().(*string)
		}
		*strsql += fmt.Sprintf(update_sql, classType, value)
	case "*int":
		value := int(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*int)
		}
		*strsql += fmt.Sprintf(update_sql, classType, strconv.FormatInt(int64(value),10))
	case "*uint":
		value := uint(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*uint)
		}
		*strsql += fmt.Sprintf(update_sql, classType, strconv.FormatUint(uint64(value),10))
	case "*struct":
		if !classVal.IsNil() {
			value := classVal.Elem().Interface()
			n, p := parseUpdateSql(value)
			noramlsql += n
			primarysql += p
		}
	case "float64":
		*strsql += fmt.Sprintf(update_sql, classType, strconv.FormatFloat(classVal.Float(), 'f', -1, 64))
	case "float32":
		*strsql += fmt.Sprintf(update_sql, classType, strconv.FormatFloat(classVal.Float(), 'f', -1, 32))
	case "bool":
		*strsql += fmt.Sprintf("%s=%t", classType, classVal.Bool())
	case "int8":
		*strsql += fmt.Sprintf(update_sql, classType, strconv.FormatInt(classVal.Int(),10))
	case "uint8":
		*strsql += fmt.Sprintf(update_sql, classType, strconv.FormatUint(classVal.Uint(),10))
	case "int16":
		*strsql += fmt.Sprintf(update_sql, classType, strconv.FormatInt(classVal.Int(),10))
	case "uint16":
		*strsql += fmt.Sprintf(update_sql, classType, strconv.FormatUint(classVal.Uint(),10))
	case "int32":
		*strsql += fmt.Sprintf(update_sql, classType, strconv.FormatInt(classVal.Int(),10))
	case "uint32":
		*strsql += fmt.Sprintf(update_sql, classType, strconv.FormatUint(classVal.Uint(), 10))
	case "int64":
		if !isDatetime(classField){
			*strsql += fmt.Sprintf(update_sql, classType, strconv.FormatInt(classVal.Int(),10))
		}else{
			*strsql += fmt.Sprintf(update_sql, classType, GetDBTimeString(classVal.Int()))
		}
	case "uint64":
		*strsql += fmt.Sprintf(update_sql, classType, strconv.FormatUint(classVal.Uint(),10))
	case "string":
		*strsql += fmt.Sprintf(update_sql, classType, classVal.String())
	case "int":
		*strsql += fmt.Sprintf(update_sql, classType, strconv.FormatInt(classVal.Int(),10))
	case "uint":
		*strsql += fmt.Sprintf(update_sql, classType, strconv.FormatUint(classVal.Uint(),10))
	case "struct":
		n, p := parseUpdateSql(classVal.Interface())
		noramlsql += n
		primarysql += p
	case "[]float64":
		value := []float64{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]float64)
		}
		for i,v := range value{
			*strsql += fmt.Sprintf(update_sqlarray, classType, i, strconv.FormatFloat(v, 'f', -1, 64))
		}
	case "[]float32":
		value := []float32{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]float32)
		}
		for i,v := range value{
			*strsql += fmt.Sprintf(update_sqlarray, classType, i, strconv.FormatFloat(float64(v), 'f', -1, 32))
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
			*strsql += fmt.Sprintf(update_sqlarray, classType, i, strconv.FormatInt(int64(v), 10))
		}
	case "[]uint8":
		value := []uint8{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]uint8)
		}
		if isBlob(classField){
			*strsql += fmt.Sprintf(update_sql, classType, classVal.Bytes())
		}else{
			for i,v := range value{
				*strsql += fmt.Sprintf(update_sqlarray, classType, i, strconv.FormatUint(uint64(v), 10))
			}
		}
	case "[]int16":
		value := []int16{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]int16)
		}
		for i,v := range value{
			*strsql += fmt.Sprintf(update_sqlarray, classType, i, strconv.FormatInt(int64(v), 10))
		}
	case "[]uint16":
		value := []uint16{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]uint16)
		}
		for i,v := range value{
			*strsql += fmt.Sprintf(update_sqlarray, classType, i, strconv.FormatUint(uint64(v), 10))
		}
	case "[]int32":
		value := []int32{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]int32)
		}
		for i,v := range value{
			*strsql += fmt.Sprintf(update_sqlarray, classType, i, strconv.FormatInt(int64(v), 10))
		}
	case "[]uint32":
		value := []uint32{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]uint32)
		}
		for i,v := range value{
			*strsql += fmt.Sprintf(update_sqlarray, classType, i, strconv.FormatUint(uint64(v), 10))
		}
	case "[]int64":
		value := []int64{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]int64)
		}
		for i,v := range value{
			if !isDatetime(classField){
				*strsql += fmt.Sprintf(update_sqlarray, classType, i, strconv.FormatInt(int64(v), 10))
			}else{
				*strsql += fmt.Sprintf(update_sqlarray, classType, i, GetDBTimeString(v))
			}
		}
	case "[]uint64":
		value := []uint64{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]uint64)
		}
		for i,v := range value{
			*strsql += fmt.Sprintf(update_sqlarray, classType, i, strconv.FormatUint(uint64(v), 10))
		}
	case "[]string":
		value := []string{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]string)
		}
		for i,v := range value{
			*strsql += fmt.Sprintf(update_sqlarray, classType, i, v)
		}
	case "[]int":
		value := []int{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]int)
		}
		for i,v := range value{
			*strsql += fmt.Sprintf(update_sqlarray, classType, i, strconv.FormatInt(int64(v), 10))
		}
	case "[]uint":
		value := []uint{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]uint)
		}
		for i,v := range value{
			*strsql += fmt.Sprintf(update_sqlarray, classType, i, strconv.FormatUint(uint64(v), 10))
		}
	case "[]struct"://no support
		for i := 0;  i < classVal.Len(); i++{
			n, p := parseUpdateSql(classVal.Index(i).Interface())
			noramlsql += n
			primarysql += p
		}
	case "[*]float64":
		for i := 0;  i < classVal.Len(); i++{
			*strsql += fmt.Sprintf(update_sqlarray, classType, i, strconv.FormatFloat(classVal.Index(i).Float(), 'f', -1, 64))
		}
	case "[*]float32":
		for i := 0;  i < classVal.Len(); i++{
			*strsql += fmt.Sprintf(update_sqlarray, classType, i, strconv.FormatFloat(classVal.Index(i).Float(), 'f', -1, 64))
		}
	case "[*]bool":
		for i := 0;  i < classVal.Len(); i++{
			*strsql += fmt.Sprintf("%s%d=%t,", classType, i, classVal.Index(i).Bool())
		}
	case "[*]int8":
		for i := 0;  i < classVal.Len(); i++{
			*strsql += fmt.Sprintf(update_sqlarray, classType, i, strconv.FormatInt(classVal.Index(i).Int(), 10))
		}
	case "[*]uint8":
		for i := 0;  i < classVal.Len(); i++{
			*strsql += fmt.Sprintf(update_sqlarray, classType, i, strconv.FormatUint(classVal.Index(i).Uint(), 10))
		}
	case "[*]int16":
		for i := 0;  i < classVal.Len(); i++{
			*strsql += fmt.Sprintf(update_sqlarray, classType, i, strconv.FormatInt(classVal.Index(i).Int(), 10))
		}
	case "[*]uint16":
		for i := 0;  i < classVal.Len(); i++{
			*strsql += fmt.Sprintf(update_sqlarray, classType, i, strconv.FormatUint(classVal.Index(i).Uint(), 10))
		}
	case "[*]int32":
		for i := 0;  i < classVal.Len(); i++{
			*strsql += fmt.Sprintf(update_sqlarray, classType, i, strconv.FormatInt(classVal.Index(i).Int(), 10))
		}
	case "[*]uint32":
		for i := 0;  i < classVal.Len(); i++{
			*strsql += fmt.Sprintf(update_sqlarray, classType, i, strconv.FormatUint(classVal.Index(i).Uint(), 10))
		}
	case "[*]int64":
		for i := 0;  i < classVal.Len(); i++{
			*strsql += fmt.Sprintf(update_sqlarray, classType, i, strconv.FormatInt(classVal.Index(i).Int(), 10))
		}
	case "[*]uint64":
		for i := 0;  i < classVal.Len(); i++{
			*strsql += fmt.Sprintf(update_sqlarray, classType, i, strconv.FormatUint(classVal.Index(i).Uint(), 10))
		}
	case "[*]string":
		for i := 0;  i < classVal.Len(); i++{
			*strsql += fmt.Sprintf(update_sqlarray, classType, i, classVal.Index(i).String())
		}
	case "[*]int":
		for i := 0;  i < classVal.Len(); i++{
			*strsql += fmt.Sprintf(update_sqlarray, classType, i, strconv.FormatInt(classVal.Index(i).Int(), 10))
		}
	case "[*]uint":
		for i := 0;  i < classVal.Len(); i++{
			*strsql += fmt.Sprintf(update_sqlarray, classType, i, strconv.FormatUint(classVal.Index(i).Uint(), 10))
		}
	case "[*]struct"://no support
		for i := 0;  i < classVal.Len(); i++{
			n, p := parseUpdateSql(classVal.Index(i).Interface())
			noramlsql += n
			primarysql += p
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
			fmt.Println("getUpdateSql type not supported", sType,  classField.Type)
			panic("getUpdateSql type not supported")
			return false, "", ""
		//}
	}
	return true, noramlsql, primarysql
}

func parseUpdateSql(obj interface{}) (string, string){
	var protoVal reflect.Value
	protoType := reflect.TypeOf(obj)
	if protoType.Kind() == reflect.Ptr {
		protoType = reflect.TypeOf(obj).Elem()
		protoVal = reflect.ValueOf(obj).Elem()
	}else if protoType.Kind() == reflect.Struct{
		protoVal = reflect.ValueOf(obj)
	}else{
		errorStr := fmt.Sprintf("parseUpdateSql no support ptr %s", protoType.Name())
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

		bRight, sqlstr, sqlprimary := getUpdateSql(protoType.Field(i), protoVal.Field(i))
		if !bRight{
			errorStr := fmt.Sprintf("parseUpdateSql type not supported %s", protoType.Name())
			panic(errorStr)
			return "",""//丢弃这个包
		}
		str += sqlstr
		primary += sqlprimary
	}
	return str,primary
}

func updateSqlStr(sqltable string, str string, primary string) string{
	index := strings.LastIndex(str, ",")
	if index!= -1{
		str = str[:index]
	}

	index = strings.LastIndex(primary, ",")
	if index!= -1{
		primary = primary[:index]
	}
	primary = strings.Replace(primary, ",", " and ", -1)
	return "update " + sqltable + " set " + str + " where "+ primary
}

//--- struct to sql
func UpdateSql(obj interface{}, sqltable string)string{
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("UpdateSql", err)
		}
	}()
	str, primary := parseUpdateSql(obj)
	return  updateSqlStr(sqltable, str, primary)
}

func UpdateSqlEx(obj interface{}, sqltable string, params ...string) string {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("UpdateSqlEx", err)
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
		if !protoVal.Field(i).CanInterface() {//private成员不能读取
			continue
		}

		sf := protoType.Field(i)
		_, exist := nameMap[getSqlName(sf)]
		if exist{
			bRight, name, value := getUpdateSql(sf, protoVal.Field(i))
			if !bRight{
				errorStr := fmt.Sprintf("UpdateSqlEx error %s", reflect.TypeOf(obj).Name())
				panic(errorStr)
				return ""//丢弃这个包
			}
			str += name
			primary += value
		}
	}
	return updateSqlStr(sqltable, str, primary)
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
	fmt.Println(db.UpdateSql(var1, "tb_test"))
	fmt.Println(db.UpdateSqlEx(var1, "tb_test", "I", "J2"))
	fmt.Println(db.LoadSql(var1, "tb_test","where playerid = 111"))
	fmt.Println(db.LoadSqlEx(var1,  "tb_test","where playerid = 111", "I", "J2",))
	fmt.Println(db.DeleteSql(var1, "tb_test"))
	fmt.Println(db.DeleteSqlEx(var1,  "tb_test", "I", "J2",))
 */


