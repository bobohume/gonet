package common

import (
	"fmt"
	"bufio"
	"os"
	"actor"
	"strings"
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
			fmt.Printf("parseConsole error [%s]", error.Error)
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
		pCmd.SendMsg(0, funcName, params...)
	}else{
		consoleError(command)
	}
}

func consoleroutine(pCmd actor.IActor) {
	command := make([]byte, 1024)
	reader := bufio.NewReader(os.Stdin)
	for {
		command, _, _ = reader.ReadLine()
		ParseConsole(pCmd, command)
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
