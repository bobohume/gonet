package base

import (
	"fmt"
	"runtime"
)

const (
	NONE_ERROR = iota
	VERSION_ERROR = iota				//版本不正确
	ACCOUNT_NOEXIST = iota			//账号不存在
)

//输出错误，跟踪代码
func TraceCode(code ...interface{}) {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	data := ""
	for _, v := range code{
		data += fmt.Sprintf("%v", v)
	}
	data += string(buf[:n])
	fmt.Printf("==> %s\n", data)
}