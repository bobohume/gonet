package orm

import (
	"encoding/json"
	"fmt"
	"reflect"
	"unsafe"
)

const (
	load_obj_sqlarrayname = "%s%d"
)

func getLoadObjSql(p *Properties, classField reflect.StructField, classVal reflect.Value, row IRow) bool {
	if !classVal.CanSet() {
		return true
	}
	classType := p.Name

	sType := p.SType
	if p.IsJson() {
		return json.Unmarshal(row.Byte(classType), classVal.Addr().Interface()) == nil
	} else if p.IsBlob() {
		for classVal.Kind() == reflect.Ptr {
			classVal = classVal.Elem()
		}

		return unMarshalBlob(row.Byte(classType), classVal) == nil
	} else if p.IsIgnore() {
		return true
	} else if p.IsTable() {
		return true
	}

	switch sType {
	case "*bool":
		value := (**bool)(unsafe.Pointer(uintptr(unsafe.Pointer(classVal.Addr().Pointer()))))
		val1 := row.Bool(classType)
		*value = &val1
	case "*string":
		value := (**string)(unsafe.Pointer(uintptr(unsafe.Pointer(classVal.Addr().Pointer()))))
		val1 := row.String(classType)
		*value = &val1
	case "*float32":
		value := (**float32)(unsafe.Pointer(uintptr(unsafe.Pointer(classVal.Addr().Pointer()))))
		val1 := float32(row.Float32(classType))
		*value = &val1
	case "*float64":
		value := (**float64)(unsafe.Pointer(uintptr(unsafe.Pointer(classVal.Addr().Pointer()))))
		val1 := row.Float64(classType)
		*value = &val1
	case "*int":
		value := (**int)(unsafe.Pointer(uintptr(unsafe.Pointer(classVal.Addr().Pointer()))))
		val1 := int(row.Int(classType))
		*value = &val1
	case "*int8":
		value := (**int8)(unsafe.Pointer(uintptr(unsafe.Pointer(classVal.Addr().Pointer()))))
		val1 := int8(row.Int(classType))
		*value = &val1
	case "*int16":
		value := (**int16)(unsafe.Pointer(uintptr(unsafe.Pointer(classVal.Addr().Pointer()))))
		val1 := int16(row.Int(classType))
		*value = &val1
	case "*int32":
		value := (**int32)(unsafe.Pointer(uintptr(unsafe.Pointer(classVal.Addr().Pointer()))))
		val1 := int32(row.Int(classType))
		*value = &val1
	case "*int64":
		value := (**int64)(unsafe.Pointer(uintptr(unsafe.Pointer(classVal.Addr().Pointer()))))
		if !p.IsDatetime() {
			val1 := int64(row.Int64(classType))
			*value = &val1
		} else {
			val1 := int64(row.Time(classType))
			*value = &val1
		}
	case "*uint":
		value := (**uint)(unsafe.Pointer(uintptr(unsafe.Pointer(classVal.Addr().Pointer()))))
		val1 := uint(row.Int(classType))
		*value = &val1
	case "*uint8":
		value := (**uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(classVal.Addr().Pointer()))))
		val1 := uint8(row.Int(classType))
		*value = &val1
	case "*uint16":
		value := (**uint16)(unsafe.Pointer(uintptr(unsafe.Pointer(classVal.Addr().Pointer()))))
		val1 := uint16(row.Int(classType))
		*value = &val1
	case "*uint32":
		value := (**uint32)(unsafe.Pointer(uintptr(unsafe.Pointer(classVal.Addr().Pointer()))))
		val1 := uint32(row.Int(classType))
		*value = &val1
	case "*uint64":
		value := (**uint64)(unsafe.Pointer(uintptr(unsafe.Pointer(classVal.Addr().Pointer()))))
		val1 := uint64(row.Int64(classType))
		*value = &val1
	case "*struct":
		if !classVal.IsNil() {
			value := classVal.Elem().Interface()
			parseLoadObjSql(value, row)
		}

	case "bool":
		classVal.SetBool(row.Bool(classType))
	case "string":
		classVal.SetString(row.String(classType))
	case "float32":
		classVal.SetFloat(row.Float64(classType))
	case "float64":
		classVal.SetFloat(row.Float64(classType))
	case "int":
		classVal.SetInt(row.Int64(classType))
	case "int8":
		classVal.SetInt(row.Int64(classType))
	case "int16":
		classVal.SetInt(row.Int64(classType))
	case "int32":
		classVal.SetInt(row.Int64(classType))
	case "int64":
		if !p.IsDatetime() {
			classVal.SetInt(row.Int64(classType))
		} else {
			classVal.SetInt(row.Time(classType))
		}
	case "uint":
		classVal.SetUint(uint64(row.Int64(classType)))
	case "uint8":
		classVal.SetUint(uint64(row.Int64(classType)))
	case "uint16":
		classVal.SetUint(uint64(row.Int64(classType)))
	case "uint32":
		classVal.SetUint(uint64(row.Int64(classType)))
	case "uint64":
		classVal.SetUint(uint64(row.Int64(classType)))
	case "struct":
		parseLoadObjSql(classVal.Addr().Interface(), row)

	case "[]bool", "[]string", "[]float32", "[]float64", "[]int", "[]int8", "[]int16",
		"[]int32", "[]int64", "[]uint", "[]uint8", "[]uint16", "[]uint32", "[]uint64":
		return parseLoadarray(classVal, classType, row)

	case "[]struct":
		for i := 0; i < classVal.Len(); i++ {
			parseLoadObjSql(classVal.Index(i).Addr().Interface(), row)
		}

	case "[*]bool", "[*]string", "[*]float32", "[*]float64", "[*]int", "[*]int8", "[*]int16",
		"[*]int32", "[*]int64", "[*]uint", "[*]uint8", "[*]uint16", "[*]uint32", "[*]uint64":
		return parseLoadarray(classVal, classType, row)

	case "[*]struct":
		for i := 0; i < classVal.Len(); i++ {
			parseLoadObjSql(classVal.Addr().Interface(), row)
		}

	case "[m]bool", "[m]string", "[m]float32", "[m]float64", "[m]int", "[m]int8", "[m]int16",
		"[m]int32", "[m]int64", "[m]uint", "[m]uint8", "[m]uint16", "[m]uint32", "[m]uint64", "[m]struct":
		return parseLoadarray(classVal, classType, row)

	default:
		fmt.Println("getLoadObjSql type not supported", sType, classField.Type)
		panic("getLoadObjSql type not supported")
		return false
		//}
	}
	return true
}

func parseLoadarray(classVal reflect.Value, classType string, row IRow) bool {
	for classVal.Kind() == reflect.Ptr {
		classVal = classVal.Elem()
	}

	return unMarshalBlob(row.Byte(classType), classVal) == nil
}

func parseLoadObjSql(obj interface{}, row IRow) bool {
	classVal, classType, table := GetTableInfo(obj)
	for i := 0; i < classType.NumField(); i++ {
		if !classVal.Field(i).CanInterface() {
			continue
		}

		p := table.Columns[i]
		bRight := getLoadObjSql(p, classType.Field(i), classVal.Field(i), row)
		if !bRight {
			errorStr := fmt.Sprintf("parseLoadObjSql type not supported %s", classType.Name())
			panic(errorStr)
			return false //丢弃这个包
		}
	}
	return true
}

// --- struct to sql
func LoadObjSql(obj interface{}, row IRow) bool {
	if row == nil {
		return false
	}

	return parseLoadObjSql(obj, row)
}
