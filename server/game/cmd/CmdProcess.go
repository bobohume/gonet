package cmd

import (
	"context"
	"fmt"
	"gonet/actor"
	"gonet/common"
	"gonet/server/message"
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

func (this *CmdProcess) Init() {
	this.Actor.Init()
	actor.MGR.RegisterActor(this)
	this.Actor.Start()
}

var (
	g_Cmd *CmdProcess
)

func Init() {
	g_Cmd = &CmdProcess{}
	g_Cmd.Init()
	common.StartConsole(g_Cmd)
	//InitWeb()
}

func (this *CmdProcess) Cpus(ctx context.Context) {
	fmt.Println(runtime.NumCPU(), " cpus and ", runtime.GOMAXPROCS(0), " in use")
}

func (this *CmdProcess) Routines(ctx context.Context) {
	fmt.Println("Current number of goroutines: ", runtime.NumGoroutine())
}

func (this *CmdProcess) Setcpus(ctx context.Context, args string) {
	n, _ := strconv.Atoi(args)
	runtime.GOMAXPROCS(n)
	fmt.Println(runtime.NumCPU(), " cpus and ", runtime.GOMAXPROCS(0), " in use")
}

func (this *CmdProcess) Startgc(ctx context.Context) {
	runtime.GC()
	fmt.Println("gc finished")
}

func (this *CmdProcess) Showrpc(ctx context.Context) {
	fmt.Printf("--------------  PACKET  -------------\n")
	for i, v := range message.Packet_CrcNamesMap {
		fmt.Printf("packetName[%s], crc[%d]\n", v, i)
	}
	fmt.Printf("--------------  PACKET  -------------\n")
}

func (this *CmdProcess) HotFix(ctx context.Context, name string) {
	fmt.Printf("--------------  PACKET  -------------\n")
}

func (this *CmdProcess) Cpus1(ctx context.Context) {
	fmt.Println(runtime.NumCPU(), " cpus111 and ", runtime.GOMAXPROCS(0), " in use")
}
