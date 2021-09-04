package main

import (
	"bufio"
	"bytes"
	"fmt"
	"gonet/base"
	"os"
	"strings"

	"github.com/tealeg/xlsx"
)

func OpenExcePb(filename string) {
	xlFile, err := xlsx.OpenFile(filename)
	if err != nil {
		fmt.Println("open [%s] error", filename)
		return
	}

	dataTypes := []int{}
	dataTypeNames := []string{}
	stream := bytes.NewBuffer([]byte{})
	filenames := strings.Split(filename, ".")
	enumKVMap := make(map[int]map[string]int) //列 key val
	enumKMap := map[string][]string{}         //列名对应key
	enumNames := []string{}                   //列名
	dataNames := []string{}
	colNames := []string{}
	{
		sheet, bEx := xlFile.Sheet["Settings_Radio"]
		if bEx {
			for i, v := range sheet.Rows {
				for i1, v1 := range v.Cells {
					if v1.String() == "" {
						continue
					}
					if i == 0 {
						enumNames = append(enumNames, v1.String())
					} else {
						enumKMap[enumNames[i1]] = append(enumKMap[enumNames[i1]], v1.String())
					}
				}
			}
		}
	}
	for page, sheet := range xlFile.Sheets {
		if page != 0 {
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
					colNames = append(colNames, cell.String())
					continue
				} else if i == COL_CLIENT_NAME {
					colName := cell.String()
					dataNames = append(dataNames, colName)
					//写proto
					if j == 0 {
						stream.WriteString("syntax = \"proto3\";\n")
						stream.WriteString("package message;\n")
						stream.WriteString("\n")
					}

					continue
				} else if i == COL_TYPE {
					coltype := strings.ToLower(cell.String())
					rd := bufio.NewReader(strings.NewReader(coltype))
					data, _, _ := rd.ReadLine()
					coltype = strings.TrimSpace(string(data))
					if coltype == "enum" {
						num := 0
						KVMap := map[string]int{}
						for data, _, _ := rd.ReadLine(); data != nil; {
							slot := strings.Split(string(data), "=")
							if len(slot) == 2 {
								KVMap[slot[0]] = base.Int(slot[1])
							}
							data, _, _ = rd.ReadLine()
						}
						keys, bEx := enumKMap[colNames[j]]
						if bEx {
							_, bEx := enumKVMap[j]
							if !bEx {
								enumKVMap[j] = make(map[string]int)
							}
							for _, v := range keys {
								val, bEx := KVMap[v]
								if bEx {
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
						dataTypeNames = append(dataTypeNames, "int32")
					case "int8":
						dataTypes = append(dataTypes, base.DType_S8)
						dataTypeNames = append(dataTypeNames, "int32")
					case "int16":
						dataTypes = append(dataTypes, base.DTypelinkClearS16)
						dataTypeNames = append(dataTypeNames, "int32")
					case "int":
						dataTypes = append(dataTypes, base.DType_S32)
						dataTypeNames = append(dataTypeNames, "int32")
					case "float":
						dataTypes = append(dataTypes, base.DType_F32)
						dataTypeNames = append(dataTypeNames, "float")
					case "float64":
						dataTypes = append(dataTypes, base.DType_F64)
						dataTypeNames = append(dataTypeNames, "double")
					case "int64":
						dataTypes = append(dataTypes, base.DType_S64)
						dataTypeNames = append(dataTypeNames, "int64")
					case "[]string":
						dataTypes = append(dataTypes, base.DType_StringArray)
						dataTypeNames = append(dataTypeNames, "repeated string")
					case "[]int8":
						dataTypes = append(dataTypes, base.DType_S8Array)
						dataTypeNames = append(dataTypeNames, "repeated int32")
					case "[]int16":
						dataTypes = append(dataTypes, base.DType_S16Array)
						dataTypeNames = append(dataTypeNames, "repeated int32")
					case "[]int":
						dataTypes = append(dataTypes, base.DType_S32Array)
						dataTypeNames = append(dataTypeNames, "repeated int32")
					case "[]float":
						dataTypes = append(dataTypes, base.DType_F32Array)
						dataTypeNames = append(dataTypeNames, "repeated float")
					case "[]float64":
						dataTypes = append(dataTypes, base.DType_F64Array)
						dataTypeNames = append(dataTypeNames, "repeated double")
					case "[]int64":
						dataTypes = append(dataTypes, base.DType_S64Array)
						dataTypeNames = append(dataTypeNames, "repeated int64")
					default:
						fmt.Printf("data [%s] [%s] col[%d] type not support in[string, enum, int8, int16, int32, float32, float64, []string, []int8, []int16, []int32, []float32, []float64]", filename, coltype, j)
						return
					}
					continue
				}

				//读取excel头部文件
				{
					//basedata
					{
						stream.WriteString("message ")
						stream.WriteString(filenames[0])
						stream.WriteString("Data\n")
						stream.WriteString("{\n")
						id := 1
						for i1, v := range dataTypeNames {
							//过滤掉不是客户端的数据
							if dataNames[i1] == "" || dataNames[i1] == "0" {
								continue
							}

							stream.WriteString(fmt.Sprintf("\t%s\t%s = %d;//%s\n", v, dataNames[i1], id, colNames[i1]))
							id++
						}
						stream.WriteString("}\n\n")
					}

					//mgr
					{
						stream.WriteString("message ")
						stream.WriteString(filenames[0])
						stream.WriteString("DataMgr\n")
						stream.WriteString("{\n")
						stream.WriteString("\trepeated int64 Keys = 1;\n")
						stream.WriteString(fmt.Sprintf("\trepeated %sData Items = 2;\n", filenames[0]))
						stream.WriteString(fmt.Sprintf("\tmap<int64, %sData> ItemsMap = 3;\n", filenames[0]))
						stream.WriteString("}\n")
					}

					//other sheet
					file, err := os.Create(filenames[0] + ".proto")
					if err == nil {
						file.Write(stream.Bytes())
						file.Close()
					}

					return
				}
			}
		}
	}
}
