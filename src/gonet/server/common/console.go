package common

import (
	"fmt"
	"bufio"
	"os"
	"gonet/actor"
	"strings"
	"time"
)

func StartConsole(pCmd actor.IActor) {
	go consoleroutine(pCmd)
}

func consoleError(buf []byte){
	fmt.Printf("Command[%s] error, try again.", string(buf))
}

 func ParseConsole(pCmd actor.IActor, command []byte) {
	defer func() {
		if err := recover(); err != nil{
			fmt.Printf("parseConsole error [%s]", err.(error).Error())
		}
	}()

	if string(command) == ""{
		return
	}

	args := strings.Split(string(command), "(")
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

	if pCmd.FindCall(funcName) != nil{
		pCmd.SendMsg(funcName, params...)
	}else{
		consoleError(command)
	}
}

//linux下面nohup &开启的时候当nohup文件占满磁盘的时候，这个就不会阻塞了，注意
func consoleroutine(pCmd actor.IActor) {
	command := make([]byte, 1024)
	for {
		reader := bufio.NewReader(os.Stdin)
		command, _, _ = reader.ReadLine()
		ParseConsole(pCmd, command)
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
