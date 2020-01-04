package db

import (
	"reflect"
	"strings"
)

func getSliceTypeString(sTypeName string) string{
	index := strings.Index(sTypeName, "]")
	if index != -1{
		sTypeName = sTypeName[index+1:]
	}

	if sTypeName == "bool" || sTypeName == "float64" || sTypeName == "float32" || sTypeName == "int8" ||
		sTypeName == "uint8" || sTypeName == "int16" || sTypeName == "uint16" || sTypeName == "int32" ||
		sTypeName == "uint32" || sTypeName == "int64" || sTypeName == "uint64" ||  sTypeName == "string"||
		sTypeName == "int" || sTypeName == "uint"{
		return "[]" + sTypeName
	}else{
		return "[]struct"
	}

	return sTypeName
}

func getArrayTypeString(sTypeName string) string{
	index := strings.Index(sTypeName, "]")
	if index != -1{
		sTypeName = sTypeName[index+1:]
	}

	if sTypeName == "bool" || sTypeName == "float64" || sTypeName == "float32" || sTypeName == "int8" ||
		sTypeName == "uint8" || sTypeName == "int16" || sTypeName == "uint16" || sTypeName == "int32" ||
		sTypeName == "uint32" || sTypeName == "int64" || sTypeName == "uint64" ||  sTypeName == "string"||
		sTypeName == "int"  || sTypeName == "uint"{
		return "[*]" + sTypeName
	}else{
		return "[*]struct"
	}

	return sTypeName
}

func getTypeString(classField reflect.StructField, classVal reflect.Value) string{
	paramType := classField.Type
	sType := ""
	if paramType.Kind() == reflect.Ptr{
		sType = "*" + paramType.Elem().Kind().String()
	}else if paramType.Kind() == reflect.Slice{
		sType = getSliceTypeString(paramType.String())
	}else if paramType.Kind() == reflect.Array{
		sType = getArrayTypeString(paramType.String())
	} else{
		sType = classField.Type.Kind().String()
	}
	return sType
}