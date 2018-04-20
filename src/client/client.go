package main

import (
	"network"
	"fmt"
	"strconv"
	"base"
)

var (
	CLIENT *network.ClientSocket
)
func main() {
	cfg := &base.Config{}
	cfg.Read("SXZ_SERVER.CFG")
	UserNetIP, UserNetPort := cfg.Get2("NetGate_WANAddress", ":")
	//UserNetIP, UserNetPort := "101.132.178.159", "31700"
	port,_ := strconv.Atoi(UserNetPort)
	var packet1 *EventProcess
	n, n1 := 0, 0
	for i:= 0; i < 1; i++{
		CLIENT = new(network.ClientSocket)
		CLIENT.Init(UserNetIP, port)
		packet := new(EventProcess)
		packet.Init(1)
		CLIENT.BindPacketFunc(packet.PacketFunc)
		CLIENT.Start()
		packet1 = packet
	}

	for {
		//time.Sleep(1000)
		packet1.LoginAccount()
		n++
		if n % 100 == 0 {
			n1++
			fmt.Println("已经运行[", n1 * 100, "]" )
		}
	}
}