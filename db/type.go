package db

import (
	"reflect"
)

func getSliceTypeString(paramType reflect.Type) string{
	sTypeName := reflect.SliceOf(paramType.Elem()).Elem().Kind().String()
	switch sTypeName{
	case "bool", "float64", "float32", "int8", "uint8", "int16", "uint16",
		"int32", "uint32", "int64", "uint64", "string", "int", "uint":
		return "[*]" + sTypeName
	}

	return "[*]struct"
}

func getArrayTypeString(paramType reflect.Type) string{
	sTypeName := reflect.SliceOf(paramType.Elem()).Elem().Kind().String()
	switch sTypeName{
	case "bool", "float64", "float32", "int8", "uint8", "int16", "uint16",
		"int32", "uint32", "int64", "uint64", "string", "int", "uint":
			return "[*]" + sTypeName
	}

	return "[*]struct"
}

func getTypeString(classField reflect.StructField, classVal reflect.Value) string{
	paramType := classField.Type
	sType := ""
	if paramType.Kind() == reflect.Ptr{
		sType = "*" + paramType.Elem().Kind().String()
	}else if paramType.Kind() == reflect.Slice{
		sType = getSliceTypeString(paramType)
	}else if paramType.Kind() == reflect.Array{
		sType = getArrayTypeString(paramType)
	} else{
		sType = classField.Type.Kind().String()
	}
	return sType
}