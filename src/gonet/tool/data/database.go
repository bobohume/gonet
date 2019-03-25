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

type(
	IDataFile interface {
		OpenExcel(filename string)
		SaveExcel(filename string)
	}
)

func OpenExcel(filename string){
	xlFile, err := xlsx.OpenFile(filename)
	if err != nil{
		fmt.Println("open [%s] error", filename)
		return
	}

	dataTypes := []int{}
	buf := make([]byte,  10 * 1024 * 1024)
	stream := base.NewBitStream(buf, 10 * 1024 * 1024)
	enmnMap := make(map[int] map[string] int)
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
		for i, row := range sheet.Rows {
			for j, cell := range row.Cells {
				if i == 0 {
					stream.WriteString(cell.String())
					continue
				}else if i == 1{
					coltype := strings.ToLower(cell.String())
					rd :=  bufio.NewReader(strings.NewReader(coltype))
					data, _, _ := rd.ReadLine()
					coltype = strings.TrimSpace(string(data))
					if coltype == "enum"{
						for data, _, _ := rd.ReadLine(); data != nil;{
							slot := strings.Split(string(data), " ")
							if len(slot) == 2{
								_, bEx := enmnMap[j]
								if !bEx{
									enmnMap[j] = make(map[string] int)
								}
								enmnMap[j][slot[0]] = base.Int(slot[1])
							}
							data, _, _ = rd.ReadLine()
						}
					}
					stream.WriteString(cell.String())
					if coltype == "string"{
						stream.WriteInt(base.DType_String, 8)
						dataTypes = append(dataTypes, base.DType_String)
					}else if coltype == "enum"{
						stream.WriteInt(base.DType_Enum, 8)
						dataTypes = append(dataTypes, base.DType_Enum)
					}else if coltype == "int8"{
						stream.WriteInt(base.DType_S8, 8)
						dataTypes = append(dataTypes, base.DType_S8)
					}else if coltype == "int16"{
						stream.WriteInt(base.DType_S16, 8)
						dataTypes = append(dataTypes, base.DType_S16)
					}else if coltype == "int"{
						stream.WriteInt(base.DType_S32, 8)
						dataTypes = append(dataTypes, base.DType_S32)
					} else if coltype == "float"{
						stream.WriteInt(base.DType_F32, 8)
						dataTypes = append(dataTypes, base.DType_F32)
					}else if coltype == "float64"{
						stream.WriteInt(base.DType_F64, 8)
						dataTypes = append(dataTypes, base.DType_F64)
					}else if coltype == "int64"{
						stream.WriteInt(base.DType_S64, 8)
						dataTypes = append(dataTypes, base.DType_S64)
					}else{
						fmt.Errorf("[%s] col[%d] type not support in[string, enum, int8, int16, int32, float32, float64]", coltype, j )
						return
					}
					continue
				}

				writeInt := func(bitnum int) {
					switch cell.Type() {
					case xlsx.CellTypeString:
						stream.WriteInt(base.Int(cell.String()), bitnum)
					case xlsx.CellTypeFormula:
						stream.WriteInt(base.Int(cell.String()), bitnum)
					case xlsx.CellTypeNumeric:
						stream.WriteInt(base.Int(cell.Value), bitnum)
					case xlsx.CellTypeBool:
						bVal := base.Bool(cell.Value)
						if bVal{
							stream.WriteInt(1, bitnum)
						}else{
							stream.WriteInt(0, bitnum)
						}
					case xlsx.CellTypeDate:
						stream.WriteInt(base.Int(cell.Value), bitnum)
					}
				}


				if dataTypes[j] == base.DType_String{
					switch cell.Type() {
					case xlsx.CellTypeString:
						stream.WriteString(cell.String())
					case xlsx.CellTypeFormula:
						stream.WriteString(cell.String())
					case xlsx.CellTypeNumeric:
						stream.WriteString(fmt.Sprintf("%d", base.Int64(cell.Value)))
					case xlsx.CellTypeBool:
						stream.WriteString(fmt.Sprintf("%v", cell.Bool()))
					case xlsx.CellTypeDate:
						stream.WriteString(cell.Value)
					}
				}else if dataTypes[j] == base.DType_Enum{
					val, bEx := enmnMap[j][cell.Value]
					if bEx{
						stream.WriteInt(val, 16)
					}else{
						stream.WriteInt(0, 16)
					}
				}else if dataTypes[j] == base.DType_S8{
					writeInt(8)
				}else if dataTypes[j] == base.DType_S16{
					writeInt(16)
				}else if dataTypes[j] == base.DType_S32{
					writeInt(32)
				}else if dataTypes[j] == base.DType_F32{
					switch cell.Type() {
					case xlsx.CellTypeString:
						stream.WriteFloat(base.Float32(cell.String()))
					case xlsx.CellTypeFormula:
						stream.WriteFloat(base.Float32(cell.String()))
					case xlsx.CellTypeNumeric:
						stream.WriteFloat(base.Float32(cell.String()))
					case xlsx.CellTypeBool:
						bVal := base.Bool(cell.Value)
						if bVal{
							stream.WriteFloat(1)
						}else{
							stream.WriteFloat(0)
						}
					case xlsx.CellTypeDate:
						stream.WriteFloat(base.Float32(cell.Value))
					}
				}else if dataTypes[j] == base.DType_F64{
					switch cell.Type() {
					case xlsx.CellTypeString:
						stream.WriteFloat64(base.Float64(cell.String()))
					case xlsx.CellTypeFormula:
						stream.WriteFloat64(base.Float64(cell.String()))
					case xlsx.CellTypeNumeric:
						stream.WriteFloat64(base.Float64(cell.String()))
					case xlsx.CellTypeBool:
						bVal := base.Bool(cell.Value)
						if bVal{
							stream.WriteFloat64(1)
						}else{
							stream.WriteFloat64(0)
						}
					case xlsx.CellTypeDate:
						stream.WriteFloat64(base.Float64(cell.Value))
					}
				}else if dataTypes[j] == base.DType_S64{
					switch cell.Type() {
					case xlsx.CellTypeString:
						stream.WriteInt64(base.Int64(cell.String()), 64)
					case xlsx.CellTypeFormula:
						stream.WriteInt64(base.Int64(cell.String()), 64)
					case xlsx.CellTypeNumeric:
						stream.WriteInt64(base.Int64(cell.Value), 64)
					case xlsx.CellTypeBool:
						bVal := base.Bool(cell.Value)
						if bVal{
							stream.WriteInt64(1, 64)
						}else{
							stream.WriteInt64(0, 64)
						}
					case xlsx.CellTypeDate:
						stream.WriteInt64(base.Int64(cell.Value), 64)
					}
				}
			}

			//头结束
			if i == 0{
				for i1 := 0; i1 < 8 - (sheet.MaxCol % 8); i1++{
					stream.WriteFlag(true)
				}
				stream.WriteBits(16, []byte{'@', '\n'})
				stream.WriteInt(sheet.MaxRow - 2, 32)
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
	enmnMap := make(map[int] map[int] string)
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
	sheet, err :=xfile.AddSheet(Sheetname)
	if err != nil{
		return
	}
	//name
	{
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
					slot := strings.Split(string(data), " ")
					if len(slot) == 2 {
						_, bEx := enmnMap[nColumnIndex]
						if !bEx {
							enmnMap[nColumnIndex] = make(map[int]string)
						}
						enmnMap[nColumnIndex][base.Int(slot[1])] = slot[0]
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
				val, bEx := enmnMap[j][fstream.ReadInt(16)]
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
			}
		}
	}

	for fstream.ReadFlag(){
		//得到记录总数
		recordNum := fstream.ReadInt(32)
		//得到列的总数
		columNum := fstream.ReadInt(32)
		sheetname := fstream.ReadString()
		sheet, err :=xfile.AddSheet(sheetname)
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
