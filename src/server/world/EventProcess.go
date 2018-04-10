package world

import (
	"actor"
	"database/sql"
	"message"
)

type (
	EventProcess struct {
		actor.Actor
		m_db *sql.DB
	}

	IEventProcess interface {
		actor.IActor
	}
)

func (this *EventProcess) Init(num int) {
	this.Actor.Init(num)
	this.m_db = SERVER.GetDB()
	this.RegisterCall("COMMON_RegisterRequest", func(caller *actor.Caller, nType int, GateId int, Ip string, Port int) {
		var ServerInfo stServerInfo
		ServerInfo.SocketId = caller.SocketId
		ServerInfo.Type = nType
		ServerInfo.GateId = GateId
		ServerInfo.Ip = Ip
		ServerInfo.Port = Port

		SERVER.GetServerMgr().SendMsg(caller.SocketId,"CONNECT", nType, GateId, Ip, Port)

		switch ServerInfo.Type {
		case int(message.SERVICE_GATESERVER):
			SERVER.GetServer().SendMsgByID(caller.SocketId, "COMMON_RegisterResponse")
		}
	})

	//断开链接
	this.RegisterCall("DISCONNECT", func(caller *actor.Caller, socketid int) {
		SERVER.GetServerMgr().SendMsg(caller.SocketId, "DISCONNECT", socketid)
	})

	this.Actor.Start()
}
