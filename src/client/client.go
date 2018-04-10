package main

import (
	"network"
	"time"
	"fmt"
)

var (
	CLIENT *network.ClientSocket
)
func main() {
	var packet1 *EventProcess
	n, n1 := 0, 0
	for i:= 0; i < 1; i++{
		CLIENT = new(network.ClientSocket)
		CLIENT.Init("192.168.1.22", 31700)
		packet := new(EventProcess)
		packet.Init(1)
		CLIENT.BindPacketFunc(packet.PacketFunc)
		CLIENT.Start()
		packet.LoginAccount()
		packet.LoginAccount()
		packet.LoginAccount()
		packet1 = packet
	}

	for {
		time.Sleep(1000)
		packet1.LoginAccount()
		packet1.LoginAccount()
		packet1.LoginAccount()
		n++
		if n % 100 == 0 {
			n1++
			fmt.Println("已经运行[", n1, "]" )
		}
	}
}