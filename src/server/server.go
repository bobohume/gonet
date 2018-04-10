package main

import (
	"server/account"
	"server/netgate"
	"time"
	"os"
	"server/world"
	"base"
)

func main() {
	args := os.Args
	if args[1] == "account"{
		account.SERVER.Init()
	}else if args[1] == "netgate"{
		netgate.SERVER.Init()
	} else if args[1] == "world"{
		world.SERVER.Init()
	}

	base.SEVERNAME = args[1]
	
	InitMgr(args[1])

	for{
		/*if args[1] == "netgate"{
			netgate.SERVER.GetAccountScoket().SendMsg("11111")
		}*/
		time.Sleep(10000)
	}
}