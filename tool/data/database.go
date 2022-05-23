package main

import (
	"bufio"
	"bytes"
	"fmt"
	"gonet/base"
	"gonet/base/vector"
	"io/ioutil"
	"os"
	"strings"

	"github.com/tealeg/xlsx"
)

const (
	COL_NAME = iota
	COL_CLIENT_NAME
	COL_VSTO
	COL_TYPE
	COL_MAX
)

const (
	ARRAY_SPLIT = "|"
)

type (
	IDataFile interface {
		OpenExcel(filename string)
		SaveExcel(filename string)
	}
)

func FILENAME(filename, sheetname, ext string) string {
	filenames := strings.Split(filename, ".")
	return filenames[0] + "_" + sheetname + ext
}

//excel第一行 中文名字
//excel第二行 客户端data下的列名
//excel第三行 插件值
//excel第四行 类型
func OpenExcel(filename string) {
	xlFile, err := xlsx.OpenFile(filename)
	if err != nil {
		fmt.Println("open [%s] error", filename)
		return
	}

	for _, sheet := range xlFile.Sheets {
		if strings.Contains(sheet.Name, "~") {
			continue
		}

		dataTypes := []int{}
		buf := make([]byte, 10*1024*1024)
		stream := base.NewBitStream(buf, 10*1024*1024)
		enumKVMap := make(map[int]map[string]int) //列 key val
		enumKMap := map[int][]string{}            //列名对应key
		//检查行列
		func() {
			if sheet.MaxRow != len(sheet.Rows) {
				fmt.Printf("data [%s] 行数不统一", filename)
				return
			}
			for i, row := range sheet.Rows {
				if sheet.MaxCol != len(row.Cells) {
					fmt.Printf("data [%s] 列数不统一,第 [%d] 行", filename, i)
					return
				}
			}
		}()

		for i, row := range sheet.Rows {
			for j, cell := range row.Cells {
				if i == COL_NAME { //excel第一列 中文名字
					stream.WriteString(cell.String())
					continue
				} else if i == COL_CLIENT_NAME { //客户端data下的列名
					stream.WriteString(cell.String())
					continue
				} else if i == COL_VSTO { //插件值
					stream.WriteString(cell.String())
					if cell.String() == "" {
						continue
					}

					enumNames := strings.Split(cell.String(), "\n")
					for _, v1 := range enumNames {
						enumKMap[j] = append(enumKMap[j], v1)
					}
					continue
				} else if i == COL_TYPE { //类型
					coltype := strings.TrimSpace(strings.ToLower(cell.String()))
					if coltype == "enum" {
						num := 0
						enumKVMap[j] = make(map[string]int)
						for _, v1 := range enumKMap[j] {
							slot := strings.Split(string(v1), "=")
							if len(slot) == 2 {
								num = base.Int(slot[1])
								v1 = slot[0]
							}

							enumKVMap[j][v1] = num
							num++
						}
					}
					//写入列名
					stream.WriteString(coltype)
					switch coltype {
					case "string":
						stream.WriteInt(base.DType_String, 8)
						dataTypes = append(dataTypes, base.DType_String)
					case "enum":
						stream.WriteInt(base.DType_Enum, 8)
						dataTypes = append(dataTypes, base.DType_Enum)
					case "int8":
						stream.WriteInt(base.DType_S8, 8)
						dataTypes = append(dataTypes, base.DType_S8)
					case "int16":
						stream.WriteInt(base.DType_S16, 8)
						dataTypes = append(dataTypes, base.DType_S16)
					case "int":
						stream.WriteInt(base.DType_S32, 8)
						dataTypes = append(dataTypes, base.DType_S32)
					case "float":
						stream.WriteInt(base.DType_F32, 8)
						dataTypes = append(dataTypes, base.DType_F32)
					case "float64":
						stream.WriteInt(base.DType_F64, 8)
						dataTypes = append(dataTypes, base.DType_F64)
					case "int64":
						stream.WriteInt(base.DType_S64, 8)
						dataTypes = append(dataTypes, base.DType_S64)

					case "[]string":
						stream.WriteInt(base.DType_StringArray, 8)
						dataTypes = append(dataTypes, base.DType_StringArray)
					case "[]int8":
						stream.WriteInt(base.DType_S8Array, 8)
						dataTypes = append(dataTypes, base.DType_S8Array)
					case "[]int16":
						stream.WriteInt(base.DType_S16Array, 8)
						dataTypes = append(dataTypes, base.DType_S16Array)
					case "[]int":
						stream.WriteInt(base.DType_S32Array, 8)
						dataTypes = append(dataTypes, base.DType_S32Array)
					case "[]float":
						stream.WriteInt(base.DType_F32Array, 8)
						dataTypes = append(dataTypes, base.DType_F32Array)
					case "[]float64":
						stream.WriteInt(base.DType_F64Array, 8)
						dataTypes = append(dataTypes, base.DType_F64Array)
					case "[]int64":
						stream.WriteInt(base.DType_S64Array, 8)
						dataTypes = append(dataTypes, base.DType_S64Array)
					default:
						fmt.Printf("data [%s] [%s] col[%d] type not support in[string, enum, int8, int16, int32, float32, float64, []string, []int8, []int16, []int32, []float32, []float64]", filename, coltype, j)
						return
					}
					continue
				}

				switch dataTypes[j] {
				case base.DType_String:
					stream.WriteString(cell.Value)
				case base.DType_Enum:
					val, bEx := enumKVMap[j][strings.ToLower(cell.Value)]
					if bEx {
						stream.WriteInt(val, 16)
					} else {
						stream.WriteInt(0, 16)
					}
				case base.DType_S8:
					stream.WriteInt(base.Int(cell.Value), 8)
				case base.DType_S16:
					stream.WriteInt(base.Int(cell.Value), 16)
				case base.DType_S32:
					stream.WriteInt(base.Int(cell.Value), 32)
				case base.DType_F32:
					stream.WriteFloat(base.Float32(cell.Value))
				case base.DType_F64:
					stream.WriteFloat64(base.Float64(cell.Value))
				case base.DType_S64:
					stream.WriteInt64(base.Int64(cell.Value), 64)

				case base.DType_StringArray:
					arr := splitArray(cell.Value)
					stream.WriteInt(len(arr), 8)
					for _, v := range arr {
						stream.WriteString(v)
					}
				case base.DType_S8Array:
					arr := splitArray(cell.Value)
					stream.WriteInt(len(arr), 8)
					for _, v := range arr {
						stream.WriteInt(base.Int(v), 8)
					}
				case base.DType_S16Array:
					arr := splitArray(cell.Value)
					stream.WriteInt(len(arr), 8)
					for _, v := range arr {
						stream.WriteInt(base.Int(v), 16)
					}
				case base.DType_S32Array:
					arr := splitArray(cell.Value)
					stream.WriteInt(len(arr), 8)
					for _, v := range arr {
						stream.WriteInt(base.Int(v), 32)
					}
				case base.DType_F32Array:
					arr := splitArray(cell.Value)
					stream.WriteInt(len(arr), 8)
					for _, v := range arr {
						stream.WriteFloat(base.Float32(v))
					}
				case base.DType_F64Array:
					arr := splitArray(cell.Value)
					stream.WriteInt(len(arr), 8)
					for _, v := range arr {
						stream.WriteFloat64(base.Float64(v))
					}
				case base.DType_S64Array:
					arr := splitArray(cell.Value)
					stream.WriteInt(len(arr), 8)
					for _, v := range arr {
						stream.WriteInt64(base.Int64(v), 64)
					}
				}
			}

			//头结束
			//前三行都写在头部
			if i == COL_VSTO {
				for i1 := 0; i1 < 8-((COL_VSTO+1)*sheet.MaxCol%8); i1++ {
					stream.WriteFlag(true)
				}
				stream.WriteBits([]byte(base.DATA_END), base.DATA_END_LENGTH<<3)
				stream.WriteInt(sheet.MaxRow-COL_MAX, 32)
				stream.WriteInt(sheet.MaxCol, 32)
				stream.WriteString(sheet.Name)
			}
		}
		//other sheet
		stream.WriteInt(0, 32)
		file, err := os.Create(FILENAME(filename, sheet.Name, ".dat"))
		if err == nil {
			file.Write(stream.GetBuffer())
			file.Close()
		}
	}
}

//excel第一列 中文名字
//excel第二列 客户端data下的列名
//excel第三行 插件值
//excel第四行 类型
func SaveExcel(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("[%s] open failed", filename)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return
	}

	rd := bufio.NewReaderSize(file, int(fileInfo.Size()))
	buf, err := ioutil.ReadAll(rd)
	if err != nil {
		return
	}
	fstream := base.NewBitStream(buf, len(buf)+10)
	hstream := base.NewBitStream(buf, len(buf)+10)
	enumKVMap := make(map[int]map[int]string)
	enumKMap := map[int][]string{} //列名对应key
	colNames := []string{}
	nLen := bytes.Index(buf, []byte(base.DATA_END))
	if nLen == -1 {
		return
	}
	fstream.SetPosition(nLen + base.DATA_END_LENGTH)
	//得到记录总数
	RecordNum := fstream.ReadInt(32)
	//得到列的总数
	ColumNum := fstream.ReadInt(32)
	Sheetname := fstream.ReadString()
	//readstep := RecordNum * ColumNum
	dataTypes := vector.NewVector()
	xfile := xlsx.NewFile()
	sheet, err := xfile.AddSheet("~" + Sheetname)
	if err != nil {
		return
	}

	for i := 0; i <= COL_VSTO; i++ {
		row := sheet.AddRow()
		for j := 0; j < ColumNum; j++ {
			cell := row.AddCell()
			val := hstream.ReadString()
			cell.SetString(val)
			if i == COL_NAME { //excel第一列 中文名字
				colNames = append(colNames, val)
				continue
			} else if i == COL_VSTO {
				if val == "" {
					continue
				}
				enumNames := strings.Split(val, "\n")
				for _, v1 := range enumNames {
					enumKMap[j] = append(enumKMap[j], v1)
				}
			}
		}
	}

	//type
	{
		row := sheet.AddRow()
		for nColumnIndex := 0; nColumnIndex < ColumNum; nColumnIndex++ {
			typeName := fstream.ReadString()
			cell := row.AddCell()
			cell.SetString(typeName)
			coltype := strings.TrimSpace(strings.ToLower(typeName))
			if coltype == "enum" {
				num := 0
				enumKVMap[nColumnIndex] = make(map[int]string)
				for _, v1 := range enumKMap[nColumnIndex] {
					slot := strings.Split(string(v1), "=")
					if len(slot) == 2 {
						num = base.Int(slot[1])
						v1 = slot[0]
					}

					enumKVMap[nColumnIndex][num] = v1
					num++
				}
			}
			nDataType := fstream.ReadInt(8)
			dataTypes.PushBack(int(nDataType))
		}
	}

	//content
	for i := 0; i < RecordNum; i++ {
		row := sheet.AddRow()
		for j := 0; j < ColumNum; j++ {
			cell := row.AddCell()
			switch dataTypes.Get(j).(int) {
			case base.DType_String:
				cell.SetString(fstream.ReadString())
			case base.DType_S8:
				cell.SetInt(fstream.ReadInt(8))
			case base.DType_S16:
				cell.SetInt(fstream.ReadInt(16))
			case base.DType_S32:
				cell.SetInt(fstream.ReadInt(32))
			case base.DType_Enum:
				val, bEx := enumKVMap[j][fstream.ReadInt(16)]
				if bEx {
					cell.SetString(val)
				} else {
					cell.SetString("")
				}
			case base.DType_F32:
				cell.SetFloat(float64(fstream.ReadFloat()))
			case base.DType_F64:
				cell.SetFloat(fstream.ReadFloat64())
			case base.DType_S64:
				cell.SetInt64(fstream.ReadInt64(64))

			case base.DType_StringArray:
				nLen := fstream.ReadInt(8)
				str := ""
				for i := 0; i < nLen; i++ {
					str += fstream.ReadString()
					if i != nLen-1 {
						str += "|"
					}
				}
				cell.SetString(str)
			case base.DType_S8Array:
				nLen := fstream.ReadInt(8)
				str := ""
				for i := 0; i < nLen; i++ {
					str += fmt.Sprintf("%d", fstream.ReadInt(8))
					if i != nLen-1 {
						str += "|"
					}
				}
				cell.SetString(str)
			case base.DType_S16Array:
				nLen := fstream.ReadInt(8)
				str := ""
				for i := 0; i < nLen; i++ {
					str += fmt.Sprintf("%d", fstream.ReadInt(16))
					if i != nLen-1 {
						str += "|"
					}
				}
				cell.SetString(str)
			case base.DType_S32Array:
				nLen := fstream.ReadInt(8)
				str := ""
				for i := 0; i < nLen; i++ {
					str += fmt.Sprintf("%d", fstream.ReadInt(32))
					if i != nLen-1 {
						str += "|"
					}
				}
				cell.SetString(str)
			case base.DType_F32Array:
				nLen := fstream.ReadInt(8)
				str := ""
				for i := 0; i < nLen; i++ {
					str += fmt.Sprintf("%f", fstream.ReadFloat())
					if i != nLen-1 {
						str += "|"
					}
				}
				cell.SetString(str)
			case base.DType_F64Array:
				nLen := fstream.ReadInt(8)
				str := ""
				for i := 0; i < nLen; i++ {
					str += fmt.Sprintf("%f", fstream.ReadFloat64())
					if i != nLen-1 {
						str += "|"
					}
				}
				cell.SetString(str)
			case base.DType_S64Array:
				nLen := fstream.ReadInt(8)
				str := ""
				for i := 0; i < nLen; i++ {
					str += fmt.Sprintf("%d", fstream.ReadInt64(64))
					if i != nLen-1 {
						str += "|"
					}
				}
				cell.SetString(str)
			}
		}
	}

	for fstream.ReadFlag() {
		//得到记录总数
		recordNum := fstream.ReadInt(32)
		//得到列的总数
		columNum := fstream.ReadInt(32)
		sheetname := fstream.ReadString()
		sheet, err := xfile.AddSheet(sheetname)
		if err != nil {
			continue
		}
		//name
		for i := 0; i < recordNum; i++ {
			row := sheet.AddRow()
			for j := 0; j < columNum; j++ {
				cell := row.AddCell()
				cell.SetString(fstream.ReadString())
			}
		}
	}
	xfile.Save(FILENAME(filename, "", "temp.xlsx"))
	return
}

func splitArray(val string) []string {
	if val == "" || val == "0" {
		return []string{}
	}
	return strings.Split(val, ARRAY_SPLIT)
}
