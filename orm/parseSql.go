package orm

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

func marshalBlob(value reflect.Value) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	enc := gob.NewEncoder(buf)
	err := enc.Encode(value.Interface())
	return buf.Bytes(), err
}

func unMarshalBlob(data []byte, value reflect.Value) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	return dec.Decode(value.Addr().Interface())
}

func parsesql(sqlData *SqlData, p *Properties, op *Op, val string) {
	switch op.sqlType {
	case SQLTYPE_INSERT:
		sqlData.Value += fmt.Sprintf("'%s',", val)
		sqlData.Name += fmt.Sprintf("`%s`,", p.Name)
	case SQLTYPE_DELETE:
		if p.IsPrimary() {
			sqlData.Key += fmt.Sprintf("`%s`='%s',", p.Name, val)
		}
	case SQLTYPE_UPDATE:
		if p.IsPrimary() {
			sqlData.Key += fmt.Sprintf("`%s`='%s',", p.Name, val)
		} else {
			sqlData.NameValue += fmt.Sprintf("`%s`='%s',", p.Name, val)
		}
	case SQLTYPE_LOAD:
		if p.IsPrimary() {
			sqlData.Key += fmt.Sprintf("`%s`='%s',", p.Name, val)
		}
		sqlData.Name += fmt.Sprintf("`%s`,", p.Name)
	case SQLTYPE_SAVE:
		sqlData.Value += fmt.Sprintf("'%s',", val)
		sqlData.Name += fmt.Sprintf("`%s`,", p.Name)
		if !p.IsPrimary() {
			sqlData.NameValue += fmt.Sprintf("`%s`='%s',", p.Name, val)
		}
	case SQLTYPE_WHERE:
		sqlData.Key += fmt.Sprintf("`%s`='%s',", p.Name, val)
	}
}

func parsesqlblob(sqlData *SqlData, p *Properties, op *Op, val []byte) {
	switch op.sqlType {
	case SQLTYPE_INSERT:
		sqlData.Value += fmt.Sprintf("'%s',", val)
		sqlData.Name += fmt.Sprintf("`%s`,", p.Name)
	case SQLTYPE_DELETE:
		break
	case SQLTYPE_UPDATE:
		if p.IsPrimary() {
			sqlData.Key += fmt.Sprintf("`%s`='%s',", p.Name, val)
		} else {
			sqlData.NameValue += fmt.Sprintf("`%s`='%s',", p.Name, val)
		}
	case SQLTYPE_LOAD:
		break
	case SQLTYPE_SAVE:
		sqlData.Value += fmt.Sprintf("'%s',", val)
		sqlData.Name += fmt.Sprintf("`%s`,", p.Name)
		if !p.IsPrimary() {
			sqlData.NameValue += fmt.Sprintf("`%s`='%s',", p.Name, val)
		}
	case SQLTYPE_WHERE:
		break
	}
}

func parsesqlarray(sqlData *SqlData, p *Properties, op *Op, classVal reflect.Value) bool {
	switch op.sqlType {
	case SQLTYPE_INSERT:
		for classVal.Kind() == reflect.Ptr {
			classVal = classVal.Elem()
		}
		data, err := marshalBlob(classVal)
		parsesqlblob(sqlData, p, op, data)
		return err == nil
	case SQLTYPE_DELETE:
		return true
	case SQLTYPE_UPDATE:
		for classVal.Kind() == reflect.Ptr {
			classVal = classVal.Elem()
		}
		data, err := marshalBlob(classVal)
		parsesqlblob(sqlData, p, op, data)
		return err == nil
	case SQLTYPE_LOAD:
		sqlData.Name += fmt.Sprintf("`%s`,", p.Name)
	case SQLTYPE_SAVE:
		for classVal.Kind() == reflect.Ptr {
			classVal = classVal.Elem()
		}
		data, err := marshalBlob(classVal)
		parsesqlblob(sqlData, p, op, data)
		return err == nil
	case SQLTYPE_WHERE:
		return true
	}
	return true
}

func parseSfSql(p *Properties, classField reflect.StructField, classVal reflect.Value, sqlData *SqlData, op *Op) bool {
	sType := p.SType
	switch op.sqlType {
	case SQLTYPE_INSERT:
		if p.IsJson() && !classVal.IsZero() {
			for classVal.Kind() == reflect.Ptr {
				classVal = classVal.Elem()
			}
			data, err := json.Marshal(classVal.Interface())
			parsesql(sqlData, p, op, string(data))
			return err == nil
		} else if p.IsBlob() && !classVal.IsZero() {
			for classVal.Kind() == reflect.Ptr {
				classVal = classVal.Elem()
			}
			data, err := marshalBlob(classVal)
			parsesqlblob(sqlData, p, op, data)
			return err == nil
		} else if p.IsIgnore() {
			return true
		} else if p.IsTable() {
			return true
		} else if !op.forceFlag && !p.IsPrimary() && !p.IsForce() && classVal.IsZero() {
			return true
		}
	case SQLTYPE_DELETE:
		//过略json
		if p.IsJson() {
			return true
		} else if p.IsBlob() {
			return true
		} else if p.IsIgnore() {
			return true
		} else if p.IsTable() {
			return true
		} else if !p.IsPrimary() {
			return true
		}
	case SQLTYPE_UPDATE:
		if p.IsJson() && !classVal.IsZero() {
			for classVal.Kind() == reflect.Ptr {
				classVal = classVal.Elem()
			}
			data, err := json.Marshal(classVal.Interface())
			parsesql(sqlData, p, op, string(data))
			return err == nil
		} else if p.IsBlob() && !classVal.IsZero() {
			for classVal.Kind() == reflect.Ptr {
				classVal = classVal.Elem()
			}
			data, err := marshalBlob(classVal)
			parsesqlblob(sqlData, p, op, data)
			return err == nil
		} else if p.IsIgnore() {
			return true
		} else if p.IsTable() {
			return true
		} else if !op.forceFlag && !p.IsPrimary() && !p.IsForce() && classVal.IsZero() {
			return true
		}
	case SQLTYPE_LOAD:
		if p.IsJson() {
			sqlData.Name += fmt.Sprintf("`%s`,", p.Name)
			return true
		} else if p.IsBlob() {
			sqlData.Name += fmt.Sprintf("`%s`,", p.Name)
			return true
		} else if p.IsIgnore() {
			return true
		} else if p.IsTable() {
			return true
		}
	case SQLTYPE_SAVE:
		if p.IsJson() && !classVal.IsZero() {
			for classVal.Kind() == reflect.Ptr {
				classVal = classVal.Elem()
			}
			data, err := json.Marshal(classVal.Interface())
			parsesql(sqlData, p, op, string(data))
			return err == nil
		} else if p.IsBlob() && !classVal.IsZero() {
			for classVal.Kind() == reflect.Ptr {
				classVal = classVal.Elem()
			}
			data, err := marshalBlob(classVal)
			parsesqlblob(sqlData, p, op, data)
			return err == nil
		} else if p.IsIgnore() {
			return true
		} else if p.IsTable() {
			return true
		} else if !op.forceFlag && !p.IsPrimary() && !p.IsForce() && classVal.IsZero() {
			return true
		}
	case SQLTYPE_WHERE:
		//过略json
		if p.IsJson() {
			return true
		} else if p.IsBlob() {
			return true
		} else if p.IsIgnore() {
			return true
		} else if p.IsTable() {
			return true
		} else if op.whereFlag && classVal.IsZero() {
			return true
		}
	}

	switch sType {
	case "*bool":
		value := bool(false)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*bool)
		}
		parsesql(sqlData, p, op, strconv.FormatBool(value))
	case "*string":
		value := string("")
		if !classVal.IsNil() {
			value = *classVal.Interface().(*string)
		}
		parsesql(sqlData, p, op, value)
	case "*float32":
		value := float32(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*float32)
		}
		parsesql(sqlData, p, op, strconv.FormatFloat(float64(value), 'f', -1, 32))
	case "*float64":
		value := float64(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*float64)
		}
		parsesql(sqlData, p, op, strconv.FormatFloat(value, 'f', -1, 64))
	case "*int":
		value := int(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*int)
		}
		parsesql(sqlData, p, op, strconv.FormatInt(int64(value), 10))
	case "*int8":
		value := int8(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*int8)
		}
		parsesql(sqlData, p, op, strconv.FormatInt(int64(value), 10))
	case "*int16":
		value := int16(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*int16)
		}
		parsesql(sqlData, p, op, strconv.FormatInt(int64(value), 10))
	case "*int32":
		value := int32(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*int32)
		}
		parsesql(sqlData, p, op, strconv.FormatInt(int64(value), 10))
	case "*int64":
		value := int64(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*int64)
		}
		if !p.IsDatetime() {
			parsesql(sqlData, p, op, strconv.FormatInt(int64(value), 10))
		} else {
			parsesql(sqlData, p, op, GetDBTimeString(int64(value)))
		}
	case "*uint":
		value := uint(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*uint)
		}
		parsesql(sqlData, p, op, strconv.FormatUint(uint64(value), 10))
	case "*uint8":
		value := uint8(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*uint8)
		}
		parsesql(sqlData, p, op, strconv.FormatUint(uint64(value), 10))
	case "*uint16":
		value := uint16(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*uint16)
		}
		parsesql(sqlData, p, op, strconv.FormatUint(uint64(value), 10))
	case "*uint32":
		value := uint32(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*uint32)
		}
		parsesql(sqlData, p, op, strconv.FormatUint(uint64(value), 10))
	case "*uint64":
		value := uint64(0)
		if !classVal.IsNil() {
			value = *classVal.Interface().(*uint64)
		}
		parsesql(sqlData, p, op, strconv.FormatUint(uint64(value), 10))
	case "*struct":
		if !classVal.IsNil() {
			value := classVal.Elem().Interface()
			parseStructSql(value, sqlData, op)
		}

	case "bool":
		parsesql(sqlData, p, op, strconv.FormatBool(classVal.Bool()))
	case "string":
		parsesql(sqlData, p, op, classVal.String())
	case "float32":
		parsesql(sqlData, p, op, strconv.FormatFloat(classVal.Float(), 'f', -1, 32))
	case "float64":
		parsesql(sqlData, p, op, strconv.FormatFloat(classVal.Float(), 'f', -1, 64))
	case "int":
		parsesql(sqlData, p, op, strconv.FormatInt(classVal.Int(), 10))
	case "int8":
		parsesql(sqlData, p, op, strconv.FormatInt(classVal.Int(), 10))
	case "int16":
		parsesql(sqlData, p, op, strconv.FormatInt(classVal.Int(), 10))
	case "int32":
		parsesql(sqlData, p, op, strconv.FormatInt(classVal.Int(), 10))
	case "int64":
		if !p.IsDatetime() {
			parsesql(sqlData, p, op, strconv.FormatInt(classVal.Int(), 10))
		} else {
			parsesql(sqlData, p, op, GetDBTimeString(classVal.Int()))
		}
	case "uint":
		parsesql(sqlData, p, op, strconv.FormatUint(classVal.Uint(), 10))
	case "uint8":
		parsesql(sqlData, p, op, strconv.FormatUint(classVal.Uint(), 10))
	case "uint16":
		parsesql(sqlData, p, op, strconv.FormatUint(classVal.Uint(), 10))
	case "uint32":
		parsesql(sqlData, p, op, strconv.FormatUint(classVal.Uint(), 10))
	case "uint64":
		parsesql(sqlData, p, op, strconv.FormatUint(classVal.Uint(), 10))
	case "struct":
		parseStructSql(classVal.Interface(), sqlData, op)

	case "[]bool", "[]string", "[]float32", "[]float64", "[]int", "[]int8", "[]int16",
		"[]int32", "[]int64", "[]uint", "[]uint8", "[]uint16", "[]uint32", "[]uint64":
		return parsesqlarray(sqlData, p, op, classVal)

	case "[]struct":
		for i := 0; i < classVal.Len(); i++ {
			parseStructSql(classVal.Index(i).Interface(), sqlData, op)
		}

	case "[*]bool", "[*]string", "[*]float32", "[*]float64", "[*]int", "[*]int8", "[*]int16",
		"[*]int32", "[*]int64", "[*]uint", "[*]uint8", "[*]uint16", "[*]uint32", "[*]uint64":
		return parsesqlarray(sqlData, p, op, classVal)

	case "[*]struct":
		for i := 0; i < classVal.Len(); i++ {
			parseStructSql(classVal.Index(i).Interface(), sqlData, op)
		}

	case "[m]bool", "[m]string", "[m]float32", "[m]float64", "[m]int", "[m]int8", "[m]int16",
		"[m]int32", "[m]int64", "[m]uint", "[m]uint8", "[m]uint16", "[m]uint32", "[m]uint64", "[m]struct":
		return parsesqlarray(sqlData, p, op, classVal)

	default:
		panic(fmt.Sprintf("parseSfSql table [%s] [%s] type[%s] not supported", sqlData.Table, classField.Name, classField.Type.String()))
		return false
		//}
	}
	return true
}

func parseStructSql(obj interface{}, sqlData *SqlData, op *Op) {
	classVal, classType, table := GetTableInfo(obj)
	for i := 0; i < classType.NumField(); i++ {
		if !classVal.Field(i).CanInterface() {
			continue
		}

		p := table.Columns[i]
		bRight := parseSfSql(p, classType.Field(i), classVal.Field(i), sqlData, op)
		if !bRight {
			errorStr := fmt.Sprintf("parseStructSql type not supported %s", classType.Name())
			panic(errorStr)
			return //丢弃这个包
		}
	}
}
