package base

import (
	"bufio"
	"os"
	"fmt"
)

const(
	DType_none    = iota
	DType_String	= iota
	DType_Enum8	= iota
	DType_Enum16	= iota
	DType_S8		= iota
	DType_S16		= iota
	DType_S32		= iota
	DType_U8		= iota
	DType_U16		= iota
	DType_U32		= iota
	DType_F32		= iota
	DType_F64		= iota
)

type(
	RData struct{
		Type	int

		String	string
		S8	int8
		S16	int16
		S32 int32
		U8  uint8
		U16 uint16
		U32 uint32
		F32 float32
		F64 float64
		Enum8 uint8
		Enum16 uint16
	}

	CDataFile struct{
		RecordNum	int//记录数量
		ColumNum	int//列数量

		fdata	*bufio.Reader
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
	this.fdata = nil
}

func (this *CDataFile) ReadDataFile(fileName string) bool{
	this.dataTypes.Clear()
	this.currentColumnIndex = 0

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("[%s] open failed", fileName)
		return false
	}

	fileInfo, err := file.Stat()
	if err != nil{
		return false
	}

	defer file.Close()
	this.fdata = bufio.NewReaderSize(file, int(fileInfo.Size()))

	for {
		tchr, _ := this.fdata.ReadByte()
		if tchr == '@'{//找到数据文件的开头
			tchr, _ = this.fdata.ReadByte()//这个是换行字符
			//fmt.Println(tchr)
			break
		}
	}
	buf := make([]byte, 4)
	this.fdata.Read(buf)//得到记录总数
	this.RecordNum = BytesToInt(buf)
	this.fdata.Read(buf)//得到列的总数
	this.ColumNum = BytesToInt(buf)

	this.readstep = this.RecordNum * this.ColumNum
	for nColumnIndex := 0; nColumnIndex < this.ColumNum; nColumnIndex++{
		nDataType, _ := this.fdata.ReadByte()
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
	if this.readstep == 0 || this.fdata == nil{
		return false
	}

	var nByte byte
	switch this.dataTypes.Get(this.currentColumnIndex).(int) {
	case DType_String:
		buf := make([]byte, 2)
		this.fdata.Read(buf)
		nLen := BytesToInt16(buf)
		buf1 := make([]byte, nLen)
		this.fdata.Read(buf1)
		pData.String = string(buf1)
		//fmt.Println(pData.String, nLen)
	case DType_U8:
		pData.U8,_ = this.fdata.ReadByte()
	case DType_S8:
		nByte,_ = this.fdata.ReadByte()
		pData.S8 = int8(nByte)
	case DType_U16:
		buf := make([]byte, 2)
		this.fdata.Read(buf)
		pData.U16 = uint16(BytesToInt16(buf))
	case DType_S16:
		buf := make([]byte, 2)
		this.fdata.Read(buf)
		pData.S16 = int16(BytesToInt16(buf))
	case DType_U32:
		buf := make([]byte, 4)
		this.fdata.Read(buf)
		pData.U32 = uint32(BytesToInt(buf))
	case DType_S32:
		buf := make([]byte, 4)
		this.fdata.Read(buf)
		pData.S32 = int32(BytesToInt(buf))
	case DType_Enum8:
		pData.Enum8,_ = this.fdata.ReadByte()
	case DType_Enum16:
		buf := make([]byte, 2)
		this.fdata.Read(buf)
		pData.Enum16 = uint16(BytesToInt(buf))
	case DType_F32:
		buf := make([]byte, 4)
		this.fdata.Read(buf)
		pData.F32 = float32(ByteToFloat32(buf))
	case DType_F64:
		buf := make([]byte, 8)
		this.fdata.Read(buf)
		pData.F64 = float64(ByteToFloat64(buf))
	}

	pData.Type = this.dataTypes.Get(this.currentColumnIndex).(int)
	this.currentColumnIndex = (this.currentColumnIndex + 1) % this.ColumNum
	this.readstep--
	return true
}
