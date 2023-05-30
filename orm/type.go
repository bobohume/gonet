package orm

import (
	"reflect"
)

func GetSliceTypeString(paramType reflect.Type) string {
	sTypeName := reflect.SliceOf(paramType.Elem()).Elem().Kind().String()
	switch sTypeName {
	case "bool", "float64", "float32", "int8", "uint8", "int16", "uint16",
		"int32", "uint32", "int64", "uint64", "string", "int", "uint":
		return "[]" + sTypeName
	}

	return "[]struct"
}

func GetArrayTypeString(paramType reflect.Type) string {
	sTypeName := reflect.SliceOf(paramType.Elem()).Elem().Kind().String()
	switch sTypeName {
	case "bool", "float64", "float32", "int8", "uint8", "int16", "uint16",
		"int32", "uint32", "int64", "uint64", "string", "int", "uint":
		return "[*]" + sTypeName
	}

	return "[*]struct"
}

func GetMapTypeString(paramType reflect.Type) string {
	for paramType.Elem().Kind() == reflect.Ptr {
		paramType = paramType.Elem()
	}
	sTypeName := (paramType.Elem()).Kind().String()
	switch sTypeName {
	case "bool", "float64", "float32", "int8", "uint8", "int16", "uint16",
		"int32", "uint32", "int64", "uint64", "string", "int", "uint":
		return "[m]" + sTypeName
	}

	return "[m]struct"
}

func GetTypeString(classField reflect.StructField) string {
	paramType := classField.Type
	sType := ""
	if paramType.Kind() == reflect.Ptr {
		sType = "*" + paramType.Elem().Kind().String()
	} else if paramType.Kind() == reflect.Slice {
		sType = GetSliceTypeString(paramType)
	} else if paramType.Kind() == reflect.Array {
		sType = GetArrayTypeString(paramType)
	} else if paramType.Kind() == reflect.Map {
		sType = GetMapTypeString(paramType)
	} else {
		sType = classField.Type.Kind().String()
	}
	return sType
}
