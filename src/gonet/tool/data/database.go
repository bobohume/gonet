package main

import (
	"gonet/base"
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"github.com/tealeg/xlsx"
)
const (
	COL_NAME = iota
	COL_CLIENT_NAME
	COL_TYPE
	COL_MAX
)

const(
	ARRAY_SPLIT = "|"
)

type(
	IDataFile interface {
		OpenExcel(filename string)
		SaveExcel(filename string)
	}
)

//excel第一行 中文名字
//excel第二行 客户端data下的列名
//excel第三行 类型
func OpenExcel(filename string){
	xlFile, err := xlsx.OpenFile(filename)
	if err != nil{
		fmt.Println("open [%s] error", filename)
		return
	}

	dataTypes := []int{}
	buf := make([]byte,  10 * 1024 * 1024)
	stream := base.NewBitStream(buf, 10 * 1024 * 1024)
	enumKVMap := make(map[int] map[string] int) //列 key val
	enumKMap := map[string] []string{}//列名对应key
	enumNames := []string{}//列名
	colNames := []string{}
	{
		sheet, bEx := xlFile.Sheet["Settings_Radio"]
		if bEx{
			for i, v := range sheet.Rows{
				for i1, v1 := range v.Cells{
					if v1.String() == ""{
						continue
					}
					if i == 0{
						enumNames = append(enumNames, v1.String())
					}else{
						enumKMap[enumNames[i1]] = append(enumKMap[enumNames[i1]], v1.String())
					}
				}
			}
		}
	}

	for page, sheet := range xlFile.Sheets{
		if page != 0{
			//other sheet
			stream.WriteFlag(true)
			stream.WriteInt(sheet.MaxRow, 32)
			stream.WriteInt(sheet.MaxCol, 32)
			stream.WriteString(sheet.Name)
			for _, row := range sheet.Rows {
				//列不统一
				for j := 0; j < sheet.MaxCol; j ++{
					if j < len(row.Cells){
						stream.WriteString(row.Cells[j].Value)
					}else{
						stream.WriteString("")
					}
				}
				/*for _, cell := range row.Cells {
					stream.WriteString(cell.Value)
				}*/
			}
			continue
		}

		//检查行列
		func(){
			if sheet.MaxRow != len(sheet.Rows){
				fmt.Printf("data [%s] 行数不统一", filename,  )
				return
			}
			for i, row := range sheet.Rows {
				if sheet.MaxCol != len(row.Cells){
					fmt.Printf("data [%s] 列数不统一,第 [%d] 行", filename,  i)
					return
				}
			}
		}()

		for i, row := range sheet.Rows {
			for j, cell := range row.Cells {
				colTypeName := cell.String()//在data解析到enum时候重新组装枚举到sheet列名
				if i == COL_NAME {//excel第一列 中文名字
					colNames = append(colNames, cell.String())
					stream.WriteString(cell.String())
					continue
				}else if i == COL_CLIENT_NAME {//客户端data下的列名
					stream.WriteString(cell.String())
					continue
				}else if i == COL_TYPE{//类型
					coltype := strings.ToLower(cell.String())
					rd :=  bufio.NewReader(strings.NewReader(coltype))
					data, _, _ := rd.ReadLine()
					coltype = strings.TrimSpace(string(data))
					if coltype == "enum"{
						num := 0
						colTypeName = "enum\n"
						KVMap := map[string] int{}
						for data, _, _ := rd.ReadLine(); data != nil;{
							slot := strings.Split(string(data), "=")
							if len(slot) == 2{
								KVMap[slot[0]] = base.Int(slot[1])
							}
							data, _, _ = rd.ReadLine()
						}
						keys, bEx := enumKMap[colNames[j]]
						if bEx{
							_, bEx := enumKVMap[j]
							if !bEx{
								enumKVMap[j] = make(map[string] int)
							}
							for _, v := range keys{
								val, bEx := KVMap[v]
								if bEx{
									num = val
								}
								enumKVMap[j][v] = num
								colTypeName += fmt.Sprintf("%s=%d\n", v, num)
								num++
							}
						}
						index := strings.LastIndex(colTypeName, ",")
						if index!= -1{
							colTypeName = colTypeName[:index]
						}
					}
					//写入列名
					stream.WriteString(colTypeName)
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
						fmt.Printf("data [%s] [%s] col[%d] type not support in[string, enum, int8, int16, int32, float32, float64, []string, []int8, []int16, []int32, []float32, []float64]", filename, coltype, j )
						return
					}
					continue
				}

				switch dataTypes[j] {
				case base.DType_String:
					stream.WriteString(cell.Value)
				case base.DType_Enum:
					val, bEx := enumKVMap[j][strings.ToLower(cell.Value)]
					if bEx{
						stream.WriteInt(val, 16)
					}else{
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
					for _, v := range arr{
						stream.WriteString(v)
					}
				case base.DType_S8Array:
					arr := splitArray(cell.Value)
					stream.WriteInt(len(arr), 8)
					for _, v := range arr{
						stream.WriteInt(base.Int(v), 8)
					}
				case base.DType_S16Array:
					arr := splitArray(cell.Value)
					stream.WriteInt(len(arr), 8)
					for _, v := range arr{
						stream.WriteInt(base.Int(v), 16)
					}
				case base.DType_S32Array:
					arr := splitArray(cell.Value)
					stream.WriteInt(len(arr), 8)
					for _, v := range arr{
						stream.WriteInt(base.Int(v), 32)
					}
				case base.DType_F32Array:
					arr := splitArray(cell.Value)
					stream.WriteInt(len(arr), 8)
					for _, v := range arr{
						stream.WriteFloat(base.Float32(v))
					}
				case base.DType_F64Array:
					arr := splitArray(cell.Value)
					stream.WriteInt(len(arr), 8)
					for _, v := range arr{
						stream.WriteFloat64(base.Float64(v))
					}
				case base.DType_S64Array:
					stream.WriteInt64(base.Int64(cell.Value), 64)
				}
			}

			//头结束
			//第一列和第二列都写在头部
			if i == COL_CLIENT_NAME{
				for i1 := 0; i1 < 8 - ((COL_CLIENT_NAME+1) * sheet.MaxCol % 8); i1++{
					stream.WriteFlag(true)
				}
				stream.WriteBits([]byte{'@', '\n'}, 16)
				stream.WriteInt(sheet.MaxRow - COL_MAX, 32)
				stream.WriteInt(sheet.MaxCol, 32)
				stream.WriteString(sheet.Name)
			}
		}
	}
	//other sheet
	filenames := strings.Split(filename, ".")
	stream.WriteInt(0, 32)
	file, err := os.Create(filenames[0] + ".dat")
	if err == nil{
		file.Write(stream.GetBuffer())
		file.Close()
	}
}

//excel第一列 中文名字
//excel第二列 客户端data下的列名
//excel第三列 类型
func SaveExcel(filename string){
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("[%s] open failed", filename)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil{
		return
	}

	rd := bufio.NewReaderSize(file, int(fileInfo.Size()))
	buf, err := ioutil.ReadAll(rd)
	if err != nil{
		return
	}
	fstream := base.NewBitStream(buf, len(buf) + 10)
	hstream := base.NewBitStream(buf, len(buf) + 10)
	enunKVMap := make(map[int] map[int] string)
	for {
		tchr := fstream.ReadInt(8)
		if tchr == '@'{//找到数据文件的开头
			tchr = fstream.ReadInt(8)//这个是换行字符
			//fmt.Println(tchr)
			break
		}
	}
	//得到记录总数
	RecordNum := fstream.ReadInt(32)
	//得到列的总数
	ColumNum := fstream.ReadInt(32)
	Sheetname := fstream.ReadString()
	//readstep := RecordNum * ColumNum
	dataTypes := base.NewVector()
	xfile := xlsx.NewFile()
	sheet, err := xfile.AddSheet(Sheetname)
	if err != nil{
		return
	}

	for i := 0; i <= COL_CLIENT_NAME; i++{
		row := sheet.AddRow()
		for j := 0; j < ColumNum; j++{
			cell := row.AddCell()
			cell.SetString(hstream.ReadString())
		}
	}

	//type
	{
		row := sheet.AddRow()
		for nColumnIndex := 0; nColumnIndex < ColumNum; nColumnIndex++ {
			typeName := fstream.ReadString()
			cell := row.AddCell()
			cell.SetString(typeName)
			coltype := strings.ToLower(typeName)
			rd := bufio.NewReader(strings.NewReader(coltype))
			data, _, _ := rd.ReadLine()
			coltype = strings.TrimSpace(string(data))
			if coltype == "enum" {
				for data, _, _ := rd.ReadLine(); data != nil; {
					slot := strings.Split(string(data), "=")
					if len(slot) == 2 {
						_, bEx := enunKVMap[nColumnIndex]
						if !bEx {
							enunKVMap[nColumnIndex] = make(map[int]string)
						}
						enunKVMap[nColumnIndex][base.Int(slot[1])] = slot[0]
					}
					data, _, _ = rd.ReadLine()
				}
			}
			nDataType := fstream.ReadInt(8)
			dataTypes.Push_back(int(nDataType))
		}
	}

	//content
	for i := 0; i < RecordNum; i++{
		row := sheet.AddRow()
		for j := 0; j < ColumNum; j++{
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
				val, bEx := enunKVMap[j][fstream.ReadInt(16)]
				if bEx{
					cell.SetString(val)
				}else{
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
				for i := 0; i < nLen; i++{
					str += fstream.ReadString()
					if i != nLen-1{
						str += "|"
					}
				}
				cell.SetString(str)
			case base.DType_S8Array:
				nLen := fstream.ReadInt(8)
				str := ""
				for i := 0; i < nLen; i++{
					str += fmt.Sprintf("%d", fstream.ReadInt(8))
					if i != nLen-1{
						str += "|"
					}
				}
				cell.SetString(str)
			case base.DType_S16Array:
				nLen := fstream.ReadInt(8)
				str := ""
				for i := 0; i < nLen; i++{
					str += fmt.Sprintf("%d", fstream.ReadInt(16))
					if i != nLen-1{
						str += "|"
					}
				}
				cell.SetString(str)
			case base.DType_S32Array:
				nLen := fstream.ReadInt(8)
				str := ""
				for i := 0; i < nLen; i++{
					str += fmt.Sprintf("%d", fstream.ReadInt(32))
					if i != nLen-1{
						str += "|"
					}
				}
				cell.SetString(str)
			case base.DType_F32Array:
				nLen := fstream.ReadInt(8)
				str := ""
				for i := 0; i < nLen; i++{
					str += fmt.Sprintf("%f", fstream.ReadFloat())
					if i != nLen-1{
						str += "|"
					}
				}
				cell.SetString(str)
			case base.DType_F64Array:
				nLen := fstream.ReadInt(8)
				str := ""
				for i := 0; i < nLen; i++{
					str += fmt.Sprintf("%f", fstream.ReadFloat64())
					if i != nLen-1{
						str += "|"
					}
				}
				cell.SetString(str)
			case base.DType_S64Array:
				nLen := fstream.ReadInt(8)
				str := ""
				for i := 0; i < nLen; i++{
					str += fmt.Sprintf("%d", fstream.ReadInt64(64))
					if i != nLen-1{
						str += "|"
					}
				}
				cell.SetString(str)
			}
		}
	}

	for fstream.ReadFlag(){
		//得到记录总数
		recordNum := fstream.ReadInt(32)
		//得到列的总数
		columNum := fstream.ReadInt(32)
		sheetname := fstream.ReadString()
		sheet, err := xfile.AddSheet(sheetname)
		if err != nil{
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
	filenames := strings.Split(filename, ".")
	xfile.Save( filenames[0]+ "_temp.xlsx")

	return
}

func splitArray(val string) []string{
	if val == "" || val == "0"{
		return []string{}
	}
	return strings.Split(val, ARRAY_SPLIT)
}
