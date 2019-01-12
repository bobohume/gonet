package world

import (
	"actor"
	"database/sql"
	"message"
	"server/common"
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
	this.RegisterCall("COMMON_RegisterRequest", func(nType int, Ip string, Port int) {
		pServerInfo := new(common.ServerInfo)
		pServerInfo.SocketId = this.GetSocketId()
		pServerInfo.Type = nType
		pServerInfo.Ip = Ip
		pServerInfo.Port = Port
		SERVER.GetServerMgr().SendMsg("CONNECT", pServerInfo)
		switch pServerInfo.Type {
		case int(message.SERVICE_GATESERVER):
			SERVER.GetServer().SendMsgByID(this.GetSocketId(), "COMMON_RegisterResponse")
		}
	})

	//断开链接
	this.RegisterCall("DISCONNECT", func(socketid int) {
		SERVER.GetServerMgr().SendMsg("DISCONNECT", socketid)
	})

	this.Actor.Start()
}
