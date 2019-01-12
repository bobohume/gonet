package base
// read ini tool
import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	//"encoding/base64"
)

const(
	STATE_NONE = iota
	STATE_SECTION = iota
	STATE_VALUE = iota

	MAX_LINE_LENGTH = 2048
	BUFFER_LENGTH = 512
	MAX_TOKEN_LENGTH = 128
	MAX_HOSTNAME_LEN = 128
	DEFAULT_CONFIG = 0
)

type(
	CfgKey struct{
		first string
		second int
	}

	SectionInfo map[string] string
	CfgInfo map[CfgKey] SectionInfo

	Config struct{
		m_cfgInfo CfgInfo
		m_config int
		m_filePath string
	}

	ICfonfig interface {
		Read(string)
		Get(string) string
		Get2(string, string)(string, string)
		Int(key string) int
		Int64(key string) int64
		Float32(key string) float32
		Float64(key string) float64
		Bool(key string) bool
		Time(key string) int64
	}
)

func Token(srcBuffer []byte, begin int, end int, toLower bool) (string, int){
	//buffer := make([]byte, len(srcBuffer))
	nlen := end - begin
	token := make([]byte, nlen)
	copy(token, srcBuffer[begin:begin+nlen])
	begin = end + 1
	//str :=strings.ToLower(string(token))
	str := string(token)
	str = strings.TrimSpace(str)
	return str, begin
}

func (this *Config) Get(key string) string{
	//key = strings.ToLower(key)
	for _,map1 := range this.m_cfgInfo{
		val,exist := map1[key];
		if (exist == true){
			return val;
		}
	}

	return "";
}

func (this *Config) Int(key string) int{
	n, _ := strconv.Atoi(this.Get(key))
	return n
}

func (this *Config) Int64(key string) int64{
	n, _ := strconv.ParseInt(this.Get(key), 0, 64)
	return n
}

func (this *Config) Float32(key string) float32{
	n, _ := strconv.ParseFloat(this.Get(key), 32)
	return float32(n)
}

func (this *Config) Float64(key string) float64{
	n, _ := strconv.ParseFloat(this.Get(key), 64)
	return n
}

func (this *Config) Bool(key string) bool{
	n, _ := strconv.ParseBool(this.Get(key))
	return n
}

func (this *Config) Time(key string) int64 {
	return GetDBTime(this.Get(key)).Unix()
}

func (this *Config) Get2(key string, sep string) (string, string){
	split := func(buf string, sep string)(string, string) {
		index := strings.Index(buf, sep)
		first := buf[:index]
		second := buf[index+1:]
		return  first,second
	}
	return  split(this.Get(key), sep)
}

func (this *Config) Read(path string)  {
	this.m_cfgInfo = make(map [CfgKey] SectionInfo)
	for i,_ := range this.m_cfgInfo{
		delete(this.m_cfgInfo, i)
	}

	if (this.m_filePath == ""){
		this.m_filePath = path
	}else{
		path = this.m_filePath
	}

	file, err := os.Open(path)
	if err != nil {
		fmt.Print("read cof error %s", err)
		return
	}

	defer file.Close()
	fileIn := bufio.NewReader(file)
	section := ""
	secCount := make(map[string] int)

	for {
		line, _, err := fileIn.ReadLine()
		//buffer1 := make([]byte, len(line)*2)
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		//nlen ,err := base64.StdEncoding.Decode(buffer1, line)
		nlen := len(line)
		//if err != nil {
		//	panic(err)
		//}

		buffer := line[:]
		state := STATE_NONE
		comment := false
		i := 0
		tokenBegin := 0
		key := ""

		InsertMap := func(){
			_,exist := this.m_cfgInfo[CfgKey{section, secCount[section]-1}]
			if exist == true{
				this.m_cfgInfo[CfgKey{section, secCount[section] - 1}][key], tokenBegin = Token(buffer, tokenBegin, i, false)
			}else{
				secotionMap := SectionInfo{}
				secotionMap[key] , tokenBegin = Token(buffer, tokenBegin, i, false)
				this.m_cfgInfo[CfgKey{section, secCount[section] - 1}] = secotionMap
			}
		}

		for i < nlen && !comment {
			switch  buffer[i]{
			case '[':
				if state == STATE_NONE{
					tokenBegin = i + 1
					state = STATE_SECTION
				}
			case ']':
				if state == STATE_SECTION{
					section, tokenBegin = Token(buffer, tokenBegin, i, false)
					if section != ""{
						secCount[section]++
						this.m_cfgInfo[CfgKey{section, secCount[section]}] =  SectionInfo{}
						state = STATE_NONE
					}
				}
			case '=':
				if state == STATE_NONE{
					key, tokenBegin = Token(buffer, tokenBegin, i,true);
					if key != ""{
						state = STATE_VALUE;
					}
				}
			case ';':
				if state == STATE_VALUE{
					if section != "" {
						InsertMap()
					}
					state = STATE_NONE;
				}
			case '/':
				if (i>1 && buffer[i-1]=='/' && state==STATE_VALUE) {
					if (section != ""){
						//fmt.Println("111111", section)
						InsertMap()
						comment = true;
						state = STATE_NONE;
					}
				}
			}
			i++
		}

		if (state == STATE_VALUE) {
			if (section != "") {
				InsertMap()
			}
			state = STATE_NONE;
		}
	}

	//fmt.Println(this.m_cfgInfo)
}
