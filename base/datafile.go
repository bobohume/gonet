package base

import (
	"bytes"
	"fmt"
	"gonet/base/vector"
	"io/ioutil"
	"os"
	"reflect"
)

const (
	DATA_END        = "data_end"
	DATA_END_LENGTH = len(DATA_END) //data结束标记
)

//datatype
const (
	DType_none        = iota
	DType_String      = iota
	DType_Enum        = iota
	DType_S8          = iota
	DType_S16         = iota
	DType_S32         = iota
	DType_F32         = iota
	DType_F64         = iota
	DType_S64         = iota
	DType_StringArray = iota
	DType_S8Array     = iota
	DType_S16Array    = iota
	DType_S32Array    = iota
	DType_F32Array    = iota
	DType_F64Array    = iota
	DType_S64Array    = iota
)

type (
	Data struct {
		fileName string
		dataType int

		str      string
		enum     int
		s8       int8
		s16      int16
		s32      int
		f32      float32
		f64      float64
		s64      int64
		strArray []string
		s8Array  []int8
		s16Array []int16
		s32Array []int
		f32Array []float32
		f64Array []float64
		s64Array []int64
	}

	DataFile struct {
		RecordNum int //记录数量
		ColumNum  int //列数量

		fileName           string //文件名
		fstream            *BitStream
		readstep           int //控制读的总数量
		dataTypes          vector.Vector[int]
		currentColumnIndex int
	}

	IDateFile interface {
		ReadDataFile(string) bool
		GetData(*Data) bool
		ReadDataInit()
	}
)

func (d *DataFile) ReadDataInit() {
	d.ColumNum = 0
	d.RecordNum = 0
	d.readstep = 0
	d.fstream = nil
}

func (d *DataFile) ReadDataFile(fileName string) bool {
	d.dataTypes.Clear()
	d.currentColumnIndex = 0
	d.fileName = fileName

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("[%s] open failed", fileName)
		return false
	}
	defer file.Close()
	buf, err := ioutil.ReadAll(file)
	if err != nil {
		return false
	}
	d.fstream = NewBitStream(buf, len(buf))
	nLen := bytes.Index(buf, []byte(DATA_END))
	if nLen == -1 {
		return false
	}
	d.fstream.SetPosition(nLen + DATA_END_LENGTH)
	//得到记录总数
	d.RecordNum = d.fstream.ReadInt(32)
	//得到列的总数
	d.ColumNum = d.fstream.ReadInt(32)
	//sheet name
	d.fstream.ReadString()

	d.readstep = d.RecordNum * d.ColumNum
	for nColumnIndex := 0; nColumnIndex < d.ColumNum; nColumnIndex++ {
		//col name
		d.fstream.ReadString()
		nDataType := d.fstream.ReadInt(8)
		d.dataTypes.PushBack(nDataType)
	}
	return true
}

/****************************
	格式:
	头文件:
	1、总记录数(int)
	2、总字段数(int)
	字段格式:
	1、字段长度(int)
	2、字读数据类型(int->2,string->1,enum->3,float->4)
	3、字段内容(int,string)
*************************/
func (d *DataFile) GetData(data *Data) bool {
	if d.readstep == 0 || d.fstream == nil {
		return false
	}

	data.fileName = d.fileName
	switch d.dataTypes.Get(d.currentColumnIndex) {
	case DType_String:
		data.str = d.fstream.ReadString()
	case DType_S8:
		data.s8 = int8(d.fstream.ReadInt(8))
	case DType_S16:
		data.s16 = int16(d.fstream.ReadInt(16))
	case DType_S32:
		data.s32 = d.fstream.ReadInt(32)
	case DType_Enum:
		data.enum = d.fstream.ReadInt(16)
	case DType_F32:
		data.f32 = d.fstream.ReadFloat()
	case DType_F64:
		data.f64 = d.fstream.ReadFloat64()
	case DType_S64:
		data.s64 = d.fstream.ReadInt64(64)

	case DType_StringArray:
		nLen := d.fstream.ReadInt(8)
		data.strArray = make([]string, nLen)
		for i := 0; i < nLen; i++ {
			data.strArray[i] = d.fstream.ReadString()
		}
	case DType_S8Array:
		nLen := d.fstream.ReadInt(8)
		data.s8Array = make([]int8, nLen)
		for i := 0; i < nLen; i++ {
			data.s8Array[i] = int8(d.fstream.ReadInt(8))
		}
	case DType_S16Array:
		nLen := d.fstream.ReadInt(8)
		data.s16Array = make([]int16, nLen)
		for i := 0; i < nLen; i++ {
			data.s16Array[i] = int16(d.fstream.ReadInt(16))
		}
	case DType_S32Array:
		nLen := d.fstream.ReadInt(8)
		data.s32Array = make([]int, nLen)
		for i := 0; i < nLen; i++ {
			data.s32Array[i] = d.fstream.ReadInt(32)
		}
	case DType_F32Array:
		nLen := d.fstream.ReadInt(8)
		data.f32Array = make([]float32, nLen)
		for i := 0; i < nLen; i++ {
			data.f32Array[i] = d.fstream.ReadFloat()
		}
	case DType_F64Array:
		nLen := d.fstream.ReadInt(8)
		data.f64Array = make([]float64, nLen)
		for i := 0; i < nLen; i++ {
			data.f64Array[i] = d.fstream.ReadFloat64()
		}
	case DType_S64Array:
		nLen := d.fstream.ReadInt(8)
		data.s64Array = make([]int64, nLen)
		for i := 0; i < nLen; i++ {
			data.s64Array[i] = d.fstream.ReadInt64(64)
		}
	}

	data.dataType = d.dataTypes.Get(d.currentColumnIndex)
	d.currentColumnIndex = (d.currentColumnIndex + 1) % d.ColumNum
	d.readstep--
	return true
}

/****************************
	Data funciton
****************************/
func (d *Data) String(datacol string) string {
	IFAssert(d.dataType == DType_String, fmt.Sprintf("read [%s] col[%s] error", d.fileName, datacol))
	return d.str
}

func (d *Data) Enum(datacol string) int {
	IFAssert(d.dataType == DType_Enum, fmt.Sprintf("read [%s] col[%s] error", d.fileName, datacol))
	return d.enum
}

func (d *Data) Int8(datacol string) int8 {
	IFAssert(d.dataType == DType_S8, fmt.Sprintf("read [%s] col[%s] error", d.fileName, datacol))
	return d.s8
}

func (d *Data) Int16(datacol string) int16 {
	IFAssert(d.dataType == DType_S16, fmt.Sprintf("read [%s] col[%s] error", d.fileName, datacol))
	return d.s16
}

func (d *Data) Int(datacol string) int {
	IFAssert(d.dataType == DType_S32, fmt.Sprintf("read [%s] col[%s] error", d.fileName, datacol))
	return d.s32
}

func (d *Data) Float32(datacol string) float32 {
	IFAssert(d.dataType == DType_F32, fmt.Sprintf("read [%s] col[%s] error", d.fileName, datacol))
	return d.f32
}

func (d *Data) Float64(datacol string) float64 {
	IFAssert(d.dataType == DType_F64, fmt.Sprintf("read [%s] col[%s] error", d.fileName, datacol))
	return d.f64
}

func (d *Data) Int64(datacol string) int64 {
	IFAssert(d.dataType == DType_S64, fmt.Sprintf("read [%s] col[%s] error", d.fileName, datacol))
	return d.s64
}

func (d *Data) StringArray(datacol string) []string {
	IFAssert(d.dataType == DType_StringArray, fmt.Sprintf("read [%s] col[%s] error", d.fileName, datacol))
	return d.strArray
}

func (d *Data) Int8Array(datacol string) []int8 {
	IFAssert(d.dataType == DType_S8Array, fmt.Sprintf("read [%s] col[%s] error", d.fileName, datacol))
	return d.s8Array
}

func (d *Data) Int16Array(datacol string) []int16 {
	IFAssert(d.dataType == DType_S16Array, fmt.Sprintf("read [%s] col[%s] error", d.fileName, datacol))
	return d.s16Array
}

func (d *Data) IntArray(datacol string) []int {
	IFAssert(d.dataType == DType_S32Array, fmt.Sprintf("read [%s] col[%s] error", d.fileName, datacol))
	return d.s32Array
}

func (d *Data) Float32Array(datacol string) []float32 {
	IFAssert(d.dataType == DType_F32Array, fmt.Sprintf("read [%s] col[%s] error", d.fileName, datacol))
	return d.f32Array
}

func (d *Data) Float64Array(datacol string) []float64 {
	IFAssert(d.dataType == DType_F64Array, fmt.Sprintf("read [%s] col[%s] error", d.fileName, datacol))
	return d.f64Array
}

func (d *Data) Int64Array(datacol string) []int64 {
	IFAssert(d.dataType == DType_S64Array, fmt.Sprintf("read [%s] col[%s] error", d.fileName, datacol))
	return d.s64Array
}

//--- struct to data
func LoadData(obj interface{}, file *DataFile) bool {
	if file == nil {
		return false
	}

	classVal, classType := getClassInfo(obj)
	for i := 0; i < classType.NumField(); i++ {
		if !classVal.Field(i).CanInterface() {
			continue
		}

		getData(classType.Field(i), classVal.Field(i), file)
	}
	return true
}

func getClassInfo(obj interface{}) (reflect.Value, reflect.Type) {
	classVal := reflect.ValueOf(obj)
	for classVal.Kind() == reflect.Ptr {
		classVal = classVal.Elem()
	}
	classType := classVal.Type()
	return classVal, classType
}

func getData(sf reflect.StructField, classVal reflect.Value, file *DataFile) {
	data := &Data{}
	colName := sf.Name

	file.GetData(data)

	switch data.dataType {
	case DType_String:
		classVal.SetString(data.String(colName))
	case DType_Enum:
		classVal.SetInt(int64(data.Enum(colName)))
	case DType_S8:
		classVal.SetInt(int64(data.Int8(colName)))
	case DType_S16:
		classVal.SetInt(int64(data.Int16(colName)))
	case DType_S32:
		classVal.SetInt(int64(data.Int(colName)))
	case DType_F32:
		classVal.SetFloat(float64(data.Float32(colName)))
	case DType_F64:
		classVal.SetFloat(float64(data.Float64(colName)))
	case DType_S64:
		classVal.SetInt(int64(data.Int64(colName)))

	case DType_StringArray:
		classVal.Set(reflect.ValueOf(data.StringArray(colName)))
	case DType_S8Array:
		classVal.Set(reflect.ValueOf(data.Int8Array(colName)))
	case DType_S16Array:
		classVal.Set(reflect.ValueOf(data.Int16Array(colName)))
	case DType_S32Array:
		classVal.Set(reflect.ValueOf(data.IntArray(colName)))
	case DType_F32Array:
		classVal.Set(reflect.ValueOf(data.Float32Array(colName)))
	case DType_F64Array:
		classVal.Set(reflect.ValueOf(data.Float64Array(colName)))
	case DType_S64Array:
		classVal.Set(reflect.ValueOf(data.Int64Array(colName)))
	}
}
