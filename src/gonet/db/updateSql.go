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

func updatesql(sqlData *SqlData, p *Properties, val string){
	if p.IsPrimary(){
		sqlData.SqlName += fmt.Sprintf("`%s`='%s',", p.Name, val)
	}else{
		sqlData.SqlValue += fmt.Sprintf("`%s`='%s',", p.Name, val)
	}
}

func updatesqlblob(sqlData *SqlData, p *Properties, val []byte){
	if p.IsPrimary(){
		sqlData.SqlName += fmt.Sprintf("`%s`='0x%s',", p.Name, val)
	}else{
		sqlData.SqlValue += fmt.Sprintf("`%s`='0x%s',", p.Name, val)
	}
}

func updatesqlarray(sqlData *SqlData, p *Properties, val string, i int){
	if p.IsPrimary() {
		sqlData.SqlName += fmt.Sprintf("`%s%d`='%s',", p.Name, i, val)
	}else if sqlData.bitMap != nil && !sqlData.bitMap.Test(i){
		return
	}else{
		sqlData.SqlValue += fmt.Sprintf("`%s%d`='%s',", p.Name, i, val)
	}
}

func getUpdateSql(classField reflect.StructField, classVal reflect.Value, sqlData *SqlData) (bool) {
	p := getProperties(classField)
	sType := base.GetTypeStringEx(classField, classVal)
	if p.IsJson(){
		data, _ := json.Marshal(classVal.Interface())
		updatesql(sqlData, p, string(data))
		return true
	}else if p.IsBlob(){
		for classVal.Kind() == reflect.Ptr {
			classVal = classVal.Elem()
		}
		data, _ := proto.Marshal(classVal.Addr().Interface().(proto.Message))
		updatesqlblob(sqlData, p, data)
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
		updatesql(sqlData, p, strconv.FormatFloat(value, 'f', -1, 64))
	case "*float32":
		value := float32(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*float32)
		}
		updatesql(sqlData, p, strconv.FormatFloat(float64(value), 'f', -1, 32))
	case "*bool":
		value := bool(false)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*bool)
		}
		updatesql(sqlData, p, strconv.FormatBool(value))
	case "*int8":
		value := int8(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*int8)
		}
		updatesql(sqlData, p, strconv.FormatInt(int64(value),10))
	case "*uint8":
		value := uint8(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*uint8)
		}
		updatesql(sqlData, p, strconv.FormatUint(uint64(value),10))
	case "*int16":
		value := int16(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*int16)
		}
		updatesql(sqlData, p, strconv.FormatInt(int64(value),10))
	case "*uint16":
		value := uint16(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*uint16)
		}
		updatesql(sqlData, p, strconv.FormatUint(uint64(value),10))
	case "*int32":
		value := int32(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*int32)
		}
		updatesql(sqlData, p, strconv.FormatInt(int64(value),10))
	case "*uint32":
		value := uint32(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*uint32)
		}
		updatesql(sqlData, p, strconv.FormatUint(uint64(value),10))
	case "*int64":
		value := int64(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*int64)
		}
		if !p.IsDatetime(){
			updatesql(sqlData, p, strconv.FormatInt(int64(value),10))
		}else{
			updatesql(sqlData, p, GetDBTimeString(int64(value)))
		}
	case "*uint64":
		value := uint64(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*uint64)
		}
		updatesql(sqlData, p, strconv.FormatUint(uint64(value),10))
	case "*string":
		value := string("")
		if !classVal.IsNil() {
			value = *classVal.Interface().(*string)
		}
		updatesql(sqlData, p, value)
	case "*int":
		value := int(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*int)
		}
		updatesql(sqlData, p, strconv.FormatInt(int64(value),10))
	case "*uint":
		value := uint(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*uint)
		}
		updatesql(sqlData, p, strconv.FormatUint(uint64(value),10))
	case "*struct":
		if !classVal.IsNil() {
			value := classVal.Elem().Interface()
			parseUpdateSql(value, sqlData)
		}
	case "float64":
		updatesql(sqlData, p, strconv.FormatFloat(classVal.Float(), 'f', -1, 64))
	case "float32":
		updatesql(sqlData, p, strconv.FormatFloat(classVal.Float(), 'f', -1, 32))
	case "bool":
		updatesql(sqlData, p, strconv.FormatBool(classVal.Bool()))
	case "int8":
		updatesql(sqlData, p, strconv.FormatInt(classVal.Int(),10))
	case "uint8":
		updatesql(sqlData, p, strconv.FormatUint(classVal.Uint(),10))
	case "int16":
		updatesql(sqlData, p, strconv.FormatInt(classVal.Int(),10))
	case "uint16":
		updatesql(sqlData, p, strconv.FormatUint(classVal.Uint(),10))
	case "int32":
		updatesql(sqlData, p, strconv.FormatInt(classVal.Int(),10))
	case "uint32":
		updatesql(sqlData, p, strconv.FormatUint(classVal.Uint(), 10))
	case "int64":
		if !p.IsDatetime(){
			updatesql(sqlData, p, strconv.FormatInt(classVal.Int(),10))
		}else{
			updatesql(sqlData, p, GetDBTimeString(classVal.Int()))
		}
	case "uint64":
		updatesql(sqlData, p, strconv.FormatUint(classVal.Uint(),10))
	case "string":
		updatesql(sqlData, p, classVal.String())
	case "int":
		updatesql(sqlData, p, strconv.FormatInt(classVal.Int(),10))
	case "uint":
		updatesql(sqlData, p, strconv.FormatUint(classVal.Uint(),10))
	case "struct":
		parseUpdateSql(classVal.Interface(), sqlData)
	case "[]float64":
		value := []float64{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]float64)
		}
		for i,v := range value{
			updatesqlarray(sqlData, p, strconv.FormatFloat(v, 'f', -1, 64), i)
		}
	case "[]float32":
		value := []float32{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]float32)
		}
		for i,v := range value{
			updatesqlarray(sqlData, p, strconv.FormatFloat(float64(v), 'f', -1, 32), i)
		}
	case "[]bool":
		value := []bool{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]bool)
		}
		for i,v := range value{
			updatesqlarray(sqlData, p, strconv.FormatBool(v), i)
		}
	case "[]int8":
		value := []int8{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]int8)
		}
		for i,v := range value{
			updatesqlarray(sqlData, p, strconv.FormatInt(int64(v), 10), i)
		}
	case "[]uint8":
		value := []uint8{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]uint8)
		}
		for i,v := range value{
			updatesqlarray(sqlData, p, strconv.FormatUint(uint64(v), 10), i)
		}
	case "[]int16":
		value := []int16{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]int16)
		}
		for i,v := range value{
			updatesqlarray(sqlData, p, strconv.FormatInt(int64(v), 10), i)
		}
	case "[]uint16":
		value := []uint16{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]uint16)
		}
		for i,v := range value{
			updatesqlarray(sqlData, p, strconv.FormatUint(uint64(v), 10), i)
		}
	case "[]int32":
		value := []int32{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]int32)
		}
		for i,v := range value{
			updatesqlarray(sqlData, p, strconv.FormatInt(int64(v), 10), i)
		}
	case "[]uint32":
		value := []uint32{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]uint32)
		}
		for i,v := range value{
			updatesqlarray(sqlData, p, strconv.FormatUint(uint64(v), 10), i)
		}
	case "[]int64":
		value := []int64{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]int64)
		}
		for i,v := range value{
			if !p.IsDatetime(){
				updatesqlarray(sqlData, p, strconv.FormatInt(int64(v), 10), i)
			}else{
				updatesqlarray(sqlData, p, GetDBTimeString(v), i)
			}
		}
	case "[]uint64":
		value := []uint64{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]uint64)
		}
		for i,v := range value{
			updatesqlarray(sqlData, p, strconv.FormatUint(uint64(v), 10), i)
		}
	case "[]string":
		value := []string{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]string)
		}
		for i,v := range value{
			updatesqlarray(sqlData, p, v, i)
		}
	case "[]int":
		value := []int{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]int)
		}
		for i,v := range value{
			updatesqlarray(sqlData, p, strconv.FormatInt(int64(v), 10), i)
		}
	case "[]uint":
		value := []uint{}
		if !classVal.IsNil() {
			value = classVal.Interface().([]uint)
		}
		for i,v := range value{
			updatesqlarray(sqlData, p, strconv.FormatUint(uint64(v), 10), i)
		}
	case "[]struct"://no support
		for i := 0;  i < classVal.Len(); i++{
			parseUpdateSql(classVal.Index(i).Interface(), sqlData)
		}
	case "[*]float64":
		for i := 0;  i < classVal.Len(); i++{
			updatesqlarray(sqlData, p, strconv.FormatFloat(classVal.Index(i).Float(), 'f', -1, 64), i)
		}
	case "[*]float32":
		for i := 0;  i < classVal.Len(); i++{
			updatesqlarray(sqlData, p, strconv.FormatFloat(classVal.Index(i).Float(), 'f', -1, 64), i)
		}
	case "[*]bool":
		for i := 0;  i < classVal.Len(); i++{
			updatesqlarray(sqlData, p, strconv.FormatBool(classVal.Index(i).Bool()), i)
		}
	case "[*]int8":
		for i := 0;  i < classVal.Len(); i++{
			updatesqlarray(sqlData, p, strconv.FormatInt(classVal.Index(i).Int(), 10), i)
		}
	case "[*]uint8":
		for i := 0;  i < classVal.Len(); i++{
			updatesqlarray(sqlData, p, strconv.FormatUint(classVal.Index(i).Uint(), 10), i)
		}
	case "[*]int16":
		for i := 0;  i < classVal.Len(); i++{
			updatesqlarray(sqlData, p, strconv.FormatInt(classVal.Index(i).Int(), 10), i)
		}
	case "[*]uint16":
		for i := 0;  i < classVal.Len(); i++{
			updatesqlarray(sqlData, p, strconv.FormatUint(classVal.Index(i).Uint(), 10), i)
		}
	case "[*]int32":
		for i := 0;  i < classVal.Len(); i++{
			updatesqlarray(sqlData, p, strconv.FormatInt(classVal.Index(i).Int(), 10), i)
		}
	case "[*]uint32":
		for i := 0;  i < classVal.Len(); i++{
			updatesqlarray(sqlData, p, strconv.FormatUint(classVal.Index(i).Uint(), 10), i)
		}
	case "[*]int64":
		for i := 0;  i < classVal.Len(); i++{
			updatesqlarray(sqlData, p, strconv.FormatInt(classVal.Index(i).Int(), 10), i)
		}
	case "[*]uint64":
		for i := 0;  i < classVal.Len(); i++{
			updatesqlarray(sqlData, p, strconv.FormatUint(classVal.Index(i).Uint(), 10), i)
		}
	case "[*]string":
		for i := 0;  i < classVal.Len(); i++{
			updatesqlarray(sqlData, p, classVal.Index(i).String(), i)
		}
	case "[*]int":
		for i := 0;  i < classVal.Len(); i++{
			updatesqlarray(sqlData, p, strconv.FormatInt(classVal.Index(i).Int(), 10), i)
		}
	case "[*]uint":
		for i := 0;  i < classVal.Len(); i++{
			updatesqlarray(sqlData, p, strconv.FormatUint(classVal.Index(i).Uint(), 10), i)
		}
	case "[*]struct"://no support
		for i := 0;  i < classVal.Len(); i++{
			parseUpdateSql(classVal.Index(i).Interface(), sqlData)
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
			return false
		//}
	}
	return true
}

func parseUpdateSql(obj interface{}, sqlData *SqlData){
	classVal := reflect.ValueOf(obj)
	for classVal.Kind() == reflect.Ptr {
		classVal = classVal.Elem()
	}
	classType := classVal.Type()

	for i := 0; i < classType.NumField(); i++{
		if !classVal.Field(i).CanInterface(){
			continue
		}

		bRight:= getUpdateSql(classType.Field(i), classVal.Field(i), sqlData)
		if !bRight{
			errorStr := fmt.Sprintf("parseUpdateSql type not supported %s", classType.Name())
			panic(errorStr)
			return//丢弃这个包
		}
	}
	return
}

func updateSqlStr(sqltable string, sqlData *SqlData) string{
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

	sqlData := &SqlData{}
	parseUpdateSql(obj, sqlData)
	return  updateSqlStr(sqltable, sqlData)
}

func UpdateSqlEx(obj interface{}, sqltable string, params ...string) string {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("UpdateSqlEx", err)
		}
	}()

	classVal := reflect.ValueOf(obj)
	for classVal.Kind() == reflect.Ptr {
		classVal = classVal.Elem()
	}
	classType := classVal.Type()

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
		if exist || p.IsPrimary(){
			sqlData.bitMap = bitMap
			bRight := getUpdateSql(sf, classVal.Field(i), sqlData)
			if !bRight{
				errorStr := fmt.Sprintf("UpdateSqlEx error %s", reflect.TypeOf(obj).Name())
				panic(errorStr)
				return ""//丢弃这个包
			}
		}
	}
	return updateSqlStr(sqltable, sqlData)
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


