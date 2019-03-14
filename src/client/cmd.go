package main

import (
	"gonet/actor"
	"github.com/golang/protobuf/proto"
	"gonet/message"
	"gonet/server/common"
)

type (
	CmdProcess struct {
		actor.Actor
	}

	ICmdProcess interface {
		actor.IActor
	}
)

func (this *CmdProcess) Init(num int) {
	this.Actor.Init(num)
	this.RegisterCall("msg", func(args string) {
		packet1 := &message.C_W_ChatMessage{PacketHead:message.BuildPacketHead( PACKET.AccountId, int(message.SERVICE_WORLDSERVER)),
			Sender:proto.Int64(PACKET.PlayerId),
			Recver:proto.Int64(0),
			MessageType:proto.Int32(int32(message.CHAT_MSG_TYPE_WORLD)),
			Message:proto.String(args),
		}
		SendPacket(packet1)
	})

	this.Actor.Start()
}

var(
	g_Cmd *CmdProcess
)

func InitCmd(){
	g_Cmd = &CmdProcess{}
	g_Cmd.Init(1000)
	common.StartConsole(g_Cmd)
}

