package main

import (
	"gonet/base"
	"fmt"
	"gonet/network"
	"os"
	"os/signal"
	"strconv"
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
	CLIENT = new(network.ClientSocket)
	CLIENT.Init(UserNetIP, port)
	PACKET = new(EventProcess)
	PACKET.Init(1)
	CLIENT.BindPacketFunc(PACKET.PacketFunc)
	PACKET.Client = CLIENT
	if CLIENT.Start(){
		PACKET.LoginAccount()
	}

	InitCmd()

	/*for i := 0; i < 10000; i++{
		client := new(network.ClientSocket)
		client.Init(UserNetIP, port)
		packet := new(EventProcess)
		packet.Init(1)
		client.BindPacketFunc(packet.PacketFunc)
		packet.Client = client
		if client.Start(){
			packet.LoginAccount()
		}
		//time.Sleep(1 * time.Microsecond)
	}*/
	//PACKET.LoginGame()
	//for{
	//	PACKET.LoginAccount()
	//}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	s := <-c
	fmt.Printf("client exit ------- signal:[%v]", s)
}