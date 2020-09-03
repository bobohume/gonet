package cmd

import (
	"context"
	"fmt"
	"gonet/actor"
	"gonet/message"
	"gonet/rpc"
	"gonet/server/common"
	"gonet/server/world/toprank"
	"runtime"
	"strconv"
)

type (
	CmdProcess struct {
		actor.Actor
	}

	ICmdProcess interface {
		actor.IActor
	}
)

func (this *CmdProcess) Init(num int) {
	this.Actor.Init(num)
	this.RegisterCall("cpus", func(ctx context.Context) {
		fmt.Println(runtime.NumCPU(), " cpus and ", runtime.GOMAXPROCS(0), " in use")
	})

	this.RegisterCall("routines", func(ctx context.Context) {
		fmt.Println("Current number of goroutines: ", runtime.NumGoroutine())
	})

	this.RegisterCall("setcpus", func(ctx context.Context, args string) {
		n, _ := strconv.Atoi(args)
		runtime.GOMAXPROCS(n)
		fmt.Println(runtime.NumCPU(), " cpus and ", runtime.GOMAXPROCS(0), " in use")
	})

	this.RegisterCall("startgc", func(ctx context.Context) {
		runtime.GC()
		fmt.Println("gc finished")
	})

	this.RegisterCall("InTopRank", func(ctx context.Context, argv0,argv1,argv2,argv3,argv4,argv5 string) {
		nType, _ := strconv.Atoi(argv0)
		id, _ := strconv.Atoi(argv1)
		name := argv2
		score, _ := strconv.Atoi(argv3)
		val0, _ := strconv.Atoi(argv4)
		val1, _ := strconv.Atoi(argv5)
		toprank.MGR().SendMsg( rpc.RpcHead{},"InTopRank", nType, int64(id), name, score, val0, val1)
	})

	this.RegisterCall("showrpc", func(ctx context.Context) {
		fmt.Printf("--------------  PACKET  -------------\n")
		for i, v := range message.Packet_CrcNamesMap{
			fmt.Printf("packetName[%s], crc[%d]\n", v, i)
		}
		fmt.Printf("--------------  PACKET  -------------\n")
	})

	this.Actor.Start()
}

var(
	g_Cmd *CmdProcess
)

func Init(){
	g_Cmd = &CmdProcess{}
	g_Cmd.Init(1000)
	common.StartConsole(g_Cmd)
	//InitWeb()
}
