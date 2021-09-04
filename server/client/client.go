package main

import (
	"fmt"
	"gonet/base"
	"gonet/common"
	"gonet/network"
	"gonet/server/message"
	"os"
	"os/signal"
)

type (
	Config struct {
		common.Server `yaml:"netgate"`
	}
)

var (
	CONF   Config
	CLIENT *network.ClientSocket
)

func main() {
	message.InitClient()
	base.ReadConf("gonet.yaml", &CONF)
	CLIENT = new(network.ClientSocket)
	CLIENT.Init(CONF.Server.Ip, CONF.Server.Port)
	PACKET = new(EventProcess)
	PACKET.Init()
	CLIENT.BindPacketFunc(PACKET.PacketFunc)
	PACKET.Client = CLIENT
	if CLIENT.Start() {
		PACKET.LoginGate()
	}

	InitCmd()

	for i := 0; i < 1000; i++ {
		client := new(network.ClientSocket)
		client.Init(CONF.Server.Ip, CONF.Server.Port)
		packet := new(EventProcess)
		packet.Init()
		client.BindPacketFunc(packet.PacketFunc)
		packet.Client = client
		if client.Start() {
			packet.LoginGate()
		}
	}
	//PACKET.LoginGame()
	//for{
	//	PACKET.LoginGate()
	//}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	s := <-c
	fmt.Printf("client exit ------- signal:[%v]", s)
}
