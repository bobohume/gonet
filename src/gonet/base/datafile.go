package base

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
)

//datatype
const(
	DType_none    = iota
	DType_String	= iota
	DType_Enum		= iota
	DType_S8		= iota
	DType_S16		= iota
	DType_S32		= iota
	DType_F32		= iota
	DType_F64		= iota
	DType_S64		= iota
)

type(
	RData struct{
		m_Type	int

		m_String	string
		m_Enum int
		m_S8	int8
		m_S16	int16
		m_S32 int
		m_F32 float32
		m_F64 float64
		m_S64 int64
	}

	CDataFile struct{
		RecordNum	int//记录数量
		ColumNum	int//列数量

		fstream		*BitStream
		readstep	int//控制读的总数量
		dataTypes   Vector
		currentColumnIndex int
	}

	IDateFile interface {
		ReadDataFile(string) bool
		GetData(*RData) bool
		ReadDataInit()
	}
)

func (this *CDataFile) ReadDataInit(){
	this.ColumNum = 0
	this.RecordNum = 0
	this.readstep = 0
	this.fstream = nil
}

func (this *CDataFile) ReadDataFile(fileName string) bool{
	this.dataTypes.Clear()
	this.currentColumnIndex = 0

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("[%s] open failed", fileName)
		return false
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil{
		return false
	}

	rd := bufio.NewReaderSize(file, int(fileInfo.Size()))
	buf, err := ioutil.ReadAll(rd)
	if err != nil{
		return false
	}
	this.fstream = NewBitStream(buf, len(buf))

	for {
		tchr := this.fstream.ReadInt(8)
		if tchr == '@'{//找到数据文件的开头
			tchr = this.fstream.ReadInt(8)//这个是换行字符
			//fmt.Println(tchr)
			break
		}
	}
	//得到记录总数
	this.RecordNum = this.fstream.ReadInt(32)
	//得到列的总数
	this.ColumNum = this.fstream.ReadInt(32)
	//sheet name
	this.fstream.ReadString()

	this.readstep = this.RecordNum * this.ColumNum
	for nColumnIndex := 0; nColumnIndex < this.ColumNum; nColumnIndex++{
		//col name
		this.fstream.ReadString()
		nDataType := this.fstream.ReadInt(8)
		this.dataTypes.Push_back(int(nDataType))
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
func (this *CDataFile) GetData(pData *RData) bool {
	if this.readstep == 0 || this.fstream == nil{
		return false
	}

	switch this.dataTypes.Get(this.currentColumnIndex).(int) {
	case DType_String:
		pData.m_String = this.fstream.ReadString()
		//fmt.Println(pData.String, nLen)
	case DType_S8:
		pData.m_S8 = int8(this.fstream.ReadInt(8))
	case DType_S16:
		pData.m_S16 = int16(this.fstream.ReadInt(16))
	case DType_S32:
		pData.m_S32 = this.fstream.ReadInt(32)
	case DType_Enum:
		pData.m_Enum = this.fstream.ReadInt(16)
	case DType_F32:
		pData.m_F32 = this.fstream.ReadFloat()
	case DType_F64:
		pData.m_F64 = this.fstream.ReadFloat64()
	case DType_S64:
		pData.m_S64 = this.fstream.ReadInt64(64)
	}

	pData.m_Type = this.dataTypes.Get(this.currentColumnIndex).(int)
	this.currentColumnIndex = (this.currentColumnIndex + 1) % this.ColumNum
	this.readstep--
	return true
}

/****************************
	RData funciton
****************************/
func (this *RData) String(dataname, datacol string) string{
	IFAssert(this.m_Type == DType_String,  fmt.Sprintf("read [%s] col[%s] error", dataname, datacol))
	return this.m_String
}

func (this *RData) Enum(dataname, datacol string) int{
	IFAssert(this.m_Type == DType_Enum,  fmt.Sprintf("read [%s] col[%s] error", dataname, datacol))
	return this.m_Enum
}

func (this *RData) Int8(dataname, datacol string) int8{
	IFAssert(this.m_Type == DType_S8,  fmt.Sprintf("read [%s] col[%s] error", dataname, datacol))
	return this.m_S8
}

func (this *RData) Int16(dataname, datacol string) int16{
	IFAssert(this.m_Type == DType_S16,  fmt.Sprintf("read [%s] col[%s] error", dataname, datacol))
	return this.m_S16
}

func (this *RData) Int(dataname, datacol string) int{
	IFAssert(this.m_Type == DType_S32,  fmt.Sprintf("read [%s] col[%s] error", dataname, datacol))
	return this.m_S32
}

func (this *RData) Float32(dataname, datacol string) float32{
	IFAssert(this.m_Type == DType_F32,  fmt.Sprintf("read [%s] col[%s] error", dataname, datacol))
	return this.m_F32
}

func (this *RData) Float64(dataname, datacol string) float64{
	IFAssert(this.m_Type == DType_F64,  fmt.Sprintf("read [%s] col[%s] error", dataname, datacol))
	return this.m_F64
}

func (this *RData) Int64(dataname, datacol string) int64{
	IFAssert(this.m_Type == DType_S64,  fmt.Sprintf("read [%s] col[%s] error", dataname, datacol))
	return this.m_S64
}

