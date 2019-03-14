package main

import (
	"gonet/base"
	"fmt"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/websocket"
	"log"
	"gonet/message"
)

var(
)

func ExampleDial() {
	origin := "http://localhost/"
	url := "ws://localhost:31700/ws"
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}

	AccountName := fmt.Sprintf("test%d", 1)
	packet1 := &message.C_A_LoginRequest{PacketHead: message.BuildPacketHead(0, int(message.SERVICE_ACCOUNTSERVER)),
		AccountName: proto.String(AccountName), BuildNo: proto.String(base.BUILD_NO), SocketId: proto.Int32(0)}
	buff := message.Encode(packet1)
	buff = base.SetTcpEnd(buff)
	if _, err := ws.Write(buff); err != nil {
		log.Fatal(err)
	}

	for{
		var msg = make([]byte, 512)
		var n int
		if n, err = ws.Read(msg); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Received: %s.\n", msg[:n])
	}
}

func main() {
	ExampleDial()

	for{
		ttt := 0
		ttt++
	}
}
