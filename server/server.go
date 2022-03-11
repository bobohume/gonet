package main

import (
	"fmt"
	"gonet/base"
	"gonet/server/db"
	"gonet/server/game"
	"gonet/server/gate"
	"gonet/server/gm"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	args := os.Args
	if args[1] == "gm"{
		gm.SERVER.Init()
	}else if args[1] == "gate"{
		gate.SERVER.Init()
	}else if args[1] == "game"{
		game.SERVER.Init()
	} else if args[1]  == "db"{
		db.SERVER.Init()
	}

	base.SEVERNAME = args[1]
	
	InitMgr(args[1])

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	s := <-c

	ExitMgr(args[1])
	fmt.Printf("server【%s】 exit ------- signal:[%v]", args[1], s)
}