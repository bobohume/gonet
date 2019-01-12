package main

import (
	"base"
	"fmt"
	"os"
	"os/signal"
	"server/account"
	"server/common"
	"server/monitor"
	"server/netgate"
	"server/world"
)

func main() {
	args := os.Args
	base.RegisterMessage(&common.ServerInfo{})
	if args[1] == "account"{
		account.SERVER.Init()
	}else if args[1] == "monitor"{
		monitor.SERVER.Init()
	}else if args[1] == "netgate"{
		netgate.SERVER.Init()
	}else if args[1] == "world"{
		world.SERVER.Init()
	}

	base.SEVERNAME = args[1]
	
	InitMgr(args[1])

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	s := <-c
	fmt.Printf("server【%s】 exit ------- signal:[%v]", args[1], s)
}