package main

import (
	"bytes"
	"fmt"
	"gonet/base"
	"os"
	"strings"

	"github.com/tealeg/xlsx"
)

const (
	FILE_GENERATE = `package data

// 自动生成代码
type(
	{StructName} struct{
{StructBody}
	}

	{MgrName} struct{
		{DataMap} map[int] *{StructName}
	}
)

var(
	{MgrHandleName} {MgrName} 
)

`
)

func OpenExceGo(filename string) {
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
		dataTypeNames := []string{}
		dataColLen := 0       //结束列数
		dataColBeginLen := -1 //开始列数
		stream := bytes.NewBuffer([]byte{})
		structName := FILENAME(filename, sheet.Name, "Data")
		structResName := FILENAME(filename, sheet.Name, "DataRes")
		enumKVMap := make(map[int]map[string]int) //列 key val
		enumKMap := map[int][]string{}            //列名对应key
		dataNames := []string{}

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
				if i == COL_NAME {
					continue
				} else if i == COL_CLIENT_NAME {
					colName := cell.String()
					dataNames = append(dataNames, colName)
					if colName != "" && colName != "0" {
						dataColLen = j
						if dataColBeginLen == -1 {
							dataColBeginLen = j
						}
					}
					continue
				} else if i == COL_VSTO {
					enumNames := strings.Split(cell.String(), "\n")
					for _, v1 := range enumNames {
						enumKMap[j] = append(enumKMap[j], v1)
					}
					continue
				} else if i == COL_TYPE {
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
						fmt.Printf("data [%s] [%s] col[%d] type not support in[string, enum, int8, int16, int32, float32, float64, []string, []int8, []int16, []int32, []float32, []float64]", filename, coltype, j)
						return
					}

					if j == dataColLen {
						//定义数据类型结构体
						structBody := ""
						for i1, v := range dataTypeNames {
							//过滤掉不是客户端的数据
							if dataNames[i1] == "" || dataNames[i1] == "0" {
								continue
							}

							structBody += fmt.Sprintf("\t\t%s\t%s\n", dataNames[i1], v)
						}
						structBody = structBody[:len(structBody)-1]
						str := FILE_GENERATE
						str = strings.Replace(str, "{StructName}", structName, -1)
						str = strings.Replace(str, "{StructBody}", structBody, -1)
						str = strings.Replace(str, "{DataMap}", "m_DataMap", -1)
						str = strings.Replace(str, "{MgrName}", structResName, -1)
						str = strings.Replace(str, "{MgrHandleName}", strings.ToUpper(structResName), -1)
						stream.WriteString(str)
						//定义初始函数
						stream.WriteString(fmt.Sprintf("func (this *%s) Init(){\n", structResName))
						//定义map
						stream.WriteString(fmt.Sprintf("\tthis.m_DataMap =   map[int] *%s{}\n", structName))
						continue
					} else {
						continue
					}
				}

				//过滤掉不是客户端的数据
				if dataNames[j] == "" || dataNames[j] == "0" {
					continue
				}

				if j == dataColBeginLen {
					//map赋值
					stream.WriteString(fmt.Sprintf("\tthis.m_DataMap[%d] =  &%s{\n", base.Int(cell.String()), structName))
				}

				switch dataTypes[j] {
				case base.DType_String:
					stream.WriteString(fmt.Sprintf("\t\t%s : \"%s\",\n", dataNames[j], cell.Value))
				case base.DType_Enum:
					val, bEx := enumKVMap[j][strings.ToLower(cell.Value)]
					if bEx {
						stream.WriteString(fmt.Sprintf("\t\t%s : %d,\n", dataNames[j], val))
					} else {
						stream.WriteString(fmt.Sprintf("\t\t%s : %d,\n", dataNames[j], 0))
					}
				case base.DType_S8:
					stream.WriteString(fmt.Sprintf("\t\t%s : %d,\n", dataNames[j], base.Int(cell.Value)))
				case base.DType_S16:
					stream.WriteString(fmt.Sprintf("\t\t%s : %d,\n", dataNames[j], base.Int(cell.Value)))
				case base.DType_S32:
					stream.WriteString(fmt.Sprintf("\t\t%s : %d,\n", dataNames[j], base.Int(cell.Value)))
				case base.DType_F32:
					stream.WriteString(fmt.Sprintf("\t\t%s : %f,\n", dataNames[j], base.Float32(cell.Value)))
				case base.DType_F64:
					stream.WriteString(fmt.Sprintf("\t\t%s : %f,\n", dataNames[j], base.Float64(cell.Value)))
				case base.DType_S64:
					stream.WriteString(fmt.Sprintf("\t\t%s : %d,\n", dataNames[j], base.Int64(cell.Value)))

				case base.DType_StringArray:
					arr := splitArray(cell.Value)
					stream.WriteString(fmt.Sprintf("\t\t%s : []string{", dataNames[j]))
					for i, v := range arr {
						stream.WriteString(fmt.Sprintf("\"%s\"", v))
						if i != len(arr)-1 {
							stream.WriteString(",")
						}
					}
					stream.WriteString("},\n")
				case base.DType_S8Array:
					arr := splitArray(cell.Value)
					stream.WriteString(fmt.Sprintf("\t\t%s : []int8{", dataNames[j]))
					for i, v := range arr {
						stream.WriteString(fmt.Sprintf("%d", base.Int(v)))
						if i != len(arr)-1 {
							stream.WriteString(",")
						}
					}
					stream.WriteString("},\n")
				case base.DType_S16Array:
					arr := splitArray(cell.Value)
					stream.WriteString(fmt.Sprintf("\t\t%s : []int16{", dataNames[j]))
					for i, v := range arr {
						stream.WriteString(fmt.Sprintf("%d", base.Int(v)))
						if i != len(arr)-1 {
							stream.WriteString(",")
						}
					}
					stream.WriteString("},\n")
				case base.DType_S32Array:
					arr := splitArray(cell.Value)
					stream.WriteString(fmt.Sprintf("\t\t%s : []int{", dataNames[j]))
					for i, v := range arr {
						stream.WriteString(fmt.Sprintf("%d", base.Int(v)))
						if i != len(arr)-1 {
							stream.WriteString(",")
						}
					}
					stream.WriteString("},\n")
				case base.DType_F32Array:
					arr := splitArray(cell.Value)
					stream.WriteString(fmt.Sprintf("\t\t%s : []float32{", dataNames[j]))
					for i, v := range arr {
						stream.WriteString(fmt.Sprintf("%f", base.Float32(v)))
						if i != len(arr)-1 {
							stream.WriteString(",")
						}
					}
					stream.WriteString("},\n")
				case base.DType_F64Array:
					arr := splitArray(cell.Value)
					stream.WriteString(fmt.Sprintf("\t\t%s : []float64{", dataNames[j]))
					for i, v := range arr {
						stream.WriteString(fmt.Sprintf("%f", base.Float64(v)))
						if i != len(arr)-1 {
							stream.WriteString(",")
						}
					}
					stream.WriteString("},\n")
				case base.DType_S64Array:
					arr := splitArray(cell.Value)
					stream.WriteString(fmt.Sprintf("\t\t%s : []int64{", dataNames[j]))
					for i, v := range arr {
						stream.WriteString(fmt.Sprintf("%d", base.Int64(v)))
						if i != len(arr)-1 {
							stream.WriteString(",")
						}
					}
					stream.WriteString("},\n")
				}

				if j == dataColLen {
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

		//文件没有可导出
		if dataColLen == 0 {
			return
		}

		//other sheet
		file, err := os.Create(FILENAME(filename, sheet.Name, ".go"))
		if err == nil {
			file.Write(stream.Bytes())
			file.Close()
		}
	}
}
