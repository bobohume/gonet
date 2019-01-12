package base

import (
	"errors"
	"fmt"
	"runtime"
)

const (
	NONE_ERROR = iota
	VERSION_ERROR = iota				//版本不正确
	ACCOUNT_NOEXIST = iota			//账号不存在
)

//输出错误，跟踪代码
func TraceCode(skip int) error {
	funcName,file,line,ok := runtime.Caller(skip)
	if(ok){
		return errors.New(fmt.Sprintf("[file:%s],[line:%d],[func:%s]", file, line, runtime.FuncForPC(funcName).Name()))
	}
	return errors.New("")
}