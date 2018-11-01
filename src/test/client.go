package main

import (
	"network"
	"time"
	"message"
	"github.com/golang/protobuf/proto"
	"base"
)

func main() {
	client := new(network.ClientSocket)
	client.Init("192.168.1.22", 21000)
	client.Start()
	packet := &message.String2Packet{PacketHead:message.BuildPacketHead(int(message.SERVICE_ACCOUNTSERVER), 0),
	String1:proto.String("test11"), String2: proto.String(base.BUILD_NO)}
	buf, _ := proto.Marshal(packet)
	client.Send(buf)

	for {
		time.Sleep(1000)
	}
}
