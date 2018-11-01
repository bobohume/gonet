package main

import (
	"base"
	"fmt"
	"network"
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
	PACKET := new(EventProcess)
	PACKET.Init(1)
	CLIENT.BindPacketFunc(PACKET.PacketFunc)
	if CLIENT.Start(){
		PACKET.LoginAccount()
	}

	InitCmd()
	//PACKET.LoginGame()
	//for{
	//	PACKET.LoginAccount()
	//}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	s := <-c
	fmt.Printf("client exit ------- signal:[%v]", s)
}