package main

import (
	"fmt"
	"gonet/base"
	"gonet/network"
	"gonet/server/message"
	"os"
	"os/signal"
	"strconv"
)

var (
	CLIENT *network.ClientSocket
)
func main() {
	message.InitClient()
	cfg := &base.Config{}
	cfg.Read("GONET_SERVER.CFG")
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
		PACKET.LoginGate()
	}

	InitCmd()
	//PACKET.LoginGame()
	//for{
	//	PACKET.LoginGate()
	//}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	s := <-c
	fmt.Printf("client exit ------- signal:[%v]", s)
}