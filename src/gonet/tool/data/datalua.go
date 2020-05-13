package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/tealeg/xlsx"
	"gonet/base"
	"os"
	"strings"
)

func OpenExceLua(filename string){
	xlFile, err := xlsx.OpenFile(filename)
	if err != nil{
		fmt.Println("open [%s] error", filename)
		return
	}

	dataTypes := []int{}
	dataColLen := 0//结束列数
	stream := bytes.NewBuffer([]byte{})
	filenames := strings.Split(filename, ".")
	enumKVMap := make(map[int] map[string] int) //列 key val
	enumKMap := map[string] []string{}//列名对应key
	enumNames := []string{}//列名
	dataNames := []string{}
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
			continue
			/*for _, row := range sheet.Rows {
				//列不统一
				for j := 0; j < sheet.MaxCol; j ++{
					if j < len(row.Cells){
						stream.WriteString(row.Cells[j].Value)
					}else{
						stream.WriteString("")
					}
				}
			}
			continue*/
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
				if i == COL_NAME {
					colNames = append(colNames, cell.String())
					continue
				} else if i == COL_CLIENT_NAME{
					colName := cell.String()
					dataNames = append(dataNames, colName)
					if j == 0{
						stream.WriteString(fmt.Sprintf("%s %s%s","local", filenames[0],"Data = {\n" ))
					}
					if colName != "" && colName != "0"{
						dataColLen = j
					}
					continue
				}else if i == COL_TYPE{
					coltype := strings.ToLower(cell.String())
					rd :=  bufio.NewReader(strings.NewReader(coltype))
					data, _, _ := rd.ReadLine()
					coltype = strings.TrimSpace(string(data))
					if coltype == "enum"{
						num := 0
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
								num++
							}
						}
					}
					if coltype == "string"{
						dataTypes = append(dataTypes, base.DType_String)
					}else if coltype == "enum"{
						dataTypes = append(dataTypes, base.DType_Enum)
					}else if coltype == "int8"{
						dataTypes = append(dataTypes, base.DType_S8)
					}else if coltype == "int16"{
						dataTypes = append(dataTypes, base.DType_S16)
					}else if coltype == "int"{
						dataTypes = append(dataTypes, base.DType_S32)
					} else if coltype == "float"{
						dataTypes = append(dataTypes, base.DType_F32)
					}else if coltype == "float64"{
						dataTypes = append(dataTypes, base.DType_F64)
					}else if coltype == "int64"{
						dataTypes = append(dataTypes, base.DType_S64)
					}else{
						fmt.Printf("data [%s] [%s] col[%d] type not support in[string, enum, int8, int16, int32, float32, float64]", filename, coltype, j )
						return
					}
					continue
				}

				if j == 0{
					stream.WriteString(fmt.Sprintf("\t[%s] = {\n", cell.Value))
				}

				writeInt := func() {
					switch cell.Type() {
					case xlsx.CellTypeString:
						stream.WriteString(fmt.Sprintf("\t\t%s = %d,\n",dataNames[j], base.Int(cell.String())))
					case xlsx.CellTypeStringFormula:
						stream.WriteString(fmt.Sprintf("\t\t%s = %d,\n",dataNames[j], base.Int(cell.String())))
					case xlsx.CellTypeNumeric:
						stream.WriteString(fmt.Sprintf("\t\t%s = %d,\n",dataNames[j], base.Int(cell.Value)))
					case xlsx.CellTypeBool:
						bVal := base.Bool(cell.Value)
						if bVal{
							stream.WriteString(fmt.Sprintf("\t\t%s = %d,\n",dataNames[j], 1))
						}else{
							stream.WriteString(fmt.Sprintf("\t\t%s = %d,\n",dataNames[j], 0))
						}
					case xlsx.CellTypeDate:
						stream.WriteString(fmt.Sprintf("\t\t%s = %d,\n",dataNames[j], base.Int(cell.Value)))
					}
				}

				//过滤掉不是客户端的数据
				if dataNames[j] == "" || dataNames[j] == "0"{
					continue
				}

				if dataTypes[j] == base.DType_String{
					switch cell.Type() {
					case xlsx.CellTypeString:
						stream.WriteString(fmt.Sprintf("\t\t%s = \"%s\",\n",dataNames[j], cell.String()))
					case xlsx.CellTypeStringFormula:
						stream.WriteString(fmt.Sprintf("\t\t%s = \"%s\",\n",dataNames[j], cell.String()))
					case xlsx.CellTypeNumeric:
						stream.WriteString(fmt.Sprintf("\t\t%s = \"%s\",\n",dataNames[j], cell.Value))
					case xlsx.CellTypeBool:
						bVal := base.Bool(cell.Value)
						if bVal{
							stream.WriteString(fmt.Sprintf("\t\t%s = \"%s\",\n",dataNames[j], "true"))
						}else{
							stream.WriteString(fmt.Sprintf("\t\t%s = \"%s\",\n",dataNames[j], "false"))
						}
					case xlsx.CellTypeDate:
						stream.WriteString(fmt.Sprintf("\t\t%s = \"%s\",\n",dataNames[j], cell.Value))
					}
				}else if dataTypes[j] == base.DType_Enum{
					val, bEx := enumKVMap[j][strings.ToLower(cell.Value)]
					if bEx{
						stream.WriteString(fmt.Sprintf("\t\t%s = %d,\n",dataNames[j], val))
					}else{
						stream.WriteString(fmt.Sprintf("\t\t%s = %d,\n",dataNames[j], 0))
					}
				}else if dataTypes[j] == base.DType_S8{
					writeInt()
				}else if dataTypes[j] == base.DType_S16{
					writeInt()
				}else if dataTypes[j] == base.DType_S32{
					writeInt()
				}else if dataTypes[j] == base.DType_F32{
					switch cell.Type() {
					case xlsx.CellTypeString:
						stream.WriteString(fmt.Sprintf("\t\t%s = %f,\n",dataNames[j], base.Float32(cell.String())))
					case xlsx.CellTypeStringFormula:
						stream.WriteString(fmt.Sprintf("\t\t%s = %f,\n",dataNames[j], base.Float32(cell.String())))
					case xlsx.CellTypeNumeric:
						stream.WriteString(fmt.Sprintf("\t\t%s = %f,\n",dataNames[j], base.Float32(cell.Value)))
					case xlsx.CellTypeBool:
						bVal := base.Bool(cell.Value)
						if bVal{
							stream.WriteString(fmt.Sprintf("\t\t%s = %f,\n",dataNames[j], 1))
						}else{
							stream.WriteString(fmt.Sprintf("\t\t%s = %f,\n",dataNames[j], 0))
						}
					case xlsx.CellTypeDate:
						stream.WriteString(fmt.Sprintf("\t\t%s = %f,\n",dataNames[j], base.Float32(cell.Value)))
					}
				}else if dataTypes[j] == base.DType_F64{
					switch cell.Type() {
					case xlsx.CellTypeString:
						stream.WriteString(fmt.Sprintf("\t\t%s = %f,\n",dataNames[j], base.Float64(cell.String())))
					case xlsx.CellTypeStringFormula:
						stream.WriteString(fmt.Sprintf("\t\t%s = %f,\n",dataNames[j], base.Float64(cell.String())))
					case xlsx.CellTypeNumeric:
						stream.WriteString(fmt.Sprintf("\t\t%s = %f,\n",dataNames[j], base.Float64(cell.Value)))
					case xlsx.CellTypeBool:
						bVal := base.Bool(cell.Value)
						if bVal{
							stream.WriteString(fmt.Sprintf("\t\t%s = %f,\n",dataNames[j], 1))
						}else{
							stream.WriteString(fmt.Sprintf("\t\t%s = %f,\n",dataNames[j], 0))
						}
					case xlsx.CellTypeDate:
						stream.WriteString(fmt.Sprintf("\t\t%s = %f,\n",dataNames[j], base.Float64(cell.Value)))
					}
				}else if dataTypes[j] == base.DType_S64{
					writeInt()
				}

				if j == dataColLen{
					stream.WriteString("\t},\n")
				}
			}
		}

		stream.WriteString("}\n")
		stream.WriteString(fmt.Sprintf("%s %s%s\n","return", filenames[0],"Data" ))
	}

	//文件没有可导出
	if dataColLen == 0{
		return
	}

	//other sheet
	file, err := os.Create(filenames[0] + ".lua")
	if err == nil{
		file.Write(stream.Bytes())
		file.Close()
	}
}

