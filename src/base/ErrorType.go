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
func TraceCode() {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	fmt.Printf("==> %s\n", string(buf[:n]))
}