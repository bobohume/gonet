package main

import (
	"bytes"
	"path/filepath"
	"strings"
)

var(
	PROTO [2][4]string = [2][4]string{
		{
			"protoc --plugin=protoc-gen-go=protoc-gen-go.exe  --go_out=../src/gonet/message  --proto_path=../src/gonet/message	",
			"::protoc --js_out=../src/gonet/message  --proto_path=../src/gonet/message	",
			"::protoc --cpp_out=../src/gonet/message/c++  --proto_path=../src/gonet/message	",
			"::protoc -o ../src/gonet/message/pb/client.pb --proto_path=../src/gonet/message	",
		},//win
		{
			"protoc --go_out=../src/gonet/message  --proto_path=../src/gonet/message	",
			"#protoc --js_out=../src/gonet/message  --proto_path=../src/gonet/message	",
			"#protoc --cpp_out=../src/gonet/message/c++  --proto_path=../src/gonet/message	",
			"#protoc -o ../src/gonet/message/pb/client.pb --proto_path=../src/gonet/message	",
		},//linux
	}
)

func main(){
	files, err := filepath.Glob("../src/gonet/message/*.proto")
	str := ""
	if err == nil{
		files1 := []string{}
		for _, v := range files{
			v = strings.Replace(v, "\\", "/", -1)
			if strings.LastIndex(v, "message.proto") != -1{
				str += v + "	"
				continue
			}
			files1 = append(files1, v)
		}

		for _, v := range files1{
			str += v + "	"
		}
	}

	index := strings.LastIndex(str, "	")
	if index!= -1{
		str = str[:index]
	}

	//生成bat文件
	{
		stream := bytes.NewBuffer([]byte{})
		file, err := os.Create("proto.bat")
		if err == nil{
			for _, v := range PROTO[0]{
				stream.WriteString(v)
				stream.WriteString(str)
				stream.WriteString("\n")
			}
			file.Write(stream.Bytes())
			file.Close()
		}
	}
	{
		stream := bytes.NewBuffer([]byte{})
		file, err := os.Create("proto.sh")
		if err == nil{
			for _, v := range PROTO[1]{
				stream.WriteString(v)
				stream.WriteString(str)
				stream.WriteString("\n")
			}
			file.Write(stream.Bytes())
			file.Close()
		}
	}
}
