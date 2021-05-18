package rpc

import (
	"reflect"
)

//params[0]:rpc.RpcHead
//params[1]:error
func Call(parmas ...interface{}) {
	/*head := *parmas[0].(*RpcHead)
	if parmas[1] == nil{
		parmas[1] = ""
	}else{
		parmas[1] = parmas[1].(error).Error()
	}*/
}

var GCall = reflect.ValueOf(Call)