package main

import (
	"flag"
	"fmt"
	"golang.org/x/net/websocket"
	"time"
)

//========================================================================================
var addr = flag.String("addr", "localhost:12300", "http service address")

type(
	Test struct {
		F float32
	}

)

func main() {
	var MM map[int] int
	MM = make(map[int] int)
	for i := 0; i <100; i++{
		MM[i] = i
	}

	for i,_ := range MM{
		i ++
		//fmt.Println(i, v)
	}

	for i,v := range MM{
		fmt.Println(i, v)
	}


	flag.Parse()

	url := "ws://"+ *addr + "/ws"
	origin := "http://localhost/"
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		fmt.Println(err)
	}

	websocket.Message.Send(ws, "hello world")

	go timeWriter(ws)

	for {
		n := 1
		n++
		var msg [512]byte
		_, err := ws.Read(msg[:])//此处阻塞，等待有数据可读
		if err != nil {
			fmt.Println("read:", err)
			return
		}

		fmt.Printf("received: %s\n", msg)
	}
}

func timeWriter(conn *websocket.Conn) {
	for {
		time.Sleep(time.Second * 2)
		websocket.Message.Send(conn, "hello world")
	}
}

