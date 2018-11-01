package db

import (
	"base"
	"fmt"
	"log"
	"reflect"
	"unsafe"
)

func getLoadObjSql(classField reflect.StructField, classVal reflect.Value, row IRow) (bool) {
	if !classVal.CanSet(){
		return true
	}
	classType := getSqlName(classField)
	/*defer func() {
		if err := recover(); err != nil {
			fmt.Println("getLoadObjSql", classType,  err)
		}
	}()*/

	sType := base.GetTypeStringEx(classField, classVal)
	//fmt.Println(classVal, classType, sType, classVal.Type().String())
	switch sType {
	case "*float64":
		if classVal.IsNil() {
			reflect.NewAt(reflect.TypeOf(float64(0)), unsafe.Pointer(classVal.Pointer()))
		}
		classVal.Elem().SetFloat(row.Flot64(classType))
	case "*float32":
		if classVal.IsNil() {
			reflect.NewAt(reflect.TypeOf(float32(0)), unsafe.Pointer(classVal.Pointer()))
		}
		classVal.Elem().SetFloat(row.Flot64(classType))
	case "*bool":
		if classVal.IsNil() {
			reflect.NewAt(reflect.TypeOf(bool(false)), unsafe.Pointer(classVal.Pointer()))
		}
		classVal.Elem().SetBool(row.Bool(classType))
	case "*int8":
		if classVal.IsNil() {
			reflect.NewAt(reflect.TypeOf(int8(0)), unsafe.Pointer(classVal.Pointer()))
		}
		classVal.Elem().SetInt(row.Int64(classType))
	case "*uint8":
		if classVal.IsNil() {
			reflect.NewAt(reflect.TypeOf(uint8(0)), unsafe.Pointer(classVal.Pointer()))
		}
		classVal.Elem().SetUint(uint64(row.Int64(classType)))
	case "*int16":
		if classVal.IsNil() {
			reflect.NewAt(reflect.TypeOf(int16(0)), unsafe.Pointer(classVal.Pointer()))
		}
		classVal.Elem().SetInt(row.Int64(classType))
	case "*uint16":
		if classVal.IsNil() {
			reflect.NewAt(reflect.TypeOf(uint16(0)), unsafe.Pointer(classVal.Pointer()))
		}
		classVal.Elem().SetUint(uint64(row.Int64(classType)))
	case "*int32":
		if classVal.IsNil() {
			reflect.NewAt(reflect.TypeOf(int32(0)), unsafe.Pointer(classVal.Pointer()))
		}
		classVal.Elem().SetInt(row.Int64(classType))
	case "*uint32":
		if classVal.IsNil() {
			reflect.NewAt(reflect.TypeOf(uint32(0)), unsafe.Pointer(classVal.Pointer()))
		}
		classVal.Elem().SetUint(uint64(row.Int64(classType)))
	case "*int64":
		if classVal.IsNil() {
			reflect.NewAt(reflect.TypeOf(int64(0)), unsafe.Pointer(classVal.Pointer()))
		}
		classVal.Elem().SetInt(row.Int64(classType))
	case "*uint64":
		if classVal.IsNil() {
			reflect.NewAt(reflect.TypeOf(uint64(0)), unsafe.Pointer(classVal.Pointer()))
		}
		classVal.Elem().SetUint(uint64(row.Int64(classType)))
	case "*string":
		if classVal.IsNil() {
			reflect.NewAt(reflect.TypeOf(string("")), unsafe.Pointer(classVal.Pointer()))
		}
		classVal.Elem().SetString(row.String(classType))
	case "*int":
		if classVal.IsNil() {
			reflect.NewAt(reflect.TypeOf(int(0)), unsafe.Pointer(classVal.Pointer()))
		}
		classVal.Elem().SetInt(row.Int64(classType))
	case "*uint":
		if classVal.IsNil() {
			reflect.NewAt(reflect.TypeOf(uint(0)), unsafe.Pointer(classVal.Pointer()))
		}
		classVal.Elem().SetUint(uint64(row.Int64(classType)))
	case "*struct":
		if !classVal.IsNil() {
			value := classVal.Elem().Interface()
			parseLoadObjSql(value, row)
		}
	case "float64":
		classVal.SetFloat(row.Flot64(classType))
	case "float32":
		classVal.SetFloat(row.Flot64(classType))
	case "bool":
		classVal.SetBool(row.Bool(classType))
	case "int8":
		classVal.SetInt(row.Int64(classType))
	case "uint8":
		classVal.SetUint(uint64(row.Int64(classType)))
	case "int16":
		classVal.SetInt(row.Int64(classType))
	case "uint16":
		classVal.SetUint(uint64(row.Int64(classType)))
	case "int32":
		classVal.SetInt(row.Int64(classType))
	case "uint32":
		classVal.SetUint(uint64(row.Int64(classType)))
	case "int64":
		if !isDatetime(classField){
			classVal.SetInt(row.Int64(classType))
		}else{
			classVal.SetInt(row.Time(classType))
		}
	case "uint64":
		classVal.SetUint((uint64(row.Int64(classType))))
	case "string":
		classVal.SetString(row.String(classType))
	case "int":
		classVal.SetInt(row.Int64(classType))
	case "uint":
		classVal.SetUint((uint64(row.Int64(classType))))
	case "struct":
		parseLoadObjSql(classVal.Interface(), row)
	case "[]float64":
		if !classVal.CanSet() {
			for i := 0; i < classVal.Len(); i++{
				classVal.Index(i).SetFloat(row.Flot64(classType))
			}
		}
	case "[]float32":
		if !classVal.CanSet() {
			for i := 0; i < classVal.Len(); i++{
				classVal.Index(i).SetFloat(row.Flot64(classType))
			}
		}
	case "[]bool":
		if !classVal.CanSet() {
			for i := 0; i < classVal.Len(); i++{
				classVal.Index(i).SetBool(row.Bool(classType))
			}
		}
	case "[]int8":
		if !classVal.CanSet() {
			for i := 0; i < classVal.Len(); i++{
				classVal.Index(i).SetInt(row.Int64(classType))
			}
		}
	case "[]uint8":
		if !classVal.CanSet() {
			for i := 0; i < classVal.Len(); i++{
				classVal.Index(i).SetUint(uint64(row.Int64(classType)))
			}
		}
	case "[]int16":
		if !classVal.CanSet() {
			for i := 0; i < classVal.Len(); i++{
				classVal.Index(i).SetInt(row.Int64(classType))
			}
		}
	case "[]uint16":
		if !classVal.CanSet() {
			for i := 0; i < classVal.Len(); i++{
				classVal.Index(i).SetUint(uint64(row.Int64(classType)))
			}
		}
	case "[]int32":
		if !classVal.CanSet() {
			for i := 0; i < classVal.Len(); i++{
				classVal.Index(i).SetInt(row.Int64(classType))
			}
		}
	case "[]uint32":
		if !classVal.CanSet() {
			for i := 0; i < classVal.Len(); i++{
				classVal.Index(i).SetUint(uint64(row.Int64(classType)))
			}
		}
	case "[]int64":
		if !classVal.CanSet() {
			for i := 0; i < classVal.Len(); i++{
				classVal.Index(i).SetInt(row.Int64(classType))
			}
		}
	case "[]uint64":
		if !classVal.CanSet() {
			for i := 0; i < classVal.Len(); i++{
				classVal.Index(i).SetUint(uint64(row.Int64(classType)))
			}
		}
	case "[]string":
		if !classVal.CanSet() {
			for i := 0; i < classVal.Len(); i++{
				classVal.Index(i).SetString(row.String(classType))
			}
		}
	case "[]int":
		if !classVal.CanSet() {
			for i := 0; i < classVal.Len(); i++{
				classVal.Index(i).SetInt(row.Int64(classType))
			}
		}
	case "[]uint":
		if !classVal.CanSet() {
			for i := 0; i < classVal.Len(); i++{
				classVal.Index(i).SetUint(uint64(row.Int64(classType)))
			}
		}
	case "[]struct":
		for i := 0;  i < classVal.Len(); i++{
			parseLoadObjSql(classVal.Index(i).Interface(), row)
		}
	case "[*]float64":
		for i := 0; i < classVal.Len(); i++{
			classVal.Index(i).SetFloat(row.Flot64(classType))
		}
	case "[*]float32":
		for i := 0; i < classVal.Len(); i++{
			classVal.Index(i).SetFloat(row.Flot64(classType))
		}
	case "[*]bool":
		for i := 0; i < classVal.Len(); i++{
			classVal.Index(i).SetBool(row.Bool(classType))
		}
	case "[*]int8":
		for i := 0; i < classVal.Len(); i++{
			classVal.Index(i).SetInt(row.Int64(classType))
		}
	case "[*]uint8":
		for i := 0; i < classVal.Len(); i++{
			classVal.Index(i).SetUint(uint64(row.Int64(classType)))
		}
	case "[*]int16":
		for i := 0; i < classVal.Len(); i++{
			classVal.Index(i).SetInt(row.Int64(classType))
		}
	case "[*]uint16":
		for i := 0; i < classVal.Len(); i++{
			classVal.Index(i).SetUint(uint64(row.Int64(classType)))
		}
	case "[*]int32":
		for i := 0; i < classVal.Len(); i++{
			classVal.Index(i).SetInt(row.Int64(classType))
		}
	case "[*]uint32":
		for i := 0; i < classVal.Len(); i++{
			classVal.Index(i).SetUint(uint64(row.Int64(classType)))
		}
	case "[*]int64":
		for i := 0; i < classVal.Len(); i++{
			classVal.Index(i).SetInt(row.Int64(classType))
		}
	case "[*]uint64":
		for i := 0; i < classVal.Len(); i++{
			classVal.Index(i).SetUint(uint64(row.Int64(classType)))
		}
	case "[*]string":
		for i := 0; i < classVal.Len(); i++{
			classVal.Index(i).SetString(row.String(classType))
		}
	case "[*]int":
		for i := 0; i < classVal.Len(); i++{
			classVal.Index(i).SetInt(row.Int64(classType))
		}
	case "[*]uint":
		for i := 0; i < classVal.Len(); i++{
			classVal.Index(i).SetUint(uint64(row.Int64(classType)))
		}
	case "[*]struct":
		for i := 0;  i < classVal.Len(); i++{
			parseLoadObjSql(classVal.Index(i).Interface(), row)
		}
	default:
		fmt.Println("getLoadObjSql type not supported", sType,  classField.Type)
		panic("getLoadObjSql type not supported")
		return false
		//}
	}
	return true
}

func parseLoadObjSql(obj interface{}, row IRow) (bool){
	var protoVal reflect.Value
	protoType := reflect.TypeOf(obj)
	if protoType.Kind() == reflect.Ptr {
		protoType = reflect.TypeOf(obj).Elem()
		protoVal = reflect.ValueOf(obj).Elem()
	}else if protoType.Kind() == reflect.Struct{
		errorStr := fmt.Sprintf("parseLoadObjSql no support struct %s", protoType.Name())
		log.Println(errorStr)
		panic(errorStr)
		return false
	} else{
		errorStr := fmt.Sprintf("parseLoadObjSql no support %s", protoType.Name())
		log.Println(errorStr)
		panic(errorStr)
		return false
	}

	for i := 0; i < protoType.NumField(); i++{
		if !protoVal.Field(i).CanInterface(){
			continue
		}

		bRight := getLoadObjSql(protoType.Field(i), protoVal.Field(i), row)
		if !bRight{
			errorStr := fmt.Sprintf("parseLoadObjSql type not supported %s", protoType.Name())
			panic(errorStr)
			return false//丢弃这个包
		}
	}
	return true
}

//--- struct to sql
func LoadObjSql(obj interface{}, row IRow)bool{
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("LoadObjSql", err)
		}
	}()

	if row == nil{
		return false
	}

	return  parseLoadObjSql(obj, row)
}



