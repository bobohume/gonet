package common

import (
	"fmt"
	"bufio"
	"gonet/base"
	"gonet/rpc"
	"os"
	"gonet/actor"
	"strings"
	"time"
)

func StartConsole(pCmd actor.IActor) {
	go consoleroutine(pCmd)
}

func consoleError(command string){
	fmt.Printf("Command[%s] error, try again.", command)
}

 func ParseConsole(pCmd actor.IActor, command string) {
	defer func() {
		if err := recover(); err != nil{
			base.TraceCode(err)
		}
	}()

	if command == ""{
		return
	}

	args := strings.Split(command, "(")
	if len(args) != 2{
		consoleError(command)
		return
	}

	funcName := args[0]
	if funcName == ""{
		return
	}

	args = strings.Split(args[1], ")")
	if len(args) != 2{
		consoleError(command)
		return
	}

	args = strings.Split(args[0], ",")
	params := make([]interface{}, 0)
	for _,v := range args{
		if v != ""{
			params = append(params, v)
		}
	}

	if pCmd.HasRpc(funcName){
		pCmd.SendMsg(rpc.RpcHead{}, funcName, params...)
	}else{
		consoleError(command)
	}
}

//linux下面nohup &开启的时候当nohup文件占满磁盘的时候，这个就不会阻塞了，注意
func consoleroutine(pCmd actor.IActor) {
	for {
		reader := bufio.NewScanner(os.Stdin)
		for reader.Scan() {
			command := reader.Text()
			ParseConsole(pCmd, command)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

/*
传递指针类型的
pCmd := &CmdProcess{}
pCmd.Init(1)
funcName := common.StartConsole
var funcName1  *func(actor.IActor)
fmt.Println(funcName, funcName1)
ponit := unsafe.Pointer(&funcName)
func1 := (*func (actor.IActor))(unsafe.Pointer(ponit))
fmt.Println(func1, funcName)
(*func1)(pCmd)*/
