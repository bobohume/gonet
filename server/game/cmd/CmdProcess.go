package cmd

import (
	"context"
	"fmt"
	"gonet/actor"
	"gonet/server/cm"
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

func (c *CmdProcess) Init() {
	c.Actor.Init()
	actor.MGR.RegisterActor(c)
	c.Actor.Start()
}

var (
	g_Cmd *CmdProcess
)

func Init() {
	g_Cmd = &CmdProcess{}
	g_Cmd.Init()
	cm.StartConsole(g_Cmd)
	//InitWeb()
}

func (c *CmdProcess) Cpus(ctx context.Context) {
	fmt.Println(runtime.NumCPU(), " cpus and ", runtime.GOMAXPROCS(0), " in use")
}

func (c *CmdProcess) Routines(ctx context.Context) {
	fmt.Println("Current number of goroutines: ", runtime.NumGoroutine())
}

func (c *CmdProcess) Setcpus(ctx context.Context, args string) {
	n, _ := strconv.Atoi(args)
	runtime.GOMAXPROCS(n)
	fmt.Println(runtime.NumCPU(), " cpus and ", runtime.GOMAXPROCS(0), " in use")
}

func (c *CmdProcess) Startgc(ctx context.Context) {
	runtime.GC()
	fmt.Println("gc finished")
}

func (c *CmdProcess) Showrpc(ctx context.Context) {
	fmt.Printf("--------------  PACKET  -------------\n")
	for i, v := range message.Packet_CrcNamesMap {
		fmt.Printf("packetName[%s], crc[%d]\n", v, i)
	}
	fmt.Printf("--------------  PACKET  -------------\n")
}

func (c *CmdProcess) Cpus1(ctx context.Context) {
	fmt.Println(runtime.NumCPU(), " cpus111 and ", runtime.GOMAXPROCS(0), " in use")
}
