package db

import (
	"gonet/base"
	"fmt"
	"log"
	"reflect"
	"unsafe"
)

const(
	load_obj_sqlarrayname = "%s%d"
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
		value :=  (**float64)(unsafe.Pointer(uintptr(unsafe.Pointer(classVal.Addr().Pointer()))))
		val1 := row.Float64(classType)
		*value = &val1
	case "*float32":
		value :=  (**float32)(unsafe.Pointer(uintptr(unsafe.Pointer(classVal.Addr().Pointer()))))
		val1 := float32(row.Float32(classType))
		*value = &val1
	case "*bool":
		value :=  (**bool)(unsafe.Pointer(uintptr(unsafe.Pointer(classVal.Addr().Pointer()))))
		val1 := row.Bool(classType)
		*value = &val1
	case "*int8":
		value :=  (**int8)(unsafe.Pointer(uintptr(unsafe.Pointer(classVal.Addr().Pointer()))))
		val1 := int8(row.Int(classType))
		*value = &val1
	case "*uint8":
		value :=  (**uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(classVal.Addr().Pointer()))))
		val1 := uint8(row.Int(classType))
		*value = &val1
	case "*int16":
		value :=  (**int16)(unsafe.Pointer(uintptr(unsafe.Pointer(classVal.Addr().Pointer()))))
		val1 := int16(row.Int(classType))
		*value = &val1
	case "*uint16":
		value :=  (**uint16)(unsafe.Pointer(uintptr(unsafe.Pointer(classVal.Addr().Pointer()))))
		val1 := uint16(row.Int(classType))
		*value = &val1
	case "*int32":
		value :=  (**int32)(unsafe.Pointer(uintptr(unsafe.Pointer(classVal.Addr().Pointer()))))
		val1 := int32(row.Int(classType))
		*value = &val1
	case "*uint32":
		value :=  (**uint32)(unsafe.Pointer(uintptr(unsafe.Pointer(classVal.Addr().Pointer()))))
		val1 := uint32(row.Int(classType))
		*value = &val1
	case "*int64":
		value :=  (**int64)(unsafe.Pointer(uintptr(unsafe.Pointer(classVal.Addr().Pointer()))))
		val1 := int64(row.Int64(classType))
		*value = &val1
	case "*uint64":
		value :=  (**uint64)(unsafe.Pointer(uintptr(unsafe.Pointer(classVal.Addr().Pointer()))))
		val1 := uint64(row.Int64(classType))
		*value = &val1
	case "*string":
		value :=  (**string)(unsafe.Pointer(uintptr(unsafe.Pointer(classVal.Addr().Pointer()))))
		val1 := row.String(classType)
		*value = &val1
	case "*int":
		value :=  (**int)(unsafe.Pointer(uintptr(unsafe.Pointer(classVal.Addr().Pointer()))))
		val1 := int(row.Int(classType))
		*value = &val1
	case "*uint":
		value :=  (**uint)(unsafe.Pointer(uintptr(unsafe.Pointer(classVal.Addr().Pointer()))))
		val1 := uint(row.Int(classType))
		*value = &val1
	case "*struct":
		if !classVal.IsNil() {
			value := classVal.Elem().Interface()
			parseLoadObjSql(value, row)
		}
	case "float64":
		classVal.SetFloat(row.Float64(classType))
	case "float32":
		classVal.SetFloat(row.Float64(classType))
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
		if classVal.CanSet() {
			for i := 0; i < classVal.Len(); i++{
				classVal.Index(i).SetFloat(row.Float64(fmt.Sprintf(load_obj_sqlarrayname, classType, i)))
			}
		}
	case "[]float32":
		if classVal.CanSet() {
			for i := 0; i < classVal.Len(); i++{
				classVal.Index(i).SetFloat(row.Float64(fmt.Sprintf(load_obj_sqlarrayname, classType, i)))
			}
		}
	case "[]bool":
		if classVal.CanSet() {
			for i := 0; i < classVal.Len(); i++{
				classVal.Index(i).SetBool(row.Bool(fmt.Sprintf(load_obj_sqlarrayname, classType, i)))
			}
		}
	case "[]int8":
		if classVal.CanSet() {
			for i := 0; i < classVal.Len(); i++{
				classVal.Index(i).SetInt(row.Int64(fmt.Sprintf(load_obj_sqlarrayname, classType, i)))
			}
		}
	case "[]uint8":
		classVal.SetBytes(row.Byte(classType)) //blo
		if classVal.CanSet() {
			if isBlob(classField){
				classVal.SetBytes(row.Byte(classType)) //blob
			}else{
				for i := 0; i < classVal.Len(); i++{
					classVal.Index(i).SetUint(uint64(row.Int64(fmt.Sprintf(load_obj_sqlarrayname, classType, i))))
				}
			}
		}
	case "[]int16":
		if classVal.CanSet() {
			for i := 0; i < classVal.Len(); i++{
				classVal.Index(i).SetInt(row.Int64(fmt.Sprintf(load_obj_sqlarrayname, classType, i)))
			}
		}
	case "[]uint16":
		if classVal.CanSet() {
			for i := 0; i < classVal.Len(); i++{
				classVal.Index(i).SetUint(uint64(row.Int64(fmt.Sprintf(load_obj_sqlarrayname, classType, i))))
			}
		}
	case "[]int32":
		if classVal.CanSet() {
			for i := 0; i < classVal.Len(); i++{
				classVal.Index(i).SetInt(row.Int64(fmt.Sprintf(load_obj_sqlarrayname, classType, i)))
			}
		}
	case "[]uint32":
		if classVal.CanSet() {
			for i := 0; i < classVal.Len(); i++{
				classVal.Index(i).SetUint(uint64(row.Int64(fmt.Sprintf(load_obj_sqlarrayname, classType, i))))
			}
		}
	case "[]int64":
		if classVal.CanSet() {
			for i := 0; i < classVal.Len(); i++{
				classVal.Index(i).SetInt(row.Int64(fmt.Sprintf(load_obj_sqlarrayname, classType, i)))
			}
		}
	case "[]uint64":
		if classVal.CanSet() {
			for i := 0; i < classVal.Len(); i++{
				classVal.Index(i).SetUint(uint64(row.Int64(fmt.Sprintf(load_obj_sqlarrayname, classType, i))))
			}
		}
	case "[]string":
		if classVal.CanSet() {
			for i := 0; i < classVal.Len(); i++{
				classVal.Index(i).SetString(row.String(fmt.Sprintf(load_obj_sqlarrayname, classType, i)))
			}
		}
	case "[]int":
		if classVal.CanSet() {
			for i := 0; i < classVal.Len(); i++{
				classVal.Index(i).SetInt(row.Int64(fmt.Sprintf(load_obj_sqlarrayname, classType, i)))
			}
		}
	case "[]uint":
		if classVal.CanSet() {
			for i := 0; i < classVal.Len(); i++{
				classVal.Index(i).SetUint(uint64(row.Int64(fmt.Sprintf(load_obj_sqlarrayname, classType, i))))
			}
		}
	case "[]struct":
		for i := 0;  i < classVal.Len(); i++{
			parseLoadObjSql(classVal.Index(i).Interface(), row)
		}
	case "[*]float64":
		for i := 0; i < classVal.Len(); i++{
			classVal.Index(i).SetFloat(row.Float64(fmt.Sprintf(load_obj_sqlarrayname, classType, i)))
		}
	case "[*]float32":
		for i := 0; i < classVal.Len(); i++{
			classVal.Index(i).SetFloat(row.Float64(fmt.Sprintf(load_obj_sqlarrayname, classType, i)))
		}
	case "[*]bool":
		for i := 0; i < classVal.Len(); i++{
			classVal.Index(i).SetBool(row.Bool(fmt.Sprintf(load_obj_sqlarrayname, classType, i)))
		}
	case "[*]int8":
		for i := 0; i < classVal.Len(); i++{
			classVal.Index(i).SetInt(row.Int64(fmt.Sprintf(load_obj_sqlarrayname, classType, i)))
		}
	case "[*]uint8":
		for i := 0; i < classVal.Len(); i++{
			classVal.Index(i).SetUint(uint64(row.Int64(fmt.Sprintf(load_obj_sqlarrayname, classType, i))))
		}
	case "[*]int16":
		for i := 0; i < classVal.Len(); i++{
			classVal.Index(i).SetInt(row.Int64(fmt.Sprintf(load_obj_sqlarrayname, classType, i)))
		}
	case "[*]uint16":
		for i := 0; i < classVal.Len(); i++{
			classVal.Index(i).SetUint(uint64(row.Int64(fmt.Sprintf(load_obj_sqlarrayname, classType, i))))
		}
	case "[*]int32":
		for i := 0; i < classVal.Len(); i++{
			classVal.Index(i).SetInt(row.Int64(fmt.Sprintf(load_obj_sqlarrayname, classType, i)))
		}
	case "[*]uint32":
		for i := 0; i < classVal.Len(); i++{
			classVal.Index(i).SetUint(uint64(row.Int64(fmt.Sprintf(load_obj_sqlarrayname, classType, i))))
		}
	case "[*]int64":
		for i := 0; i < classVal.Len(); i++{
			classVal.Index(i).SetInt(row.Int64(fmt.Sprintf(load_obj_sqlarrayname, classType, i)))
		}
	case "[*]uint64":
		for i := 0; i < classVal.Len(); i++{
			classVal.Index(i).SetUint(uint64(row.Int64(fmt.Sprintf(load_obj_sqlarrayname, classType, i))))
		}
	case "[*]string":
		for i := 0; i < classVal.Len(); i++{
			classVal.Index(i).SetString(row.String(fmt.Sprintf(load_obj_sqlarrayname, classType, i)))
		}
	case "[*]int":
		for i := 0; i < classVal.Len(); i++{
			classVal.Index(i).SetInt(row.Int64(fmt.Sprintf(load_obj_sqlarrayname, classType, i)))
		}
	case "[*]uint":
		for i := 0; i < classVal.Len(); i++{
			classVal.Index(i).SetUint(uint64(row.Int64(fmt.Sprintf(load_obj_sqlarrayname, classType, i))))
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



