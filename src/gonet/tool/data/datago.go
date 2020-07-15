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

func OpenExceGo(filename string){
	xlFile, err := xlsx.OpenFile(filename)
	if err != nil{
		fmt.Println("open [%s] error", filename)
		return
	}

	dataTypes := []int{}
	dataTypeNames := []string{}
	structName := ""
	structResName := ""
	dataColLen := 0//结束列数
	dataColBeginLen := -1//开始列数
	stream := bytes.NewBuffer([]byte{})
	filenames := strings.Split(filename, ".")
	structName = filenames[0] + "Data"
	structResName = filenames[0] + "DataRes"
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
					if colName != "" && colName != "0"{
						dataColLen = j
						if dataColBeginLen == -1{
							dataColBeginLen = j
						}
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

					switch coltype {
					case "string":
						dataTypes = append(dataTypes, base.DType_String)
						dataTypeNames = append(dataTypeNames, "string")
					case "enum":
						dataTypes = append(dataTypes, base.DType_Enum)
						dataTypeNames = append(dataTypeNames, "int")
					case "int8":
						dataTypes = append(dataTypes, base.DType_S8)
						dataTypeNames = append(dataTypeNames, "int8")
					case "int16":
						dataTypes = append(dataTypes, base.DType_S16)
						dataTypeNames = append(dataTypeNames, "int16")
					case "int":
						dataTypes = append(dataTypes, base.DType_S32)
						dataTypeNames = append(dataTypeNames, "int")
					case "float":
						dataTypes = append(dataTypes, base.DType_F32)
						dataTypeNames = append(dataTypeNames, "float")
					case "float64":
						dataTypes = append(dataTypes, base.DType_F64)
						dataTypeNames = append(dataTypeNames, "float64")
					case "int64":
						dataTypes = append(dataTypes, base.DType_S64)
						dataTypeNames = append(dataTypeNames, "int64")

					case "[]string":
						dataTypes = append(dataTypes, base.DType_StringArray)
						dataTypeNames = append(dataTypeNames, "[]string")
					case "[]int8":
						dataTypes = append(dataTypes, base.DType_S8Array)
						dataTypeNames = append(dataTypeNames, "[]int8")
					case "[]int16":
						dataTypes = append(dataTypes, base.DType_S16Array)
						dataTypeNames = append(dataTypeNames, "[]int16")
					case "[]int":
						dataTypes = append(dataTypes, base.DType_S32Array)
						dataTypeNames = append(dataTypeNames, "[]int")
					case "[]float":
						dataTypes = append(dataTypes, base.DType_F32Array)
						dataTypeNames = append(dataTypeNames, "[]float32")
					case "[]float64":
						dataTypes = append(dataTypes, base.DType_F64Array)
						dataTypeNames = append(dataTypeNames, "[]float64")
					case "[]int64":
						dataTypes = append(dataTypes, base.DType_S64Array)
						dataTypeNames = append(dataTypeNames, "[]int64")
					default:
						fmt.Printf("data [%s] [%s] col[%d] type not support in[string, enum, int8, int16, int32, float32, float64, []string, []int8, []int16, []int32, []float32, []float64]", filename, coltype, j )
						return
					}

					if j == len(row.Cells) -1 {
						stream.WriteString("package data\n")
						stream.WriteString("\n")
						stream.WriteString("type(\n")
						//定义数据类型结构体
						stream.WriteString(fmt.Sprintf("\t%s struct{\n", structName))
						for i1, v := range dataTypeNames{
							//过滤掉不是客户端的数据
							if dataNames[i1] == "" || dataNames[i1] == "0"{
								continue
							}

							stream.WriteString(fmt.Sprintf("\t\t%s\t%s\n", dataNames[i1], v))
						}
						stream.WriteString("\t}\n\n")
						//定义数据datares结构体
						stream.WriteString(fmt.Sprintf("\t%s struct{\n", structResName))
						stream.WriteString(fmt.Sprintf("\t\tm_DataMap map[int] *%s\n", structName))
						stream.WriteString("\t}\n")
						stream.WriteString(")\n")
						stream.WriteString("\n")
						//定义全局datares变量
						stream.WriteString("var(\n")
						stream.WriteString(fmt.Sprintf("\t%s\t%s\n", strings.ToUpper(structName), structResName))
						stream.WriteString(")\n")
						stream.WriteString("\n")
						//定义初始函数
						stream.WriteString(fmt.Sprintf("func (this *%s) Init(){\n", structResName))
						//定义map
						stream.WriteString(fmt.Sprintf("\tthis.m_DataMap =   map[int] *%s{}\n", structName))
						continue
					}else{
						continue
					}
				}

				//过滤掉不是客户端的数据
				if dataNames[j] == "" || dataNames[j] == "0"{
					continue
				}

				if j == dataColBeginLen{
					//map赋值
					stream.WriteString(fmt.Sprintf("\tthis.m_DataMap[%d] =  &%s{", base.Int(cell.String()), structName))
				}

				switch dataTypes[j] {
				case base.DType_String:
					stream.WriteString(fmt.Sprintf("\t\t%s : \"%s\",\n",dataNames[j], cell.Value))
				case base.DType_Enum:
					val, bEx := enumKVMap[j][strings.ToLower(cell.Value)]
					if bEx{
						stream.WriteString(fmt.Sprintf("\t\t%s : %d,\n",dataNames[j], val))
					}else{
						stream.WriteString(fmt.Sprintf("\t\t%s : %d,\n",dataNames[j], 0))
					}
				case base.DType_S8:
					stream.WriteString(fmt.Sprintf("\t\t%s : %d,\n",dataNames[j], base.Int(cell.Value)))
				case base.DType_S16:
					stream.WriteString(fmt.Sprintf("\t\t%s : %d,\n",dataNames[j], base.Int(cell.Value)))
				case base.DType_S32:
					stream.WriteString(fmt.Sprintf("\t\t%s : %d,\n",dataNames[j], base.Int(cell.Value)))
				case base.DType_F32:
					stream.WriteString(fmt.Sprintf("\t\t%s : %f,\n",dataNames[j], base.Float32(cell.Value)))
				case base.DType_F64:
					stream.WriteString(fmt.Sprintf("\t\t%s : %f,\n",dataNames[j], base.Float64(cell.Value)))
				case base.DType_S64:
					stream.WriteString(fmt.Sprintf("\t\t%s : %d,\n",dataNames[j], base.Int64(cell.Value)))

				case base.DType_StringArray:
					arr := strings.Split(cell.Value, ARRAY_SPLIT)
					stream.WriteString(fmt.Sprintf("\t\t%s : []string{", dataNames[j]))
					for i, v := range arr{
						stream.WriteString(fmt.Sprintf("\"%s\"",v))
						if i != len(arr)-1{
							stream.WriteString(",")
						}
					}
					stream.WriteString("},\n")
				case base.DType_S8Array:
					arr := strings.Split(cell.Value, ARRAY_SPLIT)
					stream.WriteString(fmt.Sprintf("\t\t%s : []int8{", dataNames[j]))
					for i, v := range arr{
						stream.WriteString(fmt.Sprintf("%d",base.Int(v)))
						if i != len(arr)-1{
							stream.WriteString(",")
						}
					}
					stream.WriteString("},\n")
				case base.DType_S16Array:
					arr := strings.Split(cell.Value, ARRAY_SPLIT)
					stream.WriteString(fmt.Sprintf("\t\t%s : []int16{", dataNames[j]))
					for i, v := range arr{
						stream.WriteString(fmt.Sprintf("%d",base.Int(v)))
						if i != len(arr)-1{
							stream.WriteString(",")
						}
					}
					stream.WriteString("},\n")
				case base.DType_S32Array:
					arr := strings.Split(cell.Value, ARRAY_SPLIT)
					stream.WriteString(fmt.Sprintf("\t\t%s : []int{", dataNames[j]))
					for i, v := range arr{
						stream.WriteString(fmt.Sprintf("%d",base.Int(v)))
						if i != len(arr)-1{
							stream.WriteString(",")
						}
					}
					stream.WriteString("},\n")
				case base.DType_F32Array:
					arr := strings.Split(cell.Value, ARRAY_SPLIT)
					stream.WriteString(fmt.Sprintf("\t\t%s : []float32{", dataNames[j]))
					for i, v := range arr{
						stream.WriteString(fmt.Sprintf("%f",base.Float32(v)))
						if i != len(arr)-1{
							stream.WriteString(",")
						}
					}
					stream.WriteString("},\n")
				case base.DType_F64Array:
					arr := strings.Split(cell.Value, ARRAY_SPLIT)
					stream.WriteString(fmt.Sprintf("\t\t%s : []float64{", dataNames[j]))
					for i, v := range arr{
						stream.WriteString(fmt.Sprintf("%f",base.Float64(v)))
						if i != len(arr)-1{
							stream.WriteString(",")
						}
					}
					stream.WriteString("},\n")
				case  base.DType_S64Array:
					arr := strings.Split(cell.Value, ARRAY_SPLIT)
					stream.WriteString(fmt.Sprintf("\t\t%s : []int64{", dataNames[j]))
					for i, v := range arr{
						stream.WriteString(fmt.Sprintf("%d",base.Int64(v)))
						if i != len(arr)-1{
							stream.WriteString(",")
						}
					}
					stream.WriteString("},\n")
				}

				if j == dataColLen{
					stream.WriteString("\t}\n")
				}
			}
		}

		//init函数执行完
		stream.WriteString("}\n")
		stream.WriteString("\n")
		//写获取函数
		stream.WriteString(fmt.Sprintf("func (this *%s) GetData(id int) *%s{\n", structResName, structName))
		stream.WriteString("\tpData, bOk := this.m_DataMap[id]\n")
		stream.WriteString("\tif bOk && pData != nil{\n")
		stream.WriteString("\t\treturn pData\n")
		stream.WriteString("\t}\n")
		stream.WriteString("\treturn nil\n")
		stream.WriteString("}\n")
	}

	//文件没有可导出
	if dataColLen == 0{
		return
	}

	//other sheet
	file, err := os.Create(filenames[0] + ".go")
	if err == nil{
		file.Write(stream.Bytes())
		file.Close()
	}
}


