package main

import (
	"context"
	"gonet/actor"
	"gonet/rpc"
	"gonet/server/cm"
	"gonet/server/message"
	"strconv"
)

type (
	CmdProcess struct {
		actor.Actor
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

func InitCmd() {
	g_Cmd = &CmdProcess{}
	g_Cmd.Init()
	cm.StartConsole(g_Cmd)
}

func (c *CmdProcess) Msg(ctx context.Context, args string) {
	packet1 := &message.ChatMessageRequest{PacketHead: message.BuildPacketHead(PACKET.PlayerId, rpc.SERVICE_GATE),
		Sender:      PACKET.PlayerId,
		Recver:      0,
		MessageType: int32(message.CHAT_MSG_TYPE_WORLD),
		Message:     (args),
	}
	SendPacket(packet1)
}

func (c *CmdProcess) Move(ctx context.Context, yaw string) {
	ya, _ := strconv.ParseFloat(yaw, 32)
	PACKET.Move(float32(ya), 100.0)
}
