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
	dataNames := []string{}
	stream := bytes.NewBuffer([]byte{})
	enmnMap := make(map[int] map[string] int)
	filenames := strings.Split(filename, ".")
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
		for i, row := range sheet.Rows {
			rowsLen := len(row.Cells) - 1
			for j, cell := range row.Cells {
				if i == 0 {
					dataNames = append(dataNames, cell.String())
					if j == 0{
						stream.WriteString(fmt.Sprintf("%s %s%s","local", filenames[0],"Data = {\n" ))
					}
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
						fmt.Errorf("[%s] col[%d] type not support in[string, enum, int8, int16, int32, float32, float64]", coltype, j )
						return
					}
					continue
				}

				writeInt := func() {
					switch cell.Type() {
					case xlsx.CellTypeString:
						stream.WriteString(fmt.Sprintf("\t\t%s = %d,\n",dataNames[j], base.Int(cell.String())))
					case xlsx.CellTypeFormula:
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

				if j == 0{
					stream.WriteString("\t{\n")
				}

				if dataTypes[j] == base.DType_String{
					switch cell.Type() {
					case xlsx.CellTypeString:
						stream.WriteString(fmt.Sprintf("\t\t%s = \"%s\",\n",dataNames[j], cell.String()))
					case xlsx.CellTypeFormula:
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
					val, bEx := enmnMap[j][cell.Value]
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
					case xlsx.CellTypeFormula:
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
					case xlsx.CellTypeFormula:
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

				if j == rowsLen{
					stream.WriteString("\t},\n")
				}
			}
		}

		stream.WriteString("}\n")
		stream.WriteString(fmt.Sprintf("%s %s%s\n","return", filenames[0],"Data" ))
	}

	/*stream.WriteString(fmt.Sprintf("%s %s%s","local", filenames[0],"DataName = {\n" ))
	for _, v := range dataNames{
		stream.WriteString(fmt.Sprintf("\t%s,\n", v))
	}
	stream.WriteString("}\n")
	stream.WriteString(fmt.Sprintf("%s %s%s\n","return", filenames[0],"DataName" ))*/
	//other sheet
	file, err := os.Create(filenames[0] + ".lua")
	if err == nil{
		file.Write(stream.Bytes())
		file.Close()
	}
}

