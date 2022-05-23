package ini

// read ini tool
import (
	"bufio"
	"fmt"
	"gonet/base"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	//"encoding/base64"
)

const (
	STATE_NONE    = iota
	STATE_SECTION = iota
	STATE_VALUE   = iota

	MAX_LINE_LENGTH  = 2048
	BUFFER_LENGTH    = 512
	MAX_TOKEN_LENGTH = 128
	MAX_HOSTNAME_LEN = 128
	DEFAULT_CONFIG   = 0
)

type (
	CfgKey struct {
		first  string
		second int
	}

	SectionInfo map[string]string
	CfgInfo     map[CfgKey]SectionInfo

	Config struct {
		cfgInfo  CfgInfo
		config   int
		filePath string
	}

	ICfonfig interface {
		Read(string)
		Get(key string) string                                    //获取key
		Get2(key string, sep string) (string, string)             //获取ip
		Get3(section string, key string, secitonId ...int) string //根据section, key, sectionid(从0开始)
		Get5(key string, sep string) []string                     //获取数组
		Get6(section string, key string, sep string) []string     //获取数组
		Int(key string) int
		Int64(key string) int64
		Float32(key string) float32
		Float64(key string) float64
		Bool(key string) bool
		Time(key string) int64
	}
)

func Token(srcBuffer []byte, begin int, end int, toLower bool) (string, int) {
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

func (c *Config) Get(key string) string {
	//key = strings.ToLower(key)
	for _, map1 := range c.cfgInfo {
		val, bEx := map1[key]
		if bEx == true {
			return val
		}
	}

	return ""
}

func (c *Config) Get2(key string, sep string) (string, string) {
	split := func(buf string, sep string) (string, string) {
		index := strings.Index(buf, sep)
		first := buf[:index]
		second := buf[index+1:]
		return first, second
	}
	ip, port := split(c.Get(key), sep)
	if ip == "0.0.0.0" {
		addrs, _ := net.InterfaceAddrs()
		for _, address := range addrs {
			if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					ip = ipnet.IP.String()
					return ip, port
				}
			}
		}
	}
	return ip, port
}

func (c *Config) Get3(seciton string, key string, sectionid ...int) string {
	//key = strings.ToLower(key)
	id := 0
	if len(sectionid) >= 1 {
		id = sectionid[0]
	}
	map1, bEx := c.cfgInfo[CfgKey{seciton, id}]
	if bEx {
		val, exist := map1[key]
		if exist == true {
			return val
		}
	}

	return ""
}

func (c *Config) Get5(key string, sep string) []string {
	return strings.Split(c.Get(key), sep)
}

func (c *Config) Get6(section string, key string, sep string) []string {
	return strings.Split(c.Get3(section, key), sep)
}

func (c *Config) Int(key string) int {
	n, _ := strconv.Atoi(c.Get(key))
	return n
}

func (c *Config) Int64(key string) int64 {
	n, _ := strconv.ParseInt(c.Get(key), 0, 64)
	return n
}

func (c *Config) Float32(key string) float32 {
	n, _ := strconv.ParseFloat(c.Get(key), 32)
	return float32(n)
}

func (c *Config) Float64(key string) float64 {
	n, _ := strconv.ParseFloat(c.Get(key), 64)
	return n
}

func (c *Config) Bool(key string) bool {
	n, _ := strconv.ParseBool(c.Get(key))
	return n
}

func (c *Config) Time(key string) int64 {
	return base.GetDBTime(c.Get(key)).Unix()
}

func (c *Config) Read(path string) {
	c.cfgInfo = make(map[CfgKey]SectionInfo)
	for i, _ := range c.cfgInfo {
		delete(c.cfgInfo, i)
	}

	if c.filePath == "" {
		c.filePath = path
	} else {
		path = c.filePath
	}

	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("read cof error %v", err)
		return
	}

	defer file.Close()
	fileIn := bufio.NewReader(file)
	section := ""
	secCount := make(map[string]int)

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

		InsertMap := func() {
			_, exist := c.cfgInfo[CfgKey{section, secCount[section]}]
			if exist == true {
				c.cfgInfo[CfgKey{section, secCount[section]}][key], tokenBegin = Token(buffer, tokenBegin, i, false)
			} else {
				secotionMap := SectionInfo{}
				secotionMap[key], tokenBegin = Token(buffer, tokenBegin, i, false)
				c.cfgInfo[CfgKey{section, secCount[section]}] = secotionMap
			}
		}

		for i < nlen && !comment {
			switch buffer[i] {
			case '[':
				if state == STATE_NONE {
					tokenBegin = i + 1
					state = STATE_SECTION
				}
			case ']':
				if state == STATE_SECTION {
					section, tokenBegin = Token(buffer, tokenBegin, i, false)
					if section != "" {
						_, bEx := secCount[section]
						if !bEx {
							secCount[section] = 0
						} else {
							secCount[section]++
						}
						c.cfgInfo[CfgKey{section, secCount[section]}] = SectionInfo{}
						state = STATE_NONE
					}
				}
			case '=':
				if state == STATE_NONE {
					key, tokenBegin = Token(buffer, tokenBegin, i, true)
					if key != "" {
						state = STATE_VALUE
					}
				}
			case ';':
				if state == STATE_VALUE {
					if section != "" {
						InsertMap()
					}
					state = STATE_NONE
				}
			case '#': //注释模块
				if state == STATE_VALUE {
					if section != "" {
						//fmt.Println("111111", section)
						InsertMap()
						comment = true
						state = STATE_NONE
					}
				}
				/*case '/':
				if (i>1 && buffer[i-1]=='/' && state==STATE_VALUE) {
					if (section != ""){
						//fmt.Println("111111", section)
						InsertMap()
						comment = true;
						state = STATE_NONE;
					}
				}*/
			}
			i++
		}

		if state == STATE_VALUE {
			if section != "" {
				InsertMap()
			}
			state = STATE_NONE
		}
	}
	//fmt.Println(c.cfgInfo)
}
